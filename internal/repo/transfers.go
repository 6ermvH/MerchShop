package repo

import (
	"context"
	"fmt"

	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/google/uuid"
)

func (r *Repo) CreateTransfer(
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

	return t, fmt.Errorf("create transfer: %w", err)
}

func (r *Repo) FindTransfersFromID(ctx context.Context, id uuid.UUID) ([]model.Transfer, error) {
	q := r.runner(ctx)

	rows, err := q.Query(ctx, `
		SELECT t.id, t.from_user_id, u.username, t.to_user_id, t.amount, t.created_at
		FROM merch_shop.transfers AS t
		JOIN merch_shop.users AS u ON t.from_user_id = u.id
		WHERE t.from_user_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get query sql: %w", err)
	}

	defer rows.Close()

	var transfers []model.Transfer

	for rows.Next() {
		var t model.Transfer
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.FromUserName, &t.ToUserID, &t.Amount, &t.CreatedAt); err != nil {
			return transfers, fmt.Errorf("check next row: %w", err)
		}

		transfers = append(transfers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check row: %w", err)
	}

	return transfers, nil
}

func (r *Repo) FindTransfersToID(ctx context.Context, id uuid.UUID) ([]model.Transfer, error) {
	q := r.runner(ctx)

	rows, err := q.Query(ctx, `
		SELECT t.id, t.from_user_id, t.to_user_id, u.username, t.amount, t.created_at
		FROM merch_shop.transfers AS t
		JOIN merch_shop.users AS u ON t.to_user_id = u.id
		WHERE t.to_user_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get query sql: %w", err)
	}

	defer rows.Close()

	var transfers []model.Transfer

	for rows.Next() {
		var t model.Transfer
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.ToUserName, &t.Amount, &t.CreatedAt); err != nil {
			return transfers, fmt.Errorf("check next row: %w", err)
		}

		transfers = append(transfers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("check row: %w", err)
	}

	return transfers, nil
}
