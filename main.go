package main

import (
	"github.com/ByPikod/go-crypto/controllers"
	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/log"
	"github.com/ByPikod/go-crypto/middleware"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/repositories"
	"github.com/ByPikod/go-crypto/services"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title           Go Crypto
// @version         1.0
// @description     Simple crypto app back-end API for educational purposes.

// @contact.name   API Support
// @contact.url    http://github.com/ByPikod
// @contact.email  admin@yahyabatulu.com

// @license.name  MIT
// @license.url   https://www.mit.edu/~amini/LICENSE.md

// @host      localhost:80
// @BasePath  /api/

// @securityDefinitions.basic  BasicAuth
func main() {

	// Load configuration from environment variables
	config := core.InitializeConfig()
	log.Info("Initialized config")

	// Initialize database with config above
	db := core.InitializeDatabase(config.Database)
	helpers.LogInfo("Initialized database")

	// Migration of the database
	helpers.LogInfo("Migrating database.")
	err := db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
	)
	if err != nil {
		panic(err)
	}

	helpers.LogInfo("Migrating completed.")

	// Create fiber app
	App := fiber.New()
	App.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"code":    200,
			"message": "Hello world!",
		})
	})

	// Initialize http server and routes.
	var (
		exchangesRepository = repositories.NewExchangesRepository()
		userRepository      = repositories.NewUserRepository(db)
		walletRepository    = repositories.NewWalletRepository(db)
	)

	// Business layer (Services)
	var (
		exchangesService = services.NewExchangeService(exchangesRepository)
		userService      = services.NewUserService(userRepository)
		walletService    = services.NewWalletService(walletRepository)
	)

	// Presentation layer (Controllers)
	var (
		exchangesController = controllers.NewExchangesController(exchangesService)
		userController      = controllers.NewUserController(userService)
		walletController    = controllers.NewWalletController(walletService, exchangesService)
	)

	// Metrics
	prometheus := fiberprometheus.New("metrics")
	prometheus.RegisterAt(App, "/metrics")
	App.Use(prometheus.Middleware)

	// Middlewares
	var (
		authMiddleware = middleware.NewAuthMiddleware(userService)
	)

	App.Use(cors.New(
		cors.Config{
			AllowHeaders:     "",
			AllowOrigins:     "*",
			AllowCredentials: true,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		},
	))

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
	App.Listen(config.Host + ":" + config.Listen)

}
