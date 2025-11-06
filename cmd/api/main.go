package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/6ermvH/MerchShop/internal/db"
	"github.com/6ermvH/MerchShop/internal/http/handlers"
	"github.com/6ermvH/MerchShop/internal/http/middleware"
	"github.com/6ermvH/MerchShop/internal/jwtutil"
	"github.com/6ermvH/MerchShop/internal/logx"
	"github.com/6ermvH/MerchShop/internal/repo"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx := logx.IntoContext(context.Background(), logx.NewSlog(logger))

	port := getenv("PORT", "8080")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("DATABASE_URL is empty; set DATABASE_URL or use docker-compose")
		os.Exit(1)
	}
	jwtSecret := getenv("JWT_SECRET", "dev-secret")
	jwtIss := getenv("JWT_ISS", "merch-shop")
	jwtAud := getenv("JWT_AUD", "merch-shop-client")

	connCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pool, err := db.NewPool(connCtx, dsn)
	if err != nil {
		logger.Error("failed to connect to Postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	repositories := repo.NewRepo(pool)
	hs := jwtutil.NewHS256(jwtSecret, jwtIss, jwtAud)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.RequestId())

	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	api := handlers.NewAPI(repositories, hs)
	api.RegisterRoutes(r)

	logger.Info("http server starting", slog.String("addr", ":"+port))
	if err := r.Run(":" + port); err != nil {
		logger.Error("http server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
