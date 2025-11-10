package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Balance      int64
	PasswordHash string
	CreatedAt    time.Time
}

type Transfer struct {
	ID           uuid.UUID
	FromUserID   uuid.UUID
	FromUserName string
	ToUserID     uuid.UUID
	ToUserName   string
	Amount       int64
	CreatedAt    time.Time
}

type Order struct {
	ID           uuid.UUID
	Count        int32
	UserID       uuid.UUID
	ProductID    uuid.UUID
	ProductTitle string
	ProductPrice int64
	CreatedAt    time.Time
}

type Product struct {
	ID    uuid.UUID
	Title string
	Price int64
}
