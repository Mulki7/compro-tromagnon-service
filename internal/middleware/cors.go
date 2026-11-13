package middleware

import (
	"compro-tromagnon-service/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCORS(cfg *config.Config) fiber.Handler {
	origins := "http://localhost:3000, http://127.0.0.1:3000"
	if cfg.ClientOrigin != "" && cfg.ClientOrigin != "http://localhost:3000" {
		origins += ", " + cfg.ClientOrigin
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}
