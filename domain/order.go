package domain

import (
	"time"
)

// Статусы заказа
const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

// Order представляет заказ в системе
type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	Items      []OrderItem `json:"items,omitempty"`
}

// OrderItem представляет элемент заказа
type OrderItem struct {
	ID        string   `json:"id"`
	OrderID   string   `json:"order_id"`
	ProductID string   `json:"product_id"`
	Product   *Product `json:"product,omitempty"`
	Quantity  int      `json:"quantity"`
	Price     float64  `json:"price"`
}

// OrderRequest представляет запрос на создание заказа
type OrderRequest struct {
	UserID string             `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

// OrderItemRequest представляет запрос на добавление элемента в заказ
type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// OrderStatusUpdateRequest представляет запрос на обновление статуса заказа
type OrderStatusUpdateRequest struct {
	Status string `json:"status"`
}
