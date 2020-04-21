package listeners

import (
	"clamp-core/config"
	"clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/streadway/amqp"
	"log"
)

type amqpListener struct {
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("[AMQP Consumer] : %s: %s", msg, err)
	}
}

func (amqpListener amqpListener) Listen() {
	go func() {
		conn, err := amqp.Dial(config.ENV.QueueConnectionStr)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			config.ENV.QueueName,
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to declare a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"clamp",
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		go func() {
			for d := range msgs {
				var res models.AsyncStepResponse
				err = json.Unmarshal(d.Body, &res)
				if err != nil {
					log.Printf("[AMQP Consumer] : Message recieved is not in proper format %s: %s", d.Body, err.Error())
				} else {
					err := binding.Validator.ValidateStruct(res)
					if err != nil {
						log.Printf("[AMQP Consumer] : Message recieved is not in proper format %s: %s", d.Body, err.Error())
					}
					log.Printf("[AMQP Consumer] : Received step completed response: %v", res)
					log.Printf("[AMQP Consumer] : Pushing step completed response to channel")
					services.AddStepResponseToResumeChannel(res)
				}
			}
		}()

		log.Printf("[AMQP Consumer] : Started listening to queue")
		<-forever
	}()
}
