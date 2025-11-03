CREATE INDEX IF NOT EXISTS transfers_from_created_idx
  ON merch_shop.transfers (from_user_id, created_at DESC) INCLUDE (amount);

CREATE INDEX IF NOT EXISTS transfers_to_created_idx
  ON merch_shop.transfers (to_user_id, created_at DESC) INCLUDE (amount);

CREATE INDEX IF NOT EXISTS orders_user_created_idx
  ON merch_shop.orders (user_id, created_at DESC) INCLUDE (product_id);