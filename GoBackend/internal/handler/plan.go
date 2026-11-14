package handler

import "log/slog"

type Plan struct {
	log *slog.Logger
}

func NewPlan(log *slog.Logger) *Plan {
	return &Plan{log: log}
}

func (h *Plan) CreatePlan() {

}
