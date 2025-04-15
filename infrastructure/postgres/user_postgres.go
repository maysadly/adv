package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"FoodStore-AdvProg2/domain"
)

type UserPostgresRepo struct{}

func NewUserPostgresRepo() *UserPostgresRepo {
	return &UserPostgresRepo{}
}

func (r *UserPostgresRepo) Save(user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `INSERT INTO users (id, username, email, full_name, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := DB.Exec(
		context.Background(),
		query,
		user.ID,
		user.Username,
		user.Email,
		user.FullName,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserPostgresRepo) FindByID(id string) (*domain.User, error) {
	query := `SELECT id, username, email, full_name, password, created_at, updated_at 
              FROM users 
              WHERE id = $1`

	var user domain.User
	err := DB.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgresRepo) FindByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, email, full_name, password, created_at, updated_at 
              FROM users 
              WHERE username = $1`

	var user domain.User
	err := DB.QueryRow(context.Background(), query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgresRepo) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, email, full_name, password, created_at, updated_at 
              FROM users 
              WHERE email = $1`

	var user domain.User
	err := DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgresRepo) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `UPDATE users 
              SET username = $1, email = $2, full_name = $3, password = $4, updated_at = $5
              WHERE id = $6`

	_, err := DB.Exec(
		context.Background(),
		query,
		user.Username,
		user.Email,
		user.FullName,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	return err
}
