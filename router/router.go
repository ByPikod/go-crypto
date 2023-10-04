package router

import (
	"github.com/ByPikod/go-crypto/controllers"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/middleware"
	"github.com/ByPikod/go-crypto/repositories"
	"github.com/ByPikod/go-crypto/services"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var App *fiber.App

// Respond with hello world and 200 status code
func helloWorld(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "Hello world!",
	})
}

func InitializeRouter(db *gorm.DB) {
	App = fiber.New()
	App.Get("/", helloWorld)

	// Repositories
	var (
		// Exchange rates
		exchangesRepository = repositories.NewExchangesRepository()
		exchangesService    = services.NewExchangeService(exchangesRepository)
		exchangesController = controllers.NewExchangesController(exchangesService)
		// User
		userRepository = repositories.NewUserRepository(db)
		userService    = services.NewUserService(userRepository)
		userController = controllers.NewUserController(userService)
		// Wallet
		walletRepository = repositories.NewWalletRepository(db)
		walletService    = services.NewWalletService(walletRepository)
		walletController = controllers.NewWalletController(walletService, exchangesService)
	)

	// Metrics
	prometheus := fiberprometheus.New("metrics")
	prometheus.RegisterAt(App, "/metrics")
	App.Use(prometheus.Middleware)

	// Middlewares
	var (
		authMiddleware = middleware.NewAuthMiddleware(userService)
	)

	// Websocket
	ws := App.Group("/ws")
	ws.Use(middleware.WebSocket) // Returns "426" if upgrade not provided.
	ws.Get("/exchange-rates", websocket.New(exchangesController.WSExchangeRates))

	// REST Api
	api := App.Group("/api")
	api.Get("/exchange-rates", exchangesController.ExchangeRates)

	user := api.Group("/user")
	user.Post("/register", middleware.Json, userController.Register)
	user.Post("/login", middleware.Json, userController.Login)
	user.Get("/me", authMiddleware.Auth, userController.Me)

	wallet := user.Group("/wallet", authMiddleware.Auth, middleware.Json)
	wallet.Post("/deposit", walletController.Deposit)
	wallet.Post("/buy", walletController.Buy)
	wallet.Post("/withdraw", walletController.Withdraw)
	wallet.Post("/sell", walletController.Sell)
	wallet.Get("/balance", walletController.Balance)

	// 404
	App.Use(helpers.NotFound)

	// Listen
	App.Listen(":8080")
}
