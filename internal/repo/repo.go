package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) SendCoins(ctx context.Context, fromUserId, toUserId uuid.UUID, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if fromUserId == toUserId {
		return fmt.Errorf("cannot transfer to self")
	}

	return r.WithTx(ctx, func(txCtx context.Context) error {
		if _, err := r.UsersRepo.AddToBalance(txCtx, fromUserId, -amount); err != nil {
			return err
		}
		if _, err := r.UsersRepo.AddToBalance(txCtx, toUserId, +amount); err != nil {
			return err
		}
		_, err := r.TransfersRepo.Create(txCtx, fromUserId, toUserId, amount)
		return err
	}, &TxOptions{Level: pgx.Serializable, MaxRetries: 10})
}

func (r *Repo) BuyProduct(ctx context.Context, userId uuid.UUID, productTitle string, count int32) error {
	if count <= 0 {
		return fmt.Errorf("count must be positive")
	}
	return r.WithTx(ctx, func(txCtx context.Context) error {
		product, err := r.ProductsRepo.FindByTitle(txCtx, productTitle)
		if err != nil {
			return err
		}
		if _, err := r.UsersRepo.AddToBalance(txCtx, userId, -product.Price*int64(count)); err != nil {
			return err
		}
		if _, err := r.OrdersRepo.Create(txCtx, userId, product.ID, count); err != nil {
			return err
		}
		return nil
	}, &TxOptions{Level: pgx.Serializable, MaxRetries: 10})
}
