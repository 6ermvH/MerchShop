package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_repo "github.com/6ermvH/MerchShop/gen/mock/repo"

	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseBearer(t *testing.T) {
	require.Equal(t, "", parseBearer(""))
	require.Equal(t, "X", parseBearer("Bearer X"))
	require.Equal(t, "", parseBearer("bearer X"))
}

func TestAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("hello-world-my-name-is-german", "merch", "merch")

	r := gin.New()
	r.GET("/x", Auth(j, repoMock), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("has code: %d, want code: %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_GoodToken_UserFound_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	user := &model.User{
		ID: id, Username: "alice", Balance: 0, CreatedAt: time.Now(),
	}

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	repoMock.EXPECT().FindUserByID(gomock.Any(), id).Return(*user, nil)

	j := jwtutil.NewHS256("hello-world-my-name-is-german", "merch", "merch")
	token, _ := j.Sign(user.ID, user.Username)

	r := gin.New()
	r.GET("/x", Auth(j, repoMock), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("has code: %d, want code: %d", w.Code, http.StatusOK)
	}
}

func TestAuth_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	user := &model.User{
		ID:        id,
		Username:  "german",
		Balance:   0,
		CreatedAt: time.Now(),
	}

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	repoMock.EXPECT().FindUserByID(gomock.Any(), id).Return(*user, repo.ErrNotFound)

	j := jwtutil.NewHS256("hello-world-my-name-is-german", "merch", "merch")
	token, _ := j.Sign(user.ID, user.Username)

	r := gin.New()
	r.GET("/x", Auth(j, repoMock), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("has code: %d, want code: %d", w.Code, http.StatusUnauthorized)
	}
}
