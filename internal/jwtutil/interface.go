package jwtutil

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT interface {
	Sign(sub uuid.UUID, name string) (string, error)
	Parse(token string) (jwt.MapClaims, error)
}
