package executors

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	amqptestserver "github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/stretchr/testify/assert"
)

const (
	testConnectionURL = "amqp://localhost:5672/test"
	testQueueName     = "test_queue"
	testExchangeName  = "test_exchange"
	testRoutingKey    = "test_routing_key"
	testContentType   = "json/application"
)

var (
	testMessageBody = map[string]interface{}{"somefield": "somevlaue"}
)

func amqpTestPublish(amqpURL, exchange, key string, contentType string, body []byte) error {
	conn, err := amqptest.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("error while dialing amqp connection: %w", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error while creating amqp channel: %w", err)
	}

	defer ch.Close()

	err = ch.Publish(exchange, key, body, wabbit.Option{"contentType": contentType})
	if err != nil {
		return fmt.Errorf("error while publishing amqp message: %w", err)
	}

	return nil
}

func popDelivery(deliveryCh <-chan wabbit.Delivery) wabbit.Delivery {
	select {
	case d := <-deliveryCh:
		return d
	default:
		return nil
	}
}

func drainDeliveries(deliveryCh <-chan wabbit.Delivery) {
	for {
		select {
		case <-deliveryCh:
		default:
			return
		}
	}
}

func TestAMQPVal_DoExecute(t *testing.T) {
	amqpPublishFunc = amqpTestPublish

	t.Run("should return error if connection fails", func(t *testing.T) {
		val := AMQPVal{
			ConnectionURL: testConnectionURL,
			QueueName:     testQueueName,
			ExchangeName:  "",
			RoutingKey:    "",
			ContentType:   testContentType,
		}
		_, err := val.DoExecute(testMessageBody, "")
		assert.Error(t, err)
	})

	amqpTestServer := amqptestserver.NewServer(testConnectionURL)
	amqpTestServer.Start()
	defer amqpTestServer.Stop()

	amqpConn, err := amqptest.Dial(testConnectionURL)
	if err != nil {
		t.Errorf("amqptest dial failed: %s", err)
		return
	}

	defer amqpConn.Close()

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		t.Errorf("amqptest channel creation failed: %s", err)
		return
	}

	defer amqpChannel.Close()

	err = amqpChannel.ExchangeDeclare(testExchangeName, "topic", wabbit.Option{})
	if err != nil {
		t.Errorf("amqptest exchange declare failed: %s", err)
		return
	}

	_, err = amqpChannel.QueueDeclare(testQueueName, wabbit.Option{})
	if err != nil {
		t.Errorf("amqptest queue declare failed: %s", err)
		return
	}

	err = amqpChannel.QueueBind(testQueueName, testRoutingKey, testExchangeName, wabbit.Option{})
	if err != nil {
		t.Errorf("amqptest queue bind failed: %s", err)
		return
	}

	deliveryCh, err := amqpChannel.Consume(testQueueName, "", wabbit.Option{})
	if err != nil {
		t.Errorf("amqptest consume failed: %s", err)
		return
	}

	t.Run("should send message to queue", func(t *testing.T) {
		val := AMQPVal{
			ConnectionURL: testConnectionURL,
			QueueName:     testQueueName,
			ExchangeName:  "",
			RoutingKey:    "",
			ContentType:   testContentType,
		}
		_, err := val.DoExecute(testMessageBody, "")
		assert.NoError(t, err)

		time.Sleep(time.Millisecond) // gives time for the delivery

		d := popDelivery(deliveryCh)
		if assert.NotNil(t, d) {
			var messageBody interface{}
			err = json.Unmarshal(d.Body(), &messageBody)
			assert.NoError(t, err)
			assert.Equal(t, testMessageBody, messageBody)
		}
	})

	drainDeliveries(deliveryCh)

	t.Run("should send message to exchange", func(t *testing.T) {
		val := AMQPVal{
			ConnectionURL: testConnectionURL,
			QueueName:     "",
			ExchangeName:  testExchangeName,
			RoutingKey:    testRoutingKey,
			ContentType:   testContentType,
		}
		_, err := val.DoExecute(testMessageBody, "")
		assert.NoError(t, err)

		time.Sleep(time.Millisecond) // gives time for the delivery

		d := popDelivery(deliveryCh)
		if assert.NotNil(t, d) {
			var messageBody interface{}
			err = json.Unmarshal(d.Body(), &messageBody)
			assert.NoError(t, err)
			assert.Equal(t, testMessageBody, messageBody)
		}
	})

	drainDeliveries(deliveryCh)
}
