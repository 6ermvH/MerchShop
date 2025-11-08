package repo

import (
	"context"
	"errors"

	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repo) FindProductByTitle(ctx context.Context, title string) (model.Product, error) {
	q := r.runner(ctx)
	var p model.Product
	err := q.QueryRow(ctx, `
		SELECT id, title, price
		FROM merch_shop.products
		WHERE lower(title) = lower($1)
	`, title).Scan(&p.ID, &p.Title, &p.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, ErrNotFound
		}
		return model.Product{}, err
	}
	return p, nil
}
