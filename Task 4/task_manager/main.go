package main

import (
	"A2SV_ProjectPhase/Task4/TaskManager/controllers"
	"A2SV_ProjectPhase/Task4/TaskManager/data"
	"A2SV_ProjectPhase/Task4/TaskManager/router"
)

func main() {
	manager := data.NewTaskManager()
	controllers.Initialize(manager)
	router := router.NewRouter()
	router.Run(":8080")
}
