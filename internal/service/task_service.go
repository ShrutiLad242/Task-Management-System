package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"task-management-system/internal/model"
	"task-management-system/internal/repository"
)

type TaskService struct {
	repo      *repository.TaskRepository
	taskQueue chan string
}

func NewTaskService(
	repo *repository.TaskRepository,
	taskQueue chan string,
) *TaskService {

	return &TaskService{
		repo:      repo,
		taskQueue: taskQueue,
	}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
	description string,
	userID string,
) (*model.Task, error) {

	task := &model.Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: description,
		Status:      model.StatusPending,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(ctx, task)

	if err != nil {
		if strings.Contains(err.Error(), "tasks_user_id_fkey") || strings.Contains(err.Error(), "SQLSTATE 23503") {
			return nil, errors.New("invalid user id: user does not exist")
		}
		return nil, err
	}

	s.taskQueue <- task.ID

	return task, nil
}

func (s *TaskService) GetTaskByID(
	ctx context.Context,
	taskID string,
	userID string,
	role model.Role,
) (*model.Task, error) {

	task, err := s.repo.GetByID(ctx, taskID)

	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	// Authorization check
	if role != model.RoleAdmin && task.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return task, nil
}

func (s *TaskService) GetAllTasks(
	ctx context.Context,
	userID string,
	role model.Role,
) ([]model.Task, error) {

	return s.repo.GetAll(ctx, userID, role)
}

func (s *TaskService) DeleteTask(
	ctx context.Context,
	taskID string,
	userID string,
	role model.Role,
) error {

	task, err := s.repo.GetByID(ctx, taskID)

	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	if role != model.RoleAdmin && task.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.Delete(ctx, taskID)
}

func (s *TaskService) CompleteTask(
	ctx context.Context,
	taskID string,
) error {

	task, err := s.repo.GetByID(ctx, taskID)

	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	if task.Status == model.StatusCompleted {
		return nil
	}

	return s.repo.UpdateStatus(ctx, taskID, model.StatusCompleted)
}
