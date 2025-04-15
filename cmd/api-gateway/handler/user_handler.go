package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"FoodStore-AdvProg2/proto/user"
)

// UserHandler обрабатывает HTTP-запросы к User API
type UserHandler struct {
	client user.UserServiceClient
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(client user.UserServiceClient) *UserHandler {
	return &UserHandler{
		client: client,
	}
}

// RegisterUser обрабатывает запрос на регистрацию нового пользователя
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"full_name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправляем запрос к gRPC сервису
	resp, err := h.client.RegisterUser(context.Background(), &user.UserRequest{
		Username: request.Username,
		Email:    request.Email,
		FullName: request.FullName,
		Password: request.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusCreated, gin.H{
		"id":        resp.User.Id,
		"username":  resp.User.Username,
		"email":     resp.User.Email,
		"full_name": resp.User.FullName,
	})
}

// AuthenticateUser обрабатывает запрос на вход пользователя
func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправляем запрос к gRPC сервису
	resp, err := h.client.Login(context.Background(), &user.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"id":        resp.User.Id,
		"username":  resp.User.Username,
		"email":     resp.User.Email,
		"full_name": resp.User.FullName,
	})
}

// GetUserProfile возвращает профиль авторизованного пользователя
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Получаем ID пользователя из контекста (устанавливается в middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Отправляем запрос к gRPC сервису
	resp, err := h.client.GetUserProfile(context.Background(), &user.UserProfileRequest{
		UserId: userID.(string),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем профиль пользователя
	c.JSON(http.StatusOK, gin.H{
		"id":        resp.User.Id,
		"username":  resp.User.Username,
		"email":     resp.User.Email,
		"full_name": resp.User.FullName,
	})
}
