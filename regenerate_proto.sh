#!/bin/bash

# Удаляем старые сгенерированные файлы
rm -f proto/inventory/*.pb.go
rm -f proto/order/*.pb.go
rm -f proto/user/*.pb.go

# Определяем путь к бинарным файлам Go
GOPATH=$(go env GOPATH)
export PATH=$PATH:$GOPATH/bin

# Обновляем зависимости
go mod tidy

# Переустанавливаем protoc-gen-go и protoc-gen-go-grpc с совместимыми версиями
echo "Installing protoc plugins..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Определяем абсолютные пути к плагинам
PROTOC_GEN_GO="$GOPATH/bin/protoc-gen-go"
PROTOC_GEN_GO_GRPC="$GOPATH/bin/protoc-gen-go-grpc"

# Проверяем существование плагинов
echo "Checking plugins..."
if [ ! -f "$PROTOC_GEN_GO" ]; then
    echo "Error: $PROTOC_GEN_GO does not exist"
    exit 1
fi

if [ ! -f "$PROTOC_GEN_GO_GRPC" ]; then
    echo "Error: $PROTOC_GEN_GO_GRPC does not exist"
    exit 1
fi

# Делаем плагины исполняемыми
echo "Making plugins executable..."
chmod +x "$PROTOC_GEN_GO"
chmod +x "$PROTOC_GEN_GO_GRPC"

echo "Generating proto files..."
# Генерируем новые файлы с явным указанием пути к плагинам
protoc \
    --plugin="protoc-gen-go=${PROTOC_GEN_GO}" \
    --plugin="protoc-gen-go-grpc=${PROTOC_GEN_GO_GRPC}" \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/inventory/inventory.proto \
    proto/order/order.proto \
    proto/user/user.proto

# Проверяем успешность генерации
if [ $? -eq 0 ]; then
    echo "Regeneration completed successfully."
else
    echo "Error during regeneration. See errors above."
fi