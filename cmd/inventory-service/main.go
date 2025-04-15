package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/proto/inventory"
	"FoodStore-AdvProg2/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	// Инициализация базы данных
	dbHost := os.Getenv("DB")
	if dbHost == "" {
		log.Fatal("DB environment variable not set")
	}
	postgres.InitDB(dbHost)
	log.Println("Connected to PostgreSQL")

	if err := postgres.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	// Создание репозитория и use case
	productRepo := postgres.NewProductPostgresRepo()
	productUC := usecase.NewProductUseCase(productRepo)

	// Настройка gRPC сервера
	port := os.Getenv("INVENTORY_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	inventoryServer := NewInventoryServiceServer(productUC)
	inventory.RegisterInventoryServiceServer(server, inventoryServer)

	// Включаем reflection для отладки
	reflection.Register(server)

	log.Printf("Inventory Service is starting on port %s...", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// InventoryServiceServer реализует gRPC сервер для Inventory Service
type InventoryServiceServer struct {
	inventory.UnimplementedInventoryServiceServer
	productUC *usecase.ProductUseCase
}

func NewInventoryServiceServer(productUC *usecase.ProductUseCase) *InventoryServiceServer {
	return &InventoryServiceServer{
		productUC: productUC,
	}
}

// GetProduct возвращает информацию о продукте по ID
func (s *InventoryServiceServer) GetProduct(ctx context.Context, req *inventory.GetProductRequest) (*inventory.GetProductResponse, error) {
	product, err := s.productUC.GetByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &inventory.GetProductResponse{
		Product: &inventory.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: int32(product.Stock),
		},
	}, nil
}

// CreateProduct создает новый продукт
func (s *InventoryServiceServer) CreateProduct(ctx context.Context, req *inventory.CreateProductRequest) (*inventory.CreateProductResponse, error) {
	product := domain.Product{
		ID:    uuid.New().String(),
		Name:  req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	}

	if err := s.productUC.Create(product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &inventory.CreateProductResponse{
		Product: &inventory.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: int32(product.Stock),
		},
	}, nil
}

// UpdateProduct обновляет существующий продукт
func (s *InventoryServiceServer) UpdateProduct(ctx context.Context, req *inventory.UpdateProductRequest) (*inventory.Product, error) {
	product := domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	}

	if err := s.productUC.Update(req.Id, product); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	updatedProduct, err := s.productUC.GetByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found after update")
	}

	return &inventory.Product{
		Id:    updatedProduct.ID,
		Name:  updatedProduct.Name,
		Price: updatedProduct.Price,
		Stock: int32(updatedProduct.Stock),
	}, nil
}

// DeleteProduct удаляет продукт по ID
func (s *InventoryServiceServer) DeleteProduct(ctx context.Context, req *inventory.DeleteProductRequest) (*inventory.DeleteProductResponse, error) {
	// Проверяем существование продукта перед удалением
	_, err := s.productUC.GetByID(req.Id)
	if err != nil {
		return &inventory.DeleteProductResponse{Success: false}, status.Error(codes.NotFound, "product not found")
	}

	if err := s.productUC.Delete(req.Id); err != nil {
		return &inventory.DeleteProductResponse{Success: false}, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &inventory.DeleteProductResponse{
		Success: true,
	}, nil
}

// ListProducts возвращает список продуктов с фильтрацией и пагинацией
func (s *InventoryServiceServer) ListProducts(ctx context.Context, req *inventory.ListProductsRequest) (*inventory.ListProductsResponse, error) {
	filter := domain.FilterParams{
		Name:     req.Filter.Name,
		MinPrice: req.Filter.MinPrice,
		MaxPrice: req.Filter.MaxPrice,
	}

	pagination := domain.PaginationParams{
		Page:    int(req.Pagination.Page),
		PerPage: int(req.Pagination.PerPage),
	}

	products, total, err := s.productUC.List(filter, pagination)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	var protoProducts []*inventory.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &inventory.Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		})
	}

	return &inventory.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
		Page:     int32(pagination.Page),
		PerPage:  int32(pagination.PerPage),
	}, nil
}

// CheckStock проверяет наличие товаров на складе
func (s *InventoryServiceServer) CheckStock(ctx context.Context, req *inventory.CheckStockRequest) (*inventory.CheckStockResponse, error) {
	for _, item := range req.Items {
		product, err := s.productUC.GetByID(item.ProductId)
		if err != nil {
			return &inventory.CheckStockResponse{
				Available:            false,
				UnavailableProductId: item.ProductId,
			}, nil
		}

		if product.Stock < int(item.Quantity) {
			return &inventory.CheckStockResponse{
				Available:            false,
				UnavailableProductId: item.ProductId,
			}, nil
		}
	}

	return &inventory.CheckStockResponse{
		Available: true,
	}, nil
}

// UpdateStock обновляет количество товаров на складе
func (s *InventoryServiceServer) UpdateStock(ctx context.Context, req *inventory.UpdateStockRequest) (*inventory.UpdateStockResponse, error) {
	for _, item := range req.Items {
		product, err := s.productUC.GetByID(item.ProductId)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "product not found: %s", item.ProductId)
		}

		product.Stock -= int(item.Quantity)
		if product.Stock < 0 {
			product.Stock = 0
		}

		if err := s.productUC.Update(product.ID, product); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
		}
	}

	return &inventory.UpdateStockResponse{
		Success: true,
	}, nil
}
