package middleware

import (
	"log"

	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/dto"

	"github.com/gofiber/fiber/v2"
)

// Recover middleware recovers from panics and returns a proper error response
func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Recovered from panic: %v", r)
				c.Status(fiber.StatusInternalServerError).JSON(
					dto.NewErrorResponse("Internal server error"),
				)
			}
		}()

		return c.Next()
	}
}
