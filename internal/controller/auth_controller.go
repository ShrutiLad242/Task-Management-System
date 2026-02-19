package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task-management-system/internal/auth"
	"task-management-system/config"
	"task-management-system/internal/model"
	"task-management-system/internal/repository"
)

type AuthController struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthController(
	userRepo *repository.UserRepository,
	cfg *config.Config,
) *AuthController {

	return &AuthController{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}

func (ac *AuthController) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx := c.Request.Context()

	user, err := ac.userRepo.GetByEmail(ctx, req.Email)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	if user == nil {

		user = &model.User{
			ID:        uuid.NewString(),
			Email:     req.Email,
			Role:      model.RoleUser,
			CreatedAt: time.Now(),
		}

		err := ac.userRepo.Create(ctx, user)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}
	}

	token, err := auth.GenerateToken(
		user.ID,
		string(user.Role),
		ac.cfg.JWTSecret,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
