package prompt

import (
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

func TestPromptGenerator_GenerateSystemMessage(t *testing.T) {
	generator := NewUnifiedPromptGenerator()
	validator := conventional.NewValidator()

	t.Run("valid validator", func(t *testing.T) {
		systemMsg, err := generator.GenerateSystemMessage(validator)
		if err != nil {
			t.Fatalf("GenerateSystemMessage() error = %v, want nil", err)
		}

		// Verify system message contains validation rules
		if !strings.Contains(systemMsg, "feat") || !strings.Contains(systemMsg, "fix") {
			t.Error("GenerateSystemMessage() should contain valid commit types")
		}

		if !strings.Contains(systemMsg, "72") {
			t.Error("GenerateSystemMessage() should contain subject length constraint")
		}

		if !strings.Contains(systemMsg, "320") {
			t.Error("GenerateSystemMessage() should contain body length constraint")
		}

		if !strings.Contains(systemMsg, "alphanumeric") {
			t.Error("GenerateSystemMessage() should contain scope format constraint")
		}

		// Verify format specification
		if !strings.Contains(systemMsg, "type(scope): subject") {
			t.Error("GenerateSystemMessage() should contain Conventional Commits format")
		}
	})

	t.Run("nil validator", func(t *testing.T) {
		systemMsg, err := generator.GenerateSystemMessage(nil)
		if err == nil {
			t.Error("GenerateSystemMessage() with nil validator should return error")
		}
		if systemMsg != "" {
			t.Error("GenerateSystemMessage() with nil validator should return empty string")
		}
		if err != ErrNilValidator {
			t.Errorf("GenerateSystemMessage() error = %v, want ErrNilValidator", err)
		}
	})
}

func TestPromptGenerator_GenerateUserMessage(t *testing.T) {
	generator := NewUnifiedPromptGenerator()

	t.Run("valid repository state", func(t *testing.T) {
		repoState := &model.RepositoryState{
			StagedFiles: []model.FileChange{
				{
					Path:   "internal/ai/openai_provider.go",
					Status: "modified",
					Diff:   "diff --git a/internal/ai/openai_provider.go b/internal/ai/openai_provider.go\n@@ -1,3 +1,5 @@\n+new code\n",
				},
			},
			UnstagedFiles: []model.FileChange{
				{
					Path:   "test/integration/prompt_test.go",
					Status: "added",
					Diff:   "diff --git a/test/integration/prompt_test.go b/test/integration/prompt_test.go\n@@ -0,0 +1,3 @@\n+test code\n",
				},
			},
		}

		userMsg, err := generator.GenerateUserMessage(repoState)
		if err != nil {
			t.Fatalf("GenerateUserMessage() error = %v, want nil", err)
		}

		// Verify user message contains repository state
		if !strings.Contains(userMsg, "openai_provider.go") {
			t.Error("GenerateUserMessage() should contain staged file path")
		}

		if !strings.Contains(userMsg, "prompt_test.go") {
			t.Error("GenerateUserMessage() should contain unstaged file path")
		}

		if !strings.Contains(userMsg, "Staged files:") {
			t.Error("GenerateUserMessage() should contain staged files section")
		}

		if !strings.Contains(userMsg, "Unstaged files:") {
			t.Error("GenerateUserMessage() should contain unstaged files section")
		}
	})

	t.Run("nil repository state", func(t *testing.T) {
		userMsg, err := generator.GenerateUserMessage(nil)
		if err == nil {
			t.Error("GenerateUserMessage() with nil repository state should return error")
		}
		if userMsg != "" {
			t.Error("GenerateUserMessage() with nil repository state should return empty string")
		}
		if err != ErrNilRepositoryState {
			t.Errorf("GenerateUserMessage() error = %v, want ErrNilRepositoryState", err)
		}
	})

	t.Run("empty repository state", func(t *testing.T) {
		repoState := &model.RepositoryState{
			StagedFiles:   []model.FileChange{},
			UnstagedFiles: []model.FileChange{},
		}

		userMsg, err := generator.GenerateUserMessage(repoState)
		if err != nil {
			t.Fatalf("GenerateUserMessage() with empty state error = %v, want nil", err)
		}

		// Should still generate a message, just with no files
		if userMsg == "" {
			t.Error("GenerateUserMessage() with empty state should return non-empty message")
		}
	})
}

func TestPromptGenerator_Consistency(t *testing.T) {
	generator := NewUnifiedPromptGenerator()
	validator := conventional.NewValidator()
	repoState := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{
				Path:   "test.go",
				Status: "modified",
				Diff:   "test diff",
			},
		},
	}

	// Generate messages multiple times
	systemMsg1, err1 := generator.GenerateSystemMessage(validator)
	if err1 != nil {
		t.Fatalf("First GenerateSystemMessage() error = %v", err1)
	}

	systemMsg2, err2 := generator.GenerateSystemMessage(validator)
	if err2 != nil {
		t.Fatalf("Second GenerateSystemMessage() error = %v", err2)
	}

	if systemMsg1 != systemMsg2 {
		t.Error("GenerateSystemMessage() should produce consistent output for same input")
	}

	userMsg1, err1 := generator.GenerateUserMessage(repoState)
	if err1 != nil {
		t.Fatalf("First GenerateUserMessage() error = %v", err1)
	}

	userMsg2, err2 := generator.GenerateUserMessage(repoState)
	if err2 != nil {
		t.Fatalf("Second GenerateUserMessage() error = %v", err2)
	}

	if userMsg1 != userMsg2 {
		t.Error("GenerateUserMessage() should produce consistent output for same input")
	}
}
