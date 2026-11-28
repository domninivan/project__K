package model

import "time"

// ScheduleDay — один день недельного шаблона work_center_schedule.
// Нет строки на день недели = выходной.
type ScheduleDay struct {
	DayOfWeek int       `json:"dayOfWeek" binding:"min=1,max=7"` // Пн=1, Вс=7
	StartTime ClockTime `json:"startTime" binding:"required"`
	EndTime   ClockTime `json:"endTime"   binding:"required"`
}

// Validate — рабочий интервал должен быть непустым (start_time < end_time).
func (d ScheduleDay) Validate() error {
	// TODO(шаг 10): сравнить StartTime.Minutes() < EndTime.Minutes()
	return nil
}

// WorkCenterScheduleResponse — ответ GET /api/work-centers/{id}/schedule.
type WorkCenterScheduleResponse struct {
	WorkCenterID int64         `json:"workCenterId"`
	Timezone     string        `json:"timezone"` // из config.Timezone
	Days         []ScheduleDay `json:"days"`
}

// ReplaceScheduleRequest — вход PUT /api/work-centers/{id}/schedule:
// полный новый массив дней (не дифф).
type ReplaceScheduleRequest struct {
	Days []ScheduleDay `json:"days" binding:"required,min=1,dive"`
}

// AvailabilityResponse — ответ GET /api/work-centers/{id}/availability.
// Календарное раскрытие недельного шаблона; занятость заказами НЕ вычитается.
type AvailabilityResponse struct {
	WorkCenterID int64                  `json:"workCenterId"`
	Timezone     string                 `json:"timezone"`
	Intervals    []AvailabilityInterval `json:"intervals"`
}

// AvailabilityInterval — рабочий интервал конкретной даты.
type AvailabilityInterval struct {
	StartAt time.Time `json:"startAt"` // RFC3339 со смещением, напр. +03:00
	EndAt   time.Time `json:"endAt"`
}
