package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// Schedule — репозиторий недельного шаблона work_center_schedule
// и выборок занятости для Ганта. Конкретные даты (availability) из
// шаблона раскрывает сервис, не SQL.
type Schedule struct{ pool *pgxpool.Pool }

func NewSchedule(pool *pgxpool.Pool) *Schedule {
	return &Schedule{pool: pool}
}

// GetWorkCenterByID — станок по id; нет -> model.ErrNotFound (для 404
// в schedule/availability эндпоинтах).
func (r *Schedule) GetWorkCenterByID(ctx context.Context, workCenterID int64) (*model.WorkCenter, error) {
	var (
		wc     model.WorkCenter
		active bool
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, code, title, active
		  FROM work_center
		 WHERE id = $1
	`, workCenterID).Scan(&wc.ID, &wc.Code, &wc.Title, &active)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get work_center by id: %w", err)
	}
	wc.Active = active
	return &wc, nil
}

// GetWorkCenterDays — недельный шаблон станка, упорядоченный по day_of_week.
// Нет строк = станок есть, но все дни выходные (пустой слайс, не ошибка).
// TIME сканируется строкой (::text даёт "09:00:00") и парсится в ClockTime —
// детерминированно, без надежды на pgx-кодек для нестандартного типа.
func (r *Schedule) GetWorkCenterDays(ctx context.Context, workCenterID int64) ([]model.ScheduleDay, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT day_of_week, start_time::text, end_time::text
		  FROM work_center_schedule
		 WHERE work_center_id = $1
		 ORDER BY day_of_week, start_time
	`, workCenterID)
	if err != nil {
		return nil, fmt.Errorf("query work_center_schedule: %w", err)
	}
	defer rows.Close()

	out := make([]model.ScheduleDay, 0)
	for rows.Next() {
		var (
			d      model.ScheduleDay
			st, en string
		)
		if err := rows.Scan(&d.DayOfWeek, &st, &en); err != nil {
			return nil, fmt.Errorf("scan work_center_schedule: %w", err)
		}
		if d.StartTime, err = model.ParseClockTime(st); err != nil {
			return nil, fmt.Errorf("parse start_time %q: %w", st, err)
		}
		if d.EndTime, err = model.ParseClockTime(en); err != nil {
			return nil, fmt.Errorf("parse end_time %q: %w", en, err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate work_center_schedule: %w", err)
	}
	return out, nil
}

// ReplaceWorkCenterDays — заменить шаблон целиком: DELETE + INSERT
// в одной транзакции (PUT-семантика, не дифф). Атомарность инкапсулирована
// здесь: это одна операция репозитория, сервису tx для неё не нужен.
func (r *Schedule) ReplaceWorkCenterDays(ctx context.Context, workCenterID int64, days []model.ScheduleDay) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin replace schedule: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		DELETE FROM work_center_schedule WHERE work_center_id = $1
	`, workCenterID); err != nil {
		return fmt.Errorf("delete work_center_schedule: %w", err)
	}
	for _, d := range days {
		if _, err := tx.Exec(ctx, `
			INSERT INTO work_center_schedule (work_center_id, day_of_week, start_time, end_time)
			VALUES ($1, $2, $3, $4)
		`, workCenterID, d.DayOfWeek, d.StartTime, d.EndTime); err != nil {
			return fmt.Errorf("insert work_center_schedule (day %d): %w", d.DayOfWeek, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit replace schedule: %w", err)
	}
	return nil
}

// GetBusyOperationsView — все операции активных планов в периоде,
// JOIN с production_plan_item, production_plan, work_center
// (ScheduledOperationView несёт orderID/orderItemID/productID/код и титул станка).
// Сервис сгруппирует по workCenterId для ScheduleResponse (Гант).
func (r *Schedule) GetBusyOperationsView(ctx context.Context, from, to time.Time) ([]model.ScheduledOperationView, error) {
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
		 WHERE pp.status IN ('PROPOSED', 'ACCEPTED')
		   AND so.finish_at > $1 AND so.start_at < $2
		 ORDER BY so.work_center_id, so.start_at, so.id
	`
	rows, err := r.pool.Query(ctx, q, from, to)
	if err != nil {
		return nil, fmt.Errorf("query busy operations view: %w", err)
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
			return nil, fmt.Errorf("scan busy operations view: %w", err)
		}
		v.ScheduledOperation.OperationType = model.OperationType(opType)
		v.ScheduledOperation.WorkCenterID = wcID
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate busy operations view: %w", err)
	}
	return out, nil
}
