package service

import (
	"context"
	"time"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// ScheduleRepository — зависимости для расписаний и Ганта
// (реализация repository.Schedule).
type ScheduleRepository interface {
	GetWorkCenterByID(ctx context.Context, workCenterID int64) (*model.WorkCenter, error)
	GetWorkCenterDays(ctx context.Context, workCenterID int64) ([]model.ScheduleDay, error)
	ReplaceWorkCenterDays(ctx context.Context, workCenterID int64, days []model.ScheduleDay) error
	GetBusyOperationsView(ctx context.Context, from, to time.Time) ([]model.ScheduledOperationView, error)
}

// Schedule — недельные шаблоны, календарная доступность и общий Гант.
type Schedule struct {
	repo     ScheduleRepository
	timezone *time.Location // cfg.Timezone, напр. Europe/Moscow
}

func NewSchedule(repo ScheduleRepository, timezone *time.Location) *Schedule {
	return &Schedule{repo: repo, timezone: timezone}
}

// WorkCenterSchedule — GET /api/work-centers/{id}/schedule.
// Отдаёт model.WorkCenterScheduleResponse (timezone из cfg).
// Нет станка -> model.ErrNotFound (404).
func (s *Schedule) WorkCenterSchedule(ctx context.Context, workCenterID int64) (*model.WorkCenterScheduleResponse, error) {
	// TODO(шаг 10)
	return nil, nil
}

// ReplaceWorkCenterSchedule — PUT /api/work-centers/{id}/schedule.
// Вход: model.ReplaceScheduleRequest (полный массив дней, после binding +
// Validate каждого дня). Нет станка -> ErrNotFound, кривой день -> 400.
func (s *Schedule) ReplaceWorkCenterSchedule(ctx context.Context, workCenterID int64, req model.ReplaceScheduleRequest) error {
	// TODO(шаг 10)
	return nil
}

// Availability — GET /api/work-centers/{id}/availability?from=&to=.
// Вход: from/to как model.Date (YYYY-MM-DD). Раскрывает недельный шаблон
// в конкретные даты в s.timezone; занятость заказами НЕ вычитается
// (правка контракта). Нет станка -> ErrNotFound.
func (s *Schedule) Availability(ctx context.Context, workCenterID int64, from, to model.Date) (*model.AvailabilityResponse, error) {
	// TODO(шаг 10)
	return nil, nil
}

// GeneralSchedule — GET /api/schedule?from=&to= (Гант).
// Вход: from/to RFC3339 со смещением. Отдаёт model.ScheduleResponse,
// УЖЕ сгруппированный по workCenterId (правка контракта):
// GetBusyOperationsView -> группировка по WorkCenterID в order.
func (s *Schedule) GeneralSchedule(ctx context.Context, from, to time.Time) (*model.ScheduleResponse, error) {
	// TODO(шаг 11)
	return nil, nil
}
