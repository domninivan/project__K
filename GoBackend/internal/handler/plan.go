package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Plan struct {
	log *slog.Logger
}

func NewPlan(log *slog.Logger) *Plan {
	return &Plan{log: log}
}

func (h *Plan) CreatePlan(c *gin.Context) {

}

func (h *Plan) GetOrderPlan(c *gin.Context) {

}

func (h *Plan) AcceptPlan(c *gin.Context) {

}

func (h *Plan) RejectPlan(c *gin.Context) {

}

func (n *Plan) CancelOrderPlan(c *gin.Context) {

}
