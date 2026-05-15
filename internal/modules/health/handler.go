package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-monolithic-boilerplate/internal/common/response"
)

var startedAt = time.Now().UTC()

type statusResponse struct {
	Status    string    `json:"status"`
	StartedAt time.Time `json:"startedAt"`
	Uptime    string    `json:"uptime"`
}

func RegisterRoutes(router gin.IRoutes) {
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, statusResponse{
			Status:    "healthy",
			StartedAt: startedAt,
			Uptime:    time.Since(startedAt).Round(time.Second).String(),
		})
	})
}
