package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-monolithic-boilerplate/internal/common/apperror"
	"github.com/example/gin-monolithic-boilerplate/internal/common/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	users := router.Group("/users")
	users.GET("", handler.List)
	users.POST("", handler.Create)
	users.GET("/:id", handler.GetByID)
	users.PATCH("/:id", handler.Update)
	users.DELETE("/:id", handler.Delete)
}

func (h *Handler) List(c *gin.Context) {
	users, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, users)
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.Validation("Invalid request body.", map[string]any{"reason": err.Error()}))
		return
	}

	user, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, user)
}

func (h *Handler) GetByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, user)
}

func (h *Handler) Update(c *gin.Context) {
	var request UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, apperror.Validation("Invalid request body.", map[string]any{"reason": err.Error()}))
		return
	}

	user, err := h.service.Update(c.Request.Context(), c.Param("id"), request)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}
