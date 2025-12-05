package service

import (
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

// ValidationService handles validation of commit messages
type ValidationService struct {
	validator conventional.MessageValidator
}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{
		validator: conventional.NewValidator(),
	}
}

// Validate validates a CommitMessage against Conventional Commits specification
func (s *ValidationService) Validate(message *model.CommitMessage) (bool, []conventional.ValidationError) {
	return s.validator.Validate(message)
}
