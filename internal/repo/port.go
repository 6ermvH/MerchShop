//go:generate mockgen -source=port.go -destination=../../gen/mock/repo/mock_repo.go -package=mock_repo -self_package=github.com/6ermvH/MerchShop/internal/repo

package repo

import (
	"context"

	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/google/uuid"
)

type MerchRepo interface { //nolint:interfacebloat
	FindUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
	FindUserByUsername(ctx context.Context, username string) (model.User, error)
	CreateUser(ctx context.Context, username, passwordHash string) (model.User, error)
	AddToBalance(ctx context.Context, userId uuid.UUID, delta int64) (model.User, error)

	FindProductByTitle(ctx context.Context, title string) (model.Product, error)

	CreateOrder(ctx context.Context, userId, productId uuid.UUID) (model.Order, error)
	FindOrdersByUserID(ctx context.Context, userId uuid.UUID) ([]model.Order, error)

	CreateTransfer(
		ctx context.Context,
		fromID, toID uuid.UUID,
		amount int64,
	) (model.Transfer, error)
	FindTransfersFromID(ctx context.Context, fromID uuid.UUID) ([]model.Transfer, error)
	FindTransfersToID(ctx context.Context, toID uuid.UUID) ([]model.Transfer, error)

	SendCoins(ctx context.Context, fromID, toID uuid.UUID, amount int64) error
	BuyProduct(ctx context.Context, userId uuid.UUID, productTitle string) error
}
