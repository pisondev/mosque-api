package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pisondev/mosque-api/internal/middleware"
	"github.com/pisondev/mosque-api/internal/module/auth"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool, log *logrus.Logger) {
	api := app.Group("/api/v1")

	// Setup Auth Module (Rute Publik)
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, log)
	authController := auth.NewController(authService, log)

	authGroup := api.Group("/auth")
	authGroup.Post("/google", authController.GoogleLogin)

	// ==========================================
	// PROTECTED ROUTES (Hanya bisa diakses jika ada JWT)
	// ==========================================

	// Semua rute di bawah "tenantGroup" ini otomatis dilindungi oleh middleware.Protected()
	tenantGroup := api.Group("/tenant", middleware.Protected())

	// Endpoint uji coba sementara untuk membuktikan middleware berfungsi
	tenantGroup.Get("/me", func(c *fiber.Ctx) error {
		// Mengambil data yang tadi disuntikkan oleh middleware
		tenantID := c.Locals("tenant_id")
		email := c.Locals("email")

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Selamat datang di area tertutup!",
			"data": fiber.Map{
				"email":     email,
				"tenant_id": tenantID,
			},
		})
	})
}
