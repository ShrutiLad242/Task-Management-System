package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"task-management-system/internal/service"
	"task-management-system/internal/model"
)

type TaskController struct {
	service *service.TaskService
}

func NewTaskController(service *service.TaskService) *TaskController {
	return &TaskController{
		service: service,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (tc *TaskController) CreateTask(c *gin.Context) {

	var req CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")

	task, err := tc.service.CreateTask(
		c.Request.Context(),
		req.Title,
		req.Description,
		userID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "invalid user id") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}


func (tc *TaskController) GetTasks(c *gin.Context) {

	userID := c.GetString("user_id")
	role := model.Role(c.GetString("role"))

	tasks, err := tc.service.GetAllTasks(
		c.Request.Context(),
		userID,
		role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}


func (tc *TaskController) GetTask(c *gin.Context) {

	taskID := c.Param("id")

	userID := c.GetString("user_id")
	role := model.Role(c.GetString("role"))

	task, err := tc.service.GetTaskByID(
		c.Request.Context(),
		taskID,
		userID,
		role,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}


func (tc *TaskController) DeleteTask(c *gin.Context) {

	taskID := c.Param("id")

	userID := c.GetString("user_id")
	role := model.Role(c.GetString("role"))

	err := tc.service.DeleteTask(
		c.Request.Context(),
		taskID,
		userID,
		role,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}
