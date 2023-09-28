package controllers

import (
	"github.com/ByPikod/go-crypto/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type (
	ExchangesController struct {
		service *services.ExchangesService
	}
)

func NewExchangesController(service *services.ExchangesService) *ExchangesController {
	controller := &ExchangesController{service: service}
	return controller
}

// @Summary		Exchange Rates
// @Description List exchange rates
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Success		200				{object}	workers.ExchangeRates
// @Security 	ApiKeyAuth
// @Router		/exchange-rates [get]
func (controller *ExchangesController) ExchangeRates(ctx *fiber.Ctx) error {
	exchangeRates := controller.service.GetExchangeRates()
	return ctx.Status(200).JSON(exchangeRates)
}

// Handle ws/exchange-rates
func (controller *ExchangesController) WSExchangeRates(ws *websocket.Conn) {
	// Listen
	ch := controller.service.AddClient(ws)
	defer controller.service.RemoveClient(ws)

	// Broadcast received message
	for {
		lastExchangeData := <-ch
		ws.WriteJSON(lastExchangeData)
	}
}
