package routes

import (
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/workers"
	"github.com/gofiber/fiber/v2"
)

func Deposit(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Amount float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx)
	}

	// Retrieve wallet
	user := ctx.Locals("user").(*models.User)
	wallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Deposit
	transaction, err := wallet.AddTransaction(
		models.TRANSACTION_TYPE_DEPOSIT,
		payload.Amount,
	)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":     true,
		"mesasge":    "OK",
		"newBalance": transaction.Balance,
	})

}

func Buy(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx)
	}

	// Validate currency
	exchanges := workers.GetExchangeRates()
	exchange, ok := exchanges.Rates[payload.Currency]
	if !ok {
		return BadRequest(ctx, "Currency not found!")
	}

	neededUSDBalance := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve USD usdWallet
	usdWallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Valide balance
	if usdWallet.Balance < neededUSDBalance {
		return ctx.Status(200).JSON(fiber.Map{
			"status":  false,
			"message": "Insufficient balance!",
			"needed":  neededUSDBalance,
			"balance": usdWallet.Balance,
		})
	}

	buyCurrencyWallet, err := user.GetOrCreateWallet(payload.Currency)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := usdWallet.AddTransaction(models.TRANSACTION_TYPE_SELL, -neededUSDBalance)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}
	buyTransaction, err := buyCurrencyWallet.AddTransaction(models.TRANSACTION_TYPE_BUY, payload.Amount)
	if err != nil {
		panic(err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":          true,
		"mesasge":         "OK!",
		"sold_currency":   "USD",
		"sold_amount":     neededUSDBalance,
		"bought_currency": payload.Currency,
		"bought_amount":   payload.Amount,
		"Balance": fiber.Map{
			"USD":            sellTransaction.Balance,
			payload.Currency: buyTransaction.Balance,
		},
	})
}

func Withdraw(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Amount float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx)
	}
	payload.Amount = -payload.Amount

	// Retrieve wallet
	user := ctx.Locals("user").(*models.User)
	wallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Deposit
	transaction, err := wallet.AddTransaction(
		models.TRANSACTION_TYPE_WITHDRAW,
		payload.Amount,
	)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":     true,
		"mesasge":    "OK",
		"newBalance": transaction.Balance,
	})

}

func Sell(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx)
	}

	// Validate currency
	exchanges := workers.GetExchangeRates()
	exchange, ok := exchanges.Rates[payload.Currency]
	if !ok {
		return BadRequest(ctx, "Currency not found!")
	}

	neededUSDBalance := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve USD usdWallet
	usdWallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Valide balance
	if usdWallet.Balance < neededUSDBalance {
		return ctx.Status(200).JSON(fiber.Map{
			"status":  false,
			"message": "Insufficient balance!",
			"needed":  neededUSDBalance,
			"balance": usdWallet.Balance,
		})
	}

	buyCurrencyWallet, err := user.GetOrCreateWallet(payload.Currency)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := usdWallet.AddTransaction(models.TRANSACTION_TYPE_SELL, -neededUSDBalance)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}
	buyTransaction, err := buyCurrencyWallet.AddTransaction(models.TRANSACTION_TYPE_BUY, payload.Amount)
	if err != nil {
		panic(err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":          true,
		"mesasge":         "OK!",
		"sold_currency":   "USD",
		"sold_amount":     neededUSDBalance,
		"bought_currency": payload.Currency,
		"bought_amount":   payload.Amount,
		"Balance": fiber.Map{
			"USD":            sellTransaction.Balance,
			payload.Currency: buyTransaction.Balance,
		},
	})
}
