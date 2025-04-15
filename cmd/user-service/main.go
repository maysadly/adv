package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/infrastructure/postgres"
	"FoodStore-AdvProg2/proto/user"
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

	// Инициализируем таблицы, если они еще не созданы
	if err := postgres.InitTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	// Создание репозитория и use case
	userRepo := postgres.NewUserPostgresRepo()
	userUC := usecase.NewUserUseCase(userRepo)

	// Настройка gRPC сервера
	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	userServer := NewUserServiceServer(userUC)
	user.RegisterUserServiceServer(server, userServer)

	// Включаем reflection для отладки
	reflection.Register(server)

	log.Printf("User Service is starting on port %s...", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// UserServiceServer реализует gRPC сервер для User Service
type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	userUC *usecase.UserUseCase
}

func NewUserServiceServer(userUC *usecase.UserUseCase) *UserServiceServer {
	return &UserServiceServer{
		userUC: userUC,
	}
}

// RegisterUser регистрирует нового пользователя
func (s *UserServiceServer) RegisterUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	// Проверяем, что пользователя с таким email или username не существует
	_, err := s.userUC.GetByUsername(req.Username)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user with username %s already exists", req.Username)
	}

	_, err = s.userUC.GetByEmail(req.Email)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", req.Email)
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Создаем пользователя
	newUser := domain.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Password: string(hashedPassword),
	}

	if err := s.userUC.Create(newUser); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &user.UserResponse{
		User: &user.User{
			Id:       newUser.ID,
			Username: newUser.Username,
			Email:    newUser.Email,
			FullName: newUser.FullName,
		},
	}, nil
}

// AuthenticateUser аутентифицирует пользователя
func (s *UserServiceServer) AuthenticateUser(ctx context.Context, req *user.AuthRequest) (*user.AuthResponse, error) {
	// Ищем пользователя по имени пользователя
	userEntity, err := s.userUC.GetByUsername(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	return &user.AuthResponse{
		User: &user.User{
			Id:       userEntity.ID,
			Username: userEntity.Username,
			Email:    userEntity.Email,
			FullName: userEntity.FullName,
		},
	}, nil
}

// GetUserProfile возвращает профиль пользователя
func (s *UserServiceServer) GetUserProfile(ctx context.Context, req *user.UserProfileRequest) (*user.UserProfile, error) {
	userEntity, err := s.userUC.GetByID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &user.UserProfile{
		User: &user.User{
			Id:       userEntity.ID,
			Username: userEntity.Username,
			Email:    userEntity.Email,
			FullName: userEntity.FullName,
		},
	}, nil
}

// Login аутентифицирует пользователя
func (s *UserServiceServer) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	userData, err := s.userUC.GetByUsername(req.Username)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	return &user.LoginResponse{
		User: &user.User{
			Id:       userData.ID,
			Username: userData.Username,
			Email:    userData.Email,
			FullName: userData.FullName,
		},
	}, nil
}
