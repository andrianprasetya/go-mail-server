package router

import (
	"time"

	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/handler"
	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/middleware"
	"github.com/andrianprasetya/go-mail-server/internal/infrastructure/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Router holds all handlers and configuration for routing
type Router struct {
	app            *fiber.App
	config         *config.Config
	contactHandler *handler.ContactHandler
	healthHandler  *handler.HealthHandler
}

// NewRouter creates a new router with all handlers
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	contactHandler *handler.ContactHandler,
	healthHandler *handler.HealthHandler,
) *Router {
	return &Router{
		app:            app,
		config:         cfg,
		contactHandler: contactHandler,
		healthHandler:  healthHandler,
	}
}

// Setup configures all routes and middleware
func (r *Router) Setup() {
	// Global middleware
	r.app.Use(middleware.Recover())
	r.app.Use(middleware.RequestLogger())

	// CORS configuration (skip in development for easier testing)
	if r.config.AppEnv != "development" {
		r.app.Use(cors.New(cors.Config{
			AllowOrigins:     stringSliceToCSV(r.config.AllowedOrigins),
			AllowMethods:     "GET,POST,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: false,
			MaxAge:           86400, // 24 hours
		}))
	}

	// Health check endpoints (no rate limiting)
	r.app.Get("/health", r.healthHandler.HealthCheck)
	r.app.Get("/ready", r.healthHandler.ReadinessCheck)

	// API routes
	api := r.app.Group("/api")

	// Contact endpoint with rate limiting (default: 2 requests per 24 hours per IP)
	contactLimiter := middleware.NewRateLimiter(middleware.RateLimiterConfig{
		Max:        r.config.RateLimit,
		Expiration: time.Duration(r.config.RateLimitExpiration) * time.Hour,
	})
	api.Post("/contact", contactLimiter, r.contactHandler.HandleContact)
}

// stringSliceToCSV converts a string slice to comma-separated string
func stringSliceToCSV(slice []string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}
