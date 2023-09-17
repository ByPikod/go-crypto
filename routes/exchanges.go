package routes

import (
	"github.com/ByPikod/go-crypto/workers"
	"github.com/gofiber/fiber/v2"
)

// Handle api/exchange-rates endpoint.
func ExchangeRates(ctx *fiber.Ctx) error {
	exchangeRates := workers.GetExchangeRates()
	return ctx.Status(200).JSON(exchangeRates)
}
