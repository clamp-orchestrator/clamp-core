package main

import (
	"clamp-core/handler"
	"clamp-core/migrations"
)

func main() {
	migrations.Migrate()
	handler.LoadHTTPRoutes()
}
