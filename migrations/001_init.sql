CREATE SCHEMA IF NOT EXISTS merch_shop;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS merch_shop.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
);

CREATE TABLE IF NOT EXISTS merch_shop.transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
);

CREATE TABLE IF NOT EXISTS merch_shop.orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
);

CREATE TABLE IF NOT EXISTS merch_shop.products (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title        VARCHAR(64) NOT NULL UNIQUE,
  price  BIGINT NOT NULL CHECK (price >= 0),
);
INSERT INTO merch_shop.products (title, price) VALUES
  ('t-shirt',     80),
  ('cup',         20),
  ('book',        50),
  ('pen',         10),
  ('powerbank',  200),
  ('hoody',      300),
  ('umbrella',   200),
  ('socks',       10),
  ('wallet',      50),
  ('pink-hoody', 500)