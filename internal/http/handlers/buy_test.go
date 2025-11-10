package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_repo "github.com/6ermvH/MerchShop/gen/mock/repo"
	"github.com/6ermvH/MerchShop/internal/http/middleware"
	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestBuyItem_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := model.User{ID: uuid.New(), Username: "german", Balance: 100}

	products := []string{"t-shirt", "cup", "book",
		"pen", "powerbank", "hoody", "umbrella",
		"socks", "wallet", "pink-hoody"}

	for _, product := range products {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Run("Buy "+product, func(t *testing.T) {
			repoMock := mock_repo.NewMockMerchRepo(ctrl)

			repoMock.EXPECT().
				FindUserByID(gomock.Any(), user.ID).
				Return(user, nil)
			repoMock.EXPECT().
				BuyProduct(gomock.Any(), user.ID, product).
				Return(nil)

			j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

			api := NewAPI(repoMock, j)
			r := gin.New()
			r.GET("/api/buy/:item", middleware.Auth(j, repoMock), api.ApiBuyItemGet)

			token, _ := j.Sign(user.ID, user.Username)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/buy/"+product, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status got %d want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
			}
		})
	}
}

func TestBuyItem_UnknownItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := model.User{ID: uuid.New(), Username: "german", Balance: 100}

	products := []string{"mackbook", "iphone", "google-pixel", "man"}

	for _, product := range products {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Run("Buy "+product, func(t *testing.T) {
			repoMock := mock_repo.NewMockMerchRepo(ctrl)

			repoMock.EXPECT().
				FindUserByID(gomock.Any(), user.ID).
				Return(user, nil)
			repoMock.EXPECT().
				BuyProduct(gomock.Any(), user.ID, product).
				Return(fmt.Errorf("Bad product"))

			j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

			api := NewAPI(repoMock, j)
			r := gin.New()
			r.GET("/api/buy/:item", middleware.Auth(j, repoMock), api.ApiBuyItemGet)

			token, _ := j.Sign(user.ID, user.Username)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/buy/"+product, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusInternalServerError {
				t.Fatalf("status got %d want %d; body=%s", w.Code, http.StatusInternalServerError, w.Body.String())
			}
		})
	}
}
