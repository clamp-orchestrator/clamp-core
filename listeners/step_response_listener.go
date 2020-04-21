package listeners

import "clamp-core/config"

type StepResponseListenerInterface interface {
	Listen()
}

var StepResponseListener StepResponseListenerInterface

func init() {
	switch config.ENV.DBDriver {
	case "postgres":
		StepResponseListener = &amqpListener{}
	}
}
