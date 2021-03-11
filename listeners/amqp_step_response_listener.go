package listeners

import "clamp-core/config"

type AMQPStepResponseListenerInterface interface {
	Listen()
}

var AMQPStepResponseListener AMQPStepResponseListenerInterface

func init() {
	switch config.ENV.QueueDriver {
	case "amqp":
		AMQPStepResponseListener = &amqpListener{}
	}
}
