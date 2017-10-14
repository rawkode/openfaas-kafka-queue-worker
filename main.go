package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	var kafkaHosts string
	if value, exists := os.LookupEnv("FAAS_KAFKA_HOSTS"); exists {
		kafkaHosts = value
	}

	var consumerGroup string
	if value, exists := os.LookupEnv("FAAS_KAFKA_CONSUMER_GROUP"); exists {
		consumerGroup = value
	}

	var topics []string
	if value, exists := os.LookupEnv("FAAS_KAFKA_TOPICS"); exists {
		topics = strings.Split(value, ",")
	}

	if kafkaHosts == "" || consumerGroup == "" || len(topics) == 0 {
		log.Println("Please ensure that you've got all the required environment variables")
		log.Println("\tFAAS_KAFKA_HOSTS")
		log.Println("\tFAAS_KAFKA_CONSUMER_GROUP")
		log.Println("\tFAAS_KAFKA_TOPICS")
		log.Fatalln("Exiting")
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               kafkaHosts,
		"group.id":                        consumerGroup,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s\n", err)
	}

	log.Printf("Created consumer: %v\n", consumer)

	err = consumer.SubscribeTopics(topics, nil)

Consumer:
	for {
		select {
		case signal := <-sigchan:
			fmt.Printf("Heh! We caught a signal (%v). Terminating", signal)
			break Consumer

		case event := <-consumer.Events():
			switch error := event.(type) {
			case *kafka.Message:
				fmt.Print("Received a message from Kafka")

			case kafka.Error:
				log.Printf("Kafka Consumer Error: %v\n", error)
				break Consumer
			}
		}
	}

	log.Println("Closing consumer and shutting down")
	consumer.Close()
}
