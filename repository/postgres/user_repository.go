package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/utils"
)

// UserRepository реализует интерфейс domain.UserRepository для PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create сохраняет нового пользователя в базе данных
func (r *UserRepository) Create(user domain.User) error {
	// Проверяем, существует ли пользователь с таким username или email
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)`
	err := r.db.QueryRow(query, user.Username, user.Email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существующего пользователя: %w", err)
	}
	if exists {
		return errors.New("пользователь с таким именем или email уже существует")
	}

	// Генерируем ID пользователя, если он не задан
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Хешируем пароль
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	// Текущее время для полей created_at и updated_at
	now := time.Now()

	// Вставляем нового пользователя
	query = `
		INSERT INTO users (id, username, email, full_name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = r.db.Exec(query, user.ID, user.Username, user.Email, user.FullName, hashedPassword, now, now)
	if err != nil {
		return fmt.Errorf("ошибка сохранения пользователя: %w", err)
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(id string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("пользователь не найден")
		}
		return domain.User{}, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// GetByUsername получает пользователя по имени пользователя
func (r *UserRepository) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("пользователь не найден")
		}
		return domain.User{}, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("пользователь не найден")
		}
		return domain.User{}, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(id string, user domain.User) error {
	// Проверяем, существует ли пользователь с таким ID
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if !exists {
		return errors.New("пользователь не найден")
	}

	// Обновляем пользователя
	query = `
		UPDATE users
		SET username = $1, email = $2, full_name = $3, updated_at = $4
		WHERE id = $5
	`
	_, err = r.db.Exec(query, user.Username, user.Email, user.FullName, time.Now(), id)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя: %w", err)
	}

	return nil
}

// Delete удаляет пользователя по ID
func (r *UserRepository) Delete(id string) error {
	// Проверяем, существует ли пользователь с таким ID
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if !exists {
		return errors.New("пользователь не найден")
	}

	// Удаляем пользователя
	query = `DELETE FROM users WHERE id = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %w", err)
	}

	return nil
}
