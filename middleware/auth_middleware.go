package middleware

import (
	"net/http"
	"strings"

	"FoodStore-AdvProg2/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "не предоставлен токен авторизации"})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена авторизации"})
			c.Abort()
			return
		}

		// Извлекаем токен
		tokenString := parts[1]

		// Проверяем и валидируем токен
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен: " + err.Error()})
			c.Abort()
			return
		}

		// Сохраняем ID пользователя в контексте для дальнейшего использования
		c.Set("userID", userID)
		c.Next()
	}
}
