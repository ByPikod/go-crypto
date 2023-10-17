package kafka

import (
	"context"
	"time"

	"github.com/ByPikod/go-crypto/tree/notifier/core"
	"github.com/ByPikod/go-crypto/tree/notifier/helpers"
	"github.com/ByPikod/go-crypto/tree/notifier/log"
	Kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type (
	Consumer struct {
		conn *Kafka.Reader
	}
)

func CreateConsumer() *Consumer {
	conn := Kafka.NewReader(Kafka.ReaderConfig{
		Brokers:        []string{core.Config.Kafka.Host},
		Topic:          "notifier",
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	})

	consumer := &Consumer{
		conn: conn,
	}

	return consumer
}

func (consumer *Consumer) ReadQueue() {
	for {
		m, err := consumer.conn.ReadMessage(context.Background())
		if err != nil {
			helpers.LogError(err.Error())
			break
		}
		log.Info(
			"Mail sent",
			zap.String("json", string(m.Value)),
		)
	}

	if err := consumer.conn.Close(); err != nil {
		helpers.LogError("failed to close reader: " + err.Error())
	}
}
