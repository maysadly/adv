# FoodStore - E-commerce платформа с микросервисной архитектурой

Реализация платформы электронной коммерции с использованием микросервисной архитектуры и gRPC для межсервисного взаимодействия.

## Архитектура

Система состоит из следующих микросервисов:

1. **API Gateway** - Принимает REST запросы от клиентов и преобразует их в gRPC вызовы к соответствующим сервисам.
2. **Inventory Service** - Управление товарами и категориями.
3. **Order Service** - Управление заказами и платежами.
4. **User Service** - Управление пользователями (регистрация, аутентификация, профили).

## Технологии

- Go (Golang)
- gRPC для межсервисного взаимодействия
- Gin для REST API в API Gateway
- PostgreSQL в качестве базы данных
- Protocol Buffers для описания API

## Запуск проекта

### Предварительные требования

- Go 1.17+
- PostgreSQL
- protoc (Protocol Buffers Compiler)
- protoc-gen-go и protoc-gen-go-grpc
- Make (опционально)

### Установка зависимостей

```bash
go mod tidy
```

### Настройка окружения

1. Скопируйте файл .env.example в .env:

```bash
cp .env.example .env
```

2. Отредактируйте .env, указав правильные настройки для подключения к базе данных.

### Генерация кода из .proto файлов

```bash
make proto
```

или без Make:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/inventory/inventory.proto proto/order/order.proto proto/user/user.proto
```

### Сборка сервисов

```bash
make build
```

### Запуск сервисов

Для запуска каждого сервиса в отдельном терминале:

```bash
# Запуск Inventory Service
make run-inventory

# Запуск Order Service
make run-order

# Запуск User Service
make run-user

# Запуск API Gateway
make run-gateway
```

Или можно запустить собранные бинарные файлы из директории bin/:

```bash
./bin/inventory-service
./bin/order-service
./bin/user-service
./bin/api-gateway
```

## API Gateway Endpoints

### Товары (Products)

- `GET /api/products` - Получить список товаров
- `GET /api/products/{id}` - Получить товар по ID
- `POST /api/products` - Создать новый товар
- `PUT /api/products/{id}` - Обновить товар
- `DELETE /api/products/{id}` - Удалить товар

### Заказы (Orders)

- `GET /api/orders` - Получить список всех заказов
- `GET /api/orders?user_id=123` - Получить заказы конкретного пользователя
- `GET /api/orders/{id}` - Получить заказ по ID
- `POST /api/orders` - Создать новый заказ
- `PATCH /api/orders/{id}` - Обновить статус заказа

### Пользователи (Users)

- `POST /api/users/register` - Зарегистрировать нового пользователя
- `POST /api/users/login` - Аутентифицировать пользователя
- `GET /api/users/profile` - Получить профиль пользователя (требует авторизации)

## Веб-интерфейс

После запуска всех сервисов веб-интерфейс доступен по адресу:

- http://localhost:8080/order - Клиентский интерфейс для создания заказов
- http://localhost:8080/admin - Административный интерфейс для управления товарами