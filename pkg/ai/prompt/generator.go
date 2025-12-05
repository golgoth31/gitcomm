package prompt

import (
	"fmt"
	"strings"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

// PromptGenerator defines the interface for generating unified AI prompts
type PromptGenerator interface {
	// GenerateSystemMessage generates the system message with validation rules
	// extracted from the MessageValidator
	GenerateSystemMessage(validator conventional.MessageValidator) (string, error)

	// GenerateUserMessage generates the user message with repository state
	// formatted for AI consumption
	GenerateUserMessage(repoState *model.RepositoryState) (string, error)
}

// UnifiedPromptGenerator implements PromptGenerator for unified prompt generation
type UnifiedPromptGenerator struct{}

// NewUnifiedPromptGenerator creates a new unified prompt generator
func NewUnifiedPromptGenerator() PromptGenerator {
	return &UnifiedPromptGenerator{}
}

// GenerateSystemMessage generates the system message with validation rules
func (g *UnifiedPromptGenerator) GenerateSystemMessage(validator conventional.MessageValidator) (string, error) {
	if validator == nil {
		return "", ErrNilValidator
	}

	// Extract validation rules from validator
	validTypes := validator.GetValidTypes()
	subjectMaxLength := validator.GetSubjectMaxLength()
	bodyMaxLength := validator.GetBodyMaxLength()
	scopeFormatDesc := validator.GetScopeFormatDescription()

	// Build system message with structured bullet points
	var sb strings.Builder

	sb.WriteString("You are a git commit message generator. When receiving a git diff, you will ONLY generate commit messages following the Conventional Commits specification.\n\n")
	sb.WriteString("Format: type(scope): subject\n\nbody\n\nfooter\n\n")
	sb.WriteString("Validation Rules:\n")

	// Type constraint
	sb.WriteString(fmt.Sprintf("• Type must be one of: %s\n", strings.Join(validTypes, ", ")))

	// Subject length constraint
	sb.WriteString(fmt.Sprintf("• Subject must be ≤%d characters\n", subjectMaxLength))

	// Body length constraint
	sb.WriteString(fmt.Sprintf("• Body must be ≤%d characters (if provided)\n", bodyMaxLength))

	// Scope format constraint
	sb.WriteString(fmt.Sprintf("• Scope must be a valid identifier (%s)\n", scopeFormatDesc))

	return sb.String(), nil
}

// GenerateUserMessage generates the user message with repository state
func (g *UnifiedPromptGenerator) GenerateUserMessage(repoState *model.RepositoryState) (string, error) {
	if repoState == nil {
		return "", ErrNilRepositoryState
	}

	var sb strings.Builder

	sb.WriteString("Generate a commit message for the following changes:\n\n")

	// Add staged files
	if len(repoState.StagedFiles) > 0 {
		sb.WriteString("Staged files:\n")
		for _, file := range repoState.StagedFiles {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", file.Path, file.Status))
			if file.Diff != "" {
				sb.WriteString(file.Diff)
				if !strings.HasSuffix(file.Diff, "\n") {
					sb.WriteString("\n")
				}
			}
		}
	}

	// Add unstaged files
	if len(repoState.UnstagedFiles) > 0 {
		if len(repoState.StagedFiles) > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("Unstaged files:\n")
		for _, file := range repoState.UnstagedFiles {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", file.Path, file.Status))
			if file.Diff != "" {
				sb.WriteString(file.Diff)
				if !strings.HasSuffix(file.Diff, "\n") {
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String(), nil
}
