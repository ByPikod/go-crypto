package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/middleware"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/workers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

const auth = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.HBNfNTMv3Jd9Wf-m3v6buHgGLQL0Srl8zwGro8JHcO4"

func init() {

	// Load configuration from environment variables
	core.InitializeConfig()
	helpers.LogInfo("Initialized config")
	// Initialize database with config above
	core.InitializeDatabase(core.Config.Database)
	helpers.LogInfo("Initialized database")

	// Migration of the database
	helpers.LogInfo("Migrating database.")
	err := core.DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
	)
	if err != nil {
		panic(err)
	}

	helpers.LogInfo("Migrating completed.")

	// Start fetching exchange rate
	workers.UpdateExchangeRate()
	helpers.LogInfo("Initialized exchange worker.")

}

func TestDeposit(t *testing.T) {
	const (
		route  = "/api/user/wallet/deposit"
		method = "POST"
	)

	app := fiber.New()
	app.Post(route, middleware.Auth, Deposit)

	tests := []struct {
		description string
		expected    int
		data        map[string]interface{}
		headers     map[string]string
	}{
		{
			description: "get OK",
			expected:    200,
			data: fiber.Map{
				"amount": 1.4,
			},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get bad request",
			expected:    400,
			data:        fiber.Map{},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get unauthorized",
			expected:    401,
			data:        fiber.Map{},
		},
	}

	for _, tt := range tests {

		data, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}

		t.Logf("\nexpected: %v\nrequest body: %v", tt.expected, string(data))

		req := httptest.NewRequest(method, route, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		for k, v := range tt.headers {
			req.Header.Set(k, v)
		}

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to send test request: %v", err)
		}

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		assert.Equalf(t, tt.expected, resp.StatusCode, tt.description)
		t.Logf("\nstatus: %v\nresponse body: %v", resp.StatusCode, string(bodyByte))
	}

}

func TestBuy(t *testing.T) {
	const (
		route  = "/api/user/wallet/buy"
		method = "POST"
	)

	app := fiber.New()
	app.Post(route, middleware.Auth, Buy)

	tests := []struct {
		description string
		expected    int
		data        map[string]interface{}
		headers     map[string]string
	}{
		{
			description: "get OK",
			expected:    200,
			data: fiber.Map{
				"amount":   0.0001,
				"currency": "BTC",
			},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get bad request",
			expected:    400,
			data:        fiber.Map{},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get unauthorized",
			expected:    401,
			data:        fiber.Map{},
		},
	}

	for _, tt := range tests {

		data, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}

		t.Logf("\nexpected: %v\nrequest body: %v", tt.expected, string(data))

		req := httptest.NewRequest(method, route, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		for k, v := range tt.headers {
			req.Header.Set(k, v)
		}

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to send test request: %v", err)
		}

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		assert.Equalf(t, tt.expected, resp.StatusCode, tt.description)
		t.Logf("\nstatus: %v\nresponse body: %v", resp.StatusCode, string(bodyByte))
	}

}

func TestSell(t *testing.T) {
	const (
		route  = "/api/user/wallet/sell"
		method = "POST"
	)

	app := fiber.New()
	app.Post(route, middleware.Auth, Sell)

	tests := []struct {
		description string
		expected    int
		data        map[string]interface{}
		headers     map[string]string
	}{
		{
			description: "get OK",
			expected:    200,
			data: fiber.Map{
				"amount":   0.0001,
				"currency": "BTC",
			},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get bad request",
			expected:    400,
			data:        fiber.Map{},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get unauthorized",
			expected:    401,
			data:        fiber.Map{},
		},
	}

	for _, tt := range tests {

		data, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}

		t.Logf("\nexpected: %v\nrequest body: %v", tt.expected, string(data))

		req := httptest.NewRequest(method, route, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		for k, v := range tt.headers {
			req.Header.Set(k, v)
		}

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to send test request: %v", err)
		}

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		assert.Equalf(t, tt.expected, resp.StatusCode, tt.description)
		t.Logf("\nstatus: %v\nresponse body: %v", resp.StatusCode, string(bodyByte))
	}

}

func TestWithdraw(t *testing.T) {
	const (
		route  = "/api/user/wallet/withdraw"
		method = "POST"
	)

	app := fiber.New()
	app.Post(route, middleware.Auth, Withdraw)

	tests := []struct {
		description string
		expected    int
		data        map[string]interface{}
		headers     map[string]string
	}{
		{
			description: "get OK",
			expected:    200,
			data: fiber.Map{
				"amount": 1.0,
			},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get bad request",
			expected:    400,
			data:        fiber.Map{},
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get unauthorized",
			expected:    401,
			data:        fiber.Map{},
		},
	}

	for _, tt := range tests {

		data, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}

		t.Logf("\nexpected: %v\nrequest body: %v", tt.expected, string(data))

		req := httptest.NewRequest(method, route, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		for k, v := range tt.headers {
			req.Header.Set(k, v)
		}

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to send test request: %v", err)
		}

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		assert.Equalf(t, tt.expected, resp.StatusCode, tt.description)
		t.Logf("\nstatus: %v\nresponse body: %v", resp.StatusCode, string(bodyByte))
	}

}

func TestBalance(t *testing.T) {
	const (
		route  = "/api/user/wallet/balance"
		method = "GET"
	)

	app := fiber.New()
	app.Get(route, middleware.Auth, Balance)

	tests := []struct {
		description string
		expected    int
		data        map[string]interface{}
		headers     map[string]string
	}{
		{
			description: "get OK",
			expected:    200,
			headers: map[string]string{
				"Authorization": "Bearer " + auth,
			},
		},
		{
			description: "get unauthorized",
			expected:    401,
			data:        fiber.Map{},
		},
	}

	for _, tt := range tests {

		data, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}

		t.Logf("\nexpected: %v\nrequest body: %v", tt.expected, string(data))

		req := httptest.NewRequest(method, route, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		for k, v := range tt.headers {
			req.Header.Set(k, v)
		}

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to send test request: %v", err)
		}

		bodyByte, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		assert.Equalf(t, tt.expected, resp.StatusCode, tt.description)
		t.Logf("\nstatus: %v\nresponse body: %v", resp.StatusCode, string(bodyByte))
	}

}
