package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"FoodStore-AdvProg2/cmd/api-gateway/handler"
	"FoodStore-AdvProg2/cmd/api-gateway/middleware"
	"FoodStore-AdvProg2/proto/inventory"
	"FoodStore-AdvProg2/proto/order"
	"FoodStore-AdvProg2/proto/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	// Запускаем отдельный сервер для главной страницы
	RunHomepageServer()

	// Подключение к сервисам через gRPC
	inventoryConn, err := connectToService("INVENTORY_SERVICE_URL", "localhost:8081")
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()

	orderConn, err := connectToService("ORDER_SERVICE_URL", "localhost:8082")
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	userConn, err := connectToService("USER_SERVICE_URL", "localhost:8083")
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	// Создание gRPC клиентов
	inventoryClient := inventory.NewInventoryServiceClient(inventoryConn)
	orderClient := order.NewOrderServiceClient(orderConn)
	userClient := user.NewUserServiceClient(userConn)

	// Создание обработчиков
	productHandler := handler.NewProductHandler(inventoryClient)
	orderHandler := handler.NewOrderHandler(orderClient)
	userHandler := handler.NewUserHandler(userClient)

	// Инициализация Gin с нуля
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Telemetry())

	// Статические файлы и шаблоны
	r.Static("/static", "./public")
	r.LoadHTMLGlob("./public/*.html")

	// Регистрируем тестовые маршруты
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test page works!")
	})

	// ВАЖНО: явная регистрация корневого маршрута
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Food Store - Вход в систему",
		})
	})

	// HTML страницы
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", nil)
	})

	// API маршруты
	api := r.Group("/api")
	{
		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PATCH("/:id", orderHandler.UpdateOrderStatus)
		}

		// User routes
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.RegisterUser)
			users.POST("/login", userHandler.AuthenticateUser)
			users.GET("/profile", middleware.AuthMiddleware(), userHandler.GetUserProfile)
		}
	}

	// Запуск сервера
	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)

	// Регистрируем корневой маршрут перед запуском
	RegisterRootHandler(r)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Функция для подключения к gRPC сервису
func connectToService(envVarName, defaultURL string) (*grpc.ClientConn, error) {
	serviceURL := os.Getenv(envVarName)
	if serviceURL == "" {
		serviceURL = defaultURL
	}

	return grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
