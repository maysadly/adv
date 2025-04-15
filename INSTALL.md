# Инструкция по установке и настройке проекта

## 1. Установка Protocol Buffers (protoc)

### Для macOS:
```bash
brew install protobuf
```

### Для Ubuntu/Debian:
```bash
sudo apt update
sudo apt install protobuf-compiler
```

### Для Windows:
Скачайте последнюю версию с [GitHub](https://github.com/protocolbuffers/protobuf/releases)
или используйте Chocolatey:
```bash
choco install protoc
```

## 2. Установка Go плагинов для Protocol Buffers

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Добавьте GOPATH/bin в PATH:
```bash
# Для Bash/Zsh
export PATH="$PATH:$(go env GOPATH)/bin"

# Для Windows PowerShell
$env:Path += ";$(go env GOPATH)\bin"
```

## 3. Установка зависимостей Go

```bash
go mod tidy
```

## 4. Генерация кода из proto-файлов

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/inventory/inventory.proto \
    proto/order/order.proto \
    proto/user/user.proto
```

## 5. Настройка PostgreSQL

### Установка PostgreSQL

#### Для macOS:
```bash
brew install postgresql
brew services start postgresql
```

#### Для Ubuntu/Debian:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### Создание базы данных и пользователя

```bash
sudo -u postgres psql

# В консоли PostgreSQL:
CREATE DATABASE foodstore;
CREATE USER foodstore_user WITH ENCRYPTED PASSWORD 'your_password_here';
GRANT ALL PRIVILEGES ON DATABASE foodstore TO foodstore_user;
\q
```

### Настройка .env файла

```bash
cp .env.example .env
```

Отредактируйте .env файл:
```
DB=postgresql://foodstore_user:your_password_here@localhost:5432/foodstore
```

## 6. Сборка и запуск проекта

### Сборка всех сервисов
```bash
mkdir -p bin
go build -o bin/inventory-service ./cmd/inventory-service
go build -o bin/order-service ./cmd/order-service
go build -o bin/user-service ./cmd/user-service
go build -o bin/api-gateway ./cmd/api-gateway
```

### Запуск сервисов (каждый в отдельном терминале)
```bash
# Терминал 1
./bin/inventory-service

# Терминал 2
./bin/order-service

# Терминал 3
./bin/user-service

# Терминал 4
./bin/api-gateway
```

После запуска всех сервисов, приложение будет доступно по адресу: http://localhost:8080