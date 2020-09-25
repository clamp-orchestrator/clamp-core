package main

import (
	"clamp-core/config"
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
	"log"
	"os"
)

func main() {
	//defer repository.CloseDB()
	os.Setenv("PORT", config.ENV.PORT)
	migrations.Migrate()
	listeners.AmqpStepResponseListener.Listen()
	listeners.KafkaStepResponseListener.Listen()
	handlers.LoadHTTPRoutes()
	log.Println("Calling listener")
}
