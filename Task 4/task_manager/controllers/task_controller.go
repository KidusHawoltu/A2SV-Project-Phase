package controllers

import (
	"A2SV_ProjectPhase/Task4/TaskManager/data"
	"A2SV_ProjectPhase/Task4/TaskManager/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var taskManager data.TaskManager

func Initialize(manager data.TaskManager) {
	taskManager = manager
}

func GetTasks(c *gin.Context) {
	tasks := taskManager.GetTasks()
	c.JSON(http.StatusOK, tasks)
}

func GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Id", "Error": err.Error()})
		return
	}
	task, err := taskManager.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
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
	task, err := taskManager.UpdateTask(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Id", "Error": err.Error()})
		return
	}

	if err := taskManager.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func AddTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect Task Body", "Error": err.Error()})
		return
	}
	if !task.Status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("There is no status called '%v'", task.Status)})
		return
	}
	newTask := taskManager.AddTask(task)
	c.JSON(http.StatusCreated, newTask)
}
