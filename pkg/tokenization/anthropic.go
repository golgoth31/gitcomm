package tokenization

import (
	"github.com/golgoth31/gitcomm/internal/model"
)

// AnthropicTokenCalculator implements tokenization for Anthropic
type AnthropicTokenCalculator struct{}

// NewAnthropicTokenCalculator creates a new Anthropic token calculator
func NewAnthropicTokenCalculator() TokenCalculator {
	return &AnthropicTokenCalculator{}
}

// Calculate estimates tokens for Anthropic (custom tokenization)
// Anthropic uses a different tokenization method than OpenAI
// Approximation: ~3.5 characters per token
func (a *AnthropicTokenCalculator) Calculate(text string) int {
	// Rough approximation for Anthropic tokenization
	return len(text) * 10 / 35 // Approximately 3.5 chars per token
}

// CalculateForRepositoryState estimates tokens for repository state
func (a *AnthropicTokenCalculator) CalculateForRepositoryState(state *model.RepositoryState) (int, error) {
	var text string
	for _, file := range state.StagedFiles {
		text += file.Path + " " + file.Status + " " + file.Diff + "\n"
	}
	for _, file := range state.UnstagedFiles {
		text += file.Path + " " + file.Status + " " + file.Diff + "\n"
	}
	return a.Calculate(text), nil
}
