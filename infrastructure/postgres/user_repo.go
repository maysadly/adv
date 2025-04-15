package postgres

import (
	"database/sql"
	"errors"
	"time"

	"FoodStore-AdvProg2/domain"
)

// UserPostgresRepo реализует интерфейс UserRepository для PostgreSQL
type UserPostgresRepo struct{}

// NewUserPostgresRepo создает новый экземпляр UserPostgresRepo
func NewUserPostgresRepo() *UserPostgresRepo {
	return &UserPostgresRepo{}
}

// Create добавляет нового пользователя в базу данных
func (r *UserPostgresRepo) Create(user domain.User) error {
	now := time.Now()

	query := `
		INSERT INTO users 
		(id, username, email, full_name, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.FullName,
		user.Password,
		now,
		now,
	)
	return err
}

// GetByID возвращает пользователя по ID
func (r *UserPostgresRepo) GetByID(id string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return user, err
}

// GetByUsername возвращает пользователя по имени пользователя
func (r *UserPostgresRepo) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return user, err
}

// GetByEmail возвращает пользователя по email
func (r *UserPostgresRepo) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, full_name, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return domain.User{}, errors.New("user not found")
	}
	return user, err
}

// Update обновляет информацию о пользователе
func (r *UserPostgresRepo) Update(id string, user domain.User) error {
	query := `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			full_name = $3,
			password = $4,
			updated_at = $5
		WHERE id = $6
	`
	_, err := db.Exec(
		query,
		user.Username,
		user.Email,
		user.FullName,
		user.Password,
		time.Now(),
		id,
	)
	return err
}

// Delete удаляет пользователя по ID
func (r *UserPostgresRepo) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
