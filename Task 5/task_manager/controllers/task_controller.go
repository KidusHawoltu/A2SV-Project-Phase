package controllers

import (
	"A2SV_ProjectPhase/Task5/TaskManager/data"
	"A2SV_ProjectPhase/Task5/TaskManager/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	taskManager data.TaskManager
}

func NewTaskController(tm data.TaskManager) *TaskController {
	return &TaskController{
		taskManager: tm,
	}
}

// Helper function for generic error responses
func sendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"error":   err.Error(),
	})
}

// Helper function for 404 Not Found responses
func sendNotFoundError(c *gin.Context, id string, err error) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": fmt.Sprintf("Task with ID '%s' not found", id),
		"error":   err.Error(),
	})
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	// Pass context from Gin to the manager
	tasks, err := tc.taskManager.GetTasks(c.Request.Context())
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tasks", err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) GetTaskById(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid Task ID format", err)
		return
	}

	// Pass context from Gin to the manager
	task, err := tc.taskManager.GetTaskById(c.Request.Context(), id)
	if err != nil {
		// Use errors.Is to check for ErrTaskNotFound
		if errors.Is(err, data.ErrTaskNotFound) {
			sendNotFoundError(c, idHex, err)
			return
		}
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve task", err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid Task ID format", err)
		return
	}

	var updatedTask models.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if !updatedTask.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid status provided: '%v'", updatedTask.Status)})
		return
	}

	// Pass context from Gin to the manager
	task, err := tc.taskManager.UpdateTask(c.Request.Context(), id, updatedTask)
	if err != nil {
		// Use errors.Is to check for ErrTaskNotFound
		if errors.Is(err, data.ErrTaskNotFound) {
			sendNotFoundError(c, idHex, err)
			return
		}
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to update task", err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid Task ID format", err)
		return
	}

	// Pass context from Gin to the manager
	if err := tc.taskManager.DeleteTask(c.Request.Context(), id); err != nil {
		// Use errors.Is to check for ErrTaskNotFound
		if errors.Is(err, data.ErrTaskNotFound) {
			sendNotFoundError(c, idHex, err)
			return
		}
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to delete task", err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (tc *TaskController) AddTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if !task.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid status provided: '%v'", task.Status)})
		return
	}

	// Pass context from Gin to the manager
	newTask, err := tc.taskManager.AddTask(c.Request.Context(), task)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to create task", err)
		return
	}
	c.JSON(http.StatusCreated, newTask)
}
