package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/http/middleware"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
)

func (api *API) ApiSendCoinPost(c *gin.Context) {
	var request openapi.SendCoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Errors: "bad payload"})
		return
	}
	if strings.TrimSpace(request.ToUser) == "" || request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Errors: "bad payload"})
		return
	}

	userRaw, ok := c.Get(middleware.CtxUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{Errors: "no user in context"})
		return
	}
	user := userRaw.(model.User)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	to, err := api.repos.FindUserByUsername(ctx, request.ToUser)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, openapi.ErrorResponse{Errors: "receiver not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "db error"})
		return
	}

	if err := api.repos.SendCoins(ctx, user.ID, to.ID, int64(request.Amount)); err != nil {
		switch {
		case strings.Contains(err.Error(), "insufficient funds"):
			c.JSON(
				http.StatusUnprocessableEntity,
				openapi.ErrorResponse{Errors: "insufficient funds"},
			)
		default:
			c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}
