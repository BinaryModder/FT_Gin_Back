package handler

import (
	"net/http"
	"strconv"

	"github.com/BinaryModder/FT_Gin_Back.git/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// CreateUser godoc
// @Summary Создать пользователя
// @Description Регистрирует нового пользователя
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	// Парсим JSON из запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	// Вызываем Service для обработки
	user, err := h.service.RegisterUser(req.Name, req.Email, req.Age)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Отправляем ответ клиенту
	c.JSON(http.StatusCreated, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

// GetUser godoc
// @Summary Получить пользователя по ID
// @Description Возвращает одного пользователя
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Description Возвращает список всех пользователей
// @Produce json
// @Success 200 {array} UserResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		}
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateUser godoc
// @Summary Обновить пользователя
// @Description Обновляет данные пользователя по ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body UpdateUserRequest true "Updated user data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	user, err := h.service.UpdateUserInfo(uint(id), req.Name, req.Email, req.Age)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Request/Response структуры

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,gt=0,lt=150"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
