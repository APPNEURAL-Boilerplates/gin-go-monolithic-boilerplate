package root

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-monolithic-boilerplate/internal/common/response"
	"github.com/example/gin-go-monolithic-boilerplate/internal/config"
)

type metadata struct {
	Name        string   `json:"name"`
	Environment string   `json:"environment"`
	Version     string   `json:"version"`
	Endpoints   []string `json:"endpoints"`
}

func RegisterRoutes(router gin.IRoutes, cfg config.Config) {
	router.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, metadata{
			Name:        cfg.AppName,
			Environment: cfg.AppEnv,
			Version:     "0.1.0",
			Endpoints: []string{
				"GET /health",
				"GET /api/v1/users",
				"POST /api/v1/users",
				"GET /api/v1/users/:id",
				"PATCH /api/v1/users/:id",
				"DELETE /api/v1/users/:id",
			},
		})
	})
}
