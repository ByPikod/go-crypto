package core

import (
	"github.com/ByPikod/go-crypto/controllers"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/middleware"
	"github.com/ByPikod/go-crypto/repositories"
	"github.com/ByPikod/go-crypto/services"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type (
	Router struct {
		app                 *fiber.App
		exchangesRepository repositories.IExchangesRepository
		userRepository      repositories.IUserRepository
		walletRepository    repositories.IWalletRepository
	}
)

// Creates a new router object.
func NewRouter(
	exchangesRepository repositories.IExchangesRepository,
	userRepository repositories.IUserRepository,
	walletRepository repositories.IWalletRepository,
) *Router {
	return &Router{
		exchangesRepository: exchangesRepository,
		userRepository:      userRepository,
		walletRepository:    walletRepository,
	}
}

// Returns the app instance.
func (router *Router) App() *fiber.App {
	return router.app
}

// Initializes the router object.
func (router *Router) Initialize() {

	// Create fiber app
	App := fiber.New()
	App.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"code":    200,
			"message": "Hello world!",
		})
	})

	// Business layer (Services)
	var (
		exchangesService = services.NewExchangeService(router.exchangesRepository)
		userService      = services.NewUserService(router.userRepository)
		walletService    = services.NewWalletService(router.walletRepository)
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
	router.app = App

}
