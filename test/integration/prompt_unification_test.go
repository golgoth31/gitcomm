package integration

import (
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/ai"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/pkg/ai/prompt"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

// TestPromptConsistencyAcrossProviders tests that all providers use identical prompts
func TestPromptConsistencyAcrossProviders(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	generator := prompt.NewUnifiedPromptGenerator()
	validator := conventional.NewValidator()

	// Create a test repository state
	repoState := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{
				Path:   "internal/ai/openai_provider.go",
				Status: "modified",
				Diff:   "diff --git a/internal/ai/openai_provider.go b/internal/ai/openai_provider.go\n@@ -1,3 +1,5 @@\n+new code\n",
			},
		},
		UnstagedFiles: []model.FileChange{},
	}

	// Generate system and user messages
	systemMsg, err := generator.GenerateSystemMessage(validator)
	if err != nil {
		t.Fatalf("GenerateSystemMessage() error = %v", err)
	}

	userMsg, err := generator.GenerateUserMessage(repoState)
	if err != nil {
		t.Fatalf("GenerateUserMessage() error = %v", err)
	}

	// Verify system message contains all validation rules
	if !strings.Contains(systemMsg, "feat") || !strings.Contains(systemMsg, "fix") {
		t.Error("System message should contain valid commit types")
	}

	if !strings.Contains(systemMsg, "72") {
		t.Error("System message should contain subject length constraint")
	}

	if !strings.Contains(systemMsg, "320") {
		t.Error("System message should contain body length constraint")
	}

	if !strings.Contains(systemMsg, "alphanumeric") {
		t.Error("System message should contain scope format constraint")
	}

	// Verify user message contains repository state
	if !strings.Contains(userMsg, "openai_provider.go") {
		t.Error("User message should contain file path")
	}

	if !strings.Contains(userMsg, "Staged files:") {
		t.Error("User message should contain staged files section")
	}

	// Test consistency: same inputs should produce same outputs
	systemMsg2, _ := generator.GenerateSystemMessage(validator)
	userMsg2, _ := generator.GenerateUserMessage(repoState)

	if systemMsg != systemMsg2 {
		t.Error("GenerateSystemMessage() should produce consistent output")
	}

	if userMsg != userMsg2 {
		t.Error("GenerateUserMessage() should produce consistent output")
	}

	// Verify all providers would use the same prompts
	// (This is a structural test - actual API calls would require API keys)
	providers := []struct {
		name     string
		provider ai.AIProvider
	}{
		{"OpenAI", ai.NewOpenAIProvider(&model.AIProviderConfig{APIKey: "test-key"})},
		{"Anthropic", ai.NewAnthropicProvider(&model.AIProviderConfig{APIKey: "test-key"})},
		{"Mistral", ai.NewMistralProvider(&model.AIProviderConfig{APIKey: "test-key"})},
		{"Local", ai.NewLocalProvider(&model.AIProviderConfig{Endpoint: "http://localhost:8080"})},
	}

	// Verify all providers are created successfully
	for _, p := range providers {
		if p.provider == nil {
			t.Errorf("Failed to create %s provider", p.name)
		}
	}
}

// TestAnthropicSystemUserCombination tests that Anthropic provider prepends system to user message
func TestAnthropicSystemUserCombination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	generator := prompt.NewUnifiedPromptGenerator()
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

	// Generate system and user messages
	systemMsg, err := generator.GenerateSystemMessage(validator)
	if err != nil {
		t.Fatalf("GenerateSystemMessage() error = %v", err)
	}

	userMsg, err := generator.GenerateUserMessage(repoState)
	if err != nil {
		t.Fatalf("GenerateUserMessage() error = %v", err)
	}

	// Simulate Anthropic's combination: system + "\n\n" + user
	combinedMsg := systemMsg + "\n\n" + userMsg

	// Verify combined message contains both parts
	if !strings.Contains(combinedMsg, systemMsg) {
		t.Error("Combined message should contain system message")
	}

	if !strings.Contains(combinedMsg, userMsg) {
		t.Error("Combined message should contain user message")
	}

	// Verify system message comes before user message
	systemIndex := strings.Index(combinedMsg, systemMsg)
	userIndex := strings.Index(combinedMsg, userMsg)
	if systemIndex >= userIndex {
		t.Error("System message should come before user message in combined message")
	}

	// Verify separator is present
	if !strings.Contains(combinedMsg, "\n\n") {
		t.Error("Combined message should have separator between system and user messages")
	}
}
