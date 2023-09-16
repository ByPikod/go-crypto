package main

import (
	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/workers"
	"github.com/ByPikod/go-crypto/workers/router"
)

func main() {
	core.InitializeConfig() // Load configuration from environment variables

	// Initialize database with config above
	core.InitializeDatabase(core.Config.Database)

	// Migration of the database
	helpers.LogInfo("Migrating database.")
	err := core.DB.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	helpers.LogInfo("Migrating completed.")

	go workers.InitializeExchangeRateWorker() // Start fetching exchange rate
	router.InitializeServer()                 // Initialize http server and routes.
}
