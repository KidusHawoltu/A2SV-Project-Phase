package router

import (
	"A2SV_ProjectPhase/Task6/TaskManager/controllers"
	"A2SV_ProjectPhase/Task6/TaskManager/middleware"
	"A2SV_ProjectPhase/Task6/TaskManager/models"

	"github.com/gin-gonic/gin"
)

func NewRouter(taskController *controllers.TaskController) *gin.Engine {
	router := gin.Default()

	tasks := router.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	tasks.GET("", taskController.GetTasks)
	tasks.GET("/:id", taskController.GetTaskById)
	adminTasks := tasks.Group("")
	adminTasks.Use(middleware.AuthorizeRole(models.RoleAdmin))
	adminTasks.PUT("/:id", taskController.UpdateTask)
	adminTasks.DELETE("/:id", taskController.DeleteTask)
	adminTasks.POST("", taskController.AddTask)

	router.POST("/register", taskController.RegisterUser)
	router.POST("/login", taskController.Login)
	return router
}
