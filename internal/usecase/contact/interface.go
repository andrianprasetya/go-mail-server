package contact

import "context"

// ContactInput represents the input for contact use case
type ContactInput struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// ContactOutput represents the output of contact use case
type ContactOutput struct {
	Success bool
	Message string
}

// UseCase defines the contact use case interface
type UseCase interface {
	// SendContact processes a contact form submission
	SendContact(ctx context.Context, input *ContactInput) (*ContactOutput, error)
}
