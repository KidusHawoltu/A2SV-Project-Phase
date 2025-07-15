package main

import (
	"library_management/controllers"
	"library_management/services"
)

func main() {
	controllers.Handler(services.NewLibrary())
}
