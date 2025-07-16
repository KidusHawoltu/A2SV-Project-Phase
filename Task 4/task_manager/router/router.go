package router

import (
	"A2SV_ProjectPhase/Task4/TaskManager/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/tasks", controllers.GetTasks)
	router.GET("/tasks/:id", controllers.GetTaskById)
	router.PUT("/tasks/:id", controllers.UpdateTask)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
	router.POST("/tasks", controllers.AddTask)
	return router
}
