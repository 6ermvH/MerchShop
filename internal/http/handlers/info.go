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

func (api *API) ApiInfoGet(c *gin.Context) {
	userRaw, ok := c.Get(middleware.CtxUserKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, openapi.ErrorResponse{Errors: "no user in context"})
		return
	}
	user := userRaw.(model.User)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	orders, err := api.repos.OrdersRepo.FindByUserId(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "db error"})
		return
	}
	inventory := makeInfoResponseInventory(orders)

	recv, err := api.repos.TransfersRepo.FindByFromId(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "db error"})
		return
	}
	coinHistoryFrom := make([]openapi.InfoResponseCoinHistoryReceivedInner, 0)
	for _, rec := range recv {
		coinHistoryFrom = append(coinHistoryFrom,
			openapi.InfoResponseCoinHistoryReceivedInner{
				FromUser: rec.FromUserName,
				Amount:   int32(rec.Amount),
			})
	}

	sent, err := api.repos.TransfersRepo.FindByToId(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, openapi.ErrorResponse{Errors: "db error"})
		return
	}
	coinHistoryTo := make([]openapi.InfoResponseCoinHistorySentInner, 0)
	for _, sen := range sent {
		coinHistoryTo = append(coinHistoryTo,
			openapi.InfoResponseCoinHistorySentInner{
				ToUser: sen.ToUserName,
				Amount: int32(sen.Amount),
			})
	}

	response := openapi.InfoResponse{
		Coins:     int32(user.Balance),
		Inventory: inventory,
		CoinHistory: openapi.InfoResponseCoinHistory{
			Received: coinHistoryFrom,
			Sent:     coinHistoryTo,
		},
	}

	c.JSON(http.StatusOK, response)
}

type productName string

func makeInfoResponseInventory(orders []model.Order) []openapi.InfoResponseInventoryInner {
	inventory := make([]openapi.InfoResponseInventoryInner, 0)
	productCounter := make(map[productName]int)

	for _, order := range orders {
		productCounter[productName(order.ProductTitle)]++
	}

	for title, count := range productCounter {
		inventory = append(inventory,
			openapi.InfoResponseInventoryInner{
				Type:     string(title),
				Quantity: int32(count),
			})
	}

	return inventory
}
