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

	api := r.Group("/api")

	api.GET("/live", h.Health.Live)
	api.GET("/ready", h.Health.Ready)

	api.GET("/work-centers", h.WorkCenter.ListWorkCenters)
	api.POST("/work-centers", h.WorkCenter.CreateWorkCenter)
	api.GET("/work-centers/:id/schedule", h.WorkCenter.GetSchedule)
	api.PUT("/work-centers/:id/schedule", h.WorkCenter.ReplaceSchedule)
	api.GET("/work-centers/:id/availability", h.WorkCenter.GetAvailability)

	api.POST("/plans", h.Plan.CreatePlan)
	api.GET("/plans/orders/:orderId", h.Plan.GetOrderPlan)
	api.POST("/plans/:planId/accept", h.Plan.AcceptPlan)
	api.POST("/plans/:planId/reject", h.Plan.RejectPlan)
	api.POST("/plans/orders/:orderId/cancel", h.Plan.CancelOrderPlan)

	api.GET("/schedule", h.Schedule.GetSchedule)

	return r
}
