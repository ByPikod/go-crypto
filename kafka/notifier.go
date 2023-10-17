package kafka

import (
	"context"
	"encoding/json"

	Kafka "github.com/segmentio/kafka-go"
)

type Notifier struct {
	kafka *Kafka.Writer
}

var notifier *Notifier

func InitializeNotifier() {
	notifier = &Notifier{
		kafka: NewProducer("notifier", 0),
	}
}

func CreateNotifier() (*Notifier, error) {
	notifier := &Notifier{}

	return notifier, nil
}

// Returns initiated notifier instance
func GetNotifier() *Notifier {
	return notifier
}

// Send mail
func (notifier *Notifier) SendMail(
	receiver string,
	subject string,
	content string,
) error {
	bytes, err := json.Marshal([]string{
		receiver, subject, content,
	})

	if err != nil {
		return err
	}

	notifier.kafka.WriteMessages(
		context.Background(),
		Kafka.Message{
			Value: bytes,
		},
	)

	return nil
}
