package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_repo "github.com/6ermvH/MerchShop/gen/mock/repo"
	"github.com/6ermvH/MerchShop/gen/openapi"
	"github.com/6ermvH/MerchShop/internal/hasher"
	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/model"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

	authReq := openapi.AuthRequest{
		Username: "german",
		Password: "password",
	}

	repoMock.EXPECT().FindUserByUsername(gomock.Any(), authReq.Username).
		Return(model.User{}, repo.ErrNotFound)

	user := model.User{
		ID:           uuid.New(),
		Username:     authReq.Username,
		Balance:      0,
		PasswordHash: authReq.Password,
		CreatedAt:    time.Now(),
	}

	repoMock.EXPECT().
		CreateUser(gomock.Any(), authReq.Username, gomock.AssignableToTypeOf("")).
		DoAndReturn(func(_ context.Context, username, passwordHash string) (model.User, error) {
			if err := hasher.CheckPassword(passwordHash, authReq.Password); err != nil {
				t.Fatalf("password hash doesn't verify against original password")
			}
			user.PasswordHash = passwordHash
			return user, nil
		})

	api := NewAPI(repoMock, j)
	r := gin.New()
	api.RegisterRoutes(r)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp openapi.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v; body: %s", err, w.Body.String())
	}
	if resp.Token == "" {
		t.Fatalf("empty token in response")
	}

}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

	authReq := openapi.AuthRequest{
		Username: "german",
		Password: "password",
	}

	hash, _ := hasher.HashPassword(authReq.Password)
	user := model.User{
		ID:           uuid.New(),
		Username:     authReq.Username,
		Balance:      0,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	repoMock.EXPECT().FindUserByUsername(gomock.Any(), authReq.Username).
		Return(user, nil)

	api := NewAPI(repoMock, j)
	r := gin.New()
	api.RegisterRoutes(r)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp openapi.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v; body: %s", err, w.Body.String())
	}
	if resp.Token == "" {
		t.Fatalf("empty token in response")
	}

}

func TestWrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

	authReq := openapi.AuthRequest{
		Username: "german",
		Password: "password",
	}

	hash, _ := hasher.HashPassword("otherPassword")
	user := model.User{
		ID:           uuid.New(),
		Username:     authReq.Username,
		Balance:      0,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	repoMock.EXPECT().FindUserByUsername(gomock.Any(), authReq.Username).
		Return(user, nil)

	api := NewAPI(repoMock, j)
	r := gin.New()
	api.RegisterRoutes(r)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: got %d, want %d; body: %s", w.Code, http.StatusUnauthorized, w.Body.String())
	}

}

func TestDbError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

	authReq := openapi.AuthRequest{
		Username: "german",
		Password: "password",
	}

	repoMock.EXPECT().FindUserByUsername(gomock.Any(), authReq.Username).
		Return(model.User{}, fmt.Errorf("db error"))

	api := NewAPI(repoMock, j)
	r := gin.New()
	api.RegisterRoutes(r)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: got %d, want %d; body: %s", w.Code, http.StatusInternalServerError, w.Body.String())
	}

}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

	authReq := openapi.AuthRequest{Username: "german", Password: "password"}

	repoMock.EXPECT().
		FindUserByUsername(gomock.Any(), authReq.Username).
		Return(model.User{}, repo.ErrNotFound)

	repoMock.EXPECT().
		CreateUser(gomock.Any(), authReq.Username, gomock.AssignableToTypeOf("")).
		Return(model.User{}, errors.New("unique_violation"))

	api := NewAPI(repoMock, j)
	r := gin.New()
	api.RegisterRoutes(r)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("unexpected status code: got %d, want %d; body: %s", w.Code, http.StatusConflict, w.Body.String())
	}
}

func TestAuth_BadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_repo.NewMockMerchRepo(ctrl)
	j := jwtutil.NewHS256("is-my-private-secret-key-hello-world", "merch", "merch")

	api := NewAPI(repoMock, j)
	r := gin.New()
	api.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader([]byte(`{bad`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d; body: %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}
