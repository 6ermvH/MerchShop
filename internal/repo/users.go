package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

func (r *Repo) FindUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	q := r.runner(ctx)
	var u model.User
	err := q.QueryRow(ctx, `
		SELECT id, username, password_hash, balance, created_at
		FROM merch_shop.users WHERE id=$1
	`, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Balance, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	return u, nil
}

func (r *Repo) FindUserByUsername(ctx context.Context, username string) (model.User, error) {
	q := r.runner(ctx)
	var u model.User
	err := q.QueryRow(ctx, `
		SELECT id, username, password_hash, balance, created_at
		FROM merch_shop.users WHERE username=$1
	`, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Balance, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	return u, nil
}

func (r *Repo) AddToBalance(
	ctx context.Context,
	userId uuid.UUID,
	delta int64,
) (model.User, error) {
	q := r.runner(ctx)

	var cur int64
	if err := q.QueryRow(ctx, `
		SELECT balance FROM merch_shop.users WHERE id=$1 FOR UPDATE
	`, userId).Scan(&cur); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	newBal := cur + delta
	if newBal < 0 {
		return model.User{}, fmt.Errorf("insufficient funds")
	}

	var u model.User
	err := q.QueryRow(ctx, `
		UPDATE merch_shop.users
		SET balance=$2
		WHERE id=$1
		RETURNING id, username, password_hash, balance, created_at
	`, userId, newBal).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Balance, &u.CreatedAt)
	return u, err
}

func (r *Repo) CreateUser(ctx context.Context, username, passwordHash string) (model.User, error) {
	q := r.runner(ctx)
	var u model.User
	err := q.QueryRow(ctx, `
		INSERT INTO merch_shop.users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash, balance, created_at
	`, username, passwordHash).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Balance, &u.CreatedAt)
	return u, err
}
