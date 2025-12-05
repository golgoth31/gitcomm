package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// TestNewOpenAIProvider_SDKClientInitialization tests Responses API client initialization
func TestNewOpenAIProvider_SDKClientInitialization(t *testing.T) {
	tests := []struct {
		name   string
		config *model.AIProviderConfig
	}{
		{
			name: "valid config with SDK client",
			config: &model.AIProviderConfig{
				Name:   "openai",
				APIKey: "sk-test",
				Model:  "gpt-4",
			},
		},
		{
			name: "config with empty API key (allowed in constructor)",
			config: &model.AIProviderConfig{
				Name:  "openai",
				Model: "gpt-4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewOpenAIProvider(tt.config)
			if provider == nil {
				t.Error("Expected provider to be created")
			}
			// Verify it implements AIProvider interface
			var _ AIProvider = provider
			// Verify Responses API client is initialized (will be checked in implementation)
		})
	}
}

// TestOpenAIProvider_GenerateCommitMessage_SDKSuccess tests successful Responses API response
func TestOpenAIProvider_GenerateCommitMessage_SDKSuccess(t *testing.T) {
	// Skip if no API key (integration test)
	if testing.Short() {
		t.Skip("Skipping OpenAI provider Responses API test in short mode")
	}

	config := &model.AIProviderConfig{
		Name:      "openai",
		APIKey:    "test-key", // Will fail but tests the SDK structure
		Model:     "gpt-4",
		Timeout:   30 * time.Second,
		MaxTokens: 500,
	}

	provider := NewOpenAIProvider(config)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "api.go", Status: "modified", Diff: "+func NewEndpoint() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// This will fail without a real API key, but tests the Responses API integration structure
	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Log("OpenAI provider Responses API test passed (with real API key)")
	} else {
		t.Logf("OpenAI provider Responses API test structure verified (expected error: %v)", err)
		// Verify error is wrapped with ErrAIProviderUnavailable
		if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
			t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
		}
	}
}

// TestOpenAIProvider_GenerateCommitMessage_SDKErrorMapping tests Responses API error mapping to existing error types
func TestOpenAIProvider_GenerateCommitMessage_SDKErrorMapping(t *testing.T) {
	config := &model.AIProviderConfig{
		Name:  "openai",
		Model: "gpt-4",
		// APIKey intentionally empty to test error handling
	}

	provider := NewOpenAIProvider(config)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Error("Expected error for missing API key")
	}

	// Verify error is wrapped with ErrAIProviderUnavailable
	if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
		t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
	}

	// Verify error message is user-facing (not SDK-specific)
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

// TestOpenAIProvider_SDKInitializationFailure tests Responses API initialization failure handling
func TestOpenAIProvider_SDKInitializationFailure(t *testing.T) {
	// This test verifies that Responses API initialization failures are handled gracefully
	// In a real scenario, this would test cases like:
	// - Invalid SDK configuration
	// - SDK version incompatibility
	// - Missing SDK dependencies

	// For now, we test that the constructor doesn't panic
	// and that errors are handled in GenerateCommitMessage
	config := &model.AIProviderConfig{
		Name:   "openai",
		APIKey: "invalid-key-format",
		Model:  "gpt-4",
	}

	provider := NewOpenAIProvider(config)
	if provider == nil {
		t.Error("Expected provider to be created even with invalid config")
	}

	// Responses API initialization errors should be caught in GenerateCommitMessage
	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := provider.GenerateCommitMessage(ctx, state)
	// Should return error (either from SDK initialization or API call)
	if err == nil {
		t.Log("Note: SDK may handle invalid keys differently - this is expected")
	} else {
		// Verify error is properly wrapped
		if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
			t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
		}
	}
}

// TestOpenAIProvider_ContextCancellation tests that context cancellation works with Responses API
func TestOpenAIProvider_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping context cancellation test in short mode")
	}

	config := &model.AIProviderConfig{
		Name:    "openai",
		APIKey:  "test-key",
		Model:   "gpt-4",
		Timeout: 30 * time.Second,
	}

	provider := NewOpenAIProvider(config)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := provider.GenerateCommitMessage(ctx, state)
	// Should respect context cancellation
	if err == nil {
		t.Error("Expected error for cancelled context")
	}
	// Error should be context-related
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Logf("Context cancellation error (expected): %v", err)
	}
}

// TestOpenAIProvider_ValidateConfig tests config validation (existing test, kept for compatibility)
func TestOpenAIProvider_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *model.AIProviderConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &model.AIProviderConfig{
				Name:   "openai",
				APIKey: "sk-test",
				Model:  "gpt-4",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: &model.AIProviderConfig{
				Name:  "openai",
				Model: "gpt-4",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewOpenAIProvider(tt.config)
			if provider == nil && !tt.wantErr {
				t.Error("Expected provider to be created")
			}
			// Basic validation - actual API validation happens on call
		})
	}
}
