package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"FoodStore-AdvProg2/proto/inventory"
)

type ProductHandler struct {
	client inventory.InventoryServiceClient
}

func NewProductHandler(client inventory.InventoryServiceClient) *ProductHandler {
	return &ProductHandler{client: client}
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	req := &inventory.GetProductRequest{
		Id: id,
	}

	resp, err := h.client.GetProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, resp.Product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "6"))

	name := c.Query("name")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	var minPrice float64
	var maxPrice float64
	if minPriceStr != "" {
		minPrice, _ = strconv.ParseFloat(minPriceStr, 64)
	}
	if maxPriceStr != "" {
		maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
	}

	req := &inventory.ListProductsRequest{
		Filter: &inventory.FilterParams{
			Name:     name,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
		},
		Pagination: &inventory.PaginationParams{
			Page:    int32(page),
			PerPage: int32(perPage),
		},
	}

	resp, err := h.client.ListProducts(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": resp.Products,
		"total":    resp.Total,
		"page":     resp.Page,
		"per_page": resp.PerPage,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var reqBody struct {
		Name  string  `json:"Name" binding:"required"`
		Price float64 `json:"Price" binding:"required,min=0"`
		Stock int     `json:"Stock" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &inventory.CreateProductRequest{
		Name:  reqBody.Name,
		Price: reqBody.Price,
		Stock: int32(reqBody.Stock),
	}

	resp, err := h.client.CreateProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp.Product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var reqBody struct {
		Name  string  `json:"Name" binding:"required"`
		Price float64 `json:"Price" binding:"required,min=0"`
		Stock int     `json:"Stock" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &inventory.UpdateProductRequest{
		Id:    id,
		Name:  reqBody.Name,
		Price: reqBody.Price,
		Stock: int32(reqBody.Stock),
	}

	resp, err := h.client.UpdateProduct(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	req := &inventory.DeleteProductRequest{
		Id: id,
	}

	resp, err := h.client.DeleteProduct(context.Background(), req)
	if err != nil {

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			}
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp})
}
