package main

import (
	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/router"
	"github.com/ByPikod/go-crypto/workers"
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
	// Initialize things router depend
	initialize()
	// Initialize http server and routes.
	router.InitializeRouter()
}

// Initialize things router depend. Seperated for test files.
func initialize() {

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
	go workers.InitializeExchangeRateWorker()
	helpers.LogInfo("Initialized exchange worker.")

}
