package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// PlannerInput — всё, что нужно чистому планировщику. Никакой БД внутри:
// данные поставляет оркестрация (service/plan.go).
type PlannerInput struct {
	Request        model.PlanningRequest                // запрос Java
	WorkCenters    []model.WorkCenter                   // активные станки (с operationTypes)
	WeeklySchedule map[int64][]model.ScheduleDay        // workCenterID -> недельный шаблон
	Busy           map[int64][]model.ScheduledOperation // workCenterID -> занятость активных планов
	Timezone       *time.Location                       // cfg.Timezone (Europe/Moscow)
}

// PlanDraft — результат планирования до записи в БД: ID проставит
// репозиторий при вставке. Plan.PlannedStartAt / PlannedFinishAt / OnTime /
// Status=PROPOSED заполнены здесь.
type PlanDraft struct {
	Plan       model.ProductionPlan
	Items      []model.ProductionPlanItem
	Operations []model.ScheduledOperation // без ID и ProductionPlanItemID до вставки

	// OperationItemIdx[i] — индекс позиции в Items, к которой относится
	// Operations[i]. Оркестрация использует его, чтобы после INSERT позиций
	// проставить ProductionPlanItemID у кусков (ID появляется только в БД).
	OperationItemIdx []int
}

// timeWindow — полуинтервал [start, end) в локальном времени Timezone.
type timeWindow struct {
	start, end time.Time
}

// plannerState — рабочие индексы, собранные один раз на вызов BuildPlan.
type plannerState struct {
	tz       *time.Location
	wcByType map[model.OperationType][]model.WorkCenter // только активные, отсортированы по ID
	sched    map[int64]map[int][][2]int                 // wc -> день недели (Пн=1) -> окна [startMin,endMin]
	busy     map[int64][]timeWindow                     // wc -> слитая занятость
}

// BuildPlan — чистая функция алгоритма (правила заказчика):
//  1. длительность операции = quantity * timePerUnitSeconds;
//  2. операция ДЕЛИМА, но ДЕТАЛЬ НЕДЕЛИМА: суммарная работа раскладывается
//     кусками по свободным слотам, каждый кусок кратен времени производства
//     одной детали; хвост отрезка короче одной детали не используется вовсе
//     (работа переносится в следующий слот); каждый кусок = отдельная запись
//     scheduled_operation (одинаковые ExternalOperationID/SequenceNumber);
//  3. горизонта нет: окна раскрываются день за днём, пока работа не набрана;
//  4. свободный слот = рабочее окно шаблона − занятость активных планов;
//  5. операции позиции идут по sequenceNumber, earliestStart следующей =
//     finish последнего куска предыдущей; позиции НЕ зависят друг от друга
//     (конкурируют только через занятость станков);
//  6. выбор станка — самое раннее окончание ВСЕЙ работы; при равном финише —
//     меньший ID станка;
//  7. ErrNoFreeWindow — если ни у одного кандидата нет рабочего дня в шаблоне
//     или ни одно окно не длиннее одной детали; ErrInvalidPlanningRequest —
//     если нет станка под тип операции;
//  8. planned_finish_at <= конец дня due_date -> onTime = true
//     (план строится всегда, onTime — только пометка в ответе).
func BuildPlan(in PlannerInput) (*PlanDraft, error) {
	if in.Timezone == nil {
		return nil, fmt.Errorf("%w: timezone is nil", model.ErrInvalidPlanningRequest)
	}
	if err := in.Request.Validate(); err != nil {
		return nil, err
	}

	st := newPlannerState(in)
	draft := &PlanDraft{
		Plan: model.ProductionPlan{
			ExternalOrderID: in.Request.OrderID,
			DueDate:         in.Request.DueDate,
			Priority:        in.Request.Priority,
			Status:          model.PlanStatusProposed,
		},
	}

	for itemIdx, item := range in.Request.Items {
		ops := sortedOps(item.Operations)
		earliest := in.Request.EarliestStartAt
		var itemStart, itemFinish time.Time

		for _, op := range ops {
			unit := time.Duration(op.TimePerUnitSeconds) * time.Second // производство одной детали
			need := time.Duration(item.Quantity) * unit
			if need <= 0 {
				return nil, fmt.Errorf("%w: operation %d has zero duration", model.ErrInvalidPlanningRequest, op.OperationID)
			}
			candidates := st.wcByType[op.OperationType]
			if len(candidates) == 0 {
				return nil, fmt.Errorf("%w: no work center for operation type %s", model.ErrInvalidPlanningRequest, op.OperationType)
			}

			best, err := st.pickCandidate(candidates, need, unit, earliest)
			if err != nil {
				return nil, err
			}
			st.occupy(best.wcID, best.segs)

			for _, seg := range best.segs {
				draft.Operations = append(draft.Operations, model.ScheduledOperation{
					ExternalOperationID: op.OperationID,
					SequenceNumber:      op.SequenceNumber,
					OperationType:       op.OperationType,
					WorkCenterID:        best.wcID,
					StartAt:             seg.start,
					FinishAt:            seg.end,
				})
				draft.OperationItemIdx = append(draft.OperationItemIdx, itemIdx)
			}
			if itemStart.IsZero() || best.segs[0].start.Before(itemStart) {
				itemStart = best.segs[0].start
			}
			if best.finish.After(itemFinish) {
				itemFinish = best.finish
			}
			earliest = best.finish
		}

		draft.Items = append(draft.Items, model.ProductionPlanItem{
			ExternalOrderItemID: item.OrderItemID,
			ExternalProductID:   item.ProductID,
			Quantity:            item.Quantity,
			PlannedStartAt:      &itemStart,
			PlannedFinishAt:     &itemFinish,
		})
	}

	st.fillHeader(draft)
	return draft, nil
}

// ---------------------------------------------------------------------------
// Внутренние шаги алгоритма
// ---------------------------------------------------------------------------

func newPlannerState(in PlannerInput) *plannerState {
	st := &plannerState{
		tz:       in.Timezone,
		wcByType: make(map[model.OperationType][]model.WorkCenter),
		sched:    make(map[int64]map[int][][2]int, len(in.WeeklySchedule)),
		busy:     make(map[int64][]timeWindow, len(in.Busy)),
	}
	for _, wc := range in.WorkCenters {
		if !wc.Active {
			continue
		}
		for _, t := range wc.OperationTypes {
			st.wcByType[model.OperationType(t)] = append(st.wcByType[model.OperationType(t)], wc)
		}
	}
	for t := range st.wcByType {
		wcs := st.wcByType[t]
		sort.Slice(wcs, func(i, j int) bool { return wcs[i].ID < wcs[j].ID })
	}
	for wcID, days := range in.WeeklySchedule {
		m := make(map[int][][2]int, len(days))
		for _, d := range days {
			m[d.DayOfWeek] = append(m[d.DayOfWeek], [2]int{d.StartTime.Minutes(), d.EndTime.Minutes()})
		}
		st.sched[wcID] = m
	}
	for wcID, ops := range in.Busy {
		ws := make([]timeWindow, 0, len(ops))
		for _, o := range ops {
			ws = append(ws, timeWindow{o.StartAt, o.FinishAt})
		}
		st.busy[wcID] = mergeWindows(ws)
	}
	return st
}

func sortedOps(ops []model.PlanningRequestOperation) []model.PlanningRequestOperation {
	out := append([]model.PlanningRequestOperation(nil), ops...)
	sort.Slice(out, func(i, j int) bool { return out[i].SequenceNumber < out[j].SequenceNumber })
	return out
}

// allocation — вариант размещения операции на конкретном станке.
type allocation struct {
	wcID   int64
	segs   []timeWindow // куски по порядку
	finish time.Time    // конец последнего куска
}

// pickCandidate — раскладывает работу на каждом станке-кандидате и
// выбирает вариант с самым ранним финишем (tie-break: меньший ID).
// Кандидат не подходит, если его окно короче одной детали (деталь неделима).
func (st *plannerState) pickCandidate(candidates []model.WorkCenter, need, unit time.Duration, earliest time.Time) (allocation, error) {
	found := false
	var best allocation
	for _, wc := range candidates {
		if len(st.sched[wc.ID]) == 0 {
			continue // шаблон пуст: все дни выходные
		}
		segs, err := st.allocateOn(wc.ID, need, unit, earliest)
		if err != nil {
			continue // ни одно окно не вмещает даже одной детали
		}
		f := segs[len(segs)-1].end
		if !found || f.Before(best.finish) {
			best = allocation{wcID: wc.ID, segs: segs, finish: f}
			found = true
		}
	}
	if !found {
		return allocation{}, model.ErrNoFreeWindow
	}
	return best, nil
}

// allocateOn — заполняет свободные слоты станка кусками, пока не наберётся
// need. Кусок кратен времени одной детали (unit): деталь неделима, поэтому
// хвост слота короче unit не используется, работа переносится дальше. Дни
// идут от даты earliest без ограничения горизонта: при непустом шаблоне, где
// хотя бы одно окно длиннее unit, оставшаяся работа каждый рабочий день
// строго уменьшается, поэтому цикл терминируется.
func (st *plannerState) allocateOn(wcID int64, need, unit time.Duration, earliest time.Time) ([]timeWindow, error) {
	// Защита от зацикливания: если ни одно окно шаблона не длиннее одной
	// детали, работа не поместится никогда.
	if st.maxWindowMinutes(wcID) < unit {
		return nil, model.ErrNoFreeWindow
	}
	segs := make([]timeWindow, 0, 1)
	remaining := need
	for day := dateIn(earliest, st.tz); remaining > 0; day = day.AddDate(0, 0, 1) {
		for _, w := range st.workWindows(day, wcID) {
			for _, free := range subtractBusy([]timeWindow{w}, st.busy[wcID]) {
				s := free.start
				if earliest.After(s) {
					s = earliest
				}
				avail := free.end.Sub(s)
				if avail <= 0 {
					continue
				}
				// Сколько ЦЕЛЫХ деталей влезает в слот; хвост < unit отбрасывается.
				unitsFit := int(avail / unit)
				if unitsFit == 0 {
					continue
				}
				unitsLeft := int(remaining / unit)
				if unitsLeft < unitsFit {
					unitsFit = unitsLeft
				}
				take := time.Duration(unitsFit) * unit
				segs = append(segs, timeWindow{s, s.Add(take)})
				remaining -= take
				if remaining == 0 {
					return segs, nil
				}
			}
		}
	}
	return segs, nil // недостижимо при непустом шаблоне с окном >= unit
}

// maxWindowMinutes — длина самого длинного рабочего окна недельного шаблона
// станка (0 — шаблон пуст).
func (st *plannerState) maxWindowMinutes(wcID int64) time.Duration {
	var max time.Duration
	for _, ivs := range st.sched[wcID] {
		for _, iv := range ivs {
			if d := time.Duration(iv[1]-iv[0]) * time.Minute; d > max {
				max = d
			}
		}
	}
	return max
}

// occupy — куски выбранного варианта становятся занятостью станка
// для последующих операций и позиций.
func (st *plannerState) occupy(wcID int64, segs []timeWindow) {
	merged := append(append([]timeWindow(nil), st.busy[wcID]...), segs...)
	st.busy[wcID] = mergeWindows(merged)
}

// workWindows — рабочие окна станка на конкретную локальную дату.
func (st *plannerState) workWindows(day time.Time, wcID int64) []timeWindow {
	dow := ((int(day.Weekday()) + 6) % 7) + 1 // Пн=1 ... Вс=7
	ws := make([]timeWindow, 0, len(st.sched[wcID][dow]))
	for _, iv := range st.sched[wcID][dow] {
		start := time.Date(day.Year(), day.Month(), day.Day(), iv[0]/60, iv[0]%60, 0, 0, st.tz)
		end := time.Date(day.Year(), day.Month(), day.Day(), iv[1]/60, iv[1]%60, 0, 0, st.tz)
		if end.After(start) {
			ws = append(ws, timeWindow{start, end})
		}
	}
	sort.Slice(ws, func(i, j int) bool { return ws[i].start.Before(ws[j].start) })
	return ws
}

// fillHeader — шапка плана: min старт / max финиш по всем операциям, onTime.
func (st *plannerState) fillHeader(draft *PlanDraft) {
	var minStart, maxFinish time.Time
	for _, op := range draft.Operations {
		if minStart.IsZero() || op.StartAt.Before(minStart) {
			minStart = op.StartAt
		}
		if op.FinishAt.After(maxFinish) {
			maxFinish = op.FinishAt
		}
	}
	draft.Plan.PlannedStartAt = &minStart
	draft.Plan.PlannedFinishAt = &maxFinish

	// onTime: финиш не позже конца календарного дня dueDate в Timezone.
	due := draft.Plan.DueDate.Time()
	dueEnd := time.Date(due.Year(), due.Month(), due.Day(), 23, 59, 59, int(time.Second-1), st.tz)
	onTime := !maxFinish.After(dueEnd)
	draft.Plan.OnTime = &onTime
}

// ---------------------------------------------------------------------------
// Чистые хелперы над интервалами
// ---------------------------------------------------------------------------

// dateIn — локальная полночь даты момента t.
func dateIn(t time.Time, loc *time.Location) time.Time {
	t = t.In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

// mergeWindows — сортирует и сливает пересекающиеся/соприкасающиеся интервалы.
func mergeWindows(in []timeWindow) []timeWindow {
	if len(in) == 0 {
		return nil
	}
	s := append([]timeWindow(nil), in...)
	sort.Slice(s, func(i, j int) bool { return s[i].start.Before(s[j].start) })
	out := make([]timeWindow, 0, len(s))
	out = append(out, s[0])
	for _, w := range s[1:] {
		last := &out[len(out)-1]
		if !w.start.After(last.end) {
			if w.end.After(last.end) {
				last.end = w.end
			}
			continue
		}
		out = append(out, w)
	}
	return out
}

// subtractBusy — вычитает занятость из окон; возвращает свободные отрезки.
// busy должен быть слит (mergeWindows).
func subtractBusy(windows, busy []timeWindow) []timeWindow {
	out := make([]timeWindow, 0, len(windows))
	for _, w := range windows {
		cur := []timeWindow{w}
		for _, b := range busy {
			next := make([]timeWindow, 0, len(cur)+1)
			for _, c := range cur {
				if !b.end.After(c.start) || !b.start.Before(c.end) {
					next = append(next, c) // не пересекаются
					continue
				}
				if b.start.After(c.start) {
					next = append(next, timeWindow{c.start, b.start})
				}
				if b.end.Before(c.end) {
					next = append(next, timeWindow{b.end, c.end})
				}
			}
			cur = next
		}
		out = append(out, cur...)
	}
	return out
}
