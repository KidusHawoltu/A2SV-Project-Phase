package routers

import (
	"A2SV_ProjectPhase/Task8/TaskManager/Delivery/controllers"
	infrastructure "A2SV_ProjectPhase/Task8/TaskManager/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupUserRouters(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/register", userController.RegisterUser)
		userRoutes.POST("/login", userController.Login)
	}
}

func SetupTaskRoutes(router *gin.Engine, taskController *controllers.TaskController, authMiddleware *infrastructure.AuthMiddleware) {
	taskRoutes := router.Group("/tasks")
	{
		authenticatedTaskRoutes := taskRoutes.Use(authMiddleware.Authenticate())
		{
			authenticatedTaskRoutes.GET("/", taskController.GetAllTasks)
			authenticatedTaskRoutes.GET("/:id", taskController.GetTaskByID)
		}

		adminTaskRoutes := taskRoutes.Group("/")
		// Apply Authenticate() FIRST, then AuthorizeAdmin()
		adminTaskRoutes.Use(authMiddleware.Authenticate(), authMiddleware.AuthorizeAdmin())
		{
			adminTaskRoutes.POST("/", taskController.CreateTask)
			adminTaskRoutes.PUT("/:id", taskController.UpdateTask)
			adminTaskRoutes.DELETE("/:id", taskController.DeleteTask)
		}
	}
}
