package kafka

import (
	"github.com/ByPikod/go-crypto/tree/crypto/core"
	Kafka "github.com/segmentio/kafka-go"
)

func NewProducer(topic string, partition int) *Kafka.Writer {
	producer := Kafka.NewWriter(Kafka.WriterConfig{
		Brokers:  []string{core.Config.Kafka.Host},
		Topic:    topic,
		Balancer: &Kafka.LeastBytes{},
	})

	return producer
}
