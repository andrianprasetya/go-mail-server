package contact

import (
	"context"
	"log"

	"github.com/andrianprasetya/go-mail-server/internal/domain/entity"
	"github.com/andrianprasetya/go-mail-server/internal/domain/repository"
)

// contactUseCase implements the UseCase interface
type contactUseCase struct {
	emailRepo repository.EmailRepository
}

// NewContactUseCase creates a new contact use case
func NewContactUseCase(emailRepo repository.EmailRepository) UseCase {
	return &contactUseCase{
		emailRepo: emailRepo,
	}
}

// SendContact processes a contact form submission
func (uc *contactUseCase) SendContact(ctx context.Context, input *ContactInput) (*ContactOutput, error) {
	// Create and validate contact entity
	contact, err := entity.NewContact(
		input.Name,
		input.Email,
		input.Subject,
		input.Message,
	)
	if err != nil {
		log.Printf("[UseCase] Validation failed: %v", err)
		return &ContactOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Send email via repository
	if err := uc.emailRepo.Send(contact); err != nil {
		log.Printf("[UseCase] Failed to send email: %v", err)
		return &ContactOutput{
			Success: false,
			Message: "Failed to send email. Please try again later.",
		}, err
	}

	log.Printf("[UseCase] Email sent successfully from %s (%s)", contact.Name, contact.Email)
	return &ContactOutput{
		Success: true,
		Message: "Email sent successfully",
	}, nil
}
