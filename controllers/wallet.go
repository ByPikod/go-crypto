package controllers

import (
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/services"
	"github.com/gofiber/fiber/v2"
)

type WalletController struct {
	service *services.WalletService
}

// Create new wallet controller
func NewWalletController(service *services.WalletService) *WalletController {
	return &WalletController{service: service}
}

// @Summary		Deposit
// @Description Add balance to the user account.
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Param		amount			body		number	true	"Amount of money to deposit."	format(float)
// @Success		200				{object}	interface{}
// @Failure		401				{object}	interface{}
// @Failure		400				{object}	interface{}
// @Security 	ApiKeyAuth
// @Router		/user/wallet/deposit [post]
func (walletController *WalletController) Deposit(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Amount float64 `json:"amount"`
	}

	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}

	if payload.Amount == 0 {
		return helpers.BadRequest(ctx)
	}

	// Retrieve wallet
	user := ctx.Locals("user").(*models.User)
	wallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	// Deposit
	transaction, err := wallet.AddTransaction(
		models.TRANSACTION_TYPE_DEPOSIT,
		payload.Amount,
	)
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":     true,
		"message":    "OK",
		"newBalance": transaction.Balance,
	})

}

// @Summary		Buy
// @Description Buy a specific amount of crypto.
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Param		amount			body		number	true	"Amount of unit to buy."	format(float)
// @Param		currency		body		string	true	"Currency to buy."
// @Success		200				{object}	interface{}
// @Failure		401				{object}	interface{}
// @Failure		400				{object}	interface{}
// @Security 	ApiKeyAuth
// @Router		/user/wallet/buy [post]
func (walletController *WalletController) Buy(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}

	// Validate currency
	exchanges := workers.GetExchangeRates()
	exchange, ok := exchanges.Rates[payload.Currency]
	if !ok {
		return helpers.BadRequest(ctx, "Currency not found!")
	}

	neededUSDBalance := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve USD Wallet
	usdWallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
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
		return helpers.InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := usdWallet.AddTransaction(models.TRANSACTION_TYPE_SELL, -neededUSDBalance)
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
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

// @Summary		Withdraw
// @Description Withdraw money from your account.
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Param		amount			body		number	true	"Amount of money to withdraw."	format(float)
// @Success		200				{object}	interface{}
// @Failure		401				{object}	interface{}
// @Failure		400				{object}	interface{}
// @Security 	ApiKeyAuth
// @Router		/user/wallet/withdraw [post]
func (walletController *WalletController) Withdraw(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Amount float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}
	if payload.Amount == 0 {
		return helpers.BadRequest(ctx)
	}

	payload.Amount = -payload.Amount

	// Retrieve wallet
	user := ctx.Locals("user").(*models.User)
	wallet, err := user.GetOrCreateWallet("USD")
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	// Deposit
	transaction, err := wallet.AddTransaction(
		models.TRANSACTION_TYPE_WITHDRAW,
		payload.Amount,
	)
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":     true,
		"message":    "OK",
		"newBalance": transaction.Balance,
	})

}

// @Summary		Sell
// @Description Sell a specific amount of crypto.
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Param		amount			body		number	true	"Amount of unit to sell."	format(float)
// @Param		currency		body		string	true	"Currency to sell."
// @Success		200				{object}	interface{}
// @Failure		401				{object}	interface{}
// @Failure		400				{object}	interface{}
// @Security 	ApiKeyAuth
// @Router		/user/wallet/sell [post]
func (walletController *WalletController) Sell(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}

	// Validate currency
	exchanges := workers.GetExchangeRates()
	exchange, ok := exchanges.Rates[payload.Currency]
	if !ok {
		return helpers.BadRequest(ctx, "Currency not found!")
	}

	additionBalanceUSD := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve requested currency wallet
	requestedWallet, err := user.GetWallet(payload.Currency)
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
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
		return helpers.InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := requestedWallet.AddTransaction(models.TRANSACTION_TYPE_SELL, -payload.Amount)
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
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

// @Summary		Balance
// @Description Retrieves the balance
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Success		200				{object}	interface{}
// @Failure		401				{object}	interface{}
// @Security 	ApiKeyAuth
// @Router		/user/wallet/balance [get]
func (walletController *WalletController) Balance(ctx *fiber.Ctx) error {

	// Fetch wallets
	user := ctx.Locals("user").(*models.User)
	err := user.PreloadWallets()
	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	// Prepare response
	respond := map[string]float64{}

	for _, wallet := range user.Wallets {
		respond[wallet.Currency] = wallet.Balance
	}

	return ctx.Status(200).JSON(respond)
}
