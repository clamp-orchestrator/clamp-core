package main

import (
	"clamp-core/handlers"
	"clamp-core/migrations"
)

func main() {
	//defer repository.CloseDB()
	migrations.Migrate()
	handlers.LoadHTTPRoutes()
}
