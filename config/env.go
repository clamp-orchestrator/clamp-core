package config

import "github.com/caarlos0/env"

var ENV = struct {

	/**
	Examples:

	CLAMP_DB_DBDRIVER     CLAMP_DB_DBCONNECTIONSTR
	=================     ===============================================================
	"postgres"            "host=localhost user=root dbname=clamp password=mypassword"

	*/
	DBDriver        string `env:"CLAMP_DB_DRIVER" envDefault:"postgres"`
	DBConnectionStr string `env:"CLAMP_DB_CONNECTION_STR" envDefault:"host=myhost user=root dbname=clamp password=mypassword"`
	/**
	Examples:

	CLAMP_QUEUE_DRIVER     CLAMP_QUEUE_CONNECTION_STR
	=================     ===============================================================
	"amqp"            "amqp://guest:guest@localhost:5672/"

	*/
	QueueDriver        string `env:"CLAMP_QUEUE_DRIVER" envDefault:"amqp"`
	QueueConnectionStr string `env:"CLAMP_QUEUE_CONNECTION_STR" envDefault:"amqp://guest:guest@localhost:5672/"`
	QueueName          string `env:"CLAMP_QUEUE_NAME" envDefault:"clamp_steps_response"`
}{}

func init() {
	err := env.Parse(&ENV)
	if err != nil {
		panic(err)
	}
}
