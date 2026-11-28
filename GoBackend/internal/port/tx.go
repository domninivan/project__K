package port

import "context"

// Tx — транзакционный примитив глазами сервиса. Сервис решает, ГДЕ
// проходит граница бизнес-операции (TxManager.WithTx), но не знает, ЧЕМ
// физически является транзакция. Реализация — pgx.Tx через адаптер.
// Методов ровно столько, сколько использует репозиторий: QueryRow
// (вставки с RETURNING id) и Query (выборки занятости внутри tx).
type Tx interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

// Row — одна строка результата (подмножество pgx.Row).
type Row interface {
	Scan(dest ...any) error
}

// Rows — курсор результата (подмножество pgx.Rows).
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}
