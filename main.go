package main

import (
	"clamp-core/config"
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
	"clamp-core/models"
	"clamp-core/repository"

	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Pinging DB...")
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
		listeners.AMQPStepResponseListener.Listen()
	}
	if config.ENV.EnableKafkaIntegration {
		listeners.KafkaStepResponseListener.Listen()
	}
	handlers.LoadHTTPRoutes()
	log.Info("Calling listener")
}
