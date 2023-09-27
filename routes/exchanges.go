package routes

import (
	"time"

	"github.com/ByPikod/go-crypto/workers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	WSExchangeClients = make(map[*websocket.Conn]chan *workers.ExchangeRates)
)

// @Summary		Exchange Rates
// @Description List exchange rates
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Success		200				{object}	workers.ExchangeRates
// @Security 	ApiKeyAuth
// @Router		/exchange-rates [get]
func ExchangeRates(ctx *fiber.Ctx) error {
	exchangeRates := workers.GetExchangeRates()
	return ctx.Status(200).JSON(exchangeRates)
}

// Handle ws/exchange-rates
func WSExchangeRates(ws *websocket.Conn) {
	// Create a channel to receive broadcasts
	ch := make(chan *workers.ExchangeRates)
	// Add client to the listeners
	WSExchangeClients[ws] = ch
	// Remove client from listeners
	defer delete(WSExchangeClients, ws)
	// Close channel
	defer close(ch)

	// Broadcast received message
	for {
		lastExchangeData := <-ch
		ws.WriteJSON(lastExchangeData)
	}
}

// Broadcast the last exchange data to all the clients with an interval.
func WSExchangeBroadcaster() {
	// Wait
	for range time.Tick(5 * time.Second) {
		lastExchangeData := workers.GetExchangeRates()
		// Broadcast the last exchange data to all the clients connected.
		for _, ch := range WSExchangeClients {
			ch <- lastExchangeData
		}
	}
}
