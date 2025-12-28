package entity

import (
	"errors"
	"regexp"
	"strings"
)

// Contact represents the contact form entity (Domain Entity)
type Contact struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// Validation constants
const (
	MaxMessageLength = 500
	MaxNameLength    = 100
	MaxSubjectLength = 200
	MaxEmailLength   = 254
)

// Validation errors
var (
	ErrNameRequired    = errors.New("name is required")
	ErrNameTooLong     = errors.New("name must be less than 100 characters")
	ErrEmailRequired   = errors.New("email is required")
	ErrEmailInvalid    = errors.New("email format is invalid")
	ErrEmailTooLong    = errors.New("email must be less than 254 characters")
	ErrSubjectRequired = errors.New("subject is required")
	ErrSubjectTooLong  = errors.New("subject must be less than 200 characters")
	ErrMessageRequired = errors.New("message is required")
	ErrMessageTooLong  = errors.New("message must be less than 500 characters")
)

// Email regex pattern (RFC 5322 simplified)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewContact creates a new Contact entity with validation
func NewContact(name, email, subject, message string) (*Contact, error) {
	contact := &Contact{
		Name:    strings.TrimSpace(name),
		Email:   strings.TrimSpace(email),
		Subject: strings.TrimSpace(subject),
		Message: strings.TrimSpace(message),
	}

	if err := contact.Validate(); err != nil {
		return nil, err
	}

	return contact, nil
}

// Validate validates the contact entity
func (c *Contact) Validate() error {
	// Validate name
	if c.Name == "" {
		return ErrNameRequired
	}
	if len(c.Name) > MaxNameLength {
		return ErrNameTooLong
	}

	// Validate email
	if c.Email == "" {
		return ErrEmailRequired
	}
	if len(c.Email) > MaxEmailLength {
		return ErrEmailTooLong
	}
	if !emailRegex.MatchString(c.Email) {
		return ErrEmailInvalid
	}

	// Validate subject
	if c.Subject == "" {
		return ErrSubjectRequired
	}
	if len(c.Subject) > MaxSubjectLength {
		return ErrSubjectTooLong
	}

	// Validate message
	if c.Message == "" {
		return ErrMessageRequired
	}
	if len(c.Message) > MaxMessageLength {
		return ErrMessageTooLong
	}

	return nil
}
