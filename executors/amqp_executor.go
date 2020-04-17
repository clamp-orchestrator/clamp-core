package executors

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
)

type AMQPVal struct {
	ConnectionURL string     `json:"connection_url" binding:"required"`
	QueueName     string     `json:"queue_name"`
	ExchangeName  string     `json:"exchange_name"`
	RoutingKey    string     `json:"routing_key"`
	ExchangeType  string     `json:"exchange_type"`
	Durable       bool       `json:"durable"`
	AutoDelete    bool       `json:"auto_delete"`
	Internal      bool       `json:"internal"`
	NoWait        bool       `json:"no_wait"`
	Exclusive     bool       `json:"exclusive"`
	Mandatory     bool       `json:"mandatory"`
	Immediate     bool       `json:"immediate"`
	ContentType   string     `json:"content_type"`
	Arguments     amqp.Table `json:"arguments"`
}

/*
{
"connection_url":"amqp://guest:guest@localhost:5672/",
"queue_name":"clamp",
"durable": true
"auto_delete":false,
"internal":false,
"no_wait":false,
"arguments":nil,
content_type":"text/plain",
"mandatory":false,
"immediate":false
}
*/

func (val AMQPVal) DoExecute(requestBody interface{}) (interface{}, error) {
	prefix := log.Prefix()
	log.SetPrefix("")
	log.Printf("%s AMQP Executor: Executing amqp %s body:%v", prefix, val.getName(), requestBody)

	conn, err := amqp.Dial(val.ConnectionURL)
	if err != nil {
		log.Printf("%s AMQP Error: %s", prefix, err.Error())
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%s AMQP Error: %s", prefix, err.Error())
		return nil, err
	}
	defer ch.Close()

	if val.ExchangeName != "" {
		return sendMessageToExchange(ch, val, requestBody, prefix)
	} else if val.QueueName != "" {
		return sendMessageToQueue(ch, val, requestBody, prefix)
	} else {
		return nil, errors.New("AMQP - queue/exchange name not specified")
	}
}

func sendMessageToQueue(ch *amqp.Channel, val AMQPVal, body interface{}, prefix string) (interface{}, error) {
	q, err := ch.QueueDeclare(
		val.QueueName,
		val.Durable,
		val.AutoDelete,
		val.Exclusive,
		val.NoWait,
		val.Arguments,
	)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: val.ContentType,
			Body:        bytes,
		})
	if err != nil {
		return nil, err
	} else {
		log.Printf("%s AMQP Executor: pushed message successfully", prefix)
	}
	return `{"msg":"AMQP-pushed message successfully"}`, nil
}

func sendMessageToExchange(ch *amqp.Channel, val AMQPVal, body interface{}, prefix string) (interface{}, error) {
	err := ch.ExchangeDeclare(
		val.ExchangeName,
		val.ExchangeType,
		val.Durable,
		val.AutoDelete,
		val.Internal,
		val.NoWait,
		val.Arguments,
	)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(body)
	err = ch.Publish(
		val.ExchangeName,
		val.RoutingKey,
		val.Mandatory,
		val.Immediate,
		amqp.Publishing{
			ContentType: val.ContentType,
			Body:        bytes,
		})
	if err != nil {
		return nil, err
	} else {
		log.Printf("%s AMQP Executor: pushed message successfully", prefix)
	}
	return `{"msg":"AMQP-pushed message successfully"}`, nil
}

func (val AMQPVal) getName() string {
	if val.QueueName != "" {
		return val.QueueName
	}
	return val.ExchangeName
}
