package middleware

import (
	"time"

	"github.com/6ermvH/MerchShop/internal/logx"
	"github.com/gin-gonic/gin"
)

func Log(lg logx.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		rid, _ := c.Get(CtxRequestId)
		l := lg.With(
			"rid", rid,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)

		ctx := logx.IntoContext(c.Request.Context(), l)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		l.Info(ctx, "request completed",
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"size", c.Writer.Size(),
			"errors", c.Errors.ByType(gin.ErrorTypeAny).String(),
		)
	}
}
