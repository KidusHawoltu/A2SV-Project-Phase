package controllers

import (
	domain "A2SV_ProjectPhase/Task8/TaskManager/Domain"
	usecases "A2SV_ProjectPhase/Task8/TaskManager/Usecases"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

func sendInternalErrorResponse(c *gin.Context, err error) {
	log.Println(err)
	sendErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred")
}

// User DTO
type UserRegisterLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Task DTOs
type CreateTaskRequest struct {
	Title       string            `json:"title" binding:"required"`
	Description string            `json:"description"`
	DueDate     time.Time         `json:"duedate" binding:"required" time_format:"2006-01-02"` // Example date format
	Status      domain.TaskStatus `json:"status" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       *string            `json:"title,omitempty"` // Pointers for optional fields
	Description *string            `json:"description,omitempty"`
	DueDate     *time.Time         `json:"duedate,omitempty" time_format:"2006-01-02"`
	Status      *domain.TaskStatus `json:"status,omitempty"`
}

// --- UserController ---

type UserController struct {
	uc *usecases.UserUseCase
}

func NewUserController(userUC *usecases.UserUseCase) *UserController {
	return &UserController{
		uc: userUC,
	}
}

func (controller *UserController) RegisterUser(c *gin.Context) {
	var userRegister UserRegisterLogin
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	registeredUser, err := controller.uc.RegisterUser(c.Request.Context(), userRegister.Username, userRegister.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameTaken) {
			sendErrorResponse(c, http.StatusConflict, err.Error())
			return
		} else if errors.Is(err, domain.ErrValidationFailed) {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		sendInternalErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":       registeredUser.Id,
		"username": registeredUser.Username,
		"role":     registeredUser.Role,
	})
}

func (controller *UserController) Login(c *gin.Context) {
	var userRegister UserRegisterLogin
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := controller.uc.Login(c.Request.Context(), userRegister.Username, userRegister.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			sendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		sendInternalErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// --- TaskController ---

type TaskController struct {
	uc *usecases.TaskUseCase
}

func NewTaskController(taskUC *usecases.TaskUseCase) *TaskController {
	return &TaskController{
		uc: taskUC,
	}
}

func (controller *TaskController) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	createdTask, err := controller.uc.CreateTask(c.Request.Context(), req.Title, req.Description, req.DueDate, req.Status)
	if err != nil {
		if errors.Is(err, domain.ErrValidationFailed) {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		sendInternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

func (controller *TaskController) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")

	task, err := controller.uc.GetTaskByID(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			sendErrorResponse(c, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrValidationFailed) {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		sendInternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (controller *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := controller.uc.GetAllTasks(c.Request.Context())
	if err != nil {
		sendInternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (controller *TaskController) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedTask, err := controller.uc.UpdateTask(
		c.Request.Context(),
		taskID,
		req.Title,       // Pass pointer for optional string
		req.Description, // Pass pointer for optional string
		req.DueDate,     // Pass pointer for optional time.Time
		req.Status,      // Pass pointer for optional TaskStatus
	)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			sendErrorResponse(c, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrValidationFailed) {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		sendInternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func (controller *TaskController) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	err := controller.uc.DeleteTask(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			sendErrorResponse(c, http.StatusNotFound, err.Error())
			return
		} else if errors.Is(err, domain.ErrValidationFailed) { // For invalid ID format
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		sendInternalErrorResponse(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
