package router

import (
	"A2SV_ProjectPhase/Task6/TaskManager/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter(taskController *controllers.TaskController) *gin.Engine {
	router := gin.Default()
	tasks := router.Group("/tasks")
	tasks.GET("", taskController.GetTasks)
	tasks.GET("/:id", taskController.GetTaskById)
	tasks.PUT("/:id", taskController.UpdateTask)
	tasks.DELETE("/:id", taskController.DeleteTask)
	tasks.POST("", taskController.AddTask)
	return router
}
