package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Schedule struct {
	log *slog.Logger
}

func NewSchedule(log *slog.Logger) *Schedule {
	return &Schedule{log: log}
}

func (h *Schedule) GetSchedule(c *gin.Context) {

}
