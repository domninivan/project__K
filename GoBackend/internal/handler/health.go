package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/domninivan/project__K/GoBackend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HealthChecker interface {
	Check(ctx context.Context) error
}

type Health struct {
	svc HealthChecker
	log *slog.Logger
}

func NewHealth(svc HealthChecker, log *slog.Logger) *Health {
	return &Health{svc: svc, log: log}
}

func (h *Health) Ready(c *gin.Context) {
	log := h.log.With("request_id", middleware.RequestIDFrom(c))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.svc.Check(ctx); err != nil {
		log.Error("readiness check failed", "err", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database unavailable",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func (h *Health) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}
