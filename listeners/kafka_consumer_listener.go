package listeners

import (
	"clamp-core/config"
	"clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
}

func (c *Consumer) Listen() {
	go func() {
		saramConfig := sarama.NewConfig()
		saramConfig.ClientID = "go-kafka-consumer"
		saramConfig.Consumer.Return.Errors = true

		brokers := config.ENV.KafkaConnectionStr

		// Create new consumer
		master, err := sarama.NewConsumer(strings.Split(brokers, ","), saramConfig)
		if err != nil {
			panic(err)
		}

		defer func() {
			if err = master.Close(); err != nil {
				panic(err)
			}
		}()

		//topics, _ := master.Topics()
		topics := config.ENV.KafkaConsumerTopicName
		consumer, errors := consume(topics, master)
		// Count how many message processed
		msgCount := 0

		// Get signnal for finish
		doneCh := make(chan struct{})
		go func() {
			for {
				select {
				case msg := <-consumer:
					msgCount++
					fmt.Println("Received messages", string("msg.Key"), string(msg.Value))
					var res models.AsyncStepResponse
					err = json.Unmarshal(msg.Value, &res)
					if err != nil {
						log.Errorf("[Kafka Consumer] : Message received is not in proper format %s: %s", msg.Value, err.Error())
					} else {
						err := binding.Validator.ValidateStruct(res)
						if err != nil {
							log.Errorf("[Kafka Consumer] : Message received is not in proper format %s: %s", msg.Value, err.Error())
						}
						log.Debugf("[Kafka Consumer] : Received step completed response: %v", res)
						log.Debug("[Kafka Consumer] : Pushing step completed response to channel")
						services.AddStepResponseToResumeChannel(&res)
					}
				case consumerError := <-errors:
					msgCount++
					fmt.Println("Received consumerError ", consumerError.Topic, string(consumerError.Partition), consumerError.Err)
					doneCh <- struct{}{}
				}
			}
		}()
		<-doneCh
	}()
}

func consume(topic string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)

	partitions, err := master.Partitions(topic)
	if err != nil {
		log.Errorf("error while retrieving partitions: %s", err)
	}

	for _, partition := range partitions {
		consumer, err := master.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("consuming topic %v partitions %v failed: %s", topic, partition, err)
		}

		fmt.Println(" Starting Kafka consumer topic ", topic)
		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					errors <- consumerError
					fmt.Println("consumerError: ", consumerError.Err)

				case msg := <-consumer.Messages():
					consumers <- msg
					fmt.Println(""+
						"Got message on topic ", topic, msg.Value)
				}
			}
		}(topic, consumer)
	}

	return consumers, errors
}
