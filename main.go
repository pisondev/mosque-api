package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/pisondev/mosque-api/internal/config"
	"github.com/pisondev/mosque-api/internal/router"
)

func main() {
	log := config.NewLogger()

	if err := godotenv.Load(); err != nil {
		log.Warn("no .env file found, falling back to system environment variables")
	}

	db := config.ConnectDB(log)
	defer db.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Setup Routes
	router.SetupRoutes(app, db, log)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Infof("server is starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
