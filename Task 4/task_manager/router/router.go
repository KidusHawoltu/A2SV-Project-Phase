package router

import (
	"A2SV_ProjectPhase/Task4/TaskManager/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter(taskController *controllers.TaskController) *gin.Engine {
	router := gin.Default()
	router.GET("/tasks", taskController.GetTasks)
	router.GET("/tasks/:id", taskController.GetTaskById)
	router.PUT("/tasks/:id", taskController.UpdateTask)
	router.DELETE("/tasks/:id", taskController.DeleteTask)
	router.POST("/tasks", taskController.AddTask)
	return router
}
