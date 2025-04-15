package domain

import (
	"time"
)

// User представляет модель пользователя
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Password  string    `json:"-"` // Не возвращаем пароль в JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository представляет интерфейс репозитория для работы с пользователями
type UserRepository interface {
	Create(user User) error
	GetByID(id string) (User, error)
	GetByUsername(username string) (User, error)
	GetByEmail(email string) (User, error)
	Update(id string, user User) error
	Delete(id string) error
}

// UserUseCase представляет интерфейс сервиса для работы с пользователями
type UserUseCase interface {
	Create(user User) error
	GetByID(id string) (User, error)
	GetByUsername(username string) (User, error)
	GetByEmail(email string) (User, error)
	Update(id string, user User) error
	Delete(id string) error
}
