package service

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type Health struct {
	repo Pinger
}

func NewHealth(repo Pinger) *Health {
	return &Health{repo: repo}
}

func (s *Health) Check(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
