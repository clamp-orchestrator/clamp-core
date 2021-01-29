package main

import (
	"clamp-core/config"
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
	"clamp-core/models"
	"clamp-core/repository"

	"log"
	"os"
)

func main() {
	log.Println("Pinging DB...")
	err := repository.GetDB().Ping()
	if err != nil {
		log.Fatalf("DB ping failed: %s", err)
	}

	var cliArgs models.CLIArguments = os.Args[1:]
	os.Setenv("PORT", config.ENV.PORT)
	migrations.Migrate()

	if cliArgs.Parse().Find("migrate-only", "no") == "yes" {
		os.Exit(0)
	}

	if config.ENV.EnableRabbitMQIntegration {
		listeners.AmqpStepResponseListener.Listen()
	}
	if config.ENV.EnableKafkaIntegration {
		listeners.KafkaStepResponseListener.Listen()
	}
	handlers.LoadHTTPRoutes()
	log.Println("Calling listener")
}
