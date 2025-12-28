package handler

import (
	"log"

	"github.com/andrianprasetya/go-mail-server/internal/delivery/http/dto"
	"github.com/andrianprasetya/go-mail-server/internal/usecase/contact"

	"github.com/gofiber/fiber/v2"
)

// ContactHandler handles contact form HTTP requests
type ContactHandler struct {
	contactUC contact.UseCase
}

// NewContactHandler creates a new contact handler
func NewContactHandler(contactUC contact.UseCase) *ContactHandler {
	return &ContactHandler{
		contactUC: contactUC,
	}
}

// HandleContact processes contact form submissions
// @Summary Submit contact form
// @Description Receives contact form data and sends email to site owner
// @Tags contact
// @Accept json
// @Produce json
// @Param request body dto.ContactRequest true "Contact form data"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /api/contact [post]
func (h *ContactHandler) HandleContact(c *fiber.Ctx) error {
	// Parse request body
	var request dto.ContactRequest
	if err := c.BodyParser(&request); err != nil {
		log.Printf("[Handler] Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.NewErrorResponse("Invalid request body"),
		)
	}

	// Create use case input
	input := &contact.ContactInput{
		Name:    request.Name,
		Email:   request.Email,
		Subject: request.Subject,
		Message: request.Message,
	}

	// Execute use case
	output, err := h.contactUC.SendContact(c.Context(), input)
	if err != nil {
		// Check if it's a validation error (client error)
		if !output.Success && output.Message != "Failed to send email. Please try again later." {
			return c.Status(fiber.StatusBadRequest).JSON(
				dto.NewErrorResponse(output.Message),
			)
		}
		// Server error
		return c.Status(fiber.StatusInternalServerError).JSON(
			dto.NewErrorResponse(output.Message),
		)
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(
		dto.NewSuccessResponse(output.Message),
	)
}
