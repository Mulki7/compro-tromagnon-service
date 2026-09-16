package main

import (
	"log"

	"compro-tromagnon-service/internal/artist"
	"compro-tromagnon-service/internal/auth"
	"compro-tromagnon-service/internal/config"
	"compro-tromagnon-service/internal/database"
	"compro-tromagnon-service/internal/hero"
	"compro-tromagnon-service/internal/live"
	"compro-tromagnon-service/internal/middleware"
	"compro-tromagnon-service/internal/news"
	"compro-tromagnon-service/internal/release"
	"compro-tromagnon-service/internal/upload"
	"compro-tromagnon-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Load Configurations
	cfg := config.LoadConfig()

	// 2. Initialize Database & Migrations
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}

	// 3. Setup Fiber Application
	app := fiber.New(fiber.Config{
		AppName: "Tromagnon Records API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return response.Error(c, code, err.Error())
		},
	})

	// 4. Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(middleware.NewCORS(cfg))

	// Serve static files from uploads folder
	app.Static("/uploads", "./uploads")

	// 5. Initialize Repositories
	authRepo := auth.NewRepository(db)
	artistRepo := artist.NewRepository(db)
	releaseRepo := release.NewRepository(db)
	newsRepo := news.NewRepository(db)
	liveRepo := live.NewRepository(db)
	heroRepo := hero.NewRepository(db)

	// 6. Initialize Services
	authService := auth.NewService(authRepo, cfg)
	artistService := artist.NewService(artistRepo)
	releaseService := release.NewService(releaseRepo)
	newsService := news.NewService(newsRepo)
	liveService := live.NewService(liveRepo)
	heroService := hero.NewService(heroRepo)

	// 7. Initialize Handlers
	authHandler := auth.NewHandler(authService)
	artistHandler := artist.NewHandler(artistService)
	releaseHandler := release.NewHandler(releaseService)
	newsHandler := news.NewHandler(newsService)
	liveHandler := live.NewHandler(liveService)
	heroHandler := hero.NewHandler(heroService)
	uploadHandler := upload.NewHandler(cfg, "general")

	// 8. Register Routes
	api := app.Group("/api")

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, fiber.Map{
			"status":  "healthy",
			"service": "Tromagnon Records API",
		}, "API is running smoothly")
	})

	// ── Public Routes ────────────────────────────────────────────────
	api.Get("/artists", artistHandler.GetPublicArtists)
	api.Get("/artists/:slug", artistHandler.GetPublicArtistBySlug)

	api.Get("/releases", releaseHandler.GetPublicReleases)
	api.Get("/releases/:slug", releaseHandler.GetPublicReleaseBySlug)

	api.Get("/news", newsHandler.GetPublicNews)
	api.Get("/news/:slug", newsHandler.GetPublicNewsBySlug)

	api.Get("/live", liveHandler.GetPublicLiveEvents)

	// Hero banners — public (untuk slideshow homepage)
	api.Get("/hero-banners", heroHandler.GetBanners)

	// ── Admin Routes ─────────────────────────────────────────────────
	admin := api.Group("/admin")

	// Auth (tidak perlu token)
	admin.Post("/login", authHandler.Login)

	// Protected Admin Routes
	adminProtected := admin.Group("", middleware.JWTAuth(cfg))
	adminProtected.Get("/me", authHandler.GetMe)

	// Admin Artists
	adminProtected.Get("/artists", artistHandler.GetAdminArtists)
	adminProtected.Get("/artists/:id", artistHandler.GetAdminArtistByID)
	adminProtected.Post("/artists", artistHandler.CreateArtist)
	adminProtected.Put("/artists/:id", artistHandler.UpdateArtist)
	adminProtected.Patch("/artists/:id", artistHandler.UpdateArtist)
	adminProtected.Delete("/artists/:id", artistHandler.DeleteArtist)

	// Admin Releases
	adminProtected.Get("/releases", releaseHandler.GetAdminReleases)
	adminProtected.Get("/releases/:id", releaseHandler.GetAdminReleaseByID)
	adminProtected.Post("/releases", releaseHandler.CreateRelease)
	adminProtected.Put("/releases/:id", releaseHandler.UpdateRelease)
	adminProtected.Patch("/releases/:id", releaseHandler.UpdateRelease)
	adminProtected.Delete("/releases/:id", releaseHandler.DeleteRelease)

	// Admin News
	adminProtected.Get("/news", newsHandler.GetAdminNews)
	adminProtected.Get("/news/:id", newsHandler.GetAdminNewsByID)
	adminProtected.Post("/news", newsHandler.CreateNews)
	adminProtected.Put("/news/:id", newsHandler.UpdateNews)
	adminProtected.Patch("/news/:id", newsHandler.UpdateNews)
	adminProtected.Delete("/news/:id", newsHandler.DeleteNews)

	// Admin Live Events
	adminProtected.Get("/live", liveHandler.GetAdminLiveEvents)
	adminProtected.Get("/live/:id", liveHandler.GetAdminLiveEventByID)
	adminProtected.Post("/live", liveHandler.CreateLiveEvent)
	adminProtected.Put("/live/:id", liveHandler.UpdateLiveEvent)
	adminProtected.Patch("/live/:id", liveHandler.UpdateLiveEvent)
	adminProtected.Delete("/live/:id", liveHandler.DeleteLiveEvent)

	// Admin Hero Banners (slideshow homepage)
	adminProtected.Get("/hero-banners", heroHandler.GetBanners)
	adminProtected.Get("/hero-banners/:id", heroHandler.GetBannerByID)
	adminProtected.Post("/hero-banners", heroHandler.CreateBanner)
	adminProtected.Put("/hero-banners/:id", heroHandler.UpdateBanner)
	adminProtected.Patch("/hero-banners/:id", heroHandler.UpdateBanner)
	adminProtected.Delete("/hero-banners/:id", heroHandler.DeleteBanner)

	// Admin File Upload (utility)
	adminProtected.Post("/upload", uploadHandler.UploadFile)

	// 9. Start Server
	listenAddr := ":" + cfg.Port
	log.Printf("Server is starting on http://localhost%s\n", listenAddr)
	if err := app.Listen(listenAddr); err != nil {
		log.Fatalf("Server shutdown with error: %v", err)
	}
}
