package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-monolithic-boilerplate/internal/common/apperror"
	"github.com/example/gin-go-monolithic-boilerplate/internal/common/response"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered",
			"panic", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"request_id", GetRequestID(c),
		)
		response.Error(c, apperror.New(http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred.", nil))
	})
}
