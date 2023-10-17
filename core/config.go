package core

import (
	"os"

	"github.com/ByPikod/go-crypto/tree/notifier/helpers"
	"github.com/joho/godotenv"
)

type (
	Configuration struct {
		Loki  *LokiInfo
		Kafka *KafkaInfo
	}
	// LokiInfo struct contains authentication information for the logging database.
	LokiInfo struct {
		Host string
	}
	KafkaInfo struct {
		Host string
	}
)

// It will be nil if config haven't been initialized.
var Config *Configuration

func init() {
	Config = InitializeConfig()
}

func or(x string, y string) string {
	if x == "" {
		return y
	}
	return x
}

// Initializes config and makes Config variable above ready to use by loading environment variables.
// ".env" is supported.
func InitializeConfig() *Configuration {

	err := godotenv.Load()
	if err != nil {
		helpers.LogError(`File ".env" not found or cannot parsed: ` + err.Error())
	}

	loki := &LokiInfo{
		Host: or(os.Getenv("LOKI_HOST"), "http://loki:3100"),
	}
	kafka := &KafkaInfo{
		Host: or(os.Getenv("KAFKA_HOST"), "kafka:9092"),
	}

	config := Configuration{
		Loki:  loki,
		Kafka: kafka,
	}

	return &config

}
