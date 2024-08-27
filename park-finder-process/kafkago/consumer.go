package kafkago

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type Consumer struct {
	PoolSize        chan int
	ready           chan bool
	MethodsCallback func(*sarama.ConsumerMessage)
	ConsumerGroup   sarama.ConsumerGroup
}

func NewConsumer(kafkaBroker []string, topic string, assignor string, groupID string, methodsCallback func(*sarama.ConsumerMessage), poolSize int) {
	log.Println("Starting a new Sarama consumer")
	config := newConfig()

	oldest := true

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:

		log.Printf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		log.Println("Consumer OffsetOldest.")
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		log.Println("Consumer OffsetNewest.")
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	consumer := Consumer{
		PoolSize:        make(chan int, poolSize),
		ready:           make(chan bool),
		MethodsCallback: methodsCallback,
	}

	ctx, cancel := context.WithCancel(context.Background())

	client, err := sarama.NewConsumerGroup(kafkaBroker, groupID, config)
	if err != nil {
		log.Printf("Error creating consumer group client: %v", err)
		return
	}

	consumer.ConsumerGroup = client

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			if err := client.Consume(ctx, strings.Split(topic, ","), &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up

	log.Printf("Consumer is running with group: %s | on broker: %s", groupID, kafkaBroker)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}

}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// fmt.Println("ConsumeClaim: 2m9h1nf7-default")
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for msg := range claim.Messages() {

		consumer.PoolSize <- 1 // will block if there is MAX ints in sem
		log.Printf("Message claimed: key = %s, value = %s, timestamp = %v, topic = %s", string(msg.Key), string(msg.Value), msg.Timestamp, msg.Topic)

		go func(msg *sarama.ConsumerMessage) {

			fmt.Printf("Partition: %d, Offset: %d\n", msg.Partition, msg.Offset)
			consumer.MethodsCallback(msg)
			session.MarkMessage(msg, "")

			<-consumer.PoolSize // removes an int from sem, allowing another to proceed

		}(msg)
	}

	return nil
}
