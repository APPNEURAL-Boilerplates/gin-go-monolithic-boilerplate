package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-monolithic-boilerplate/internal/common/apperror"
	"github.com/example/gin-go-monolithic-boilerplate/internal/common/response"
	"github.com/example/gin-go-monolithic-boilerplate/internal/config"
	"github.com/example/gin-go-monolithic-boilerplate/internal/http/middleware"
	"github.com/example/gin-go-monolithic-boilerplate/internal/modules/health"
	"github.com/example/gin-go-monolithic-boilerplate/internal/modules/root"
	"github.com/example/gin-go-monolithic-boilerplate/internal/modules/users"
)

func NewRouter(cfg config.Config, logger *slog.Logger) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(
		middleware.RequestID(),
		middleware.SecurityHeaders(),
		middleware.CORS(cfg.AllowedOrigins),
		middleware.Recovery(logger),
		middleware.Logger(logger),
	)

	root.RegisterRoutes(router, cfg)
	health.RegisterRoutes(router)

	api := router.Group("/api/v1")
	userRepository := users.NewInMemoryRepository()
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)
	users.RegisterRoutes(api, userHandler)

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, apperror.New(
			http.StatusNotFound,
			"NOT_FOUND",
			"The requested endpoint was not found.",
			nil,
		))
	})

	router.NoMethod(func(c *gin.Context) {
		response.Error(c, apperror.New(
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"The requested HTTP method is not allowed for this endpoint.",
			nil,
		))
	})

	return router
}
