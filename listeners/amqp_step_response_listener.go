package listeners

import "clamp-core/config"

type AmqpStepResponseListenerInterface interface {
	Listen()
}

var AmqpStepResponseListener AmqpStepResponseListenerInterface

func init() {
	switch config.ENV.QueueDriver {
	case "amqp":
		AmqpStepResponseListener = &amqpListener{}
	}
}
