package service

import (
	"fmt"
	"strings"

	"github.com/golgoth31/gitcomm/internal/model"
)

// FormattingService handles formatting of commit messages
type FormattingService struct{}

// NewFormattingService creates a new formatting service
func NewFormattingService() *FormattingService {
	return &FormattingService{}
}

// Format formats a CommitMessage according to Conventional Commits specification
func (s *FormattingService) Format(message *model.CommitMessage) string {
	var parts []string

	// Format header: type(scope): subject
	header := message.Type
	if message.Scope != "" {
		header = fmt.Sprintf("%s(%s)", header, message.Scope)
	}
	header = fmt.Sprintf("%s: %s", header, message.Subject)
	parts = append(parts, header)

	// Add blank line before body if body exists
	if message.Body != "" {
		parts = append(parts, "")
		parts = append(parts, message.Body)
	}

	// Add blank line before footer if footer exists
	if message.Footer != "" {
		parts = append(parts, "")
		parts = append(parts, message.Footer)
	}

	// Note: Signoff is handled separately during commit creation
	result := strings.Join(parts, "\n")
	return result
}
