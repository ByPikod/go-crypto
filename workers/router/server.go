package router

import (
	"github.com/ByPikod/go-crypto/middleware"
	"github.com/ByPikod/go-crypto/routes"
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

func InitializeServer() {
	App = fiber.New()
	api := App.Group("/api")
	user := api.Group("/user")

	// Handle API requests
	App.Get("/", helloWorld)
	api.Get("/exchange-rates", routes.ExchangeRates)

	user.Post("/register", middleware.Json, routes.Register)
	user.Post("/login", middleware.Json, routes.Login)
	user.Post("/me", routes.Register)

	// 404
	App.Use(routes.NotFound)

	// Listen
	App.Listen(":80")
}
