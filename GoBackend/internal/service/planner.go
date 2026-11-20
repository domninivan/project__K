package service

import (
	"time"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// PlannerInput — всё, что нужно чистому планировщику. Никакой БД внутри:
// данные поставляет оркестрация (service/plan.go, шаг 9).
type PlannerInput struct {
	Request        model.PlanningRequest                // запрос Java
	WorkCenters    []model.WorkCenter                   // активные станки (с operationTypes)
	WeeklySchedule map[int64][]model.ScheduleDay        // workCenterID -> недельный шаблон
	Busy           map[int64][]model.ScheduledOperation // workCenterID -> занятость активных планов
	Timezone       *time.Location                       // cfg.Timezone (Europe/Moscow)
	HorizonDays    int                                  // cfg.PlanningHorizonDays (90)
}

// PlanDraft — результат планирования до записи в БД: ID проставит
// репозиторий при вставке. Plan.PlannedStartAt / PlannedFinishAt / OnTime /
// Status=PROPOSED заполнены здесь.
type PlanDraft struct {
	Plan       model.ProductionPlan
	Items      []model.ProductionPlanItem
	Operations []model.ScheduledOperation // без ID и ProductionPlanItemID до вставки
}

// BuildPlan — чистая функция алгоритма из доки 03:
//  1. длительность операции = quantity * timePerUnitSeconds;
//  2. операция неделима — целиком в одно непрерывное свободное окно;
//  3. свободное окно = рабочие интервалы шаблона − занятость активных планов;
//  4. операции позиции идут по sequenceNumber, earliestStart следующей =
//     finish предыдущей; позиции не зависят друг от друга;
//  5. выбор станка/окна — самое раннее окончание;
//  6. нет окна в горизонте HorizonDays -> model.ErrNoFreeWindow;
//  7. planned_finish_at <= due_date -> onTime = true (план всё равно строится).
func BuildPlan(in PlannerInput) (*PlanDraft, error) {
	// TODO(шаг 8): алгоритм + table-driven тесты (раздел 8 доки 03)
	return nil, nil
}
