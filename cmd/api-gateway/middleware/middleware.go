package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger - middleware для логирования запросов
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Выполнение следующих обработчиков
		c.Next()

		// После завершения запроса
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[API-GATEWAY] %s | %d | %s | %s | %s",
			method,
			statusCode,
			latency,
			clientIP,
			path,
		)
	}
}

// Telemetry - middleware для телеметрии
func Telemetry() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Выполнение следующих обработчиков
		c.Next()

		// После завершения запроса
		duration := time.Since(startTime)

		// Добавление заголовков телеметрии
		c.Header("X-Response-Time", duration.String())
	}
}

// AuthMiddleware - middleware для проверки аутентификации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")

		// Простая проверка формата (должен быть "Bearer token")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// В реальном приложении здесь должна быть проверка JWT токена
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Извлекаем ID пользователя из токена
		// В простейшем случае считаем, что ID - это первая часть токена до дефиса
		userID := strings.Split(token, "-")[0]

		// Сохраняем ID пользователя в контексте
		c.Set("user_id", userID)

		// Переходим к следующим обработчикам
		c.Next()
	}
}
