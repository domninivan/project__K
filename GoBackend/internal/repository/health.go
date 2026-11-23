package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Health struct {
	pool *pgxpool.Pool
}

func NewHealth(pool *pgxpool.Pool) *Health {
	return &Health{pool: pool}
}

func (r *Health) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}
