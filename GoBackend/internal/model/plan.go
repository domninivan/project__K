package model

import "time"

type PlanStatus string

const (
	PlanStatusProposed  PlanStatus = "PROPOSED"
	PlanStatusAccepted  PlanStatus = "ACCEPTED"
	PlanStatusRejected  PlanStatus = "REJECTED"
	PlanStatusCancelled PlanStatus = "CANCELLED"
)

func (s PlanStatus) OccupiesCapacity() bool {
	return s == PlanStatusProposed || s == PlanStatusAccepted
}

func (s PlanStatus) IsValid() bool {
	switch s {
	case PlanStatusProposed, PlanStatusAccepted, PlanStatusRejected, PlanStatusCancelled:
		return true
	}
	return false
}

var allowedTransitions = map[PlanStatus][]PlanStatus{
	PlanStatusProposed:  {PlanStatusAccepted, PlanStatusRejected, PlanStatusCancelled},
	PlanStatusAccepted:  {PlanStatusCancelled},
	PlanStatusRejected:  {},
	PlanStatusCancelled: {},
}

func (s PlanStatus) CanTransitionTo(next PlanStatus) bool {
	for _, t := range allowedTransitions[s] {
		if t == next {
			return true
		}
	}
	return false
}

type PlanPriority string

const (
	PlanPriorityLow    PlanPriority = "LOW"
	PlanPriorityMedium PlanPriority = "MEDIUM"
	PlanPriorityHigh   PlanPriority = "HIGH"
)

func (p PlanPriority) IsValid() bool {
	switch p {
	case PlanPriorityLow, PlanPriorityMedium, PlanPriorityHigh:
		return true
	}
	return false
}

type OperationType string

const (
	OpCutting    OperationType = "CUTTING"
	OpTurning    OperationType = "TURNING"
	OpMilling    OperationType = "MILLING"
	OpGrinding   OperationType = "GRINDING"
	OpDrilling   OperationType = "DRILLING"
	OpInspection OperationType = "INSPECTION"
)

func (t OperationType) IsValid() bool {
	switch t {
	case OpCutting, OpTurning, OpMilling, OpGrinding, OpDrilling, OpInspection:
		return true
	}
	return false
}

type ProductionPlan struct {
	ID              int64
	ExternalOrderID int64
	DueDate         Date
	Priority        PlanPriority
	PlannedStartAt  *time.Time
	PlannedFinishAt *time.Time
	OnTime          *bool
	Status          PlanStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ProductionPlanItem struct {
	ID                  int64
	ProductionPlanID    int64
	ExternalOrderItemID int64
	ExternalProductID   int64
	Quantity            int
	PlannedStartAt      *time.Time
	PlannedFinishAt     *time.Time
}

type ScheduledOperation struct {
	ID                   int64
	ProductionPlanItemID int64
	ExternalOperationID  int64
	SequenceNumber       int
	OperationType        OperationType
	WorkCenterID         int64
	StartAt              time.Time
	FinishAt             time.Time
}

type ScheduledOperationView struct {
	ScheduledOperation ScheduledOperation
	PlanID             int64
	OrderID            int64
	OrderItemID        int64
	ProductID          int64
	WorkCenterCode     string
	WorkCenterTitle    string
}

type PlanningRequest struct {
	OrderID         int64                 `json:"orderId"         binding:"required,min=1"`
	EarliestStartAt time.Time             `json:"earliestStartAt" binding:"required"`
	DueDate         Date                  `json:"dueDate"         binding:"required"`
	Priority        PlanPriority          `json:"priority"        binding:"required"`
	Items           []PlanningRequestItem `json:"items"           binding:"required,min=1,dive"`
}

type PlanningRequestItem struct {
	OrderItemID int64                      `json:"orderItemId" binding:"required,min=1"`
	ProductID   int64                      `json:"productId"   binding:"required,min=1"`
	Quantity    int                        `json:"quantity"    binding:"required,min=1"`
	Operations  []PlanningRequestOperation `json:"operations"  binding:"required,min=1,dive"`
}

type PlanningRequestOperation struct {
	OperationID        int64         `json:"operationId"        binding:"required,min=1"`
	SequenceNumber     int           `json:"sequenceNumber"     binding:"required,min=1"`
	OperationType      OperationType `json:"operationType"      binding:"required"`
	TimePerUnitSeconds int           `json:"timePerUnitSeconds" binding:"required,min=1"`
}

func (r PlanningRequest) Validate() error {
	if !r.Priority.IsValid() {
		return ErrInvalidPlanningRequest
	}
	if !r.DueDate.IsZero() && !r.EarliestStartAt.IsZero() {
		earliestDay := NewDateFromTime(r.EarliestStartAt)
		if r.DueDate.Time().Before(earliestDay.Time()) {
			return ErrInvalidPlanningRequest
		}
	}
	seenItems := make(map[int64]struct{}, len(r.Items))
	for _, item := range r.Items {
		if _, dup := seenItems[item.OrderItemID]; dup {
			return ErrInvalidPlanningRequest
		}
		seenItems[item.OrderItemID] = struct{}{}

		seenOps := make(map[int64]struct{}, len(item.Operations))
		prevSeq := 0
		for _, op := range item.Operations {
			if !op.OperationType.IsValid() {
				return ErrInvalidPlanningRequest
			}
			if op.SequenceNumber <= prevSeq {
				return ErrInvalidPlanningRequest
			}
			prevSeq = op.SequenceNumber
			if _, dup := seenOps[op.OperationID]; dup {
				return ErrInvalidPlanningRequest
			}
			seenOps[op.OperationID] = struct{}{}
		}
	}
	return nil
}

func NewDateFromTime(t time.Time) Date {
	return Date{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())}
}

type ProductionPlanResponse struct {
	PlanID          int64              `json:"planId"`
	OrderID         int64              `json:"orderId"`
	PlannedStartAt  *time.Time         `json:"plannedStartAt"`
	PlannedFinishAt *time.Time         `json:"plannedFinishAt"`
	DueDate         Date               `json:"dueDate"`
	OnTime          bool               `json:"onTime"`
	Status          PlanStatus         `json:"status"`
	Items           []PlanItemResponse `json:"items"`
}

type PlanItemResponse struct {
	OrderItemID     int64                   `json:"orderItemId"`
	ProductID       int64                   `json:"productId"`
	Quantity        int                     `json:"quantity"`
	PlannedStartAt  *time.Time              `json:"plannedStartAt"`
	PlannedFinishAt *time.Time              `json:"plannedFinishAt"`
	Operations      []PlanOperationResponse `json:"operations"`
}

type PlanOperationResponse struct {
	ScheduledOperationID int64         `json:"scheduledOperationId"`
	OperationID          int64         `json:"operationId"`
	SequenceNumber       int           `json:"sequenceNumber"`
	OperationType        OperationType `json:"operationType"`
	WorkCenterID         int64         `json:"workCenterId"`
	WorkCenterCode       string        `json:"workCenterCode"`
	WorkCenterTitle      string        `json:"workCenterTitle"`
	StartAt              time.Time     `json:"startAt"`
	FinishAt             time.Time     `json:"finishAt"`
}

type ScheduleResponse struct {
	From        time.Time                    `json:"from"`
	To          time.Time                    `json:"to"`
	WorkCenters []ScheduleWorkCenterResponse `json:"workCenters"`
}

type ScheduleWorkCenterResponse struct {
	WorkCenterID    int64                       `json:"workCenterId"`
	WorkCenterCode  string                      `json:"workCenterCode"`
	WorkCenterTitle string                      `json:"workCenterTitle"`
	Operations      []ScheduleOperationResponse `json:"operations"`
}

type ScheduleOperationResponse struct {
	ScheduledOperationID int64         `json:"scheduledOperationId"`
	OrderID              int64         `json:"orderId"`
	OrderItemID          int64         `json:"orderItemId"`
	ProductID            int64         `json:"productId"`
	OperationType        OperationType `json:"operationType"`
	SequenceNumber       int           `json:"sequenceNumber"`
	StartAt              time.Time     `json:"startAt"`
	FinishAt             time.Time     `json:"finishAt"`
}

type PlanStatusResponse struct {
	PlanID    int64      `json:"planId"`
	OrderID   int64      `json:"orderId"`
	Status    PlanStatus `json:"status"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
