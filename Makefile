.PHONY: proto build run clean

# Генерация proto файлов
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/inventory/inventory.proto \
		proto/order/order.proto \
		proto/user/user.proto

# Сборка всех сервисов
build: proto
	go build -o bin/api-gateway ./cmd/api-gateway
	go build -o bin/inventory-service ./cmd/inventory-service
	go build -o bin/order-service ./cmd/order-service
	go build -o bin/user-service ./cmd/user-service

# Запуск каждого сервиса по отдельности
run-inventory:
	go run ./cmd/inventory-service/main.go

run-order:
	go run ./cmd/order-service/main.go

run-user:
	go run ./cmd/user-service/main.go

run-gateway:
	go run ./cmd/api-gateway/main.go

# Очистка собранных бинарных файлов
clean:
	rm -rf bin/*