package tokenization

import (
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
)

func TestCalculateTokens(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		text     string
		wantMin  int
		wantMax  int
	}{
		{
			name:     "OpenAI short text",
			provider: "openai",
			text:     "Hello world",
			wantMin:  1,
			wantMax:  10,
		},
		{
			name:     "Anthropic short text",
			provider: "anthropic",
			text:     "Hello world",
			wantMin:  1,
			wantMax:  10,
		},
		{
			name:     "Fallback short text",
			provider: "unknown",
			text:     "Hello world",
			wantMin:  1,
			wantMax:  10,
		},
		{
			name:     "OpenAI long text",
			provider: "openai",
			text:     string(make([]byte, 1000)),
			wantMin:  100,
			wantMax:  500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewTokenCalculator(tt.provider)
			got := calc.Calculate(tt.text)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("CalculateTokens() = %d, want between %d and %d", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestTokenCalculator_CalculateForRepositoryState(t *testing.T) {
	calc := NewTokenCalculator("openai")

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "file1.go", Status: "modified", Diff: "diff content here"},
			{Path: "file2.go", Status: "added", Diff: "new file content"},
		},
		UnstagedFiles: []model.FileChange{
			{Path: "file3.go", Status: "modified", Diff: "more diff"},
		},
	}

	got, err := calc.CalculateForRepositoryState(state)
	if err != nil {
		t.Errorf("CalculateForRepositoryState() = %v, want nil", err)
	}
	if got <= 0 {
		t.Errorf("CalculateForRepositoryState() = %d, want > 0", got)
	}
}
