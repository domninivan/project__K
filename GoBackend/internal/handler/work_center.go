package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

type WorkCenterService interface {
	List(ctx context.Context) ([]model.WorkCenter, error)
	Create(ctx context.Context, req model.CreateWorkCenterRequest) (*model.WorkCenter, error)
}

type WorkCenter struct {
	svc WorkCenterService
}

func NewWorkCenter(svc WorkCenterService) *WorkCenter {
	return &WorkCenter{svc: svc}
}

// GET /api/v1/work-centers
func (h *WorkCenter) ListWorkCenters(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list work centers"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /api/v1/work-centers
func (h *WorkCenter) CreateWorkCenter(c *gin.Context) {
	var req model.CreateWorkCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wc, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": "work center with this code already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create work center"})
		}
		return
	}
	c.JSON(http.StatusCreated, wc)
}
