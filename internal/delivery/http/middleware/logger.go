package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger logs incoming requests with timing information
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get real IP (supports proxy)
		ip := c.Get("X-Real-IP")
		if ip == "" {
			ip = c.Get("X-Forwarded-For")
		}
		if ip == "" {
			ip = c.IP()
		}

		// Log request details
		log.Printf(
			"[HTTP] %s %s | %d | %v | IP: %s | UA: %s",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
			ip,
			c.Get("User-Agent"),
		)

		return err
	}
}
