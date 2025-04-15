package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Глобальное хранилище заказов и продуктов для демо
var (
	orders    = make(map[string]gin.H)
	ordersMux sync.RWMutex
	products  = []gin.H{
		{
			"id":    "p1",
			"name":  "Продукт 1",
			"price": 100.0,
			"stock": 10,
		},
		{
			"id":    "p2",
			"name":  "Продукт 2",
			"price": 200.0,
			"stock": 5,
		},
		{
			"id":    "p3",
			"name":  "Продукт 3",
			"price": 300.0,
			"stock": 15,
		},
	}
	productsMux sync.RWMutex
)

func main() {
	// Создаём чистый экземпляр gin без дополнительных middleware
	r := gin.New()
	r.Use(gin.Logger())

	// Определение абсолютного пути к проекту
	projectRoot := "/Users/tleukhanmakhmutov/Desktop/Study/FoodStore-AdvProg2-1"

	// Загружаем шаблоны и статические файлы из абсолютных путей
	r.LoadHTMLGlob(projectRoot + "/public/*.html")

	// Регистрируем пути для статических файлов
	r.Static("/static", projectRoot+"/public")

	// Маршрут для перенаправления с /login на /
	r.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	// Добавляем обработчики для /order и /admin для тестирования
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", gin.H{
			"title": "Food Store - Заказы",
		})
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", gin.H{
			"title": "Food Store - Админ панель",
		})
	})

	// Регистрируем корневой маршрут
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":      "Food Store - Вход в систему",
			"apiBaseUrl": "http://localhost:8079", // Используем текущий сервер для API
		})
	})

	// Добавляем маршрут для страницы регистрации
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":      "Food Store - Регистрация",
			"apiBaseUrl": "http://localhost:8079",
		})
	})

	// Обрабатываем запрос логина локально
	r.POST("/api/users/login", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Login attempt: username=%s, password=%s", request.Username, request.Password)

		// Упрощенная проверка для демо - принимаем любые учетные данные
		// В реальном приложении здесь был бы запрос к user-service
		if request.Username == "admin" {
			c.JSON(http.StatusOK, gin.H{
				"id":        "admin-id",
				"username":  request.Username,
				"email":     "admin@example.com",
				"full_name": "Admin User",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"id":        "user-" + request.Username,
				"username":  request.Username,
				"email":     request.Username + "@example.com",
				"full_name": "Regular User",
			})
		}
	})

	// Добавляем обработчик запроса регистрации
	r.POST("/api/users/register", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			FullName string `json:"full_name"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Простая эмуляция регистрации
		c.JSON(http.StatusCreated, gin.H{
			"id":        "user-" + request.Username,
			"username":  request.Username,
			"email":     request.Email,
			"full_name": request.FullName,
		})
	})

	// Заглушки для API продуктов
	r.GET("/api/products", func(c *gin.Context) {
		// Эмуляция списка продуктов
		productsMux.RLock()
		defer productsMux.RUnlock()

		c.JSON(http.StatusOK, gin.H{
			"products": products,
			"total":    len(products),
			"page":     1,
			"per_page": 10,
		})
	})

	// Заглушка для создания продукта
	r.POST("/api/products", func(c *gin.Context) {
		var request struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
			Stock int     `json:"stock"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newProduct := gin.H{
			"id":    fmt.Sprintf("new-product-%d", time.Now().UnixNano()),
			"name":  request.Name,
			"price": request.Price,
			"stock": request.Stock,
		}

		productsMux.Lock()
		products = append(products, newProduct)
		productsMux.Unlock()

		c.JSON(http.StatusCreated, newProduct)
	})

	// Заглушки для API заказов
	r.GET("/api/orders", func(c *gin.Context) {
		userID := c.Query("user_id")

		ordersMux.RLock()
		defer ordersMux.RUnlock()

		result := []gin.H{}
		for _, order := range orders {
			// Если указан user_id и он не совпадает, пропускаем
			if userID != "" && order["user_id"] != userID {
				continue
			}
			result = append(result, order)
		}

		c.JSON(http.StatusOK, result)
	})

	// Получение конкретного заказа
	r.GET("/api/orders/:id", func(c *gin.Context) {
		orderID := c.Param("id")

		ordersMux.RLock()
		order, exists := orders[orderID]
		ordersMux.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	// Создание заказа
	r.POST("/api/orders", func(c *gin.Context) {
		var request struct {
			UserID string `json:"user_id" binding:"required"`
			Items  []struct {
				ProductID string `json:"product_id" binding:"required"`
				Quantity  int    `json:"quantity" binding:"required,min=1"`
			} `json:"items" binding:"required,min=1"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создаём новый заказ
		orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())
		orderItems := []gin.H{}
		totalAmount := 0.0

		// Получаем информацию о продуктах и создаём элементы заказа
		for _, item := range request.Items {
			// В реальном приложении здесь был бы запрос к сервису продуктов
			// Сейчас используем фиксированные цены
			var price float64
			var productName string

			switch item.ProductID {
			case "p1":
				price = 100.0
				productName = "Продукт 1"
			case "p2":
				price = 200.0
				productName = "Продукт 2"
			case "p3":
				price = 300.0
				productName = "Продукт 3"
			default:
				price = 150.0
				productName = "Неизвестный продукт"
			}

			itemTotal := price * float64(item.Quantity)
			totalAmount += itemTotal

			orderItems = append(orderItems, gin.H{
				"id":         fmt.Sprintf("item-%d", time.Now().UnixNano()),
				"product_id": item.ProductID,
				"quantity":   item.Quantity,
				"price":      price,
				"product": gin.H{
					"id":    item.ProductID,
					"name":  productName,
					"price": price,
				},
			})
		}

		// Формируем объект заказа
		order := gin.H{
			"id":           orderID,
			"user_id":      request.UserID,
			"total_amount": totalAmount,
			"status":       "pending",
			"created_at":   time.Now().Format(time.RFC3339),
			"items":        orderItems,
		}

		// Сохраняем заказ
		ordersMux.Lock()
		orders[orderID] = order
		ordersMux.Unlock()

		c.JSON(http.StatusCreated, gin.H{
			"id":      orderID,
			"message": "Order created successfully",
		})
	})

	// Обновление статуса заказа
	r.PATCH("/api/orders/:id", func(c *gin.Context) {
		orderID := c.Param("id")

		var request struct {
			Status string `json:"status" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем валидность статуса
		if request.Status != "pending" && request.Status != "completed" && request.Status != "cancelled" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		ordersMux.Lock()
		order, exists := orders[orderID]
		if !exists {
			ordersMux.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		// Обновляем статус заказа
		order["status"] = request.Status
		ordersMux.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"id":      orderID,
			"status":  request.Status,
			"message": "Order status updated successfully",
		})
	})

	// Настраиваем прокси для остальных API-запросов на основной сервер
	mainServer, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(mainServer)

	// Проксируем API-запросы, за исключением login
	r.NoRoute(func(c *gin.Context) {
		// Если путь начинается с /api/ и это не /api/users/login - проксируем
		if strings.HasPrefix(c.Request.URL.Path, "/api/") &&
			c.Request.URL.Path != "/api/users/login" {
			proxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		// В остальных случаях возвращаем 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
	})

	// Обработка перенаправлений после входа
	r.GET("/redirect", func(c *gin.Context) {
		role := c.Query("role")
		if role == "admin" {
			c.Redirect(http.StatusFound, "http://localhost:8080/admin")
		} else {
			c.Redirect(http.StatusFound, "http://localhost:8080/order")
		}
	})

	// Запускаем на порту 8079
	port := "8079"
	log.Printf("Homepage server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Homepage server failed: %v", err)
	}
}
