package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CtxUserKey = "auth_user"
)

func Auth(hs jwtutil.JWT, repository repo.MerchRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := parseBearer(c.GetHeader("Authorization"))
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		claims, err := hs.Parse(raw)
		if err != nil || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		name, ok := claims["name"]
		if !ok || name == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: no name"})
			return
		}

		sub, ok := claims["sub"]
		if !ok || sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: no sub"})
			return
		}

		userId, err := uuid.Parse(sub.(string))
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid token: bad user_id"},
			)
			return
		}

		user, err := repository.FindUserByID(ctx, userId)
		if err != nil && err != repo.ErrNotFound {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		// TODO: add Validate token
		if user.Username != name {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid token: name not match"},
			)
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
