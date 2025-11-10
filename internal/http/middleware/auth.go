package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const CtxUserKey = "auth_user"

func Auth(hs jwtutil.JWT, repository repo.MerchRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := parseBearer(c.GetHeader("Authorization"))
		if raw == "" {
			unauth(c, "missing token")

			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second) //nolint:mnd
		defer cancel()

		name, userID, ok := parseClaims(hs, raw)
		if !ok {
			unauth(c, "invalid token")

			return
		}

		user, ok, err := findUser(ctx, repository, userID)
		if err != nil {
			internal(c)

			return
		}

		if !ok {
			unauthJSON(c, openapi.ErrorResponse{Errors: "invalid token"})

			return
		}

		if user.Username != name {
			unauth(c, "invalid token: name not match")

			return
		}

		c.Set(CtxUserKey, user)
		c.Next()
	}
}

func parseBearer(h string) string {
	const p = "Bearer "
	if len(h) > len(p) && h[:len(p)] == p {
		return h[len(p):]
	}

	return ""
}

func parseClaims(hs jwtutil.JWT, raw string) (string, uuid.UUID, bool) {
	claims, err := hs.Parse(raw)
	if err != nil || claims == nil {
		return "", uuid.Nil, false
	}

	nameV, _ := claims["name"].(string)
	if strings.TrimSpace(nameV) == "" {
		return "", uuid.Nil, false
	}

	subV, _ := claims["sub"].(string)

	id, err := uuid.Parse(subV)
	if err != nil {
		return "", uuid.Nil, false
	}

	return nameV, id, true
}

func findUser(
	ctx context.Context,
	r repo.MerchRepo,
	id uuid.UUID,
) (model.User, bool, error) {
	user, err := r.FindUserByID(ctx, id)

	switch {
	case err == nil:
		return user, true, nil
	case errors.Is(err, repo.ErrNotFound):
		return model.User{}, false, nil
	default:
		return model.User{}, false, fmt.Errorf("find user: %w", err)
	}
}

func unauth(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
}

func unauthJSON(c *gin.Context, v any) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, v)
}

func internal(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db error"})
}
