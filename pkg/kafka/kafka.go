package kafka

import (
	"github.com/IBM/sarama"
	"github.com/shani34/book-management-system/config"
	"fmt"
)

var Producer sarama.SyncProducer

func InitKafkaProducer() error {
	cfg := config.Get().Kafka
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll

	var err error
	Producer, err = sarama.NewSyncProducer(cfg.Brokers, kafkaConfig)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}
	return nil
}

func PublishEvent(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := Producer.SendMessage(msg)
	return err
}