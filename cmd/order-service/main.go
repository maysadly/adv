package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/proto/inventory"
	"FoodStore-AdvProg2/proto/order"
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

	// Подключение к Inventory Service
	inventoryServiceURL := os.Getenv("INVENTORY_SERVICE_URL")
	if inventoryServiceURL == "" {
		inventoryServiceURL = "localhost:8081"
	}

	inventoryConn, err := grpc.Dial(inventoryServiceURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := inventory.NewInventoryServiceClient(inventoryConn)

	// Создание репозитория и use case
	orderRepo := postgres.NewOrderPostgresRepo()
	productRepo := postgres.NewProductPostgresRepo()
	orderUC := usecase.NewOrderUseCase(orderRepo, productRepo)

	// Настройка gRPC сервера
	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	orderServer := NewOrderServiceServer(orderUC, inventoryClient)
	order.RegisterOrderServiceServer(server, orderServer)

	// Включаем reflection для отладки
	reflection.Register(server)

	log.Printf("Order Service is starting on port %s...", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// OrderServiceServer реализует gRPC сервер для Order Service
type OrderServiceServer struct {
	order.UnimplementedOrderServiceServer
	orderUC         *usecase.OrderUseCase
	inventoryClient inventory.InventoryServiceClient
}

func NewOrderServiceServer(orderUC *usecase.OrderUseCase, inventoryClient inventory.InventoryServiceClient) *OrderServiceServer {
	return &OrderServiceServer{
		orderUC:         orderUC,
		inventoryClient: inventoryClient,
	}
}

// CreateOrder создает новый заказ
func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// Проверяем наличие товаров через Inventory Service
	var inventoryItems []*inventory.OrderItem
	for _, item := range req.Items {
		inventoryItems = append(inventoryItems, &inventory.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	checkResp, err := s.inventoryClient.CheckStock(ctx, &inventory.CheckStockRequest{
		Items: inventoryItems,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check stock: %v", err)
	}

	if !checkResp.Available {
		return nil, status.Errorf(codes.FailedPrecondition, "product %s is out of stock", checkResp.UnavailableProductId)
	}

	// Преобразуем запрос в домен
	var orderItems []domain.OrderItemRequest
	for _, item := range req.Items {
		orderItems = append(orderItems, domain.OrderItemRequest{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
		})
	}

	orderRequest := domain.OrderRequest{
		UserID: req.UserId,
		Items:  orderItems,
	}

	// Создаем заказ
	orderID, err := s.orderUC.CreateOrder(orderRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Обновляем сток через Inventory Service
	_, err = s.inventoryClient.UpdateStock(ctx, &inventory.UpdateStockRequest{
		Items: inventoryItems,
	})

	if err != nil {
		// В реальном приложении нужна компенсационная транзакция
		log.Printf("Failed to update stock after order creation: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update stock")
	}

	return &order.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

// GetOrder возвращает информацию о заказе по ID
func (s *OrderServiceServer) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	domainOrder, err := s.orderUC.GetOrderByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	var orderItems []*order.OrderItem
	for _, item := range domainOrder.Items {
		var productInfo *order.ProductInfo
		if item.Product != nil {
			productInfo = &order.ProductInfo{
				Id:    item.Product.ID,
				Name:  item.Product.Name,
				Price: item.Product.Price,
				Stock: int32(item.Product.Stock),
			}
		}

		orderItems = append(orderItems, &order.OrderItem{
			Id:        item.ID,
			OrderId:   item.OrderID,
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
			Product:   productInfo,
		})
	}

	return &order.GetOrderResponse{
		Order: &order.Order{
			Id:         domainOrder.ID,
			UserId:     domainOrder.UserID,
			TotalPrice: domainOrder.TotalPrice,
			Status:     domainOrder.Status,
			CreatedAt:  timestamppb.New(domainOrder.CreatedAt),
			Items:      orderItems,
		},
	}, nil
}

// UpdateOrderStatus обновляет статус заказа
func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	err := s.orderUC.UpdateOrderStatus(req.Id, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return &order.UpdateOrderStatusResponse{
		Success: true,
	}, nil
}

// GetUserOrders возвращает заказы пользователя
func (s *OrderServiceServer) GetUserOrders(ctx context.Context, req *order.GetUserOrdersRequest) (*order.GetUserOrdersResponse, error) {
	domainOrders, err := s.orderUC.GetOrdersByUserID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user orders: %v", err)
	}

	var orders []*order.Order
	for _, o := range domainOrders {
		orders = append(orders, &order.Order{
			Id:         o.ID,
			UserId:     o.UserID,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedAt:  timestamppb.New(o.CreatedAt),
		})
	}

	return &order.GetUserOrdersResponse{
		Orders: orders,
	}, nil
}

// GetAllOrders возвращает все заказы
func (s *OrderServiceServer) GetAllOrders(ctx context.Context, req *order.GetAllOrdersRequest) (*order.GetAllOrdersResponse, error) {
	domainOrders, err := s.orderUC.GetAllOrders()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all orders: %v", err)
	}

	var orders []*order.Order
	for _, o := range domainOrders {
		orders = append(orders, &order.Order{
			Id:         o.ID,
			UserId:     o.UserID,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedAt:  timestamppb.New(o.CreatedAt),
		})
	}

	return &order.GetAllOrdersResponse{
		Orders: orders,
	}, nil
}
