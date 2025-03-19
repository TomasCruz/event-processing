package kafkaqueue

import (
	"log"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/ports"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaProducer struct {
	config config.Config
	kp     *kafka.Producer
	topic  string
}

func InitProducer(c config.Config) (ports.AsyncMsgProducer, error) {
	kp, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": c.KafkaURL,
	})
	if err != nil {
		log.Fatal(err, "failed to create Kafka producer")
	}

	return &kafkaProducer{
		kp:     kp,
		config: c,
	}, nil
}

func (k *kafkaProducer) Close() error {
	return nil
}

func (k *kafkaProducer) SetTopic(topic string) error {
	k.topic = topic
	return nil
}

func (k *kafkaProducer) SendAsyncMsg(msg []byte) error {
	if err := k.kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil); err != nil {
		log.Printf("failed to produce %s event: %v\n", k.topic, err)
		return err
	}

	return nil
}
