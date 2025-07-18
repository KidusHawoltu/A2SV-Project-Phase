package controllers

import (
	"A2SV_ProjectPhase/Task5/TaskManager/data"
	"A2SV_ProjectPhase/Task5/TaskManager/models"
	"fmt"
	"net/http"
	"strings" // For checking error messages

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

func (tc *TaskController) GetTasks(c *gin.Context) {
	// Pass context from Gin to the manager
	tasks, err := tc.taskManager.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to retrieve tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) GetTaskById(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Task ID format", "error": err.Error()})
		return
	}

	// Pass context from Gin to the manager
	task, err := tc.taskManager.GetTaskById(c.Request.Context(), id)
	if err != nil {
		// Check for specific "not found" error message from the data layer
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID '%s' not found", idHex), "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to retrieve task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Task ID format", "error": err.Error()})
		return
	}

	var updatedTask models.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	if !updatedTask.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid status provided: '%v'", updatedTask.Status)})
		return
	}

	// Pass context from Gin to the manager
	task, err := tc.taskManager.UpdateTask(c.Request.Context(), id, updatedTask)
	if err != nil {
		// Check for specific "not found" error message from the data layer
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID '%s' not found", idHex), "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	idHex := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Task ID format", "error": err.Error()})
		return
	}

	// Pass context from Gin to the manager
	if err := tc.taskManager.DeleteTask(c.Request.Context(), id); err != nil {
		// Check for specific "not found" error message from the data layer
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Task with ID '%s' not found", idHex), "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (tc *TaskController) AddTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}
	if !task.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid status provided: '%v'", task.Status)})
		return
	}

	// Pass context from Gin to the manager
	newTask, err := tc.taskManager.AddTask(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, newTask)
}
