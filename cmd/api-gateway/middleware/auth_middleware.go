package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authenticate проверяет JWT токен и устанавливает ID пользователя в контекст
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Проверяем формат токена
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// В реальном приложении здесь должна быть проверка подлинности JWT
		// и извлечение ID пользователя из токена
		// Пример:
		// userID, err := validateToken(tokenString)
		// if err != nil {
		//    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		//    c.Abort()
		//    return
		// }

		// Для демонстрации просто используем сам токен как ID пользователя
		// В реальном приложении этот код должен быть заменен на проверку JWT
		userID := extractUserIDFromToken(tokenString)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Устанавливаем ID пользователя в контекст
		c.Set("user_id", userID)
		c.Next()
	}
}

// extractUserIDFromToken извлекает ID пользователя из токена (пример)
func extractUserIDFromToken(token string) string {
	// В реальном приложении здесь должна быть валидация JWT токена
	// и извлечение идентификатора пользователя из полезной нагрузки токена

	// Для демонстрации просто извлекаем часть токена до дефиса
	// предполагая, что токен имеет формат "user_id-timestamp"
	parts := strings.Split(token, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
