package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/domninivan/project__K/GoBackend/internal/middleware"
	"github.com/domninivan/project__K/GoBackend/internal/model"
)

type WorkCenterService interface {
	List(ctx context.Context) ([]model.WorkCenter, error)
	Create(ctx context.Context, req model.CreateWorkCenterRequest) (*model.WorkCenter, error)
}

type WorkCenter struct {
	svc WorkCenterService
	log *slog.Logger
}

func NewWorkCenter(svc WorkCenterService, log *slog.Logger) *WorkCenter {
	return &WorkCenter{svc: svc, log: log}
}

// GET /api/v1/work-centers
func (h *WorkCenter) ListWorkCenters(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		log.Error("list work centers failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list work centers"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /api/v1/work-centers
func (h *WorkCenter) CreateWorkCenter(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	var req model.CreateWorkCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid request body", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wc, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrConflict):
			log.Warn("work center code conflict", "err", err)
			c.JSON(http.StatusConflict, gin.H{"error": "work center with this code already exists"})
		default:
			log.Error("create work center failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create work center"})
		}
		return
	}
	c.JSON(http.StatusCreated, wc)
}
