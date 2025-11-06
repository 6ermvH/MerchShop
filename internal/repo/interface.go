package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Runner interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Tx interface {
	Runner
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DB interface {
	Runner
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}
