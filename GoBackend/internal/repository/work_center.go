package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/domninivan/project__K/GoBackend/internal/model"
)

type WorkCenter struct{ pool *pgxpool.Pool }

func NewWorkCenter(pool *pgxpool.Pool) *WorkCenter {
	return &WorkCenter{pool: pool}
}

func (r *WorkCenter) List(ctx context.Context) ([]model.WorkCenter, error) {
	const q = `
		SELECT wc.id, wc.code, wc.title, wc.active,
		       COALESCE(
		           array_agg(wcot.operation_type ORDER BY wcot.operation_type)
		           FILTER (WHERE wcot.operation_type IS NOT NULL),
		           '{}'
		       )
		  FROM work_center wc
		  LEFT JOIN work_center_operation_type wcot
		         ON wcot.work_center_id = wc.id
		 GROUP BY wc.id, wc.code, wc.title, wc.active
		 ORDER BY wc.id
	`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query work_centers: %w", err)
	}
	defer rows.Close()

	out := make([]model.WorkCenter, 0)
	for rows.Next() {
		var wc model.WorkCenter
		if err := rows.Scan(&wc.ID, &wc.Code, &wc.Title, &wc.Active, &wc.OperationTypes); err != nil {
			return nil, fmt.Errorf("scan work_center: %w", err)
		}
		out = append(out, wc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate work_centers: %w", err)
	}
	return out, nil
}

func (r *WorkCenter) Create(ctx context.Context, req model.CreateWorkCenterRequest) (*model.WorkCenter, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var wc model.WorkCenter
	err = tx.QueryRow(ctx, `
		INSERT INTO work_center (code, title, active)
		VALUES ($1, $2, $3)
		RETURNING id, code, title, active
	`, req.Code, req.Title, req.IsActiveOrDefault()).Scan(
		&wc.ID, &wc.Code, &wc.Title, &wc.Active,
	)
	if err != nil {
		return nil, mapPgError(err, "insert work_center")
	}

	if len(req.OperationTypes) > 0 {
		if _, err = tx.Exec(ctx, `
			INSERT INTO work_center_operation_type (work_center_id, operation_type)
			SELECT $1, unnest($2::varchar[])
		`, wc.ID, req.OperationTypes); err != nil {
			return nil, mapPgError(err, "insert operation_types")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	wc.OperationTypes = req.OperationTypes
	return &wc, nil
}

func mapPgError(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return fmt.Errorf("%w: %s", model.ErrConflict, pgErr.ConstraintName)
	}
	return fmt.Errorf("%s: %w", op, err)
}
