package router

import (
	"log/slog"

	"github.com/domninivan/project__K/GoBackend/internal/handler"
	"github.com/domninivan/project__K/GoBackend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	WorkCenter *handler.WorkCenter
	Plan       *handler.Plan
	Schedule   *handler.Schedule
	Health     *handler.Health
}

func New(h Handlers, logger *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logging(logger), middleware.Recovery(logger))

	v1 := r.Group("/api/v1")

	v1.GET("/live", h.Health.Live)
	v1.GET("/ready", h.Health.Ready)

	v1.GET("/work-centers", h.WorkCenter.ListWorkCenters)
	v1.POST("/work-centers", h.WorkCenter.CreateWorkCenter)
	v1.GET("/work-centers/:id/schedule", h.WorkCenter.GetSchedule)
	v1.PUT("/work-centers/:id/schedule", h.WorkCenter.ReplaceSchedule)
	v1.GET("/work-centers/:id/availability", h.WorkCenter.GetAvailability)

	v1.POST("/plans", h.Plan.CreatePlan)
	v1.GET("/plans/orders/:orderId", h.Plan.GetOrderPlan)
	v1.POST("/plans/:planId/accept", h.Plan.AcceptPlan)
	v1.POST("/plans/:planId/reject", h.Plan.RejectPlan)
	v1.POST("/plans/orders/:orderId/cancel", h.Plan.CancelOrderPlan)

	v1.GET("/schedule", h.Schedule.GetSchedule)

	return r
}
