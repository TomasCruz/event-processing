package config

type Config struct {
	Port               string
	DBURL              string
	KafkaURL           string
	EventCreatedTopic  string
	ApiLayerAPIKey     string
	FreeCurrencyAPIKey string
}
