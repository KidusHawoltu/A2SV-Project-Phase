package main

import (
	"A2SV_ProjectPhase/Task4/TaskManager/controllers"
	"A2SV_ProjectPhase/Task4/TaskManager/data"
	"A2SV_ProjectPhase/Task4/TaskManager/router"
)

func main() {
	manager := data.NewTaskManager()
	controller := controllers.NewTaskController(manager)
	router := router.NewRouter(controller)
	router.Run(":8080")
}
