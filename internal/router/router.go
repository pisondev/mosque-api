package router

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/pisondev/mosque-api/docs"
	"github.com/pisondev/mosque-api/internal/constant"
	"github.com/pisondev/mosque-api/internal/middleware"
	"github.com/pisondev/mosque-api/internal/module/auth"
	"github.com/pisondev/mosque-api/internal/module/community"
	"github.com/pisondev/mosque-api/internal/module/engagement"
	"github.com/pisondev/mosque-api/internal/module/finance"
	"github.com/pisondev/mosque-api/internal/module/management"
	"github.com/pisondev/mosque-api/internal/module/worship"
	"github.com/pisondev/mosque-api/internal/storage"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool, log *logrus.Logger) {
	api := app.Group("/api/v1")
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// Setup Auth Module (Rute Publik)
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, log, auth.NewResendSender(log))
	authController := auth.NewController(authService, log)

	authGroup := api.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)
	authGroup.Post("/forgot-password", authController.ForgotPassword)
	authGroup.Post("/reset-password", authController.ResetPassword)
	authGroup.Post("/google", authController.GoogleLogin)

	// ==========================================
	// PROTECTED ROUTES (Hanya bisa diakses jika ada JWT)
	// ==========================================

	// Semua rute di bawah "tenantGroup" ini otomatis dilindungi oleh middleware.Protected()
	tenantGroup := api.Group("/tenant", middleware.Protected())

	managementRepo := management.NewRepository(db)
	managementService := management.NewService(managementRepo, log)
	managementController := management.NewController(managementService, log)

	worshipRepo := worship.NewRepository(db)
	worshipService := worship.NewService(worshipRepo, log)
	worshipController := worship.NewController(worshipService, log)

	communityRepo := community.NewRepository(db)
	communityService := community.NewService(communityRepo, log)
	communityController := community.NewController(communityService, log)

	engagementRepo := engagement.NewRepository(db)
	engagementService := engagement.NewService(engagementRepo, log)
	engagementController := engagement.NewController(engagementService, log)

	// ==========================================
	// INIT FINANCE MODULE
	// ==========================================
	financeRepo := finance.NewRepository(db)
	financeService := finance.NewService(financeRepo, log)
	financeController := finance.NewController(financeService, log)

	// Storage / Upload
	store := storage.New()
	tenantGroup.Post("/upload", management.UploadHandler(store, db))
	tenantGroup.Get("/storage-quota", management.StorageQuotaHandler(store, db))

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
	// Endpoint SaaS State untuk Frontend
	tenantGroup.Get("/billing-status", managementController.GetBillingStatus)
	tenantGroup.Get("/subscription/plans", financeController.ListSubscriptionPlans)
	tenantGroup.Post("/subscription/checkout", financeController.CreateSubscriptionCheckout)
	tenantGroup.Post("/subscription/quote", financeController.GetSubscriptionQuote)
	tenantGroup.Post("/subscription/checkout-v2", financeController.CreateSubscriptionCheckoutFromQuote)
	tenantGroup.Post("/subscription/activate-free", financeController.ActivateFreePlan)
	tenantGroup.Get("/subscription/transactions", financeController.ListSubscriptionTransactions)
	tenantGroup.Get("/subscription/transactions/active", financeController.GetActiveSubscriptionTransaction)
	tenantGroup.Post("/subscription/transactions/:id/cancel", financeController.CancelSubscriptionTransaction)
	tenantGroup.Get("/subscription/transactions/:id", financeController.GetSubscriptionTransaction)

	tenantGroup.Get("/prayer-time-settings", worshipController.GetPrayerTimeSettings)
	tenantGroup.Put("/prayer-time-settings", worshipController.UpsertPrayerTimeSettings)
	tenantGroup.Get("/prayer-times-daily", worshipController.ListPrayerTimesDaily)
	tenantGroup.Post("/prayer-times-daily", worshipController.CreatePrayerTimesDaily)
	tenantGroup.Get("/prayer-times-daily/:id", worshipController.GetPrayerTimesDaily)
	tenantGroup.Put("/prayer-times-daily/:id", worshipController.UpdatePrayerTimesDaily)
	tenantGroup.Delete("/prayer-times-daily/:id", worshipController.DeletePrayerTimesDaily)
	tenantGroup.Get("/prayer-duties", worshipController.ListPrayerDuties)
	tenantGroup.Post("/prayer-duties", worshipController.CreatePrayerDuty)
	tenantGroup.Get("/prayer-duties/:id", worshipController.GetPrayerDuty)
	tenantGroup.Put("/prayer-duties/:id", worshipController.UpdatePrayerDuty)
	tenantGroup.Delete("/prayer-duties/:id", worshipController.DeletePrayerDuty)
	tenantGroup.Get("/special-days", worshipController.ListSpecialDays)
	tenantGroup.Post("/special-days", worshipController.CreateSpecialDay)
	tenantGroup.Get("/special-days/:id", worshipController.GetSpecialDay)
	tenantGroup.Put("/special-days/:id", worshipController.UpdateSpecialDay)
	tenantGroup.Delete("/special-days/:id", worshipController.DeleteSpecialDay)
	tenantGroup.Get("/prayer-calendar", worshipController.GetPrayerCalendar)

	tenantGroup.Get("/events", communityController.ListEvents)
	tenantGroup.Post("/events", communityController.CreateEvent)
	tenantGroup.Get("/events/:id", communityController.GetEvent)
	tenantGroup.Put("/events/:id", communityController.UpdateEvent)
	tenantGroup.Patch("/events/:id/status", communityController.UpdateEventStatus)
	tenantGroup.Delete("/events/:id", communityController.DeleteEvent)
	tenantGroup.Get("/gallery/albums", communityController.ListGalleryAlbums)
	tenantGroup.Post("/gallery/albums", communityController.CreateGalleryAlbum)
	tenantGroup.Get("/gallery/albums/:id", communityController.GetGalleryAlbum)
	tenantGroup.Put("/gallery/albums/:id", communityController.UpdateGalleryAlbum)
	tenantGroup.Delete("/gallery/albums/:id", communityController.DeleteGalleryAlbum)
	tenantGroup.Get("/gallery/items", communityController.ListGalleryItems)
	tenantGroup.Post("/gallery/items", communityController.CreateGalleryItem)
	tenantGroup.Get("/gallery/items/:id", communityController.GetGalleryItem)
	tenantGroup.Put("/gallery/items/:id", communityController.UpdateGalleryItem)
	tenantGroup.Delete("/gallery/items/:id", communityController.DeleteGalleryItem)
	tenantGroup.Get("/management-members", communityController.ListManagementMembers)
	tenantGroup.Post("/management-members", communityController.CreateManagementMember)
	tenantGroup.Get("/management-members/:id", communityController.GetManagementMember)
	tenantGroup.Put("/management-members/:id", communityController.UpdateManagementMember)
	tenantGroup.Delete("/management-members/:id", communityController.DeleteManagementMember)

	tenantGroup.Get("/static-payment-methods", engagementController.ListStaticPaymentMethods)
	tenantGroup.Post("/static-payment-methods", engagementController.CreateStaticPaymentMethod)
	tenantGroup.Get("/static-payment-methods/:id", engagementController.GetStaticPaymentMethod)
	tenantGroup.Put("/static-payment-methods/:id", engagementController.UpdateStaticPaymentMethod)
	tenantGroup.Delete("/static-payment-methods/:id", engagementController.DeleteStaticPaymentMethod)
	tenantGroup.Get("/social-links", engagementController.ListSocialLinks)
	tenantGroup.Post("/social-links", engagementController.CreateSocialLink)
	tenantGroup.Get("/social-links/:id", engagementController.GetSocialLink)
	tenantGroup.Put("/social-links/:id", engagementController.UpdateSocialLink)
	tenantGroup.Delete("/social-links/:id", engagementController.DeleteSocialLink)
	tenantGroup.Get("/external-links", engagementController.ListExternalLinks)
	tenantGroup.Post("/external-links", engagementController.CreateExternalLink)
	tenantGroup.Get("/external-links/:id", engagementController.GetExternalLink)
	tenantGroup.Put("/external-links/:id", engagementController.UpdateExternalLink)
	tenantGroup.Delete("/external-links/:id", engagementController.DeleteExternalLink)
	tenantGroup.Get("/feature-catalog", engagementController.ListFeatureCatalog)
	tenantGroup.Get("/website-features", engagementController.ListWebsiteFeatures)
	tenantGroup.Put("/website-features/:feature_id", engagementController.UpsertWebsiteFeature)
	tenantGroup.Patch("/website-features/bulk", engagementController.BulkUpsertWebsiteFeatures)

	// ==========================================
	// FINANCE & DONATION ROUTES (PROTECTED)
	// ==========================================
	// Semua fitur donasi digital butuh paket PRO++ atau MAX+++

	tenantGroup.Get("/pg-config", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.GetPGConfig)
	tenantGroup.Put("/pg-config", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.UpsertPGConfig)

	tenantGroup.Get("/campaigns", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.ListCampaigns)
	tenantGroup.Post("/campaigns", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.CreateCampaign)
	tenantGroup.Get("/campaigns/:id", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.GetCampaign)
	tenantGroup.Put("/campaigns/:id", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.UpdateCampaign)

	tenantGroup.Get("/campaigns/:id/transactions", middleware.RequireFeature(db, constant.FeaturePGDigital), financeController.ListTransactions)

	publicGroup := api.Group("/public/:hostname")
	publicGroup.Get("/events", communityController.ListPublicEvents)
	publicGroup.Get("/gallery/albums", communityController.ListPublicGalleryAlbums)
	publicGroup.Get("/gallery/items", communityController.ListPublicGalleryItems)
	publicGroup.Get("/management-members", communityController.ListPublicManagementMembers)
	publicGroup.Get("/static-payment-methods", engagementController.ListPublicStaticPaymentMethods)
	publicGroup.Get("/social-links", engagementController.ListPublicSocialLinks)
	publicGroup.Get("/external-links", engagementController.ListPublicExternalLinks)

	publicGroup.Get("/campaigns", financeController.ListPublicCampaigns)
	publicGroup.Get("/campaigns/:slug", financeController.GetPublicCampaignBySlug)
	publicGroup.Get("/campaigns/:id/donors", financeController.ListPublicDonors)
	publicGroup.Post("/campaigns/:id/donate", financeController.CreateDonation)

	// ==========================================
	// WEBHOOK ROUTES (NO AUTH, NO HOSTNAME)
	// ==========================================
	webhookGroup := api.Group("/webhook")
	webhookGroup.Post("/midtrans", financeController.MidtransWebhook)
}
