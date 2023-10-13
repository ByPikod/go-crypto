package controllers

import (
	"github.com/ByPikod/go-crypto/tree/crypto/helpers"
	"github.com/ByPikod/go-crypto/tree/crypto/log"
	"github.com/ByPikod/go-crypto/tree/crypto/models"
	"github.com/ByPikod/go-crypto/tree/crypto/services"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WalletController struct {
	service          *services.WalletService
	exchangesService *services.ExchangesService
}

// Create new wallet controller
func NewWalletController(
	service *services.WalletService,
	exchangesService *services.ExchangesService,
) *WalletController {
	return &WalletController{
		service:          service,
		exchangesService: exchangesService,
	}
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
func (controller *WalletController) Deposit(ctx *fiber.Ctx) error {

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
	wallet, err := controller.service.GetOrCreateWallet(user.ID, "USD")
	if err != nil {
		log.ControllerError("Deposit", err)
		return helpers.InternalServerError(ctx)
	}

	// Deposit
	transaction, err := controller.service.AddTransaction(
		wallet,
		models.TRANSACTION_TYPE_DEPOSIT,
		payload.Amount,
	)
	if err != nil {
		log.ControllerError("Deposit", err)
		return helpers.InternalServerError(ctx)
	}

	log.Info(
		"Successfull deposit",
		zap.Uint("user_id", user.ID),
		zap.String("user_mail", user.Mail),
		zap.Float64("amount", payload.Amount),
	)

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
func (controller *WalletController) Buy(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}

	// Validate currency
	exchange, ok := controller.exchangesService.GetCurrency(payload.Currency)
	if !ok {
		return helpers.BadRequest(ctx, "Currency not found!")
	}

	neededUSDBalance := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve USD Wallet
	usdWallet, err := controller.service.GetOrCreateWallet(user.ID, "USD")
	if err != nil {
		log.ControllerError("Buy", err)
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

	buyCurrencyWallet, err := controller.service.GetOrCreateWallet(user.ID, payload.Currency)
	if err != nil {
		log.ControllerError("Buy", err)
		return helpers.InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := controller.service.AddTransaction(
		usdWallet,
		models.TRANSACTION_TYPE_SELL,
		-neededUSDBalance,
	)

	if err != nil {
		log.ControllerError("Buy", err)
		return helpers.InternalServerError(ctx)
	}

	buyTransaction, err := controller.service.AddTransaction(
		buyCurrencyWallet,
		models.TRANSACTION_TYPE_BUY,
		payload.Amount,
	)
	if err != nil {
		panic(err)
	}

	log.Info(
		"Successfull buying",
		zap.Uint("user_ids", user.ID),
		zap.String("user_mail", user.Mail),
		zap.String("sold_currency", "USD"),
		zap.Float64("sold_amount", neededUSDBalance),
		zap.String("bought_currency", payload.Currency),
		zap.Float64("bought_amount", payload.Amount),
		zap.Float64("left", sellTransaction.Balance),
		zap.Float64("new_balance", buyTransaction.Balance),
	)

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
func (controller *WalletController) Withdraw(ctx *fiber.Ctx) error {

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
	wallet, err := controller.service.GetOrCreateWallet(user.ID, "USD")
	if err != nil {
		log.ControllerError("Withdraw", err)
		return helpers.InternalServerError(ctx)
	}

	// Deposit
	transaction, err := controller.service.AddTransaction(
		wallet,
		models.TRANSACTION_TYPE_WITHDRAW,
		payload.Amount,
	)
	if err != nil {
		log.ControllerError("Withdraw", err)
		return helpers.InternalServerError(ctx)
	}

	log.Info(
		"Successfull withdaraw",
		zap.Uint("user_id", user.ID),
		zap.String("user_mail", user.Mail),
		zap.Float64("amount", payload.Amount),
	)

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
func (controller *WalletController) Sell(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}

	// Validate currency
	exchange, ok := controller.exchangesService.GetCurrency(payload.Currency)
	if !ok {
		return helpers.BadRequest(ctx, "Currency not found!")
	}

	additionBalanceUSD := payload.Amount / exchange
	user := ctx.Locals("user").(*models.User)

	// Retrieve requested currency wallet
	requestedWallet, err := controller.service.GetWallet(user.ID, payload.Currency)
	if err != nil {
		log.ControllerError("Sell", err)
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
	usdWallet, err := controller.service.GetOrCreateWallet(user.ID, "USD")
	if err != nil {
		log.ControllerError("Sell", err)
		return helpers.InternalServerError(ctx)
	}

	// Add transactions
	sellTransaction, err := controller.service.AddTransaction(
		requestedWallet,
		models.TRANSACTION_TYPE_SELL,
		-payload.Amount,
	)
	if err != nil {
		log.ControllerError("Sell", err)
		return helpers.InternalServerError(ctx)
	}

	buyTransaction, err := controller.service.AddTransaction(
		usdWallet,
		models.TRANSACTION_TYPE_BUY,
		additionBalanceUSD,
	)
	if err != nil {
		panic(err)
	}

	log.Info(
		"Successfull selling",
		zap.Uint("user_ids", user.ID),
		zap.String("user_mail", user.Mail),
		zap.String("sold_currency", payload.Currency),
		zap.Float64("sold_amount", payload.Amount),
		zap.String("bought_currency", "USD"),
		zap.Float64("bought_amount", additionBalanceUSD),
		zap.Float64("left", sellTransaction.Balance),
		zap.Float64("new_balance", buyTransaction.Balance),
	)

	return ctx.Status(200).JSON(fiber.Map{
		"status":          true,
		"message":         "OK!",
		"sold_currency":   payload.Currency,
		"sold_amount":     payload.Amount,
		"bought_currency": "USD",
		"bought_amount":   additionBalanceUSD,
		"Balance": fiber.Map{
			payload.Currency: sellTransaction.Balance,
			"USD":            buyTransaction.Balance,
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
func (controller *WalletController) Balance(ctx *fiber.Ctx) error {

	// Fetch wallets
	user := ctx.Locals("user").(*models.User)
	err := controller.service.LoadWallets(user)
	if err != nil {
		log.ControllerError("Balance", err)
		return helpers.InternalServerError(ctx)
	}

	// Prepare response
	respond := map[string]float64{}

	for _, wallet := range user.Wallets {
		respond[wallet.Currency] = wallet.Balance
	}

	return ctx.Status(200).JSON(respond)
}
