package kafka

import (
	"context"
	"golang-api/repository"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type Consumer struct {
	ready chan bool
}

var wg *sync.WaitGroup

// Consume funtion to listen
func Consume() {
	KafkaServer := "localhost:9092"
	KafkaTopic := "messageTopic"
	config := sarama.NewConfig()
	config.Version = sarama.V2_4_0_0
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup([]string{KafkaServer}, "notification", config)
	if err != nil {
		panic(err)
	}
	consumer := Consumer{
		ready: make(chan bool),
	}
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(KafkaTopic, ","), &consumer); err != nil {
				log.Panic().Caller().Err(err).Msg("Error from consumer")
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Info().Msg("Sarama consumer up and running!...")
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Error().Caller().Msg("terminating: context cancelled")
	case <-sigterm:
		log.Error().Caller().Msg("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panic().Caller().Err(err).Msg("Error closing client")
	}
	<-consumer.ready
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info().Msg("Kafka notification")
		repository.UpdateCache()
		sess.MarkMessage(msg, "")
	}
	return nil
}
