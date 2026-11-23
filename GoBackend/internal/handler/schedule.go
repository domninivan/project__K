package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/domninivan/project__K/GoBackend/internal/middleware"
	"github.com/domninivan/project__K/GoBackend/internal/model"
)

// ScheduleService — контракт расписаний/Ганта глазами хендлера
// (реализация service.Schedule).
type ScheduleService interface {
	WorkCenterSchedule(ctx context.Context, workCenterID int64) (*model.WorkCenterScheduleResponse, error)
	ReplaceWorkCenterSchedule(ctx context.Context, workCenterID int64, req model.ReplaceScheduleRequest) error
	Availability(ctx context.Context, workCenterID int64, from, to model.Date) (*model.AvailabilityResponse, error)
	GeneralSchedule(ctx context.Context, from, to time.Time) (*model.ScheduleResponse, error)
}

type Schedule struct {
	svc ScheduleService
	log *slog.Logger
}

func NewSchedule(svc ScheduleService, log *slog.Logger) *Schedule {
	return &Schedule{svc: svc, log: log}
}

// GetSchedule — GET /api/schedule?from=&to= (Гант по всем станкам).
// from/to — RFC3339 со смещением, напр. 2026-11-21T00:00:00+03:00.
func (h *Schedule) GetSchedule(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	from, ok := queryTime(c, log, "from")
	if !ok {
		return
	}
	to, ok := queryTime(c, log, "to")
	if !ok {
		return
	}

	resp, err := h.svc.GeneralSchedule(c.Request.Context(), from, to)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidDate):
			log.Warn("invalid schedule period", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Error("get schedule failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load schedule"})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

// queryTime — обязательный RFC3339 query-параметр; кривой формат -> 400.
func queryTime(c *gin.Context, log *slog.Logger, name string) (time.Time, bool) {
	raw := c.Query(name)
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing " + name + " query parameter (RFC3339)"})
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		log.Warn("invalid query parameter", "name", name, "value", raw, "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name + " query parameter: " + err.Error()})
		return time.Time{}, false
	}
	return t, true
}
