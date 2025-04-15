package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"FoodStore-AdvProg2/domain"
)

// UserUseCase содержит бизнес-логику для работы с пользователями
type UserUseCase struct {
	repo domain.UserRepository
}

// NewUserUseCase создает новый экземпляр UserUseCase
func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) RegisterUser(username, email, fullName, password string) (*domain.User, string, error) {
	// Проверка существования пользователя
	existingUser, err := uc.repo.FindByUsername(username)
	if (err == nil && existingUser != nil) || (err == nil && existingUser != nil) {
		return nil, "", errors.New("username already taken")
	}

	existingUser, err = uc.repo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("error hashing password: %w", err)
	}

	user := &domain.User{
		ID:       uuid.New().String(),
		Username: username,
		Email:    email,
		FullName: fullName,
		Password: string(hashedPassword),
	}

	if err := uc.repo.Save(user); err != nil {
		return nil, "", fmt.Errorf("error saving user: %w", err)
	}

	// Создаем токен (в продакшене лучше использовать JWT)
	token := generateToken(user.ID)

	return user, token, nil
}

func (uc *UserUseCase) AuthenticateUser(username, password string) (*domain.User, string, error) {
	// Поиск пользователя по имени
	user, err := uc.repo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// Создаем токен
	token := generateToken(user.ID)

	return user, token, nil
}

func (uc *UserUseCase) GetUserProfile(userID string) (*domain.User, error) {
	user, err := uc.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Create создает нового пользователя
func (uc *UserUseCase) Create(user domain.User) error {
	return uc.repo.Create(user)
}

// GetByID возвращает пользователя по ID
func (uc *UserUseCase) GetByID(id string) (domain.User, error) {
	return uc.repo.GetByID(id)
}

// GetByUsername возвращает пользователя по имени пользователя
func (uc *UserUseCase) GetByUsername(username string) (domain.User, error) {
	return uc.repo.GetByUsername(username)
}

// GetByEmail возвращает пользователя по email
func (uc *UserUseCase) GetByEmail(email string) (domain.User, error) {
	return uc.repo.GetByEmail(email)
}

// Update обновляет информацию о пользователе
func (uc *UserUseCase) Update(id string, user domain.User) error {
	return uc.repo.Update(id, user)
}

// Delete удаляет пользователя
func (uc *UserUseCase) Delete(id string) error {
	return uc.repo.Delete(id)
}

// Вспомогательная функция для генерации токена
func generateToken(userID string) string {
	// В реальном приложении здесь должна быть реализация JWT
	return fmt.Sprintf("%s-%d", userID, time.Now().UnixNano())
}
