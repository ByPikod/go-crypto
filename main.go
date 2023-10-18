package main

import (
	"github.com/ByPikod/go-crypto/tree/notifier/helpers"
	"github.com/ByPikod/go-crypto/tree/notifier/notifier"
)

func main() {
	helpers.LogInfo("Initializing Notifier")
	consumer := notifier.CreateConsumer()
	consumer.ReadQueue(notifier.ConsumeMessage)
}
