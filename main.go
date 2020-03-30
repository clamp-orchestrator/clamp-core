package main

import (
	"clamp-core/handlers"
	"clamp-core/migrations"
	"clamp-core/repository"
)

func main() {
	defer repository.CloseDB()
	migrations.Migrate()
	handlers.LoadHTTPRoutes()
}
