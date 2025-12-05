package tokenization

import (
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// TikTokenCalculator implements tokenization for OpenAI using tiktoken
type TikTokenCalculator struct{}

// NewTikTokenCalculator creates a new OpenAI token calculator
func NewTikTokenCalculator() TokenCalculator {
	return &TikTokenCalculator{}
}

// Calculate estimates tokens using OpenAI's tiktoken (approximation)
// TODO: Integrate actual tiktoken library when Go bindings are available
// For now, using a character-based approximation: ~4 characters per token
func (t *TikTokenCalculator) Calculate(text string) int {
	// Rough approximation: OpenAI tokens are typically ~4 characters
	// This is a simplified version - actual tiktoken would be more accurate
	return len(text) / 4
}

// CalculateForRepositoryState estimates tokens for repository state
func (t *TikTokenCalculator) CalculateForRepositoryState(state *model.RepositoryState) (int, error) {
	var text string
	for _, file := range state.StagedFiles {
		utils.Logger.Debug().Msgf("Staged file: %+v", file)
		text += file.Path + " " + file.Status + " " + file.Diff + "\n"
	}
	for _, file := range state.UnstagedFiles {
		utils.Logger.Debug().Msgf("Unstaged file: %+v", file)
		text += file.Path + " " + file.Status + " " + file.Diff + "\n"
	}
	return t.Calculate(text), nil
}
