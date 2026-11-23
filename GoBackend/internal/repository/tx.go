package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/domninivan/project__K/GoBackend/internal/port"
)

// txAdapter — оборачивает pgx.Tx в port.Tx: сервис видит только контракт,
// pgx не утекает за пределы repository. pgx.Row / pgx.Rows реализуют
// port.Row / port.Rows структурно (надмножества методов).
type txAdapter struct{ tx pgx.Tx }

func (a txAdapter) QueryRow(ctx context.Context, sql string, args ...any) port.Row {
	return a.tx.QueryRow(ctx, sql, args...)
}

func (a txAdapter) Query(ctx context.Context, sql string, args ...any) (port.Rows, error) {
	rows, err := a.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// NewTx — оборачивает pgx.Tx в контракт port.Tx. Для тестов и сценариев,
// где транзакция открыта вручную (repository-интеграционные тесты пишут
// в БД через Begin/Rollback напрямую, без TxManager).
func NewTx(tx pgx.Tx) port.Tx { return txAdapter{tx: tx} }

// poolAdapter — тот же контракт port.Tx поверх пула: чтения вне транзакции
// (снимок занятости для планировщика) идут через него, чтобы repository
// везде работал с одним интерфейсом. Прямое присваивание пула port.Tx
// невозможно: Query возвращает pgx.Rows, а контракту нужен port.Rows —
// адаптер конвертирует возвращаемые значения.
type poolAdapter struct{ pool *pgxpool.Pool }

func (p poolAdapter) QueryRow(ctx context.Context, sql string, args ...any) port.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p poolAdapter) Query(ctx context.Context, sql string, args ...any) (port.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// TxManager — реализация контракта транзакций сервиса поверх pgxpool.Pool.
// Сервис вызывает WithTx и через возвращаемую ошибку решает судьбу
// транзакции: nil -> COMMIT, ошибка -> ROLLBACK. Сервис не может забыть
// откат — rollback выполняется адаптером на любом не-nil выходе из fn.
type TxManager struct{ pool *pgxpool.Pool }

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithTx — открывает транзакцию, передаёт её в fn как port.Tx, коммитит
// при nil. Отдельная ошибка коммита (обрыв связи, deadlock) возвращается
// как есть — данные при этом откатаны сервером вместе с транзакцией.
func (m *TxManager) WithTx(ctx context.Context, fn func(tx port.Tx) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	if err := fn(txAdapter{tx: tx}); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
