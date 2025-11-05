package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Runner interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type txKey struct{}

type Repo struct {
	DB            *pgxpool.Pool
	UsersRepo     *UsersRepo
	OrdersRepo    *OrdersRepo
	TransfersRepo *TransfersRepo
	ProductsRepo  *ProductsRepo
}

func NewRepo(db *pgxpool.Pool) *Repo {
	r := &Repo{DB: db}
	r.UsersRepo = NewUsersRepo(db)
	r.OrdersRepo = NewOrdersRepo(db)
	r.TransfersRepo = NewTransfersRepo(db)
	r.ProductsRepo = NewProductsRepo(db)
	return r
}

func (r *Repo) runner(ctx context.Context) Runner {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return tx
	}
	return r.DB
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
	level := pgx.ReadCommitted
	retries := 5
	timeout := time.Duration(0)
	if opts != nil {
		level = opts.Level
		if opts.MaxRetries > 0 {
			retries = opts.MaxRetries
		}
		if opts.AttemptTimeout > 0 {
			timeout = opts.AttemptTimeout
		}
	}
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return fn(ctx)
	}

	for attempt := 0; ; attempt++ {
		inner := ctx
		var cancel context.CancelFunc
		if timeout > 0 {
			inner, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		tx, err := r.DB.BeginTx(inner, pgx.TxOptions{IsoLevel: level})
		if err != nil {
			return err
		}

		err = fn(context.WithValue(inner, txKey{}, tx))
		if err != nil {
			_ = tx.Rollback(inner)
			if isSerialization(err) && attempt < retries && level == pgx.Serializable {
				continue
			}
			return err
		}
		if err := tx.Commit(inner); err != nil {
			if isSerialization(err) && attempt < retries && level == pgx.Serializable {
				continue
			}
			return err
		}
		return nil
	}
}

func isSerialization(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "40001"
}
