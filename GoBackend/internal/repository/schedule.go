package repository

import (
	"context"
	"time"

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
// в schedule/availability эндпоинтах). Позже можно перенести в WorkCenter.
func (r *Schedule) GetWorkCenterByID(ctx context.Context, workCenterID int64) (*model.WorkCenter, error) {
	// TODO(шаг 10): SELECT id, code, title, active FROM work_center WHERE id = $1
	return nil, nil
}

// GetWorkCenterDays — недельный шаблон станка, упорядоченный по day_of_week.
// Нет строк = станок есть, но все дни выходные (пустой слайс, не ошибка).
func (r *Schedule) GetWorkCenterDays(ctx context.Context, workCenterID int64) ([]model.ScheduleDay, error) {
	// TODO(шаг 10): SELECT day_of_week, start_time, end_time
	//   FROM work_center_schedule WHERE work_center_id = $1 ORDER BY day_of_week
	return nil, nil
}

// ReplaceWorkCenterDays — заменить шаблон целиком: DELETE + INSERT
// в одной транзакции (PUT-семантика, не дифф).
func (r *Schedule) ReplaceWorkCenterDays(ctx context.Context, workCenterID int64, days []model.ScheduleDay) error {
	// TODO(шаг 10): BEGIN; DELETE FROM work_center_schedule WHERE work_center_id=$1;
	//   INSERT ... unnest; COMMIT
	return nil
}

// GetBusyOperationsView — все операции активных планов в периоде,
// JOIN с production_plan_item, production_plan, work_center
// (ScheduledOperationView несёт orderID/orderItemID/productID/код и титул станка).
// Сервис сгруппирует по workCenterId для ScheduleResponse (Гант).
func (r *Schedule) GetBusyOperationsView(ctx context.Context, from, to time.Time) ([]model.ScheduledOperationView, error) {
	// TODO(шаг 11): SELECT so.id, so.production_plan_item_id, ..., wc.code, wc.title
	//   FROM scheduled_operation so
	//   JOIN production_plan_item ppi ON ...
	//   JOIN production_plan pp ON ...
	//   JOIN work_center wc ON ...
	//  WHERE pp.status IN ('PROPOSED','ACCEPTED')
	//    AND so.finish_at > $1 AND so.start_at < $2
	//  ORDER BY so.work_center_id, so.start_at
	return nil, nil
}
