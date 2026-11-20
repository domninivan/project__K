package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// Plan — репозиторий таблиц production_plan / production_plan_item /
// scheduled_operation. Методы записи принимают tx: сервис открывает
// транзакцию и передаёт её сюда (шаг 9).
type Plan struct{ pool *pgxpool.Pool }

func NewPlan(pool *pgxpool.Pool) *Plan {
	return &Plan{pool: pool}
}

// ---------------------------------------------------------------------------
// Запись плана (внутри транзакции)
// ---------------------------------------------------------------------------

// Insert — вставка шапки плана, возвращает сгенерированный id.
// Рекомендуется INSERT ... RETURNING id, created_at, updated_at.
func (r *Plan) Insert(ctx context.Context, tx pgx.Tx, plan *model.ProductionPlan) (int64, error) {
	// TODO(шаг 7): INSERT INTO production_plan
	//   (external_order_id, due_date, priority, planned_start_at, planned_finish_at, on_time, status)
	//   VALUES (...) RETURNING id
	return 0, nil
}

// InsertItem — вставка позиции плана, возвращает id.
// UNIQUE (production_plan_id, external_order_item_id) — конфликт -> ErrConflict.
func (r *Plan) InsertItem(ctx context.Context, tx pgx.Tx, item *model.ProductionPlanItem) (int64, error) {
	// TODO(шаг 7): INSERT INTO production_plan_item ... RETURNING id
	return 0, nil
}

// InsertOperation — вставка плановой операции, возвращает id
// (это будет scheduledOperationId в ответе).
func (r *Plan) InsertOperation(ctx context.Context, tx pgx.Tx, op *model.ScheduledOperation) (int64, error) {
	// TODO(шаг 7): INSERT INTO scheduled_operation ... RETURNING id
	return 0, nil
}

// ---------------------------------------------------------------------------
// Чтение
// ---------------------------------------------------------------------------

// GetByID — шапка плана по id. Нет строки -> model.ErrNotFound.
func (r *Plan) GetByID(ctx context.Context, planID int64) (*model.ProductionPlan, error) {
	// TODO(шаг 7): SELECT ... FROM production_plan WHERE id = $1
	return nil, nil
}

// GetByOrderID — план Java-заказа (индекс production_plan(external_order_id)).
// Нет плана -> model.ErrNotFound (handler отдаст 404).
func (r *Plan) GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error) {
	// TODO(шаг 7): SELECT ... WHERE external_order_id = $1 ORDER BY id DESC LIMIT 1
	return nil, nil
}

// GetItems — позиции плана (для сборки ProductionPlanResponse).
func (r *Plan) GetItems(ctx context.Context, planID int64) ([]model.ProductionPlanItem, error) {
	// TODO(шаг 7): SELECT ... FROM production_plan_item WHERE production_plan_id = $1
	return nil, nil
}

// GetOperations — операции плана в JOIN с позицией и станком
// (ScheduledOperationView): из них собираются PlanItemResponse.Operations
// и response для Ганта.
func (r *Plan) GetOperations(ctx context.Context, planID int64) ([]model.ScheduledOperationView, error) {
	// TODO(шаг 7): JOIN scheduled_operation + production_plan_item + work_center
	return nil, nil
}

// ---------------------------------------------------------------------------
// Занятость (для планировщика)
// ---------------------------------------------------------------------------

// GetBusyOperations — занятые интервалы станка в периоде.
// Учитываются только активные планы (status IN ('PROPOSED','ACCEPTED')).
// Вызывается до вставки результата; при параллельных расчётах сервис
// повторно проверяет занятость внутри транзакции (правило доки 01).
func (r *Plan) GetBusyOperations(ctx context.Context, workCenterID int64, from, to time.Time) ([]model.ScheduledOperation, error) {
	// TODO(шаг 7): SELECT ... WHERE work_center_id = $1
	//   AND finish_at > $2 AND start_at < $3
	//   AND status IN ('PROPOSED','ACCEPTED')  -- через JOIN с production_plan
	return nil, nil
}

// ---------------------------------------------------------------------------
// Машина состояний
// ---------------------------------------------------------------------------

// UpdateStatus — атомарная смена статуса:
// UPDATE production_plan SET status=$to, updated_at=now()
//
//	WHERE id=$1 AND status=$from RETURNING ...
//
// 0 затронутых строк -> model.ErrConflict (handler отдаст 409).
func (r *Plan) UpdateStatus(ctx context.Context, planID int64, from, to model.PlanStatus) (*model.ProductionPlan, error) {
	// TODO(шаг 12)
	return nil, nil
}

// CancelActiveByOrderID — cancel по заказу: отменяет активный план
// (status IN ('PROPOSED','ACCEPTED') -> 'CANCELLED') одним UPDATE.
// Нет активного плана -> model.ErrNotFound.
func (r *Plan) CancelActiveByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error) {
	// TODO(шаг 12): UPDATE ... WHERE external_order_id=$1
	//   AND status IN ('PROPOSED','ACCEPTED') RETURNING ...
	return nil, nil
}
