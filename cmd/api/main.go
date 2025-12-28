package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/handler"
	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/router"
	"github.com/andrianprasetya/go-mail-server/internal/infrastructure/config"
	"github.com/andrianprasetya/go-mail-server/internal/infrastructure/email"
	"github.com/andrianprasetya/go-mail-server/internal/usecase/contact"

	"github.com/gofiber/fiber/v2"
)

// Version is set at build time
var Version = "1.0.0"

func main() {
	log.Println("🚀 Starting Contact Form API...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Initialize infrastructure layer
	emailRepo := email.NewSMTPRepository(cfg)

	// Initialize use case layer
	contactUC := contact.NewContactUseCase(emailRepo)

	// Initialize delivery layer (handlers)
	contactHandler := handler.NewContactHandler(contactUC)
	healthHandler := handler.NewHealthHandler(Version)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Contact Form API",
		DisableStartupMessage: cfg.IsProduction(),
		ErrorHandler:          customErrorHandler,
	})

	// Setup router
	r := router.NewRouter(app, cfg, contactHandler, healthHandler)
	r.Setup()

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("🛑 Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("❌ Error shutting down server: %v", err)
		}
	}()

	// Log startup info
	log.Printf("📧 SMTP: %s:%d", cfg.SMTPHost, cfg.SMTPPort)
	log.Printf("🔒 Rate Limit: %d requests/minute", cfg.RateLimit)
	log.Printf("🌐 Allowed Origins: %v", cfg.AllowedOrigins)
	log.Printf("✅ Server listening on port %s", cfg.AppPort)

	// Start server
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// customErrorHandler handles global errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("[Error] %v", err)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
