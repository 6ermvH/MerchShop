package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

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

type txKey struct{}

func (r *Repo) runner(ctx context.Context) Runner {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return tx
	}

	return r.db
}

type TxOptions struct {
	Level          pgx.TxIsoLevel
	MaxRetries     int
	AttemptTimeout time.Duration
}

func (r *Repo) WithTx(
	ctx context.Context,
	fn func(txCtx context.Context) error,
	opts *TxOptions,
) error {
	level, retries, timeout := normalizeTxOpts(opts)

	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return fn(ctx)
	}

	for attempt := 0; attempt <= retries; attempt++ {
		inner, cancel := withMaybeTimeout(ctx, timeout)
		err := r.doTxAttempt(inner, fn, level)

		cancel()

		if shouldRetry(err, level, attempt, retries) {
			continue
		}

		return err
	}

	return nil
}

func (r *Repo) doTxAttempt(
	ctx context.Context,
	fn func(txCtx context.Context) error,
	level pgx.TxIsoLevel,
) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: level})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("tx fn: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func normalizeTxOpts(opts *TxOptions) (pgx.TxIsoLevel, int, time.Duration) {
	level := pgx.ReadCommitted
	retries := 5
	timeout := time.Duration(0)

	if opts == nil {
		return level, retries, timeout
	}

	if opts.Level != pgx.TxIsoLevel(strconv.Itoa(0)) {
		level = opts.Level
	}

	if opts.MaxRetries > 0 {
		retries = opts.MaxRetries
	}

	if opts.AttemptTimeout > 0 {
		timeout = opts.AttemptTimeout
	}

	return level, retries, timeout
}

func withMaybeTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, d)
}

func shouldRetry(err error, level pgx.TxIsoLevel, attempt, retries int) bool {
	return err != nil && isSerialization(err) && level == pgx.Serializable && attempt < retries
}

func isSerialization(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == "40001"
}
