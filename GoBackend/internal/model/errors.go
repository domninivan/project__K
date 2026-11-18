package model

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")

	ErrInvalidDate             = errors.New("invalid date: not a real calendar date")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrPlanNotReschedulable    = errors.New("order has accepted plan")
	ErrNoFreeWindow            = errors.New("no free window within planning horizon")
	ErrInvalidPlanningRequest  = errors.New("invalid planning request")
)
