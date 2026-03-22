package router

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/pisondev/mosque-api/docs"
	"github.com/pisondev/mosque-api/internal/middleware"
	"github.com/pisondev/mosque-api/internal/module/auth"
	"github.com/pisondev/mosque-api/internal/module/management"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool, log *logrus.Logger) {
	api := app.Group("/api/v1")
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

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

	managementRepo := management.NewRepository(db)
	managementService := management.NewService(managementRepo, log)
	managementController := management.NewController(managementService, log)

	tenantGroup.Get("/me", managementController.TenantMe)
	tenantGroup.Patch("/setup", managementController.SetupTenant)
	tenantGroup.Get("/domains", managementController.ListDomains)
	tenantGroup.Post("/domains", managementController.CreateDomain)
	tenantGroup.Patch("/domains/:id", managementController.UpdateDomain)
	tenantGroup.Delete("/domains/:id", managementController.DeleteDomain)

	tenantGroup.Get("/profile", managementController.GetProfile)
	tenantGroup.Put("/profile", managementController.UpsertProfile)

	tenantGroup.Get("/tags", managementController.ListTags)
	tenantGroup.Post("/tags", managementController.CreateTag)
	tenantGroup.Patch("/tags/:id", managementController.UpdateTag)
	tenantGroup.Delete("/tags/:id", managementController.DeleteTag)

	tenantGroup.Get("/posts", managementController.ListPosts)
	tenantGroup.Post("/posts", managementController.CreatePost)
	tenantGroup.Get("/posts/:id", managementController.GetPost)
	tenantGroup.Put("/posts/:id", managementController.UpdatePost)
	tenantGroup.Patch("/posts/:id/status", managementController.UpdatePostStatus)
	tenantGroup.Delete("/posts/:id", managementController.DeletePost)

	tenantGroup.Get("/static-pages", managementController.ListStaticPages)
	tenantGroup.Put("/static-pages/:slug", managementController.UpsertStaticPage)
}
