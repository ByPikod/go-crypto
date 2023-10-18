package notifier

import (
	"encoding/json"

	"github.com/ByPikod/go-crypto/tree/notifier/helpers"
	"github.com/ByPikod/go-crypto/tree/notifier/log"
	Kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type (
	MessageBody struct {
		Receiver string
		Subject  string
		Body     string
	}
	MessageBodyRaw []string
)

// Called when a message is received.
func ConsumeMessage(message Kafka.Message) {
	body, err := ParseMessage(message)
	if err != nil {
		helpers.LogError("Failed to parse message: " + err.Error())
	}

	SendMail(body.Receiver, body.Subject, body.Body)
}

// Parse a message from Kafka.
func ParseMessage(message Kafka.Message) (*MessageBody, error) {

	var bodyRaw MessageBodyRaw
	helpers.LogInfo(string(message.Value))
	err := json.Unmarshal(message.Value, &bodyRaw)
	if err != nil {
		return nil, err
	}

	body := MessageBody{
		Receiver: bodyRaw[0],
		Subject:  bodyRaw[1],
		Body:     bodyRaw[2],
	}

	return &body, nil

}

// Send mail
func SendMail(receiver string, subject string, body string) {
	log.Info(
		"Mail sent",
		zap.String("receiver", receiver),
		zap.String("subject", subject),
		zap.String("body", body),
	)
}
