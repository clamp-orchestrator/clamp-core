package main

import (
	"clamp-core/config"
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
	"clamp-core/models"
	"log"
	"os"
)

func main() {
	var cliArgs models.CLIArguments = os.Args[1:]
	os.Setenv("PORT", config.ENV.PORT)
	migrations.Migrate()

	if cliArgs.Parse().Find("migrate-only", "no") == "yes" {
		os.Exit(0)
	}

	listeners.AmqpStepResponseListener.Listen()
	listeners.KafkaStepResponseListener.Listen()
	handlers.LoadHTTPRoutes()
	log.Println("Calling listener")
}
