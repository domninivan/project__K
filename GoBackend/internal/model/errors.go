package model

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")

	ErrInvalidDate             = errors.New("invalid date: not a real calendar date")
	ErrInvalidScheduleDay      = errors.New("invalid schedule day: startTime must be before endTime")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrPlanNotReschedulable    = errors.New("order has accepted plan")
	ErrNoFreeWindow            = errors.New("no work center with working days in schedule")
	ErrInvalidPlanningRequest  = errors.New("invalid planning request")
)
