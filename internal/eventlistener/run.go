package eventlistener

import (
	"log"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/database"
	"github.com/TomasCruz/event-processing/internal/kafkaqueue"
)

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
	kc, err := kafkaqueue.InitConsumer(c)
	if err != nil {
		log.Fatal(err, "failed to create Kafka consumer")
	}

	// subscribe to topic
	err = kc.SubscribeTopic(c.EventCreatedTopic)
	if err != nil {
		log.Fatal(err, "failed to subscribe to topic")
	}

	// run
	pData := enricherSvc{
		config:   c,
		db:       db,
		consumer: kc,
	}
	pData.run()
}
