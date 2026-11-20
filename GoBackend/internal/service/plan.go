package service

import (
	"context"

	"github.com/domninivan/project__K/GoBackend/internal/config"
	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// Зависимости оркестрации плана. Реализации — repository.Plan,
// repository.WorkCenter, repository.Schedule.
type PlanRepository interface {
	Insert(ctx context.Context, tx planTx, plan *model.ProductionPlan) (int64, error)
	InsertItem(ctx context.Context, tx planTx, item *model.ProductionPlanItem) (int64, error)
	InsertOperation(ctx context.Context, tx planTx, op *model.ScheduledOperation) (int64, error)
	GetByID(ctx context.Context, planID int64) (*model.ProductionPlan, error)
	GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error)
	GetItems(ctx context.Context, planID int64) ([]model.ProductionPlanItem, error)
	GetOperations(ctx context.Context, planID int64) ([]model.ScheduledOperationView, error)
	UpdateStatus(ctx context.Context, planID int64, from, to model.PlanStatus) (*model.ProductionPlan, error)
	CancelActiveByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlan, error)
}

// WorkCenterSourceForPlan — минимум, который оркестрации нужен от справочника.
type WorkCenterSourceForPlan interface {
	List(ctx context.Context) ([]model.WorkCenter, error)
}

type planTx = any // TODO(шаг 9): заменить на pgx.Tx после правки интерфейса репозитория

// Plan — оркестрация: транзакция -> чистый planner -> запись -> коммит.
type Plan struct {
	repo      PlanRepository
	wcSources WorkCenterSourceForPlan
	cfg       *config.Config
}

func NewPlan(repo PlanRepository, wcSources WorkCenterSourceForPlan, cfg *config.Config) *Plan {
	return &Plan{repo: repo, wcSources: wcSources, cfg: cfg}
}

// Create — POST /api/plans.
// Вход: model.PlanningRequest (после binding + req.Validate()).
// Отдаёт: полный model.ProductionPlanResponse (правка контракта).
// Ошибки: ErrInvalidPlanningRequest (400), ErrNoFreeWindow (422/409),
// ErrConflict (409), прочие (500).
// Алгоритм: BEGIN -> собран planner input -> planner.BuildPlan -> Insert/InsertItem/
// InsertOperation -> повторная проверка занятости в tx -> COMMIT.
func (s *Plan) Create(ctx context.Context, req model.PlanningRequest) (*model.ProductionPlanResponse, error) {
	// TODO(шаг 9)
	return nil, nil
}

// GetByOrderID — GET /api/plans/orders/{orderId}.
// Отдаёт: полный model.ProductionPlanResponse; нет плана -> model.ErrNotFound (404).
func (s *Plan) GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlanResponse, error) {
	// TODO(шаг 11): GetByOrderID + GetItems + GetOperations -> сборка ответа
	return nil, nil
}

// Accept — POST /api/plans/{planId}/accept. PROPOSED -> ACCEPTED.
// Отдаёт model.PlanStatusResponse. Недопустимый переход -> ErrConflict (409).
func (s *Plan) Accept(ctx context.Context, planID int64) (*model.PlanStatusResponse, error) {
	// TODO(шаг 12): проверка CanTransitionTo + атомарный UpdateStatus
	return nil, nil
}

// Reject — POST /api/plans/{planId}/reject. PROPOSED -> REJECTED.
// После reject операции плана не занимают станки (OccupiesCapacity=false).
func (s *Plan) Reject(ctx context.Context, planID int64) (*model.PlanStatusResponse, error) {
	// TODO(шаг 12)
	return nil, nil
}

// CancelByOrder — POST /api/plans/orders/{orderId}/cancel.
// PROPOSED/ACCEPTED -> CANCELLED (по orderId, не по planId).
func (s *Plan) CancelByOrder(ctx context.Context, orderID int64) (*model.PlanStatusResponse, error) {
	// TODO(шаг 12): CancelActiveByOrderID
	return nil, nil
}

// buildPlanResponse — внутренний маппинг домена в ProductionPlanResponse
// (общий для Create и GetByOrderID).
func (s *Plan) buildPlanResponse(ctx context.Context, plan *model.ProductionPlan) (*model.ProductionPlanResponse, error) {
	// TODO(шаг 11)
	return nil, nil
}
