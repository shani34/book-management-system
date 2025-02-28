package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"os"
)

var Producer sarama.SyncProducer

func InitKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	var err error
	Producer, err = sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKER")}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
}

func PublishEvent(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := Producer.SendMessage(msg)
	return err
}