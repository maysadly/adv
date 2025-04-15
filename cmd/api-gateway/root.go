package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Регистрирует корневой маршрут напрямую для решения проблемы с маршрутизацией
func RegisterRootHandler(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Food Store - Вход в систему",
		})
	})

	// Добавим также простой тестовый маршрут
	r.GET("/welcome", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Food Store API!")
	})
}
