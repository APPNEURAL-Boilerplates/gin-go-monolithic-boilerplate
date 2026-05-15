package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-monolithic-boilerplate/internal/common/id"
	"github.com/example/gin-go-monolithic-boilerplate/internal/common/requestid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestid.Header)
		if requestID == "" {
			requestID = id.New()
		}

		c.Set(requestid.ContextKey, requestID)
		c.Header(requestid.Header, requestID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	return requestid.Get(c)
}
