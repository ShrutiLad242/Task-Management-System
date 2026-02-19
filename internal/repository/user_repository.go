package repository

import (
	"context"
	"database/sql"
	"errors"

	"task-management-system/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {

	query := `
	INSERT INTO users (id, email, password, role, created_at)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Password,
		user.Role,
		user.CreatedAt,
	)

	return err
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	query := `
	SELECT id, email, password, role, created_at
	FROM users
	WHERE email = $1
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id string,
) (*model.User, error) {

	query := `
	SELECT id, email, password, role, created_at
	FROM users
	WHERE id = $1
	`

	user := &model.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
