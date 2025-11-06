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

func (r *TransfersRepo) FindByFromId(ctx context.Context, id uuid.UUID) ([]model.Transfer, error) {
	q := r.runner(ctx)
	rows, err := q.Query(ctx, `
		SELECT t.id, t.from_user_id, u.username, t.to_user_id, t.amount, t.created_at
		FROM merch_shop.transfers AS t
		JOIN merch_shop.users AS u ON t.from_user_id = u.id
		WHERE t.from_user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []model.Transfer
	for rows.Next() {
		var t model.Transfer
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.FromUserName, &t.ToUserID, &t.Amount, &t.CreatedAt); err != nil {
			return transfers, err
		}
		transfers = append(transfers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *TransfersRepo) FindByToId(ctx context.Context, id uuid.UUID) ([]model.Transfer, error) {
	q := r.runner(ctx)
	rows, err := q.Query(ctx, `
		SELECT t.id, t.from_user_id, t.to_user_id, u.username, t.amount, t.created_at
		FROM merch_shop.transfers AS t
		JOIN merch_shop.users AS u ON t.to_user_id = u.id
		WHERE t.to_user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []model.Transfer
	for rows.Next() {
		var t model.Transfer
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.ToUserName, &t.Amount, &t.CreatedAt); err != nil {
			return transfers, err
		}
		transfers = append(transfers, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transfers, nil
}
