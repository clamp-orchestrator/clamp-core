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
	DBConnectionStr string `env:"CLAMP_DB_CONNECTION_STR" envDefault:"host=localhost:5432 user=clamp dbname=clampdev password=clamppass"`
	/**
	Examples:

	CLAMP_QUEUE_DRIVER     CLAMP_QUEUE_CONNECTION_STR
	=================     ===============================================================
	"amqp"            "amqp://guest:guest@localhost:5672/"

	*/
	QueueDriver        string `env:"CLAMP_QUEUE_DRIVER" envDefault:"amqp"`
	QueueConnectionStr string `env:"CLAMP_QUEUE_CONNECTION_STR" envDefault:"amqp://clamp:clampdev!@54.190.25.178:5672/"`
	QueueName          string `env:"CLAMP_QUEUE_NAME" envDefault:"clamp_steps_response"`
	/**
	Examples:

	CLAMP_QUEUE_DRIVER     CLAMP_QUEUE_CONNECTION_STR
	=================     ===============================================================
	"kafka"            "amqp://guest:guest@localhost:5672/"

	*/
	KafkaDriver            string `env:"CLAMP_KAFKA_DRIVER" envDefault:"kafka"`
	KafkaConnectionStr     string `env:"CLAMP_KAFKA_CONNECTION_STR" envDefault:"54.190.25.178:9092"`
	KafkaTopicName         string `env:"CLAMP_KAFKA_TOPIC_NAME" envDefault:"clamp_steps_response"`
	KafkaConsumerTopicName string `env:"CLAMP_KAFKA_TOPIC_NAME" envDefault:"clamp_consumer_topic"`
	/**
	System Defaults
	*/
	AsyncStepType string   `env:"ASYNC_STEP_TYPE" envDefault:"ASYNC"`
	SyncStepType  string   `env:"SYNC_STEP_TYPE" envDefault:"SYNC"`
	PORT          string   `env:"APP_PORT" envDefault:"8080"`
	AllowOrigins  []string `env:"ALLOW_ORIGINS" envDefault:"http://localhost:3000,http://http://54.149.76.62"`
}{}

func init() {
	err := env.Parse(&ENV)
	if err != nil {
		panic(err)
	}
}
