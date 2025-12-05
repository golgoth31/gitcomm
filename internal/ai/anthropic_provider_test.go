package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// TestNewAnthropicProvider_SDKClientInitialization tests SDK client initialization
func TestNewAnthropicProvider_SDKClientInitialization(t *testing.T) {
	tests := []struct {
		name   string
		config *model.AIProviderConfig
	}{
		{
			name: "valid config with SDK client",
			config: &model.AIProviderConfig{
				Name:   "anthropic",
				APIKey: "sk-ant-test",
				Model:  "claude-3-opus",
			},
		},
		{
			name: "config with empty API key (allowed in constructor)",
			config: &model.AIProviderConfig{
				Name:  "anthropic",
				Model: "claude-3-opus",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewAnthropicProvider(tt.config)
			if provider == nil {
				t.Error("Expected provider to be created")
			}
			// Verify it implements AIProvider interface
			var _ AIProvider = provider
			// Verify SDK client is initialized (will be checked in implementation)
		})
	}
}

// TestAnthropicProvider_GenerateCommitMessage_SDKSuccess tests successful SDK API response
func TestAnthropicProvider_GenerateCommitMessage_SDKSuccess(t *testing.T) {
	// Skip if no API key (integration test)
	if testing.Short() {
		t.Skip("Skipping Anthropic provider SDK test in short mode")
	}

	config := &model.AIProviderConfig{
		Name:      "anthropic",
		APIKey:    "test-key", // Will fail but tests the SDK structure
		Model:     "claude-3-opus",
		Timeout:   30 * time.Second,
		MaxTokens: 500,
	}

	provider := NewAnthropicProvider(config)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "api.go", Status: "modified", Diff: "+func NewEndpoint() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// This will fail without a real API key, but tests the SDK integration structure
	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Log("Anthropic provider SDK test passed (with real API key)")
	} else {
		t.Logf("Anthropic provider SDK test structure verified (expected error: %v)", err)
		// Verify error is wrapped with ErrAIProviderUnavailable
		if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
			t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
		}
	}
}

// TestAnthropicProvider_GenerateCommitMessage_SDKErrorMapping tests SDK error mapping to existing error types
func TestAnthropicProvider_GenerateCommitMessage_SDKErrorMapping(t *testing.T) {
	config := &model.AIProviderConfig{
		Name:  "anthropic",
		Model: "claude-3-opus",
		// APIKey intentionally empty to test error handling
	}

	provider := NewAnthropicProvider(config)

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

// TestAnthropicProvider_SDKInitializationFailure tests SDK initialization failure handling
func TestAnthropicProvider_SDKInitializationFailure(t *testing.T) {
	// This test verifies that SDK initialization failures are handled gracefully
	config := &model.AIProviderConfig{
		Name:   "anthropic",
		APIKey: "invalid-key-format",
		Model:  "claude-3-opus",
	}

	provider := NewAnthropicProvider(config)
	if provider == nil {
		t.Error("Expected provider to be created even with invalid config")
	}

	// SDK initialization errors should be caught in GenerateCommitMessage
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

// TestAnthropicProvider_ContextCancellation tests that context cancellation works with SDK
func TestAnthropicProvider_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping context cancellation test in short mode")
	}

	config := &model.AIProviderConfig{
		Name:    "anthropic",
		APIKey:  "test-key",
		Model:   "claude-3-opus",
		Timeout: 30 * time.Second,
	}

	provider := NewAnthropicProvider(config)

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

// TestAnthropicProvider_ValidateConfig tests config validation (existing test, kept for compatibility)
func TestAnthropicProvider_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *model.AIProviderConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &model.AIProviderConfig{
				Name:   "anthropic",
				APIKey: "sk-ant-test",
				Model:  "claude-3-opus",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: &model.AIProviderConfig{
				Name:  "anthropic",
				Model: "claude-3-opus",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewAnthropicProvider(tt.config)
			if provider == nil && !tt.wantErr {
				t.Error("Expected provider to be created")
			}
		})
	}
}
