package main

import (
	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/repositories"
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
	helpers.LogInfo("Initialized config")

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

	// Initialize http server and routes.
	router := core.NewRouter(
		repositories.NewExchangesRepository(),
		repositories.NewUserRepository(db, config.AuthSecret),
		repositories.NewWalletRepository(db),
	)

	router.Initialize()
	router.App().Listen(config.Host + ":" + config.Listen)

}
