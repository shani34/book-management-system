package kafka

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
)
var Producer *kafka.Writer
func InitKafkaProducer() {
    TOPIC_NAME := "TOPIC_NAME"

    keypair, err := tls.LoadX509KeyPair("service.cert", "service.key")
    if err != nil {
        log.Fatalf("Failed to load access key and/or access certificate: %s", err)
    }

    caCert, err := ioutil.ReadFile("ca.pem")
    if err != nil {
        log.Fatalf("Failed to read CA certificate file: %s", err)
    }

    caCertPool := x509.NewCertPool()
    ok := caCertPool.AppendCertsFromPEM(caCert)
    if !ok {
        log.Fatalf("Failed to parse CA certificate file: %s", err)
    }

    dialer := &kafka.Dialer{
        Timeout:   10 * time.Second,
        DualStack: true,
        TLS: &tls.Config{
            Certificates: []tls.Certificate{keypair},
            RootCAs:      caCertPool,
        },
    }

    // init producer
    Producer = kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{"kafka-245868dd-shani-kafka-test-007.c.aivencloud.com:24869"},
        Topic:   TOPIC_NAME,
        Dialer:  dialer,
    })

    Producer.Close()
}

func PublishEvent(topic string,  message []byte)error{
	return Producer.WriteMessages(context.Background(), kafka.Message{Value: []byte(message)})
}