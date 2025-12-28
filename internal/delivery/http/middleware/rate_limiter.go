package middleware

import (
	"time"

	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Max        int
	Expiration time.Duration
}

// NewRateLimiter creates a new rate limiting middleware
func NewRateLimiter(cfg RateLimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.Max,
		Expiration: cfg.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use X-Forwarded-For header if behind a proxy, otherwise use IP
			forwarded := c.Get("X-Forwarded-For")
			if forwarded != "" {
				return forwarded
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				dto.NewErrorResponse("Too many requests. Please try again later."),
			)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// DefaultRateLimiter creates a rate limiter with default settings (10 req/min)
func DefaultRateLimiter() fiber.Handler {
	return NewRateLimiter(RateLimiterConfig{
		Max:        10,
		Expiration: 1 * time.Minute,
	})
}

// StrictRateLimiter creates a stricter rate limiter (5 req/5 min)
func StrictRateLimiter() fiber.Handler {
	return NewRateLimiter(RateLimiterConfig{
		Max:        5,
		Expiration: 5 * time.Minute,
	})
}
