package repo

import (
	"context"

	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransfersRepo struct{ db *pgxpool.Pool }

func NewTransfersRepo(db *pgxpool.Pool) *TransfersRepo { return &TransfersRepo{db: db} }

func (r *TransfersRepo) runner(ctx context.Context) Runner {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok && tx != nil {
		return tx
	}
	return r.db
}

func (r *TransfersRepo) Create(
	ctx context.Context,
	fromID, toID uuid.UUID,
	amount int64,
) (model.Transfer, error) {
	q := r.runner(ctx)
	var t model.Transfer
	err := q.QueryRow(ctx, `
		INSERT INTO merch_shop.transfers (from_user_id, to_user_id, amount)
		VALUES ($1, $2, $3)
		RETURNING id, from_user_id, to_user_id, amount, created_at
	`, fromID, toID, amount).Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.CreatedAt)
	return t, err
}
