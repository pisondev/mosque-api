package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pisondev/mosque-api/internal/module/auth"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool, log *logrus.Logger) {
	api := app.Group("/api/v1")

	// Setup Auth Module
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, log)
	authController := auth.NewController(authService, log)

	// Auth Routes
	authGroup := api.Group("/auth")
	authGroup.Post("/google", authController.GoogleLogin)
}
