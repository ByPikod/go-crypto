package main

import (
	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/workers"
)

func main() {
	// Initializations
	core.InitializeConfig()                       // Load configuration from environment variables
	core.InitializeDatabase(core.Config.Database) // Initialize database with config above
	go workers.InitializeExchangeRateWorker()     // Start fetching exchange rate
	core.InitializeServer()                       // Initialize http server and routes.
}
