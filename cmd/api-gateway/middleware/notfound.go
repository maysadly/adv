package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler перехватывает все запросы к несуществующим маршрутам
// и специально обрабатывает корневой путь "/"
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Если это запрос к корневому пути
		if c.Request.URL.Path == "/" {
			// Проверяем, авторизован ли пользователь
			userID, userIDErr := c.Cookie("user_id")
			userRole, userRoleErr := c.Cookie("user_role")

			// Если пользователь авторизован, перенаправляем его
			if userIDErr == nil && userRoleErr == nil && userID != "" {
				if userRole == "admin" {
					c.Redirect(http.StatusFound, "/admin")
					return
				} else {
					c.Redirect(http.StatusFound, "/order")
					return
				}
			}

			// Если не авторизован, показываем страницу входа
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Food Store - Вход в систему",
			})
			return
		}

		// Для остальных несуществующих маршрутов возвращаем 404
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Page not found",
		})
	}
}
