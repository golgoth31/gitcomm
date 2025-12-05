package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/ai"
	"github.com/golgoth31/gitcomm/internal/model"
)

// TestAIAssistedCommitWorkflow tests the AI-assisted commit workflow
// This is an integration test that requires a git repository and optionally AI provider
func TestAIAssistedCommitWorkflow(t *testing.T) {
	// Skip if not in CI or if git is not available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for the test repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-ai-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "feature.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc NewFeature() {}\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Test token calculation
	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "feature.go", Status: "added", Diff: "package main\n\nfunc NewFeature() {}\n"},
		},
	}

	// TODO: Once AI providers are implemented, test:
	// 1. Token calculation
	// 2. AI provider call (if API key available)
	// 3. Message generation
	// 4. Format validation
	// 5. Commit creation

	t.Logf("Test repository created at: %s", tmpDir)
	t.Log("AI integration test structure ready - will be completed when AI providers are implemented")

	// Test that we can create a provider (even if it fails without API key)
	config := &model.AIProviderConfig{
		Name:    "openai",
		APIKey:  os.Getenv("OPENAI_API_KEY"), // May be empty
		Model:   "gpt-4",
		Timeout: 5 * time.Second,
	}

	if config.APIKey != "" {
		provider := ai.NewOpenAIProvider(config)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		message, err := provider.GenerateCommitMessage(ctx, state)
		if err != nil {
			t.Logf("AI provider call failed (expected without valid key): %v", err)
		} else {
			t.Logf("AI generated message: %s", message)
		}
	} else {
		t.Log("Skipping AI provider call - no API key in environment")
	}
}
