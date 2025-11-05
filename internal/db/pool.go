package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustPool(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	return pool
}
