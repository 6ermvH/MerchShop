package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	mock_repo "github.com/6ermvH/MerchShop/gen/mock/repo"
	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/http/middleware"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func withUser(u model.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxUserKey, u)
		c.Next()
	}
}

func TestInfo_NoUserInContext_401(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	api := NewAPI(repoMock, nil)

	r := gin.New()
	r.GET("/api/info", api.ApiInfoGet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestInfo_FindOrdersError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := model.User{ID: uuid.New(), Username: "u", Balance: 100, CreatedAt: time.Now()}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	repoMock.EXPECT().
		FindOrdersByUserID(gomock.Any(), user.ID).
		Return(nil, errors.New("db"))

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.GET("/api/info", withUser(user), api.ApiInfoGet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestInfo_FindTransfersFromError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := model.User{ID: uuid.New(), Username: "u", Balance: 100}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindOrdersByUserID(gomock.Any(), user.ID).
		Return([]model.Order{{ProductTitle: "coffee", ProductPrice: 50, CreatedAt: time.Now()}}, nil)

	repoMock.EXPECT().
		FindTransfersFromID(gomock.Any(), user.ID).
		Return(nil, errors.New("db"))

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.GET("/api/info", withUser(user), api.ApiInfoGet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestInfo_FindTransfersToError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := model.User{ID: uuid.New(), Username: "u", Balance: 100}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindOrdersByUserID(gomock.Any(), user.ID).
		Return([]model.Order{{ProductTitle: "coffee", ProductPrice: 50, CreatedAt: time.Now()}}, nil)

	repoMock.EXPECT().
		FindTransfersFromID(gomock.Any(), user.ID).
		Return([]model.Transfer{}, nil)

	repoMock.EXPECT().
		FindTransfersToID(gomock.Any(), user.ID).
		Return(nil, errors.New("db"))

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.GET("/api/info", withUser(user), api.ApiInfoGet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestInfo_OK_Aggregates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := model.User{ID: uuid.New(), Username: "u", Balance: 130}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	orders := []model.Order{
		{ProductTitle: "coffee", ProductPrice: 50},
		{ProductTitle: "coffee", ProductPrice: 50},
		{ProductTitle: "tea", ProductPrice: 30},
	}
	repoMock.EXPECT().
		FindOrdersByUserID(gomock.Any(), user.ID).
		Return(orders, nil)

	recv := []model.Transfer{
		{FromUserName: "alice", Amount: 10},
		{FromUserName: "bob", Amount: 5},
	}
	repoMock.EXPECT().
		FindTransfersFromID(gomock.Any(), user.ID).
		Return(recv, nil)

	sent := []model.Transfer{
		{ToUserName: "charlie", Amount: 7},
	}
	repoMock.EXPECT().
		FindTransfersToID(gomock.Any(), user.ID).
		Return(sent, nil)

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.GET("/api/info", withUser(user), api.ApiInfoGet)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp openapi.InfoResponse

	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Equal(t, int32(user.Balance), resp.Coins)

	sort.Slice(
		resp.Inventory,
		func(i, j int) bool { return resp.Inventory[i].Type < resp.Inventory[j].Type },
	)
	require.Equal(t, []openapi.InfoResponseInventoryInner{
		{Type: "coffee", Quantity: 2},
		{Type: "tea", Quantity: 1},
	}, resp.Inventory)

	sort.Slice(resp.CoinHistory.Received, func(i, j int) bool {
		return resp.CoinHistory.Received[i].FromUser < resp.CoinHistory.Received[j].FromUser
	})
	require.Equal(t, []openapi.InfoResponseCoinHistoryReceivedInner{
		{FromUser: "alice", Amount: 10},
		{FromUser: "bob", Amount: 5},
	}, resp.CoinHistory.Received)

	require.Equal(t, []openapi.InfoResponseCoinHistorySentInner{
		{ToUser: "charlie", Amount: 7},
	}, resp.CoinHistory.Sent)
}
