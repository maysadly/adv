package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RunHomepageServer запускает отдельный сервер для корневого маршрута
func RunHomepageServer() {
	// Создаём чистый экземпляр gin без дополнительных middleware
	r := gin.New()
	r.Use(gin.Logger())

	// Загружаем шаблоны
	r.LoadHTMLGlob("./public/*.html")

	// Регистрируем только корневой маршрут
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Food Store - Вход в систему",
		})
	})

	// Запускаем в отдельной горутине, чтобы не блокировать основной сервер
	go func() {
		// Используем порт 8079 - убедитесь, что он свободен
		port := "8079"
		log.Printf("Homepage server starting on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Printf("Homepage server failed: %v", err)
		}
	}()

	// Добавляем перенаправление в основной router
	log.Println("Added redirect from root path to homepage server")
}
