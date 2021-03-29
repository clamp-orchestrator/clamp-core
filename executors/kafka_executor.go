package executors

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

var newSyncProducerFunc func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) = sarama.NewSyncProducer

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

	syncProducer, err := newSyncProducerFunc([]string{val.ConnectionURL}, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating kafka sync producer: %w", err)
	}

	requestJSONBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling kafka executor request: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: val.TopicName,
		Value: sarama.StringEncoder(requestJSONBytes),
	}

	_, _, err = syncProducer.SendMessage(msg)
	if err != nil {
		return nil, fmt.Errorf("error while sending kafka message: %w", err)
	}

	log.Debugf("%s Kafka Executor: pushed message successfully", prefix)
	return nil, nil
}
