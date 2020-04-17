package main

import (
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
)

func main() {
	//defer repository.CloseDB()
	migrations.Migrate()
	handlers.LoadHTTPRoutes()
	listeners.StepResponseListener.Listen()
}
