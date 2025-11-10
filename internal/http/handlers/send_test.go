package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_repo "github.com/6ermvH/MerchShop/gen/mock/repo"
	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSendCoin_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewAPI(mock_repo.NewMockMerchRepo(ctrl), nil)
	r := gin.New()
	r.POST("/api/sendCoin", api.ApiSendCoinPost)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendCoin_BadPayload_EmptyOrNonPositive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewAPI(mock_repo.NewMockMerchRepo(ctrl), nil)
	r := gin.New()
	r.POST("/api/sendCoin", withUser(model.User{ID: uuid.New(), Username: "me"}), api.ApiSendCoinPost)

	cases := []openapi.SendCoinRequest{
		{ToUser: "", Amount: 10},
		{ToUser: "   ", Amount: 10},
		{ToUser: "alice", Amount: 0},
		{ToUser: "alice", Amount: -5},
	}
	for i, cse := range cases {
		body, _ := json.Marshal(cse)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equalf(t, http.StatusBadRequest, w.Code, "case %d: %+v", i, cse)
	}
}

func TestSendCoin_NoUserInContext_401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := NewAPI(mock_repo.NewMockMerchRepo(ctrl), nil)
	r := gin.New()
	r.POST("/api/sendCoin", api.ApiSendCoinPost)

	body, _ := json.Marshal(openapi.SendCoinRequest{ToUser: "alice", Amount: 10})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSendCoin_ReceiverNotFound_404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	me := model.User{ID: uuid.New(), Username: "me", Balance: 100, CreatedAt: time.Now()}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindUserByUsername(gomock.Any(), "alice").
		Return(model.User{}, repo.ErrNotFound)

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.POST("/api/sendCoin", withUser(me), api.ApiSendCoinPost)

	body, _ := json.Marshal(openapi.SendCoinRequest{ToUser: "alice", Amount: 10})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestSendCoin_FindUserDBError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	me := model.User{ID: uuid.New(), Username: "me"}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindUserByUsername(gomock.Any(), "alice").
		Return(model.User{}, errors.New("db error"))

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.POST("/api/sendCoin", withUser(me), api.ApiSendCoinPost)

	body, _ := json.Marshal(openapi.SendCoinRequest{ToUser: "alice", Amount: 10})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSendCoin_InsufficientFunds_422(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	me := model.User{ID: uuid.New(), Username: "me"}
	to := model.User{ID: uuid.New(), Username: "alice"}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindUserByUsername(gomock.Any(), "alice").
		Return(to, nil)

	repoMock.EXPECT().
		SendCoins(gomock.Any(), me.ID, to.ID, int64(100)).
		Return(errors.New("insufficient funds: balance=50, need=100"))

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.POST("/api/sendCoin", withUser(me), api.ApiSendCoinPost)

	body, _ := json.Marshal(openapi.SendCoinRequest{ToUser: "alice", Amount: 100})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestSendCoin_OtherSendError_500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	me := model.User{ID: uuid.New(), Username: "me"}
	to := model.User{ID: uuid.New(), Username: "alice"}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindUserByUsername(gomock.Any(), "alice").
		Return(to, nil)

	repoMock.EXPECT().
		SendCoins(gomock.Any(), me.ID, to.ID, int64(5)).
		Return(errors.New("deadlock detected"))

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.POST("/api/sendCoin", withUser(me), api.ApiSendCoinPost)

	body, _ := json.Marshal(openapi.SendCoinRequest{ToUser: "alice", Amount: 5})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSendCoin_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	me := model.User{ID: uuid.New(), Username: "me"}
	to := model.User{ID: uuid.New(), Username: "alice"}
	repoMock := mock_repo.NewMockMerchRepo(ctrl)

	repoMock.EXPECT().
		FindUserByUsername(gomock.Any(), "alice").
		Return(to, nil)

	repoMock.EXPECT().
		SendCoins(gomock.Any(), me.ID, to.ID, int64(10)).
		Return(nil)

	api := NewAPI(repoMock, nil)
	r := gin.New()
	r.POST("/api/sendCoin", withUser(me), api.ApiSendCoinPost)

	body, _ := json.Marshal(openapi.SendCoinRequest{ToUser: "alice", Amount: 10})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
