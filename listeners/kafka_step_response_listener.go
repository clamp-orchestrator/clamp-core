package listeners

import "clamp-core/config"

type KafkaStepResponseListenerInterface interface {
	Listen()
}

var KafkaStepResponseListener KafkaStepResponseListenerInterface

func init() {
	switch config.ENV.KafkaDriver {
	case "kafka":
		KafkaStepResponseListener = &Consumer{}
	}
}
