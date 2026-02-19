package worker

import (
	"context"
	"log"
	"time"

	"task-management-system/internal/model"
	"task-management-system/internal/repository"
)

type TaskWorker struct {
	repo  *repository.TaskRepository
	queue chan string
	delay time.Duration
}

func NewTaskWorker(
	repo *repository.TaskRepository,
	queue chan string,
	delay time.Duration,
) *TaskWorker {

	return &TaskWorker{
		repo:  repo,
		queue: queue,
		delay: delay,
	}
}

func (w *TaskWorker) Start(ctx context.Context) {

	log.Println("Task worker started")

	for {
		select {

		case taskID := <-w.queue:
			go w.processTask(taskID)

		case <-ctx.Done():

			log.Println("Task worker stopped")
			return
		}
	}
}
func (w *TaskWorker) processTask(taskID string) {

	log.Printf("Processing task: %s\n", taskID)
	time.Sleep(w.delay)

	ctx := context.Background()

	task, err := w.repo.GetByID(ctx, taskID)

	if err != nil {

		log.Println("Worker DB error:", err)
		return
	}

	if task == nil {

		log.Println("Task not found:", taskID)
		return
	}

	if task.Status == model.StatusPending ||
		task.Status == model.StatusInProgress {

		err := w.repo.UpdateStatus(
			ctx,
			taskID,
			model.StatusCompleted,
		)

		if err != nil {

			log.Println("Failed to update task:", err)
			return
		}

		log.Printf("Task auto-completed: %s\n", taskID)
	}
}
