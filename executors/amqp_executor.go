package executors

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// AMQPVal : Rabbitmq configuration details
type AMQPVal struct {
	ConnectionURL string `json:"connection_url" binding:"required"`
	QueueName     string `json:"queue_name"`
	ExchangeName  string `json:"exchange_name"`
	RoutingKey    string `json:"routing_key"`
	ContentType   string `json:"content_type"`
	ReplyTo       string `json:"reply_to"`
}

var amqpPublishFunc = func(amqpURL, exchange, key string, contentType string, body []byte) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("error while dialing amqp connection: %w", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error while creating amqp channel: %w", err)
	}

	defer ch.Close()

	err = ch.Publish(exchange, key, false, false, amqp.Publishing{ContentType: contentType, Body: body})
	if err != nil {
		return fmt.Errorf("error while publishing amqp message: %w", err)
	}

	return nil
}

// DoExecute : Connection to Rabbitmq and sending message into Exchange
func (val *AMQPVal) DoExecute(requestBody interface{}, prefix string) (interface{}, error) {
	log.Debugf("%s AMQP Executor: Executing amqp %s body:%v", prefix, val.getName(), requestBody)

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling amqp request body: %w", err)
	}

	if val.ExchangeName != "" {
		err = amqpPublishFunc(val.ConnectionURL, val.ExchangeName, val.RoutingKey, val.ContentType, body)
	} else if val.QueueName != "" {
		err = amqpPublishFunc(val.ConnectionURL, "", val.QueueName, val.ContentType, body)
	} else {
		err = errors.New("AMQP - queue/exchange name not specified")
	}

	if err == nil {
		log.Debugf("%s AMQP Executor: pushed message successfully", prefix)
	}

	return nil, err
}

func (val *AMQPVal) getName() string {
	if val.QueueName != "" {
		return val.QueueName
	}
	return val.ExchangeName
}
