package router

import (
	"log/slog"

	"github.com/domninivan/project__K/GoBackend/internal/handler"
	"github.com/domninivan/project__K/GoBackend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	WorkCenter *handler.WorkCenter
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

	return r
}
