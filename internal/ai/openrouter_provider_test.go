package ai

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// TestNewOpenRouterProvider_ClientInitialization tests client initialization
func TestNewOpenRouterProvider_ClientInitialization(t *testing.T) {
	tests := []struct {
		name   string
		config *model.AIProviderConfig
	}{
		{
			name: "valid config with SDK client",
			config: &model.AIProviderConfig{
				Name:   "openrouter",
				APIKey: "sk-or-test-12345",
				Model:  "openrouter/auto",
			},
		},
		{
			name: "config with empty API key (allowed in constructor)",
			config: &model.AIProviderConfig{
				Name:  "openrouter",
				Model: "openrouter/auto",
			},
		},
		{
			name: "config with custom endpoint",
			config: &model.AIProviderConfig{
				Name:     "openrouter",
				APIKey:   "sk-or-test-12345",
				Endpoint: "https://openrouter.ai/api/v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewOpenRouterProvider(tt.config)
			if provider == nil {
				t.Fatal("Expected provider to be created")
			}
			// Verify it implements AIProvider interface
			var _ AIProvider = provider
		})
	}
}

// TestOpenRouterProvider_GenerateCommitMessage_MissingAPIKey tests error handling when no API key
func TestOpenRouterProvider_GenerateCommitMessage_MissingAPIKey(t *testing.T) {
	config := &model.AIProviderConfig{
		Name:  "openrouter",
		Model: "openrouter/auto",
	}

	provider := NewOpenRouterProvider(config)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Fatal("Expected error for missing API key")
	}

	// Verify error is wrapped with ErrAIProviderUnavailable
	if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
		t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
	}
}

// TestOpenRouterProvider_GenerateCommitMessage_SDKErrorMapping tests SDK error mapping to existing error types
func TestOpenRouterProvider_GenerateCommitMessage_SDKErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantSubstr string
	}{
		{
			name:       "authentication error maps to API key invalid",
			statusCode: http.StatusUnauthorized,
			body:       `{"error":{"message":"Invalid API key provided"}}`,
			wantSubstr: "API key invalid",
		},
		{
			name:       "rate limit error maps to rate limit exceeded",
			statusCode: http.StatusTooManyRequests,
			body:       `{"error":{"message":"Rate limit reached, please retry"}}`,
			wantSubstr: "rate limit exceeded",
		},
		{
			name:       "server error maps to generic unavailable",
			statusCode: http.StatusInternalServerError,
			body:       `{"error":{"message":"Internal server error"}}`,
			wantSubstr: "AI provider unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			config := &model.AIProviderConfig{
				Name:     "openrouter",
				APIKey:   "sk-or-test",
				Model:    "openrouter/auto",
				Endpoint: server.URL,
				Timeout:  5 * time.Second,
			}

			provider := NewOpenRouterProvider(config)

			state := &model.RepositoryState{
				StagedFiles: []model.FileChange{
					{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := provider.GenerateCommitMessage(ctx, state)
			if err == nil {
				t.Fatalf("Expected error for status %d", tt.statusCode)
			}

			if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
				t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
			}

			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Errorf("Expected error containing %q, got: %v", tt.wantSubstr, err)
			}
		})
	}
}

// TestOpenRouterProvider_ContextCancellation tests that context cancellation works with Chat Completions API
func TestOpenRouterProvider_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping context cancellation test in short mode")
	}

	config := &model.AIProviderConfig{
		Name:    "openrouter",
		APIKey:  "sk-or-test-12345",
		Model:   "openrouter/auto",
		Timeout: 30 * time.Second,
	}

	provider := NewOpenRouterProvider(config)

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
		t.Fatal("Expected error for cancelled context")
	}
	// Error must preserve its context identity
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context-related error, got: %v", err)
	}
}
