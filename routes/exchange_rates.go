package routes

import (
	"github.com/ByPikod/go-crypto/workers"
	"github.com/gofiber/fiber/v2"
)

func ExchangeRates(ctx *fiber.Ctx) error {
	exchangeRates := workers.GetExchangeRates()
	return ctx.Status(200).JSON(exchangeRates)
}
