package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/http/middleware"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/gin-gonic/gin"
)

// TODO: add logs, error wrapping
func (api *API) ApiBuyItemGet(c *gin.Context) {
	product := c.Param("item")
	if product == "" {
		c.JSON(http.StatusBadRequest, openapi.ErrorResponse{Errors: "empty item"})
		return
	}

	userRaw, ok := c.Get(middleware.CtxUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{Errors: "no user in context"})
		return
	}
	user, _ := userRaw.(model.User)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := api.repos.BuyProduct(ctx, user.ID, product); err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "db error"})
		return
	}

	c.Status(http.StatusOK)
}
