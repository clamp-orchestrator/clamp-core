package config

import "github.com/caarlos0/env"

// ENV is a config that is loaded at the application start up. The values here can be overridden
// by setting environment variables before the process starts up. An example of overriding
// the value would be `PORT` has a default value of `env:"APP_PORT" envDefault:"8080"`.
// If environment variable APP_PORT is set to 9000 then the value 9000 will be used as the port for the application
var ENV = struct {

	/**
	Examples:

	CLAMP_DB_DBDRIVER     CLAMP_DB_DBCONNECTIONSTR
	=================     ===============================================================
	"postgres"            "host=localhost user=root dbname=clamp password=mypassword"

	*/
	DBDriver        string `env:"CLAMP_DB_DRIVER" envDefault:"postgres"`
	DBConnectionStr string `env:"CLAMP_DB_CONNECTION_STR" envDefault:"host=localhost:5432 user=postgres dbname=clampdev password=mysecretpassword"`
	/**
	Examples:

	CLAMP_QUEUE_DRIVER     CLAMP_QUEUE_CONNECTION_STR
	=================     ===============================================================
	"amqp"            "amqp://guest:guest@localhost:5672/"

	*/
	QueueDriver        string `env:"CLAMP_QUEUE_DRIVER" envDefault:"amqp"`
	QueueConnectionStr string `env:"CLAMP_QUEUE_CONNECTION_STR" envDefault:"amqp://clamp:clampdev!@localhost:5672/"`
	QueueName          string `env:"CLAMP_QUEUE_NAME" envDefault:"clamp_steps_response"`
	/**
	Examples:

	CLAMP_QUEUE_DRIVER     CLAMP_QUEUE_CONNECTION_STR
	=================     ===============================================================
	"kafka"            "amqp://guest:guest@localhost:5672/"

	*/
	KafkaDriver            string `env:"CLAMP_KAFKA_DRIVER" envDefault:"kafka"`
	KafkaConnectionStr     string `env:"CLAMP_KAFKA_CONNECTION_STR" envDefault:"localhost:9092"`
	KafkaTopicName         string `env:"CLAMP_KAFKA_TOPIC_NAME" envDefault:"clamp_steps_response"`
	KafkaConsumerTopicName string `env:"CLAMP_KAFKA_TOPIC_NAME" envDefault:"clamp_consumer_topic"`
	/**
	System Defaults
	*/
	PORT         string   `env:"APP_PORT" envDefault:"8080"`
	AllowOrigins []string `env:"ALLOW_ORIGINS" envDefault:"http://localhost:3000"`

	EnableKafkaIntegration    bool `env:"ENABLE_KAFKA_INTEGRATION" envDefault:"false"`
	EnableRabbitMQIntegration bool `env:"ENABLE_AMQP_INTEGRATION" envDefault:"false"`

	EnableSQLQueriesLog bool `env:"ENABLE_SQL_QUERIES_LOG" envDefault:"false"`
}{}

func init() {
	err := env.Parse(&ENV)
	if err != nil {
		panic(err)
	}
}
