package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pisondev/mosque-api/internal/config"
	"github.com/pisondev/mosque-api/internal/router"
)

func main() {
	log := config.NewLogger()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn("no .env file found, falling back to system environment variables")
	}

	// Initialize Database
	db := config.ConnectDB(log)
	defer db.Close()

	// Initialize Fiber App
	app := fiber.New(fiber.Config{})

	// Setup Routes
	router.SetupRoutes(app, db, log)

	// Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Infof("server is starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
