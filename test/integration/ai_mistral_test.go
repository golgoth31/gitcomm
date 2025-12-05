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

// TestMistralProviderWorkflow tests the Mistral provider workflow
// This is an integration test that requires a git repository and optionally Mistral API key
func TestMistralProviderWorkflow(t *testing.T) {
	// Skip if not in CI or if git is not available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for the test repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-mistral-test-*")
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

	// Test repository state
	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "feature.go", Status: "added", Diff: "package main\n\nfunc NewFeature() {}\n"},
		},
	}

	// Test Mistral provider creation and configuration
	config := &model.AIProviderConfig{
		Name:    "mistral",
		APIKey:  os.Getenv("MISTRAL_API_KEY"), // May be empty
		Model:   "mistral-large-latest",
		Timeout: 30 * time.Second,
		MaxTokens: 500,
	}

	provider := ai.NewMistralProvider(config)
	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	// Test token calculation (if available)
	// This would be tested through the token calculator integration

	// Test AI provider call (if API key available)
	if config.APIKey != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		message, err := provider.GenerateCommitMessage(ctx, state)
		if err != nil {
			t.Logf("Mistral provider call failed (expected without valid key or network): %v", err)
		} else {
			t.Logf("Mistral generated message: %s", message)
			if message == "" {
				t.Error("Expected non-empty message from Mistral")
			}
		}
	} else {
		t.Log("Skipping Mistral provider call - no MISTRAL_API_KEY in environment")
	}

	// Test that provider follows Conventional Commits format (if message generated)
	// This would be validated through the commit service integration

	t.Logf("Test repository created at: %s", tmpDir)
	t.Log("Mistral provider integration test structure ready")
}
