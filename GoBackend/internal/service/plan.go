package service

import (
	"context"
	"fmt"
	"time"

	"github.com/domninivan/project__K/GoBackend/internal/config"
	"github.com/domninivan/project__K/GoBackend/internal/model"
	"github.com/domninivan/project__K/GoBackend/internal/port"
)

// Зависимости оркестрации плана. Реализации — repository.Plan,
// repository.WorkCenter, repository.Schedule, repository.TxManager.
type PlanRepository interface {
	Insert(ctx context.Context, tx port.Tx, plan *model.ProductionPlan) (int64, error)
	InsertItem(ctx context.Context, tx port.Tx, item *model.ProductionPlanItem) (int64, error)
	InsertOperation(ctx context.Context, tx port.Tx, op *model.ScheduledOperation) (int64, error)
	GetBusyOperations(ctx context.Context, tx port.Tx, workCenterID int64, from time.Time) ([]model.ScheduledOperation, error)
	GetByID(ctx context.Context, planID int64) (*model.ProductionPlan, error)
	GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error)
	GetItems(ctx context.Context, planID int64) ([]model.ProductionPlanItem, error)
	GetOperations(ctx context.Context, planID int64) ([]model.ScheduledOperationView, error)
	UpdateStatus(ctx context.Context, planID int64, from, to model.PlanStatus) (*model.ProductionPlan, error)
	CancelActiveByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error)
}

// WorkCenterSourceForPlan — минимум, который оркестрации нужен от справочника.
type WorkCenterSourceForPlan interface {
	List(ctx context.Context) ([]model.WorkCenter, error)
}

// ScheduleSourceForPlan — недельные шаблоны станков для PlannerInput.
type ScheduleSourceForPlan interface {
	GetWorkCenterDays(ctx context.Context, workCenterID int64) ([]model.ScheduleDay, error)
}

// TxManager — контракт транзакций глазами сервиса (реализация —
// repository.TxManager поверх pgxpool). Сервис знает ГДЕ граница
// бизнес-операции: nil из fn -> COMMIT, ошибка -> ROLLBACK (rollback
// выполняет адаптер, сервис не может забыть). ЧЕМ является транзакция
// (pgx.Tx) сервис не знает — внутрь приходит port.Tx.
type TxManager interface {
	WithTx(ctx context.Context, fn func(tx port.Tx) error) error
}

// Plan — оркестрация: транзакция -> чистый planner -> запись -> коммит.
type Plan struct {
	repo      PlanRepository
	wcSources WorkCenterSourceForPlan
	sched     ScheduleSourceForPlan
	db        TxManager
	cfg       *config.Config
}

func NewPlan(repo PlanRepository, wcSources WorkCenterSourceForPlan,
	sched ScheduleSourceForPlan, db TxManager, cfg *config.Config) *Plan {
	return &Plan{repo: repo, wcSources: wcSources, sched: sched, db: db, cfg: cfg}
}

// Create — POST /api/plans.
// Вход: model.PlanningRequest (после binding + req.Validate()).
// Отдаёт: полный model.ProductionPlanResponse.
// Ошибки: ErrInvalidPlanningRequest (400), ErrNoFreeWindow (409),
// ErrConflict (409), ErrNotFound (404), прочие (500).
// Алгоритм: снимок мира -> BuildPlan -> BEGIN -> повторная проверка занятости
// в tx -> Insert/InsertItem/InsertOperation -> COMMIT.
func (s *Plan) Create(ctx context.Context, req model.PlanningRequest) (*model.ProductionPlanResponse, error) {
	// Ранняя валидация: 400 до того, как что-то записано.
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// --- Снимок мира для чистого планировщика ---
	wcs, err := s.wcSources.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list work centers: %w", err)
	}

	in := PlannerInput{
		Request:        req,
		WorkCenters:    wcs,
		WeeklySchedule: make(map[int64][]model.ScheduleDay, len(wcs)),
		Busy:           make(map[int64][]model.ScheduledOperation, len(wcs)),
		Timezone:       s.cfg.Timezone,
	}
	for _, wc := range wcs {
		days, err := s.sched.GetWorkCenterDays(ctx, wc.ID)
		if err != nil {
			return nil, fmt.Errorf("load schedule of wc %d: %w", wc.ID, err)
		}
		in.WeeklySchedule[wc.ID] = days

		busy, err := s.repo.GetBusyOperations(ctx, nil, wc.ID, req.EarliestStartAt)
		if err != nil {
			return nil, fmt.Errorf("load busy ops of wc %d: %w", wc.ID, err)
		}
		in.Busy[wc.ID] = busy
	}

	draft, err := BuildPlan(in)
	if err != nil {
		return nil, err // ErrInvalidPlanningRequest | ErrNoFreeWindow
	}

	// --- Запись: транзакция, перепроверка занятости, вставка ---
	// Вся вставка атомарна: nil из fn -> COMMIT, ошибка (в т.ч. конфликт
	// перепроверки) -> ROLLBACK, в БД не остаётся даже шапки плана.
	err = s.db.WithTx(ctx, func(tx port.Tx) error {
		// Повторная проверка занятости ВНУТРИ транзакции: между снимком
		// выше и этим моментом параллельный расчёт мог закоммитить план,
		// занимающий те же слоты. FK/CHECK такой конфликт не ловят
		// (пересечение интервалов в SQL-ограничениях не выражается),
		// поэтому ловим здесь; ошибка откатит всё через WithTx.
		if conflict, err := s.recheckBusy(ctx, tx, draft); err != nil {
			return err
		} else if conflict != "" {
			return fmt.Errorf("%w: %s", model.ErrConflict, conflict)
		}

		if _, err := s.repo.Insert(ctx, tx, &draft.Plan); err != nil {
			return err
		}
		itemIDByItemIdx := make(map[int]int64, len(draft.Items))
		for i := range draft.Items {
			it := &draft.Items[i]
			it.ProductionPlanID = draft.Plan.ID
			if _, err := s.repo.InsertItem(ctx, tx, it); err != nil {
				return err
			}
			itemIDByItemIdx[i] = it.ID
		}
		for i := range draft.Operations {
			op := &draft.Operations[i]
			op.ProductionPlanItemID = itemIDByItemIdx[draft.OperationItemIdx[i]]
			if _, err := s.repo.InsertOperation(ctx, tx, op); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Сборка ответа — после COMMIT: репозиторий читает уже закоммиченные
	// строки через пул, tx ей не нужен.
	return s.buildPlanResponse(ctx, &draft.Plan)
}

// recheckBusy — сравнивает куски draft с занятостью активных планов,
// перечитанной ВНУТРИ транзакции tx. Возвращает непустое описание
// конфликта (для ошибки 409) или "". Считывание идёт через tx — проверка
// должна видеть данные в том же снапшоте, куда затем пишем; свои куски
// ещё не вставлены, поэтому на себя они не указывают.
func (s *Plan) recheckBusy(ctx context.Context, tx port.Tx, draft *PlanDraft) (string, error) {
	if len(draft.Operations) == 0 {
		return "", nil
	}
	// Разбираем куски по станкам и ищем min start (нижняя граница выборки).
	segsByWC := make(map[int64][]model.ScheduledOperation)
	var earliest time.Time
	for _, op := range draft.Operations {
		segsByWC[op.WorkCenterID] = append(segsByWC[op.WorkCenterID], op)
		if earliest.IsZero() || op.StartAt.Before(earliest) {
			earliest = op.StartAt
		}
	}
	for wcID, segs := range segsByWC {
		busy, err := s.repo.GetBusyOperations(ctx, tx, wcID, earliest)
		if err != nil {
			return "", fmt.Errorf("recheck busy ops of wc %d: %w", wcID, err)
		}
		for _, b := range busy {
			for _, seg := range segs {
				if b.StartAt.Before(seg.FinishAt) && seg.StartAt.Before(b.FinishAt) {
					return fmt.Sprintf("wc %d: planned %s..%s overlaps %s..%s (plan item %d)",
						wcID,
						seg.StartAt.Format(time.RFC3339), seg.FinishAt.Format(time.RFC3339),
						b.StartAt.Format(time.RFC3339), b.FinishAt.Format(time.RFC3339),
						b.ProductionPlanItemID), nil
				}
			}
		}
	}
	return "", nil
}

// GetByOrderID — GET /api/plans/orders/{orderId}.
// Отдаёт: полный model.ProductionPlanResponse; нет плана -> model.ErrNotFound (404).
func (s *Plan) GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlanResponse, error) {
	plan, err := s.repo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return s.buildPlanResponse(ctx, plan)
}

// Accept — POST /api/plans/{planId}/accept. PROPOSED -> ACCEPTED.
// Недопустимый переход или проигранная гонка -> ErrConflict (409).
func (s *Plan) Accept(ctx context.Context, planID int64) (*model.PlanStatusResponse, error) {
	return s.transition(ctx, planID, model.PlanStatusAccepted)
}

// Reject — POST /api/plans/{planId}/reject. PROPOSED -> REJECTED.
// После reject операции плана не занимают станки (OccupiesCapacity=false).
func (s *Plan) Reject(ctx context.Context, planID int64) (*model.PlanStatusResponse, error) {
	return s.transition(ctx, planID, model.PlanStatusRejected)
}

// CancelByOrder — POST /api/plans/orders/{orderId}/cancel.
// PROPOSED/ACCEPTED -> CANCELLED (по orderId, не по planId).
// Нет активного плана -> ErrNotFound.
func (s *Plan) CancelByOrder(ctx context.Context, orderID int64) (*model.PlanStatusResponse, error) {
	plan, err := s.repo.CancelActiveByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return &model.PlanStatusResponse{
		PlanID:    plan.ID,
		OrderID:   plan.ExternalOrderID,
		Status:    plan.Status,
		UpdatedAt: plan.UpdatedAt,
	}, nil
}

// transition — общий сценарий accept/reject: правила переходов проверяет
// домен (CanTransitionTo), атомарность гонки гарантирует CAS в репозитории.
// Двойная проверка не избыточна: CanTransitionTo отвечает за БИЗНЕС-запрет
// (409 с осмысленным "нельзя"), CAS — за ГОНКУ (двое приняли один план).
func (s *Plan) transition(ctx context.Context, planID int64, to model.PlanStatus) (*model.PlanStatusResponse, error) {
	current, err := s.repo.GetByID(ctx, planID)
	if err != nil {
		return nil, err // ErrNotFound -> 404
	}
	if !current.Status.CanTransitionTo(to) {
		return nil, fmt.Errorf("%w: plan %d cannot transition %s -> %s",
			model.ErrConflict, planID, current.Status, to)
	}
	updated, err := s.repo.UpdateStatus(ctx, planID, current.Status, to)
	if err != nil {
		return nil, err // ErrConflict: статус изменился параллельно
	}
	return &model.PlanStatusResponse{
		PlanID:    updated.ID,
		OrderID:   updated.ExternalOrderID,
		Status:    updated.Status,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

// buildPlanResponse — внутренний маппинг домена в ProductionPlanResponse
// (общий для Create и GetByOrderID): позиции плана + операции в JOIN
// (ScheduledOperationView), сгруппированные по позициям.
func (s *Plan) buildPlanResponse(ctx context.Context, plan *model.ProductionPlan) (*model.ProductionPlanResponse, error) {
	items, err := s.repo.GetItems(ctx, plan.ID)
	if err != nil {
		return nil, fmt.Errorf("load items of plan %d: %w", plan.ID, err)
	}
	ops, err := s.repo.GetOperations(ctx, plan.ID)
	if err != nil {
		return nil, fmt.Errorf("load operations of plan %d: %w", plan.ID, err)
	}

	opsByItem := make(map[int64][]model.ScheduledOperationView, len(items))
	for _, v := range ops {
		opsByItem[v.ScheduledOperation.ProductionPlanItemID] =
			append(opsByItem[v.ScheduledOperation.ProductionPlanItemID], v)
	}

	onTime := plan.OnTime != nil && *plan.OnTime
	resp := &model.ProductionPlanResponse{
		PlanID:          plan.ID,
		OrderID:         plan.ExternalOrderID,
		PlannedStartAt:  plan.PlannedStartAt,
		PlannedFinishAt: plan.PlannedFinishAt,
		DueDate:         plan.DueDate,
		OnTime:          onTime,
		Status:          plan.Status,
		Items:           make([]model.PlanItemResponse, 0, len(items)),
	}
	for _, it := range items {
		pi := model.PlanItemResponse{
			OrderItemID:     it.ExternalOrderItemID,
			ProductID:       it.ExternalProductID,
			Quantity:        it.Quantity,
			PlannedStartAt:  it.PlannedStartAt,
			PlannedFinishAt: it.PlannedFinishAt,
			Operations:      make([]model.PlanOperationResponse, 0, len(opsByItem[it.ID])),
		}
		for _, v := range opsByItem[it.ID] {
			pi.Operations = append(pi.Operations, model.PlanOperationResponse{
				ScheduledOperationID: v.ScheduledOperation.ID,
				OperationID:          v.ScheduledOperation.ExternalOperationID,
				SequenceNumber:       v.ScheduledOperation.SequenceNumber,
				OperationType:        v.ScheduledOperation.OperationType,
				WorkCenterID:         v.ScheduledOperation.WorkCenterID,
				WorkCenterCode:       v.WorkCenterCode,
				WorkCenterTitle:      v.WorkCenterTitle,
				StartAt:              v.ScheduledOperation.StartAt,
				FinishAt:             v.ScheduledOperation.FinishAt,
			})
		}
		resp.Items = append(resp.Items, pi)
	}
	return resp, nil
}
