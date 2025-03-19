package generator

import (
	"log"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/database"
	"github.com/TomasCruz/event-processing/internal/kafkaqueue"
)

// infra layer, instantiating deps and running presenter
func Run() {
	// populate configuration
	c, err := config.ConfigFromEnvVars()
	if err != nil {
		log.Fatal("failed to read environment variables", err)
	}

	// init DB
	db, err := database.New(c)
	if err != nil {
		log.Fatal(err, "failed to initialize database")
	}

	// Kafka consumer
	kp, err := kafkaqueue.InitProducer(c)
	if err != nil {
		log.Fatal(err, "failed to create Kafka producer")
	}

	// set Kafka topic
	topic := c.EventCreatedTopic
	err = kp.SetTopic(topic)
	if err != nil {
		log.Fatal(err, "failed to set Kafka producer topic")
	}

	// run
	genSvc := generatorPresenter{
		topic: topic,
		db:    db,
		p:     kp,
	}
	genSvc.run()
}
