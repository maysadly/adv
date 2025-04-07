package usecase

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/repository"
	"errors"
)

type OrderUseCase struct {
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
}

func NewOrderUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderUseCase {
	return &OrderUseCase{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
	}
}

func (uc *OrderUseCase) CreateOrder(orderReq domain.OrderRequest) (string, error) {
	if len(orderReq.Items) == 0 {
		return "", errors.New("заказ должен содержать хотя бы один товар")
	}

	// Подготовка заказа
	var totalPrice float64
	var orderItems []domain.OrderItem

	for _, item := range orderReq.Items {
		if item.Quantity <= 0 {
			return "", errors.New("количество товара должно быть больше 0")
		}

		product, err := uc.ProductRepo.FindByID(item.ProductID)
		if err != nil {
			return "", errors.New("товар не найден")
		}

		if product.Stock < item.Quantity {
			return "", errors.New("недостаточное количество товара")
		}

		// Создаем элемент заказа
		orderItem := domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)

		// Обновляем общую цену
		totalPrice += product.Price * float64(item.Quantity)
	}

	order := domain.Order{
		UserID:     orderReq.UserID,
		TotalPrice: totalPrice,
		Status:     domain.OrderStatusPending,
	}

	// Сохраняем заказ и его элементы
	orderID, err := uc.OrderRepo.Save(order, orderItems)
	if err != nil {
		return "", err
	}

	// Обновляем складские запасы
	for _, item := range orderReq.Items {
		product, _ := uc.ProductRepo.FindByID(item.ProductID)
		product.Stock -= item.Quantity
		uc.ProductRepo.Update(product.ID, product)
	}

	return orderID, nil
}

func (uc *OrderUseCase) GetOrderByID(id string) (domain.Order, error) {
	order, items, err := uc.OrderRepo.FindByID(id)
	if err != nil {
		return domain.Order{}, err
	}

	order.Items = items
	return order, nil
}

func (uc *OrderUseCase) UpdateOrderStatus(id string, status string) error {
	// Валидация статуса
	if status != domain.OrderStatusPending &&
		status != domain.OrderStatusCompleted &&
		status != domain.OrderStatusCancelled {
		return errors.New("неверный статус заказа")
	}

	// Получаем текущий заказ
	order, _, err := uc.OrderRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Если отменяем заказ, возвращаем товары на склад
	if status == domain.OrderStatusCancelled && order.Status != domain.OrderStatusCancelled {
		_, items, _ := uc.OrderRepo.FindByID(id)
		for _, item := range items {
			product, _ := uc.ProductRepo.FindByID(item.ProductID)
			product.Stock += item.Quantity
			uc.ProductRepo.Update(product.ID, product)
		}
	}

	return uc.OrderRepo.UpdateStatus(id, status)
}

func (uc *OrderUseCase) GetOrdersByUserID(userID string) ([]domain.Order, error) {
	return uc.OrderRepo.FindByUserID(userID)
}

func (uc *OrderUseCase) GetAllOrders() ([]domain.Order, error) {
	return uc.OrderRepo.FindAll()
}
