package config

import (
	"github.com/joho/godotenv"
)

func ConfigFromEnvVars() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	port, err := readAndCheckIntEnvVar("EVENT_PROCESSING_WEB_PORT")
	if err != nil {
		return Config{}, err
	}

	dbURL, err := readAndCheckEnvVar("EVENT_PROCESSING_DB_URL")
	if err != nil {
		return Config{}, err
	}

	kafkaURL, err := readAndCheckEnvVar("EVENT_PROCESSING_KAFKA_BROKER")
	if err != nil {
		return Config{}, err
	}

	eventCreatedTopic, err := readAndCheckEnvVar("EVENT_PROCESSING_KAFKA_TOPIC_EVENT_CREATED")
	if err != nil {
		return Config{}, err
	}

	apiLayerAPIKey, err := readAndCheckEnvVar("EVENT_PROCESSING_APILAYER_API_KEY")
	if err != nil {
		return Config{}, err
	}

	freeCurrencyAPIKey, err := readAndCheckEnvVar("EVENT_PROCESSING_FREECURRENCY_API_KEY")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:               port,
		DBURL:              dbURL,
		KafkaURL:           kafkaURL,
		EventCreatedTopic:  eventCreatedTopic,
		ApiLayerAPIKey:     apiLayerAPIKey,
		FreeCurrencyAPIKey: freeCurrencyAPIKey,
	}, nil
}
