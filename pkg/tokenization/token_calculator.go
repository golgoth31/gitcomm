package tokenization

import (
	"github.com/golgoth31/gitcomm/internal/model"
)

// TokenCalculator defines the interface for token calculation
type TokenCalculator interface {
	// Calculate estimates tokens for a given text
	Calculate(text string) int

	// CalculateForRepositoryState estimates tokens for repository state
	CalculateForRepositoryState(state *model.RepositoryState) (int, error)
}

// NewTokenCalculator creates a token calculator for the specified provider
func NewTokenCalculator(provider string) TokenCalculator {
	switch provider {
	case "openai":
		return NewTikTokenCalculator()
	case "anthropic":
		return NewAnthropicTokenCalculator()
	case "mistral":
		return NewFallbackTokenCalculator()
	default:
		return NewFallbackTokenCalculator()
	}
}
