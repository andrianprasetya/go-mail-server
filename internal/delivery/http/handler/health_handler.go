package handler

import (
	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/dto"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	version string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
	}
}

// HealthCheck returns server health status
// @Summary Health check
// @Description Returns server health status
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(dto.HealthResponse{
		Status:  "healthy",
		Service: "contact-form-api",
		Version: h.version,
	})
}

// ReadinessCheck returns server readiness status
// @Summary Readiness check
// @Description Returns server readiness status for Kubernetes
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(dto.HealthResponse{
		Status:  "ready",
		Service: "contact-form-api",
		Version: h.version,
	})
}
