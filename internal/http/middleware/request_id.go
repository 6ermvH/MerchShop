package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CtxRequestId = "request_id"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Set(CtxRequestId, id)
		c.Next()
	}
}
