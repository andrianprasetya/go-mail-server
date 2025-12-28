package repository

import "github.com/andrianprasetya/go-mail-server/internal/domain/entity"

// EmailRepository defines the interface for email operations (Domain Layer)
// This interface is implemented by infrastructure layer (SMTP, SendGrid, etc.)
type EmailRepository interface {
	// Send sends an email based on contact information
	Send(contact *entity.Contact) error
}
