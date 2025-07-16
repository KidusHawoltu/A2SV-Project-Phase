package controllers

import (
	"A2SV_ProjectPhase/Task4/TaskManager/data"
	"A2SV_ProjectPhase/Task4/TaskManager/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskManager data.TaskManager
}

func NewController(manager data.TaskManager) *TaskController {
	return &TaskController{
		taskManager: manager,
	}
}

func (taskController *TaskController) GetTasks(c *gin.Context) {
	tasks := taskController.taskManager.GetTasks()
	c.JSON(http.StatusOK, tasks)
}

func (taskController *TaskController) GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Id", "Error": err.Error()})
		return
	}
	task, err := taskController.taskManager.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (taskController *TaskController) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Id", "Error": err.Error()})
		return
	}
	var updatedTask models.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Task Body", "Error": err.Error()})
		return
	}
	if !updatedTask.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("There is no status called '%v'", updatedTask.Status)})
		return
	}
	task, err := taskController.taskManager.UpdateTask(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (taskController *TaskController) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Id", "Error": err.Error()})
		return
	}

	if err := taskController.taskManager.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (taskController *TaskController) AddTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Task Body", "Error": err.Error()})
		return
	}
	if !task.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("There is no status called '%v'", task.Status)})
		return
	}
	newTask := taskController.taskManager.AddTask(task)
	c.JSON(http.StatusCreated, newTask)
}
