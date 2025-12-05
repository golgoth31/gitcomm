package tokenization

import (
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// FallbackTokenCalculator implements character-based token estimation
type FallbackTokenCalculator struct{}

// NewFallbackTokenCalculator creates a new fallback token calculator
func NewFallbackTokenCalculator() TokenCalculator {
	return &FallbackTokenCalculator{}
}

// Calculate estimates tokens using character-based fallback
// Uses the common approximation: ~4 characters per token
func (f *FallbackTokenCalculator) Calculate(text string) int {
	return len(text) / 4
}

// CalculateForRepositoryState estimates tokens for repository state
func (f *FallbackTokenCalculator) CalculateForRepositoryState(state *model.RepositoryState) (int, error) {
	var text string
	for _, file := range state.StagedFiles {
		utils.Logger.Debug().Msgf("Staged file: %s", file.Diff)
		text += file.Path + " " + file.Status + " " + file.Diff + "\n"
	}
	for _, file := range state.UnstagedFiles {
		utils.Logger.Debug().Msgf("Unstaged file: %s", file.Diff)
		text += file.Path + " " + file.Status + " " + file.Diff + "\n"
	}
	return f.Calculate(text), nil
}
