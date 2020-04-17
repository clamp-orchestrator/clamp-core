package main

import (
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
	"log"
)

func main() {
	//defer repository.CloseDB()
	migrations.Migrate()
	listeners.StepResponseListener.Listen()
	handlers.LoadHTTPRoutes()
	log.Println("Calling listener")
}
