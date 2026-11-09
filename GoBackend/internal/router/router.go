package router

import (
	"github.com/domninivan/project__K/GoBackend/internal/handler"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	WorkCenter *handler.WorkCenter
	Health     *handler.Health
}

func New(h Handlers) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/live", h.Health.Live)
	v1.GET("/ready", h.Health.Ready)

	v1.GET("/work-centers", h.WorkCenter.ListWorkCenters)

	v1.POST("/work-centers", h.WorkCenter.CreateWorkCenter)

	return r
}
