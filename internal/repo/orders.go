package repo

import (
	"context"

	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/google/uuid"
)

func (r *Repo) CreateOrder(
	ctx context.Context,
	userId, productId uuid.UUID,
) (model.Order, error) {
	q := r.runner(ctx)
	var o model.Order
	err := q.QueryRow(ctx, `
		INSERT INTO merch_shop.orders (user_id, product_id, count)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, product_id, created_at
	`, userId, productId).Scan(&o.ID, &o.UserID, &o.ProductID, &o.CreatedAt)
	return o, err
}

func (r *Repo) FindOrdersByUserID(ctx context.Context, userId uuid.UUID) ([]model.Order, error) {
	q := r.runner(ctx)
	rows, err := q.Query(ctx, `
		SELECT o.id, o.user_id, o.product_id, o.created_at,
		       p.title AS product_title, p.price AS product_price
		FROM merch_shop.orders AS o
		JOIN merch_shop.products AS p ON p.id = o.product_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.CreatedAt, &o.ProductTitle, &o.ProductPrice); err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
