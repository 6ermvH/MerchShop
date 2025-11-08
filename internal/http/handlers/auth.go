package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/hasher"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
)

// TODO: add logs
func (api *API) ApiAuthPost(c *gin.Context) {
	var request openapi.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Errors: "bad payload"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	user, err := api.repos.FindUserByUsername(ctx, request.Username)
	switch err {
	case nil:
		if err := hasher.CheckPassword(user.PasswordHash, request.Password); err != nil {
			c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{Errors: "invalid credentials"})
			return
		}
	case repo.ErrNotFound:
		hash, err := hasher.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "hash error"})
			return
		}
		user, err = api.repos.CreateUser(ctx, request.Username, hash)
		if err != nil {
			c.JSON(http.StatusConflict, openapi.ErrorResponse{Errors: "username already exists"})
			return
		}
	default:
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "db error"})
		return
	}

	tok, err := api.hs.Sign(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "sign JWT"})
		return
	}

	c.JSON(http.StatusOK, openapi.AuthResponse{
		Token: tok,
	})
}
