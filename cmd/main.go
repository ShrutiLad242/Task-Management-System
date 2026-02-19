package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"task-management-system/config"
	"task-management-system/internal/db"
	"task-management-system/internal/repository"
	"task-management-system/internal/service"
	"task-management-system/internal/worker"
	"task-management-system/internal/controller"
	"task-management-system/internal/middleware"

)

func main() {

	// Load config
	cfg := config.Load()

	// Connect DB and run migrations
	database := db.NewPostgres(cfg)
	defer database.Close()

	// Create repositories
	taskRepo := repository.NewTaskRepository(database)

	// Create task queue channel
	taskQueue := make(chan string, 100)

	// Create service
	taskService := service.NewTaskService(taskRepo, taskQueue)

	taskController := controller.NewTaskController(taskService)

	userRepo := repository.NewUserRepository(database)

	authController := controller.NewAuthController(
		userRepo,
		cfg,
	)



	// Create worker
	taskWorker := worker.NewTaskWorker(
		taskRepo,
		taskQueue,
		cfg.AutoCompleteDelay,
	)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Start worker
	go taskWorker.Start(ctx)

	// Create Gin router
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/login", authController.Login)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))

	protected.POST("/tasks", taskController.CreateTask)
	protected.GET("/tasks", taskController.GetTasks)
	protected.GET("/tasks/:id", taskController.GetTask)
	protected.DELETE("/tasks/:id", taskController.DeleteTask)

	go func() {

		log.Println("Server running on port", cfg.Port)

		err := router.Run(":" + cfg.Port)

		if err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)

	signal.Notify(stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutting down...")

	cancel()
}
