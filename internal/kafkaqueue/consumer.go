package kafkaqueue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/ports"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaConsumer struct {
	kc *kafka.Consumer
}

func InitConsumer(config config.Config) (ports.AsyncMsgConsumer, error) {
	kc, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaURL,
		"group.id":          "group.id",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	consumer := kafkaConsumer{kc: kc}

	return &consumer, nil
}

func (k *kafkaConsumer) SubscribeTopic(topic string) error {
	// Subscribe to the Kafka topic
	err := k.kc.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	return nil
}

func (k *kafkaConsumer) Close() error {
	return nil
}

func (k *kafkaConsumer) Consume(ctx context.Context, msgChan chan<- []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := k.kc.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if err.(kafka.Error).Code() != kafka.ErrTimedOut {
				log.Printf("Consumer error: %v (%v)\n", err, msg)
			}
			continue
		}

		msgChan <- msg.Value
	}
}
