package service

import (
	"context"
	"fmt"
	"sort"
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
	if _, err := s.repo.GetWorkCenterByID(ctx, workCenterID); err != nil {
		return nil, err // ErrNotFound -> 404; станок проверяем до чтения дней
	}
	days, err := s.repo.GetWorkCenterDays(ctx, workCenterID)
	if err != nil {
		return nil, fmt.Errorf("load schedule of wc %d: %w", workCenterID, err)
	}
	if days == nil {
		days = []model.ScheduleDay{} // JSON [] вместо null: шаблон пуст = все дни выходные
	}
	return &model.WorkCenterScheduleResponse{
		WorkCenterID: workCenterID,
		Timezone:     s.timezone.String(),
		Days:         days,
	}, nil
}

// ReplaceWorkCenterSchedule — PUT /api/work-centers/{id}/schedule.
// Вход: model.ReplaceScheduleRequest (полный массив дней). Нет станка ->
// ErrNotFound (404), кривой день -> model.ErrInvalidScheduleDay (400).
// Валидация дней здесь, а не только в binding: dive-правило ловит форму
// (обязательные поля), а бизнес-правило start<end — домен.
func (s *Schedule) ReplaceWorkCenterSchedule(ctx context.Context, workCenterID int64, req model.ReplaceScheduleRequest) error {
	if _, err := s.repo.GetWorkCenterByID(ctx, workCenterID); err != nil {
		return err
	}
	for _, d := range req.Days {
		if err := d.Validate(); err != nil {
			return err // ErrInvalidScheduleDay -> 400
		}
	}
	if err := s.repo.ReplaceWorkCenterDays(ctx, workCenterID, req.Days); err != nil {
		return fmt.Errorf("replace schedule of wc %d: %w", workCenterID, err)
	}
	return nil
}

// Availability — GET /api/work-centers/{id}/availability?from=&to=.
// Раскрывает недельный шаблон в конкретные даты [from..to] включительно
// в s.timezone; занятость заказами НЕ вычитается (правка контракта).
// Нет станка -> ErrNotFound (404); from > to -> ErrInvalidDate (400).
func (s *Schedule) Availability(ctx context.Context, workCenterID int64, from, to model.Date) (*model.AvailabilityResponse, error) {
	if from.Time().After(to.Time()) {
		return nil, model.ErrInvalidDate
	}
	if _, err := s.repo.GetWorkCenterByID(ctx, workCenterID); err != nil {
		return nil, err
	}
	days, err := s.repo.GetWorkCenterDays(ctx, workCenterID)
	if err != nil {
		return nil, fmt.Errorf("load schedule of wc %d: %w", workCenterID, err)
	}

	// Шаблон по дням недели: ScheduleDay.DayOfWeek = Пн..Вс (1..7),
	// time.Weekday = Вс..Сб (0..6), переводим в Пн=1..Вс=7.
	byDow := make(map[int][]model.ScheduleDay, 7)
	for _, d := range days {
		byDow[d.DayOfWeek] = append(byDow[d.DayOfWeek], d)
	}

	intervals := make([]model.AvailabilityInterval, 0)
	for day := from.Time(); !day.After(to.Time()); day = day.AddDate(0, 0, 1) {
		// Полночь даты переводим в таймзону проекта: Date хранится в UTC,
		// а рабочие интервалы должны быть в локальном времени (RFC3339
		// со смещением, напр. +03:00).
		local := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, s.timezone)
		dow := (int(local.Weekday()) + 6) % 7 // Пн=0..Вс=6
		for _, d := range byDow[dow+1] {
			intervals = append(intervals, model.AvailabilityInterval{
				StartAt: local.Add(time.Duration(d.StartTime.Minutes()) * time.Minute),
				EndAt:   local.Add(time.Duration(d.EndTime.Minutes()) * time.Minute),
			})
		}
	}
	return &model.AvailabilityResponse{
		WorkCenterID: workCenterID,
		Timezone:     s.timezone.String(),
		Intervals:    intervals,
	}, nil
}

// GeneralSchedule — GET /api/schedule?from=&to= (Гант).
// Отдаёт model.ScheduleResponse, УЖЕ сгруппированный по workCenterId
// (правка контракта): GetBusyOperationsView -> группировка. Учитываются
// только операции активных планов (PROPOSED/ACCEPTED) — фильтр в SQL.
func (s *Schedule) GeneralSchedule(ctx context.Context, from, to time.Time) (*model.ScheduleResponse, error) {
	if from.After(to) {
		return nil, model.ErrInvalidDate
	}
	views, err := s.repo.GetBusyOperationsView(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("load schedule operations: %w", err)
	}

	// Группировка по станкам; SQL отдаёт по work_center_id, start_at,
	// но порядок станков гарантируем сами (сортировка по id), а операции
	// внутри станка идут в порядке выборки.
	resp := &model.ScheduleResponse{From: from, To: to}
	order := make([]int64, 0)
	byWC := make(map[int64]*model.ScheduleWorkCenterResponse)
	for _, v := range views {
		wc, ok := byWC[v.ScheduledOperation.WorkCenterID]
		if !ok {
			wc = &model.ScheduleWorkCenterResponse{
				WorkCenterID:    v.ScheduledOperation.WorkCenterID,
				WorkCenterCode:  v.WorkCenterCode,
				WorkCenterTitle: v.WorkCenterTitle,
				Operations:      make([]model.ScheduleOperationResponse, 0),
			}
			byWC[wc.WorkCenterID] = wc
			order = append(order, wc.WorkCenterID)
		}
		wc.Operations = append(wc.Operations, model.ScheduleOperationResponse{
			ScheduledOperationID: v.ScheduledOperation.ID,
			OrderID:              v.OrderID,
			OrderItemID:          v.OrderItemID,
			ProductID:            v.ProductID,
			OperationType:        v.ScheduledOperation.OperationType,
			SequenceNumber:       v.ScheduledOperation.SequenceNumber,
			StartAt:              v.ScheduledOperation.StartAt,
			FinishAt:             v.ScheduledOperation.FinishAt,
		})
	}
	sort.Slice(order, func(i, j int) bool { return order[i] < order[j] })
	for _, id := range order {
		resp.WorkCenters = append(resp.WorkCenters, *byWC[id])
	}
	if resp.WorkCenters == nil {
		resp.WorkCenters = []model.ScheduleWorkCenterResponse{} // JSON [] вместо null
	}
	return resp, nil
}
