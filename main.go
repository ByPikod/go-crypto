package main

import (
	"github.com/ByPikod/go-crypto/tree/notifier/helpers"
	"github.com/ByPikod/go-crypto/tree/notifier/kafka"
)

func main() {
	helpers.LogInfo("Initializing Notifier")
	consumer := kafka.CreateConsumer()
	consumer.ReadQueue()
}
