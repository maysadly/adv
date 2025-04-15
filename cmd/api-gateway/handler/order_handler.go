package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"FoodStore-AdvProg2/proto/order"
)

type OrderHandler struct {
	client order.OrderServiceClient
}

func NewOrderHandler(client order.OrderServiceClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var reqBody struct {
		UserID string                   `json:"user_id" binding:"required"`
		Items  []*order.CreateOrderItem `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &order.CreateOrderRequest{
		UserId: reqBody.UserID,
		Items:  reqBody.Items,
	}

	resp, err := h.client.CreateOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": resp.OrderId})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	req := &order.GetOrderRequest{
		Id: id,
	}

	resp, err := h.client.GetOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, resp.Order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var reqBody struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &order.UpdateOrderStatusRequest{
		Id:     id,
		Status: reqBody.Status,
	}

	resp, err := h.client.UpdateOrderStatus(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.Query("user_id")

	var orders []*order.Order

	if userID != "" {

		userOrdersReq := &order.GetUserOrdersRequest{
			UserId: userID,
		}
		resp, err := h.client.GetUserOrders(context.Background(), userOrdersReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = resp.Orders
	} else {

		allOrdersReq := &order.GetAllOrdersRequest{}
		resp, err := h.client.GetAllOrders(context.Background(), allOrdersReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = resp.Orders
	}

	c.JSON(http.StatusOK, orders)
}
