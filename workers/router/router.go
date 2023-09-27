package router

import (
	"github.com/ByPikod/go-crypto/middleware"
	"github.com/ByPikod/go-crypto/routes"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

// Respond with hello world and 200 status code
func helloWorld(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"code":    200,
		"message": "Hello world!",
	})
}

func InitializeRouter() {
	App = fiber.New()
	App.Get("/", helloWorld)
	apiRoutes(App)
	wsRoutes(App)

	// 404
	App.Use(routes.NotFound)

	// Listen
	App.Listen(":80")
}

// Route: /ws/
func wsRoutes(parent fiber.Router) {

	ws := parent.Group("/ws")
	ws.Use(middleware.WebSocket) // Returns "426" if upgrade not provided.

	ws.Get("/exchange-rates", websocket.New(routes.WSExchangeRates))
	go routes.WSExchangeBroadcaster() // Broadcast exchange data with an interval.

}

// Route: /api/
func apiRoutes(parent fiber.Router) {
	api := parent.Group("/api")
	api.Get("/exchange-rates", routes.ExchangeRates)
	userRoutes(api)
}

// Route: /api/user/
func userRoutes(parent fiber.Router) {
	user := parent.Group("/user")
	user.Post("/register", middleware.Json, routes.Register)
	user.Post("/login", middleware.Json, routes.Login)
	user.Get("/me", middleware.Auth, routes.Me)
	walletRoutes(user)
}

// Route: /api/user/wallet
func walletRoutes(parent fiber.Router) {
	wallet := parent.Group("/wallet", middleware.Auth, middleware.Json)
	wallet.Get("/deposit", routes.Deposit)
	wallet.Get("/buy", routes.Buy)
	wallet.Get("/withdraw", routes.Withdraw)
	wallet.Get("/sell", routes.Sell)
	wallet.Get("/balance", routes.Balance)
}
