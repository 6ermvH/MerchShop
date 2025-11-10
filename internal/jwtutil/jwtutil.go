package jwtutil

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type HS256 struct {
	secret   []byte
	iss, aud string
}

func NewHS256(secret, iss, aud string) *HS256 {
	return &HS256{[]byte(secret), iss, aud}
}

func (hs *HS256) Sign(sub uuid.UUID, name string) (string, error) {
	claims := jwt.MapClaims{
		"sub": sub.String(), "name": name,
		"iss": hs.iss, "aud": hs.aud,
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tok, err := token.SignedString(hs.secret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return tok, nil
}

var (
	errUnexpectedAlg = errors.New("unexpected alg")
	errBadClaims     = errors.New("bad claims type")
)

func (hs *HS256) Parse(token string) (jwt.MapClaims, error) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errUnexpectedAlg
		}

		return hs.secret, nil
	}, jwt.WithIssuer(hs.iss), jwt.WithAudience(hs.aud))
	if err != nil || !tok.Valid {
		return nil, fmt.Errorf("invalid JWT: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT: %w", errBadClaims)
	}

	return claims, nil
}
