package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db DB
}

func NewRepo(db DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SendCoins(ctx context.Context, fromUserId, toUserId uuid.UUID, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if fromUserId == toUserId {
		return fmt.Errorf("cannot transfer to self")
	}

	return r.WithTx(ctx, func(txCtx context.Context) error {
		if _, err := r.AddToBalance(txCtx, fromUserId, -amount); err != nil {
			return err
		}
		if _, err := r.AddToBalance(txCtx, toUserId, +amount); err != nil {
			return err
		}
		_, err := r.CreateTransfer(txCtx, fromUserId, toUserId, amount)
		return err
	}, &TxOptions{Level: pgx.Serializable, MaxRetries: 10})
}

func (r *Repo) BuyProduct(
	ctx context.Context,
	userId uuid.UUID,
	productTitle string,
) error {
	return r.WithTx(ctx, func(txCtx context.Context) error {
		product, err := r.FindProductByTitle(txCtx, productTitle)
		if err != nil {
			return err
		}
		if _, err := r.AddToBalance(txCtx, userId, -product.Price); err != nil {
			return err
		}
		if _, err := r.CreateOrder(txCtx, userId, product.ID); err != nil {
			return err
		}
		return nil
	}, &TxOptions{Level: pgx.Serializable, MaxRetries: 10})
}
