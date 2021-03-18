package executors

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// KafkaVal : Kafka configurations details
type KafkaVal struct {
	ConnectionURL string `json:"connection_url" binding:"required"`
	TopicName     string `json:"topic_name"`
	ContentType   string `json:"content_type"`
	ReplyTo       string `json:"reply_to"`
}

// DoExecute : Connecting to Kakfa URL and producing a message to Topic
func (val *KafkaVal) DoExecute(requestBody interface{}, prefix string) (interface{}, error) {
	log.Debugf("%s Kafka Executor: Executing kafka %s body:%v", prefix, val.TopicName, requestBody)
	syncProducer, err := sarama.NewSyncProducer([]string{val.ConnectionURL}, nil)
	//asyncProducer, err := sarama.NewAsyncProducer([]string{val.ConnectionURL}, nil)
	if err != nil {
		log.Errorf("%s Kafka Error: %s", prefix, err.Error())
		return nil, err
	}

	requestJSONBytes, _ := json.Marshal(requestBody)
	msg := &sarama.ProducerMessage{
		Topic: val.TopicName,
		Value: sarama.StringEncoder(requestJSONBytes),
	}

	_, _, err = syncProducer.SendMessage(msg)
	if err != nil {
		panic(err)
	}
	log.Debugf("%s Kafka Executor: pushed message successfully", prefix)

	//asyncProducer.Input() <- &sarama.ProducerMessage{
	//	Topic: config.ENV.KafkaTopicName,
	//	Value: sarama.StringEncoder(requestJSONBytes),
	//}
	//log.Debugf("%s Kafka Executor: pushed async message successfully", prefix)
	return nil, nil
}
