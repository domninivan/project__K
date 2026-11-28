package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/domninivan/project__K/GoBackend/internal/model"
	"github.com/domninivan/project__K/GoBackend/internal/port"
)

// Plan — репозиторий таблиц production_plan / production_plan_item /
// scheduled_operation. Методы записи принимают tx port.Tx (контракт сервиса,
// физически — pgx.Tx через txAdapter из tx.go): сервис открывает транзакцию
// через TxManager и передаёт её сюда; методы чтения идут через пул.
type Plan struct{ pool *pgxpool.Pool }

func NewPlan(pool *pgxpool.Pool) *Plan {
	return &Plan{pool: pool}
}

// planColumns — колонки шапки плана в порядке scanPlanInto; используется
// во всех SELECT/RETURNING, чтобы сканы не разъезжались.
const planColumns = `
	id, external_order_id, due_date, priority,
	planned_start_at, planned_finish_at, on_time, status,
	created_at, updated_at
`

// scanPlanInto — сканирует шапку плана через промежуточные переменные:
// DATE сканируется в time.Time (полночь UTC) и заворачивается в model.Date,
// перечисления сканируются в string и конвертируются в именованные типы.
// Нет строки (pgx.ErrNoRows) -> model.ErrNotFound.
func scanPlanInto(row pgx.Row) (*model.ProductionPlan, error) {
	var (
		p        model.ProductionPlan
		due      time.Time
		priority string
		status   string
	)
	err := row.Scan(&p.ID, &p.ExternalOrderID, &due, &priority,
		&p.PlannedStartAt, &p.PlannedFinishAt, &p.OnTime, &status,
		&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan production_plan: %w", err)
	}
	p.DueDate = model.NewDateFromTime(due)
	p.Priority = model.PlanPriority(priority)
	p.Status = model.PlanStatus(status)
	return &p, nil
}

// ---------------------------------------------------------------------------
// Запись плана (внутри транзакции)
// ---------------------------------------------------------------------------

// Insert — вставка шапки плана, возвращает сгенерированный id.
// created_at/updated_at проставляет БД (DEFAULT now()), сканируем их обратно
// в структуру. Статус ожидается уже валидным (CHECK в схеме страхует).
func (r *Plan) Insert(ctx context.Context, tx port.Tx, plan *model.ProductionPlan) (int64, error) {
	err := tx.QueryRow(ctx, `
		INSERT INTO production_plan
			(external_order_id, due_date, priority,
			 planned_start_at, planned_finish_at, on_time, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`,
		plan.ExternalOrderID, plan.DueDate.Time(), plan.Priority,
		plan.PlannedStartAt, plan.PlannedFinishAt, plan.OnTime, plan.Status,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("insert production_plan: %w", err)
	}
	return plan.ID, nil
}

// InsertItem — вставка позиции плана, возвращает id.
// Нарушение UNIQUE (production_plan_id, external_order_item_id) -> ErrConflict.
func (r *Plan) InsertItem(ctx context.Context, tx port.Tx, item *model.ProductionPlanItem) (int64, error) {
	err := tx.QueryRow(ctx, `
		INSERT INTO production_plan_item
			(production_plan_id, external_order_item_id, external_product_id,
			 quantity, planned_start_at, planned_finish_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		item.ProductionPlanID, item.ExternalOrderItemID, item.ExternalProductID,
		item.Quantity, item.PlannedStartAt, item.PlannedFinishAt,
	).Scan(&item.ID)
	if err != nil {
		return 0, mapPgError(err, "insert production_plan_item")
	}
	return item.ID, nil
}

// InsertOperation — вставка плановой операции (одного куска), возвращает id
// (это будет scheduledOperationId в ответе).
func (r *Plan) InsertOperation(ctx context.Context, tx port.Tx, op *model.ScheduledOperation) (int64, error) {
	err := tx.QueryRow(ctx, `
		INSERT INTO scheduled_operation
			(production_plan_item_id, external_operation_id, sequence_number,
			 operation_type, work_center_id, start_at, finish_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`,
		op.ProductionPlanItemID, op.ExternalOperationID, op.SequenceNumber,
		op.OperationType, op.WorkCenterID, op.StartAt, op.FinishAt,
	).Scan(&op.ID)
	if err != nil {
		return 0, fmt.Errorf("insert scheduled_operation: %w", err)
	}
	return op.ID, nil
}

// ---------------------------------------------------------------------------
// Чтение
// ---------------------------------------------------------------------------

// GetByID — шапка плана по id. Нет строки -> model.ErrNotFound.
func (r *Plan) GetByID(ctx context.Context, planID int64) (*model.ProductionPlan, error) {
	p, err := scanPlanInto(r.pool.QueryRow(ctx, `
		SELECT `+planColumns+`
		  FROM production_plan
		 WHERE id = $1
	`, planID))
	if err != nil {
		return nil, fmt.Errorf("get production_plan by id: %w", err)
	}
	return p, nil
}

// GetByOrderID — план Java-заказа. У заказа может быть несколько планов
// (перепланирование создаёт новую строку, старую отменяет сервис), поэтому
// берём самый свежий по id. Нет плана -> model.ErrNotFound (404).
func (r *Plan) GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error) {
	p, err := scanPlanInto(r.pool.QueryRow(ctx, `
		SELECT `+planColumns+`
		  FROM production_plan
		 WHERE external_order_id = $1
		 ORDER BY id DESC
		 LIMIT 1
	`, orderID))
	if err != nil {
		return nil, fmt.Errorf("get production_plan by order: %w", err)
	}
	return p, nil
}

// GetItems — позиции плана (для сборки ProductionPlanResponse).
func (r *Plan) GetItems(ctx context.Context, planID int64) ([]model.ProductionPlanItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, production_plan_id, external_order_item_id, external_product_id,
		       quantity, planned_start_at, planned_finish_at
		  FROM production_plan_item
		 WHERE production_plan_id = $1
		 ORDER BY id
	`, planID)
	if err != nil {
		return nil, fmt.Errorf("query production_plan_item: %w", err)
	}
	defer rows.Close()

	out := make([]model.ProductionPlanItem, 0)
	for rows.Next() {
		var it model.ProductionPlanItem
		if err := rows.Scan(&it.ID, &it.ProductionPlanID, &it.ExternalOrderItemID,
			&it.ExternalProductID, &it.Quantity, &it.PlannedStartAt, &it.PlannedFinishAt); err != nil {
			return nil, fmt.Errorf("scan production_plan_item: %w", err)
		}
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate production_plan_item: %w", err)
	}
	return out, nil
}

// GetOperations — операции плана в JOIN с позицией, планом и станком
// (ScheduledOperationView): из них собираются PlanItemResponse.Operations
// и ответ Ганта. Куски одной операции несут одинаковые external_operation_id
// и sequence_number, поэтому сортировка: позиция -> последовательность -> старт.
func (r *Plan) GetOperations(ctx context.Context, planID int64) ([]model.ScheduledOperationView, error) {
	const q = `
		SELECT so.id, so.production_plan_item_id, so.external_operation_id,
		       so.sequence_number, so.operation_type, so.work_center_id,
		       so.start_at, so.finish_at,
		       pp.id, pp.external_order_id,
		       ppi.external_order_item_id, ppi.external_product_id,
		       wc.code, wc.title
		  FROM scheduled_operation so
		  JOIN production_plan_item ppi ON ppi.id = so.production_plan_item_id
		  JOIN production_plan pp      ON pp.id = ppi.production_plan_id
		  JOIN work_center wc          ON wc.id = so.work_center_id
		 WHERE ppi.production_plan_id = $1
		 ORDER BY ppi.id, so.sequence_number, so.start_at, so.id
	`
	rows, err := r.pool.Query(ctx, q, planID)
	if err != nil {
		return nil, fmt.Errorf("query scheduled_operation view: %w", err)
	}
	defer rows.Close()

	out := make([]model.ScheduledOperationView, 0)
	for rows.Next() {
		var (
			v      model.ScheduledOperationView
			opType string
			wcID   int64
		)
		if err := rows.Scan(&v.ScheduledOperation.ID, &v.ScheduledOperation.ProductionPlanItemID,
			&v.ScheduledOperation.ExternalOperationID, &v.ScheduledOperation.SequenceNumber,
			&opType, &wcID,
			&v.ScheduledOperation.StartAt, &v.ScheduledOperation.FinishAt,
			&v.PlanID, &v.OrderID,
			&v.OrderItemID, &v.ProductID,
			&v.WorkCenterCode, &v.WorkCenterTitle); err != nil {
			return nil, fmt.Errorf("scan scheduled_operation view: %w", err)
		}
		v.ScheduledOperation.OperationType = model.OperationType(opType)
		v.ScheduledOperation.WorkCenterID = wcID
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scheduled_operation view: %w", err)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Занятость (для планировщика)
// ---------------------------------------------------------------------------

// GetBusyOperations — все занятые интервалы станка начиная с from.
// tx может быть nil (чтение вне транзакции, снимок для планировщика)
// или транзакцией оркестрации — перепроверка занятости идёт в той же tx,
// в которой потом вставляется план (правило доки 01). Пул реализует
// port.Tx структурно, поэтому "nil -> читать из пула" выражается одной
// заменой exec-объекта.
// Правой границы нет намеренно: у планировщика нет горизонта, он раскрывает
// окна день за днём, поэтому занятость на любую будущую дату блокирует слоты
// и обязана учитываться. Учитываются только активные планы
// (status IN ('PROPOSED','ACCEPTED')) — OccupiesCapacity(); отменённые и
// отклонённые станок не занимают.
func (r *Plan) GetBusyOperations(ctx context.Context, tx port.Tx, workCenterID int64, from time.Time) ([]model.ScheduledOperation, error) {
	q := port.Tx(poolAdapter{pool: r.pool})
	if tx != nil {
		q = tx
	}
	const sql = `
		SELECT so.id, so.production_plan_item_id, so.external_operation_id,
		       so.sequence_number, so.operation_type, so.work_center_id,
		       so.start_at, so.finish_at
		  FROM scheduled_operation so
		  JOIN production_plan_item ppi ON ppi.id = so.production_plan_item_id
		  JOIN production_plan pp      ON pp.id = ppi.production_plan_id
		 WHERE so.work_center_id = $1
		   AND so.finish_at > $2
		   AND pp.status IN ('PROPOSED', 'ACCEPTED')
		 ORDER BY so.start_at
	`
	rows, err := q.Query(ctx, sql, workCenterID, from)
	if err != nil {
		return nil, fmt.Errorf("query busy operations: %w", err)
	}
	defer rows.Close()

	out := make([]model.ScheduledOperation, 0)
	for rows.Next() {
		var (
			o      model.ScheduledOperation
			opType string
		)
		if err := rows.Scan(&o.ID, &o.ProductionPlanItemID, &o.ExternalOperationID,
			&o.SequenceNumber, &opType, &o.WorkCenterID, &o.StartAt, &o.FinishAt); err != nil {
			return nil, fmt.Errorf("scan busy operation: %w", err)
		}
		o.OperationType = model.OperationType(opType)
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate busy operations: %w", err)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Машина состояний
// ---------------------------------------------------------------------------

// UpdateStatus — атомарная смена статуса (CAS: compare-and-swap в одном
// UPDATE). WHERE id AND status=$from гарантирует, что из двух одновременных
// переходов выиграет ровно один; 0 строк -> model.ErrConflict (409).
// Правила допустимых переходов (CanTransitionTo) проверяет сервис, SQL знает
// только "из какого статуса".
func (r *Plan) UpdateStatus(ctx context.Context, planID int64, from, to model.PlanStatus) (*model.ProductionPlan, error) {
	p, err := scanPlanInto(r.pool.QueryRow(ctx, `
		UPDATE production_plan
		   SET status = $2, updated_at = now()
		 WHERE id = $1 AND status = $3
		RETURNING `+planColumns+`
	`, planID, to, from))
	if errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("%w: plan %d is not in status %s", model.ErrConflict, planID, from)
	}
	if err != nil {
		return nil, fmt.Errorf("update production_plan status: %w", err)
	}
	return p, nil
}

// CancelActiveByOrderID — cancel по заказу: отменяет ВСЕ активные планы
// заказа (status IN ('PROPOSED','ACCEPTED') -> 'CANCELLED') одним UPDATE
// и возвращает самый свежий из них. Отменяются все, а не один: инвариант
// "не более одного активного плана на заказ" поддерживает сервис, но при
// его нарушении отмена обязана освободить занятость целиком.
// Нет активного плана -> model.ErrNotFound.
func (r *Plan) CancelActiveByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error) {
	const q = `
		WITH upd AS (
			UPDATE production_plan
			   SET status = 'CANCELLED', updated_at = now()
			 WHERE external_order_id = $1
			   AND status IN ('PROPOSED', 'ACCEPTED')
			RETURNING ` + planColumns + `
		)
		SELECT ` + planColumns + ` FROM upd ORDER BY id DESC LIMIT 1
	`
	p, err := scanPlanInto(r.pool.QueryRow(ctx, q, orderID))
	if errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("%w: no active plan for order %d", model.ErrNotFound, orderID)
	}
	if err != nil {
		return nil, fmt.Errorf("cancel production_plan by order: %w", err)
	}
	return p, nil
}
