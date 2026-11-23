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
	svc         WorkCenterService
	scheduleSvc ScheduleService
	log         *slog.Logger
}

func NewWorkCenter(svc WorkCenterService, scheduleSvc ScheduleService, log *slog.Logger) *WorkCenter {
	return &WorkCenter{svc: svc, scheduleSvc: scheduleSvc, log: log}
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

// GetSchedule — GET /api/work-centers/{id}/schedule: недельный шаблон.
func (h *WorkCenter) GetSchedule(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	resp, err := h.scheduleSvc.WorkCenterSchedule(c.Request.Context(), id)
	if err != nil {
		h.respondScheduleError(c, log, "get work center schedule", err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ReplaceSchedule — PUT /api/work-centers/{id}/schedule: полная замена
// недельного шаблона (дни валидирует binding + домен в сервисе).
func (h *WorkCenter) ReplaceSchedule(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	var req model.ReplaceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid request body", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.scheduleSvc.ReplaceWorkCenterSchedule(c.Request.Context(), id, req); err != nil {
		h.respondScheduleError(c, log, "replace work center schedule", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workCenterId": id, "days": req.Days})
}

// GetAvailability — GET /api/work-centers/{id}/availability?from=&to=:
// раскрытие шаблона в конкретные даты (from/to = YYYY-MM-DD).
func (h *WorkCenter) GetAvailability(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	id, ok := pathID(c, "id")
	if !ok {
		return
	}
	from, ok := queryDate(c, log, "from")
	if !ok {
		return
	}
	to, ok := queryDate(c, log, "to")
	if !ok {
		return
	}
	resp, err := h.scheduleSvc.Availability(c.Request.Context(), id, from, to)
	if err != nil {
		h.respondScheduleError(c, log, "get availability", err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// queryDate — обязательный query-параметр даты YYYY-MM-DD; кривой -> 400.
func queryDate(c *gin.Context, log *slog.Logger, name string) (model.Date, bool) {
	raw := c.Query(name)
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing " + name + " query parameter (YYYY-MM-DD)"})
		return model.Date{}, false
	}
	d, err := model.ParseDate(raw)
	if err != nil {
		log.Warn("invalid query parameter", "name", name, "value", raw, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name + " query parameter: " + err.Error()})
		return model.Date{}, false
	}
	return d, true
}

// respondScheduleError — маппинг доменных ошибок расписаний на HTTP:
// кривой день/период -> 400, нет станка -> 404, прочее -> 500.
func (h *WorkCenter) respondScheduleError(c *gin.Context, log *slog.Logger, action string, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidScheduleDay),
		errors.Is(err, model.ErrInvalidDate):
		log.Warn(action+" rejected", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrNotFound):
		log.Warn(action+" rejected", "err", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "work center not found"})
	default:
		log.Error(action+" failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process schedule"})
	}
}
