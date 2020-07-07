package listeners

import (
	"clamp-core/config"
	"clamp-core/models"
	"clamp-core/services"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin/binding"
	"log"
	"os"
	"os/signal"
	"strings"
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
			if err := master.Close(); err != nil {
				panic(err)
			}
		}()

		//topics, _ := master.Topics()
		topics := config.ENV.KafkaConsumerTopicName
		consumer, errors := consume(topics, master)

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

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
						log.Printf("[AMQP Consumer] : Message recieved is not in proper format %s: %s", msg.Value, err.Error())
					} else {
						err := binding.Validator.ValidateStruct(res)
						if err != nil {
							log.Printf("[AMQP Consumer] : Message recieved is not in proper format %s: %s", msg.Value, err.Error())
						}
						log.Printf("[AMQP Consumer] : Received step completed response: %v", res)
						log.Printf("[AMQP Consumer] : Pushing step completed response to channel")
						services.AddStepResponseToResumeChannel(res)
					}
				case consumerError := <-errors:
					msgCount++
					fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
					doneCh <- struct{}{}
				case <-signals:
					fmt.Println("Interrupt is detected")
					doneCh <- struct{}{}
				}
			}
			fmt.Println("Processed", msgCount, "messages")
		}()
		<-doneCh
	}()
}

func consume(topic string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)

	partitions, _ := master.Partitions(topic)
	for _, partition := range partitions {
		consumer, err := master.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if nil != err {
			fmt.Printf("Topic %v Partitions: %v", topic, partition)
			panic(err)
		}
		fmt.Println(" Start consuming topic ", topic)
		go func(topic string, consumer sarama.PartitionConsumer) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					errors <- consumerError
					fmt.Println("consumerError: ", consumerError.Err)

				case msg := <-consumer.Messages():
					consumers <- msg
					fmt.Println("" +
						"Got message on topic ", topic, msg.Value)
				}
			}
		}(topic, consumer)
	}

	return consumers, errors
}