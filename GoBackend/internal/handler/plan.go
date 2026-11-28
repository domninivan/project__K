package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/domninivan/project__K/GoBackend/internal/middleware"
	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// PlanService — контракт оркестрации планов глазами хендлера
// (реализация service.Plan).
type PlanService interface {
	Create(ctx context.Context, req model.PlanningRequest) (*model.ProductionPlanResponse, error)
	GetByOrderID(ctx context.Context, orderID int64) (*model.ProductionPlanResponse, error)
	Accept(ctx context.Context, planID int64) (*model.PlanStatusResponse, error)
	Reject(ctx context.Context, planID int64) (*model.PlanStatusResponse, error)
	CancelByOrder(ctx context.Context, orderID int64) (*model.PlanStatusResponse, error)
}

type Plan struct {
	svc PlanService
	log *slog.Logger
}

func NewPlan(svc PlanService, log *slog.Logger) *Plan {
	return &Plan{svc: svc, log: log}
}

// mapPlanError — единый маппинг доменных ошибок плана на HTTP-коды:
// 400 кривой запрос, 404 нет объекта, 409 конфликт домена/гонки, 500 прочее.
func mapPlanError(err error) (int, string) {
	switch {
	case errors.Is(err, model.ErrInvalidPlanningRequest),
		errors.Is(err, model.ErrInvalidDate),
		errors.Is(err, model.ErrInvalidScheduleDay):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound, "plan not found"
	case errors.Is(err, model.ErrConflict),
		errors.Is(err, model.ErrNoFreeWindow),
		errors.Is(err, model.ErrPlanNotReschedulable),
		errors.Is(err, model.ErrInvalidStatusTransition):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}

// respondPlanError — лог + JSON-ответ по доменной ошибке. 4xx — warn
// (виноват клиент), 5xx — error (виноваты мы).
func (h *Plan) respondPlanError(c *gin.Context, log *slog.Logger, action string, err error) {
	code, msg := mapPlanError(err)
	if code >= 500 {
		log.Error(action+" failed", "err", err)
	} else {
		log.Warn(action+" rejected", "err", err)
	}
	c.JSON(code, gin.H{"error": msg})
}

// pathID — парсинг числового параметра пути (/plans/:planId, /orders/:orderId).
func pathID(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name + " path parameter"})
		return 0, false
	}
	return id, true
}

// CreatePlan — POST /api/plans. Успех -> 201 + ProductionPlanResponse.
func (h *Plan) CreatePlan(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	var req model.PlanningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid request body", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		h.respondPlanError(c, log, "create plan", err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetOrderPlan — GET /api/plans/orders/{orderId}. Успех -> 200.
func (h *Plan) GetOrderPlan(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	orderID, ok := pathID(c, "orderId")
	if !ok {
		return
	}
	resp, err := h.svc.GetByOrderID(c.Request.Context(), orderID)
	if err != nil {
		h.respondPlanError(c, log, "get plan by order", err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// AcceptPlan — POST /api/plans/{planId}/accept. PROPOSED -> ACCEPTED.
func (h *Plan) AcceptPlan(c *gin.Context) {
	h.transition(c, "accept", h.svc.Accept)
}

// RejectPlan — POST /api/plans/{planId}/reject. PROPOSED -> REJECTED.
func (h *Plan) RejectPlan(c *gin.Context) {
	h.transition(c, "reject", h.svc.Reject)
}

// transition — общий каркас accept/reject: парсинг planId -> вызов перехода.
func (h *Plan) transition(c *gin.Context, action string,
	call func(ctx context.Context, planID int64) (*model.PlanStatusResponse, error)) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	planID, ok := pathID(c, "planId")
	if !ok {
		return
	}
	resp, err := call(c.Request.Context(), planID)
	if err != nil {
		h.respondPlanError(c, log, action+" plan", err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CancelOrderPlan — POST /api/plans/orders/{orderId}/cancel.
func (h *Plan) CancelOrderPlan(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	orderID, ok := pathID(c, "orderId")
	if !ok {
		return
	}
	resp, err := h.svc.CancelByOrder(c.Request.Context(), orderID)
	if err != nil {
		h.respondPlanError(c, log, "cancel plan by order", err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
