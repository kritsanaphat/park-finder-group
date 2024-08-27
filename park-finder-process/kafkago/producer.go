package kafkago

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type Producer struct {
	DataCollector sarama.SyncProducer
}

type IProducer interface {
	ProduceMessage(topic string, keyName []byte, value []byte) error
}

func NewProducerProvider() IProducer {
	fmt.Println("Starting create a producer connection....")
	config := newConfig()
	kafkaBroker := strings.Split(os.Getenv("CLOUDKARAFKA_BROKERS"), ",")

	log.Println("Starting create producer on : ", os.Getenv("CLOUDKARAFKA_BROKERS"))

	producer, err := sarama.NewSyncProducer(kafkaBroker, config)
	if err != nil {
		log.Fatal("Failed to start Sarama producer: ", err)
	}
	return &Producer{DataCollector: producer}
}

func (p *Producer) ProduceMessage(topic string, keyName []byte, value []byte) error {

	prefix := os.Getenv("CLOUDKARAFKA_TOPIC_PREFIX")

	topic = prefix + topic

	log.Println("Topic: ", topic)

	partition, offset, err := p.DataCollector.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(keyName),
		Value: sarama.ByteEncoder(value),
	})

	if err != nil {
		log.Fatal("Fail to producer message: ", err)
		return err
	} else {
		fmt.Printf("Your data is stored in partition: %d | offset: %d \n", partition, offset)
		return nil
	}
}
