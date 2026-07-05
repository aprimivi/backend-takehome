package repository

import (
	"context"
	"database/sql"
	"errors"

	model "app/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES (?, ?, ?)
	`, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return model.User{}, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return model.User{}, err
	}

	return r.GetByID(ctx, insertedID)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}
