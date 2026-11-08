package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/domninivan/project__K/GoBackend/internal/model"
	"github.com/domninivan/project__K/GoBackend/internal/service"
)

type WorkCenterHandler struct{ svc *service.WorkCenterService }

func NewWorkCenterHandler(svc *service.WorkCenterService) *WorkCenterHandler {
	return &WorkCenterHandler{svc: svc}
}

// GET /api/v1/work-centers
func (h *WorkCenterHandler) ListWorkCentersHandler(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list work centers"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /api/v1/work-centers
func (h *WorkCenterHandler) CreateWorkCenterHandler(c *gin.Context) {
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
