package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"            json:"id"`
	Username     string    `db:"username"      json:"username"`
	Balance      int64     `db:"balance"       json:"balance"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

type Transfer struct {
	ID           uuid.UUID `db:"id"           json:"id"`
	FromUserID   uuid.UUID `db:"from_user_id" json:"from_user_id"`
	FromUserName string
	ToUserID     uuid.UUID `db:"to_user_id"   json:"to_user_id"`
	ToUserName   string
	Amount       int64     `db:"amount"       json:"amount"`
	CreatedAt    time.Time `db:"created_at"   json:"created_at"`
}

type Order struct {
	ID           uuid.UUID `db:"id"         json:"id"`
	Count        int32
	UserID       uuid.UUID `db:"user_id"    json:"user_id"`
	ProductID    uuid.UUID `db:"product_id" json:"product_id"`
	ProductTitle string
	ProductPrice int64
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Product struct {
	ID    uuid.UUID `db:"id"    json:"id"`
	Title string    `db:"title" json:"title"`
	Price int64     `db:"price" json:"price"`
}
