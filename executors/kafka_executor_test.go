package executors

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestKafkaVal_DoExecute(t *testing.T) {
	assert := assert.New(t)

	testKafkaConnectionURL := "localhost:61234"
	testTopicName := "topic_test"
	testMessageBody := make(map[string]interface{})
	testMessageBody["payload"] = "data"

	t.Run("ErrorOnConnectionFailure", func(t *testing.T) {
		val := KafkaVal{
			ConnectionURL: testKafkaConnectionURL,
			TopicName:     testTopicName,
			ContentType:   "text/plain",
		}
		_, err := val.DoExecute(testMessageBody, "")
		assert.Error(err)
	})

	mockSyncProducer := mocks.NewSyncProducer(t, &sarama.Config{})
	newSyncProducerFunc = func(_ []string, config *sarama.Config) (sarama.SyncProducer, error) {
		return mockSyncProducer, nil
	}

	t.Run("SuccessfulSendMessage", func(t *testing.T) {
		val := KafkaVal{
			ConnectionURL: testKafkaConnectionURL,
			TopicName:     testTopicName,
			ContentType:   "text/plain",
		}

		message := make(map[string]interface{})
		mockSyncProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(msg []byte) error {
			json.Unmarshal(msg, &message)
			return nil
		})

		_, err := val.DoExecute(testMessageBody, "")
		assert.NoError(err)
		assert.Equal(testMessageBody, message)
	})

	t.Run("SendMessageFailure", func(t *testing.T) {
		val := KafkaVal{
			ConnectionURL: testKafkaConnectionURL,
			TopicName:     testTopicName,
			ContentType:   "text/plain",
		}

		mockSyncProducer.ExpectSendMessageAndFail(errors.New("internal error"))

		_, err := val.DoExecute(testMessageBody, "")
		assert.Error(err)
	})
}
