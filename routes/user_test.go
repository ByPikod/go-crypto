package routes

import (
	"testing"

	"github.com/ByPikod/go-crypto/testhttp"
	"github.com/gofiber/fiber/v2"
)

func TestRegister(t *testing.T) {
	base := testhttp.New().
		SetRoute("/api/user/register")

	base.GiveBirth().SetBodyJSON(fiber.Map{
		"name":      "Yahya",
		"firstName": "Batulu",
		"mail":      "admin@yahyabatulu.com",
		"password":  "demo123",
	}).Execute(func() {

	})

}
