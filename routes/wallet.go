package routes

import (
	"github.com/ByPikod/go-crypto/core"
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
		"message":    "OK",
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

	// Retrieve USD Wallet
	usdWallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Validate balance
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
		"message":         "OK!",
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
		"message":    "OK",
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

	additionBalanceUSD := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve requested currency wallet
	requestedWallet, err := user.GetWallet(payload.Currency)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}
	if requestedWallet == nil {
		return ctx.Status(200).JSON(fiber.Map{
			"status":  false,
			"message": "Insufficient balance!",
			"needed":  payload.Amount,
			"balance": 0,
		})
	}

	// Validate balance
	if requestedWallet.Balance < payload.Amount {
		return ctx.Status(200).JSON(fiber.Map{
			"status":  false,
			"message": "Insufficient balance!",
			"needed":  payload.Amount,
			"balance": requestedWallet.Balance,
		})
	}

	// Retrieve USD Wallet
	usdWallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := requestedWallet.AddTransaction(models.TRANSACTION_TYPE_SELL, -payload.Amount)
	if err != nil {
		helpers.LogError(err.Error())
		return InternalServerError(ctx)
	}
	buyTransaction, err := usdWallet.AddTransaction(models.TRANSACTION_TYPE_BUY, additionBalanceUSD)
	if err != nil {
		panic(err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":          true,
		"message":         "OK!",
		"sold_currency":   payload.Currency,
		"sold_amount":     payload.Amount,
		"bought_currency": "USD",
		"bought_amount":   additionBalanceUSD,
		"Balance": fiber.Map{
			"USD":            sellTransaction.Balance,
			payload.Currency: buyTransaction.Balance,
		},
	})
}

func Balance(ctx *fiber.Ctx) error {

	// Fetch wallets
	user := ctx.Locals("user").(*models.User)
	res := core.DB.Preload("Wallets").Find(user)
	if res.Error != nil {
		helpers.LogError(res.Error.Error())
		return InternalServerError(ctx)
	}

	// Prepare response
	respond := map[string]float64{}

	for _, wallet := range user.Wallets {
		respond[wallet.Currency] = wallet.Balance
	}

	return ctx.Status(200).JSON(respond)
}
