package conventional

import "github.com/golgoth31/gitcomm/internal/model"

// MessageValidator defines the interface for validating Conventional Commits messages
type MessageValidator interface {
	// Validate validates a CommitMessage against the Conventional Commits specification
	Validate(message *model.CommitMessage) (bool, []ValidationError)
	// GetValidTypes returns the list of valid commit types
	GetValidTypes() []string
	// GetSubjectMaxLength returns the maximum allowed length for commit message subject
	GetSubjectMaxLength() int
	// GetBodyMaxLength returns the maximum allowed length for commit message body
	GetBodyMaxLength() int
	// GetScopeFormatDescription returns a human-readable description of valid scope format
	GetScopeFormatDescription() string
}

// ValidationError represents a validation error with a specific field and message
type ValidationError struct {
	Field   string
	Message string
}

// Validator implements MessageValidator
type Validator struct{}

// NewValidator creates a new Conventional Commits validator
func NewValidator() MessageValidator {
	return &Validator{}
}

// Validate validates a CommitMessage against the Conventional Commits specification
func (v *Validator) Validate(message *model.CommitMessage) (bool, []ValidationError) {
	var errors []ValidationError

	// Validate type
	if !isValidType(message.Type) {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: "type must be one of: feat, fix, docs, style, refactor, test, chore, version",
		})
	}

	// Validate subject
	if message.Subject == "" {
		errors = append(errors, ValidationError{
			Field:   "subject",
			Message: "subject cannot be empty",
		})
	} else if len(message.Subject) > 72 {
		errors = append(errors, ValidationError{
			Field:   "subject",
			Message: "subject must be ≤72 characters",
		})
	}

	// Validate body
	if message.Body != "" && len(message.Body) > 320 {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "body must be ≤320 characters",
		})
	}

	// Validate scope (if provided)
	if message.Scope != "" && !isValidScope(message.Scope) {
		errors = append(errors, ValidationError{
			Field:   "scope",
			Message: "scope must be a valid identifier (alphanumeric, hyphens, underscores)",
		})
	}

	return len(errors) == 0, errors
}

// isValidType checks if the type is a valid Conventional Commits type
func isValidType(t string) bool {
	validTypes := []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "version"}
	for _, vt := range validTypes {
		if t == vt {
			return true
		}
	}
	return false
}

// isValidScope checks if the scope is a valid identifier
func isValidScope(scope string) bool {
	if scope == "" {
		return true
	}
	// Simple validation: alphanumeric, hyphens, underscores
	for _, r := range scope {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// GetValidTypes returns the list of valid commit types
func (v *Validator) GetValidTypes() []string {
	return []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "version"}
}

// GetSubjectMaxLength returns the maximum allowed length for commit message subject
func (v *Validator) GetSubjectMaxLength() int {
	return 72
}

// GetBodyMaxLength returns the maximum allowed length for commit message body
func (v *Validator) GetBodyMaxLength() int {
	return 320
}

// GetScopeFormatDescription returns a human-readable description of valid scope format
func (v *Validator) GetScopeFormatDescription() string {
	return "alphanumeric, hyphens, underscores only"
}
