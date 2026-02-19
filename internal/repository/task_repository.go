package repository

import (
	"context"
	"database/sql"
	"errors"
	"task-management-system/internal/model"
	"time"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {

	query := `
	INSERT INTO tasks (id, title, description, status, user_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UserID,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {

	query := `
	SELECT id, title, description, status, user_id, created_at, updated_at
	FROM tasks
	WHERE id = $1
	`

	task := &model.Task{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return task, nil
}
func (r *TaskRepository) Delete(ctx context.Context, id string) error {

	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status model.TaskStatus,
) error {

	query := `
	UPDATE tasks
	SET status = $1, updated_at = $2
	WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		time.Now(),
		id,
	)

	return err
}

func (r *TaskRepository) GetAll(ctx context.Context, userID string, role model.Role) ([]model.Task, error) {

	var query string
	var rows *sql.Rows
	var err error

	if role == model.RoleAdmin || userID == "" {
		query = `
		SELECT id, title, description, status, user_id, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
		`
		rows, err = r.db.QueryContext(ctx, query)
	} else {
		query = `
		SELECT id, title, description, status, user_id, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
		`
		rows, err = r.db.QueryContext(ctx, query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		task := model.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
