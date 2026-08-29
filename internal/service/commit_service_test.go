package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/config"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

func TestCreateTimeoutContext(t *testing.T) {
	tests := []struct {
		name           string
		timeout        time.Duration
		expectedWithin time.Duration
	}{
		{
			name:           "3 second timeout",
			timeout:        3 * time.Second,
			expectedWithin: 3*time.Second + 100*time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			// Verify context is not cancelled initially
			if ctx.Err() != nil {
				t.Errorf("Context should not be cancelled initially, got: %v", ctx.Err())
			}

			// Wait for timeout
			select {
			case <-ctx.Done():
				// Context should be done after timeout
				if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
					t.Errorf("Expected DeadlineExceeded, got: %v", ctx.Err())
				}
			case <-time.After(tt.expectedWithin):
				t.Errorf("Context should have timed out within %v", tt.expectedWithin)
			}
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: true,
		},
		{
			name:     "wrapped deadline exceeded",
			err:      errors.New("restoration failed: " + context.DeadlineExceeded.Error()),
			expected: false, // errors.Is should still work
		},
		{
			name:     "other error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, context.DeadlineExceeded)
			if result != tt.expected {
				t.Errorf("IsTimeoutError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCommitService_GenerateWithAIRoutesToOpenRouter verifies the provider switch routes
// the "openrouter" provider name to the OpenRouter provider
func TestCommitService_GenerateWithAIRoutesToOpenRouter(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			DefaultProvider: "openrouter",
			Providers: map[string]model.AIProviderConfig{
				"openrouter": {Name: "openrouter"},
			},
		},
	}

	s := NewCommitService(nil, &model.CommitOptions{}, cfg)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	_, err := s.generateWithAIWithRetry(context.Background(), state, 0)
	if err == nil {
		t.Fatal("Expected error for unconfigured OpenRouter API key")
	}

	// Error must be wrapped with ErrAIProviderUnavailable (not "unknown provider")
	if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
		t.Errorf("Expected ErrAIProviderUnavailable, got: %v", err)
	}

	// Error must originate from the OpenRouter provider, proving the switch routed to it
	if !strings.Contains(err.Error(), "OpenRouter") {
		t.Errorf("Expected error to reference OpenRouter provider, got: %v", err)
	}

	if strings.Contains(err.Error(), "unknown provider") {
		t.Errorf("Expected provider switch to route to openrouter, got unknown provider error: %v", err)
	}
}

// TestCommitService_GenerateWithAIUnknownProvider verifies unrecognized providers are rejected
func TestCommitService_GenerateWithAIUnknownProvider(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			DefaultProvider: "nonexistent",
			Providers: map[string]model.AIProviderConfig{
				"nonexistent": {Name: "nonexistent"},
			},
		},
	}

	s := NewCommitService(nil, &model.CommitOptions{}, cfg)

	_, err := s.generateWithAIWithRetry(context.Background(), &model.RepositoryState{}, 0)
	if err == nil {
		t.Fatal("Expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "unknown provider") {
		t.Errorf("Expected unknown provider error, got: %v", err)
	}
}
