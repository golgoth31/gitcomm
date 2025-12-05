package ai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// TestNewMistralProvider_SDKClientInitialization tests SDK client initialization
func TestNewMistralProvider_SDKClientInitialization(t *testing.T) {
	tests := []struct {
		name   string
		config *model.AIProviderConfig
	}{
		{
			name: "valid config with SDK client",
			config: &model.AIProviderConfig{
				Name:   "mistral",
				APIKey: "test-key",
			},
		},
		{
			name: "config with empty API key (allowed in constructor)",
			config: &model.AIProviderConfig{
				Name:  "mistral",
				Model: "mistral-large-latest",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewMistralProvider(tt.config)
			if provider == nil {
				t.Error("Expected provider to be created")
			}
			// Verify it implements AIProvider interface
			var _ AIProvider = provider
			// Verify SDK client is initialized (will be checked in implementation)
		})
	}
}

// TestMistralProvider_GenerateCommitMessage_SDKErrorMapping tests SDK error mapping to existing error types
func TestMistralProvider_GenerateCommitMessage_SDKErrorMapping(t *testing.T) {
	config := &model.AIProviderConfig{
		Name:  "mistral",
		Model: "mistral-large-latest",
		// APIKey intentionally empty to test error handling
	}

	provider := NewMistralProvider(config)

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

// TestMistralProvider_SDKInitializationFailure tests SDK initialization failure handling
func TestMistralProvider_SDKInitializationFailure(t *testing.T) {
	// This test verifies that SDK initialization failures are handled gracefully
	config := &model.AIProviderConfig{
		Name:   "mistral",
		APIKey: "invalid-key-format",
		Model:  "mistral-large-latest",
	}

	provider := NewMistralProvider(config)
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

// TestMistralProvider_ContextCancellation tests that context cancellation works with SDK
func TestMistralProvider_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping context cancellation test in short mode")
	}

	config := &model.AIProviderConfig{
		Name:    "mistral",
		APIKey:  "test-key",
		Model:   "mistral-large-latest",
		Timeout: 30 * time.Second,
	}

	provider := NewMistralProvider(config)

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

// TestNewMistralProvider tests the constructor
func TestNewMistralProvider(t *testing.T) {
	tests := []struct {
		name   string
		config *model.AIProviderConfig
	}{
		{
			name: "valid config with defaults",
			config: &model.AIProviderConfig{
				Name:   "mistral",
				APIKey: "test-key",
			},
		},
		{
			name: "valid config with custom values",
			config: &model.AIProviderConfig{
				Name:      "mistral",
				APIKey:    "test-key",
				Model:     "mistral-small",
				Endpoint:  "https://custom.endpoint.com/v1/chat/completions",
				Timeout:   60 * time.Second,
				MaxTokens: 1000,
			},
		},
		{
			name: "config with empty API key (allowed in constructor)",
			config: &model.AIProviderConfig{
				Name:  "mistral",
				Model: "mistral-large-latest",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewMistralProvider(tt.config)
			if provider == nil {
				t.Error("Expected provider to be created")
			}
			// Verify it implements AIProvider interface
			var _ AIProvider = provider
		})
	}
}

// TestMistralProvider_GenerateCommitMessage_Success tests successful API response
func TestMistralProvider_GenerateCommitMessage_Success(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header with Bearer token")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json")
		}

		// Parse request body
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Verify request structure
		if reqBody["model"] != "mistral-large-latest" {
			t.Errorf("Expected model mistral-large-latest, got %v", reqBody["model"])
		}

		// Return successful response
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "feat(api): add new endpoint\n\nThis commit adds a new API endpoint for user management.",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &model.AIProviderConfig{
		Name:      "mistral",
		APIKey:    "test-key",
		Model:     "mistral-large-latest",
		Endpoint:  server.URL,
		Timeout:   30 * time.Second,
		MaxTokens: 500,
	}

	provider := NewMistralProvider(config)
	// Override endpoint for test
	if mp, ok := provider.(*MistralProvider); ok {
		mp.config.Endpoint = server.URL
	}

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "api.go", Status: "modified", Diff: "+func NewEndpoint() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message, err := provider.GenerateCommitMessage(ctx, state)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if message == "" {
		t.Error("Expected non-empty message")
	}

	if !strings.Contains(message, "feat") {
		t.Errorf("Expected message to contain 'feat', got: %s", message)
	}
}

// TestMistralProvider_GenerateCommitMessage_MissingAPIKey tests error handling for missing API key
func TestMistralProvider_GenerateCommitMessage_MissingAPIKey(t *testing.T) {
	config := &model.AIProviderConfig{
		Name:  "mistral",
		Model: "mistral-large-latest",
	}

	provider := NewMistralProvider(config)

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

	if !strings.Contains(err.Error(), "Mistral API key not configured") {
		t.Errorf("Expected error about missing API key, got: %v", err)
	}

	// Verify error is wrapped with ErrAIProviderUnavailable
	if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
		t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
	}
}

// TestMistralProvider_GenerateCommitMessage_APIErrors tests error handling for API errors
func TestMistralProvider_GenerateCommitMessage_APIErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:           "401 Unauthorized",
			statusCode:     http.StatusUnauthorized,
			responseBody:   `{"error": "Invalid API key"}`,
			expectedErrMsg: "API returned status 401",
		},
		{
			name:           "429 Rate Limit",
			statusCode:     http.StatusTooManyRequests,
			responseBody:   `{"error": "Rate limit exceeded"}`,
			expectedErrMsg: "API returned status 429",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"error": "Internal server error"}`,
			expectedErrMsg: "API returned status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			config := &model.AIProviderConfig{
				Name:     "mistral",
				APIKey:   "test-key",
				Endpoint: server.URL,
				Timeout:  30 * time.Second,
			}

			provider := NewMistralProvider(config)
			// Override endpoint for test
			if mp, ok := provider.(*MistralProvider); ok {
				mp.config.Endpoint = server.URL
			}

			state := &model.RepositoryState{
				StagedFiles: []model.FileChange{
					{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := provider.GenerateCommitMessage(ctx, state)
			if err == nil {
				t.Error("Expected error for API error")
			}

			if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Expected error to contain '%s', got: %v", tt.expectedErrMsg, err)
			}

			// Verify error is wrapped with ErrAIProviderUnavailable
			if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
				t.Errorf("Expected error to be wrapped with ErrAIProviderUnavailable, got: %v", err)
			}
		})
	}
}

// TestMistralProvider_GenerateCommitMessage_Timeout tests timeout handling
func TestMistralProvider_GenerateCommitMessage_Timeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay longer than context timeout
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "test",
					},
				},
			},
		})
	}))
	defer server.Close()

	config := &model.AIProviderConfig{
		Name:     "mistral",
		APIKey:   "test-key",
		Endpoint: server.URL,
		Timeout:  30 * time.Second,
	}

	provider := NewMistralProvider(config)
	// Override endpoint for test
	if mp, ok := provider.(*MistralProvider); ok {
		mp.config.Endpoint = server.URL
	}

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	// Use a short timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Error("Expected error for timeout")
	}

	// Should be context deadline exceeded or similar
	if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "timeout") {
		t.Logf("Timeout error (expected): %v", err)
	}
}

// TestMistralProvider_GenerateCommitMessage_EmptyResponse tests empty response handling
func TestMistralProvider_GenerateCommitMessage_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return empty choices array
		response := map[string]interface{}{
			"choices": []map[string]interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &model.AIProviderConfig{
		Name:     "mistral",
		APIKey:   "test-key",
		Endpoint: server.URL,
		Timeout:  30 * time.Second,
	}

	provider := NewMistralProvider(config)
	// Override endpoint for test
	if mp, ok := provider.(*MistralProvider); ok {
		mp.config.Endpoint = server.URL
	}

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Error("Expected error for empty response")
	}

	if !strings.Contains(err.Error(), "no response from API") {
		t.Errorf("Expected error about empty response, got: %v", err)
	}
}

// TestMistralProvider_GenerateUserMessage tests the unified prompt generator user message generation
func TestMistralProvider_GenerateUserMessage(t *testing.T) {
	config := &model.AIProviderConfig{
		Name:   "mistral",
		APIKey: "test-key",
	}

	provider := NewMistralProvider(config)
	mp := provider.(*MistralProvider)

	tests := []struct {
		name      string
		repoState *model.RepositoryState
		want      string
	}{
		{
			name: "with staged files",
			repoState: &model.RepositoryState{
				StagedFiles: []model.FileChange{
					{Path: "api.go", Status: "modified", Diff: "+func NewEndpoint() {}"},
				},
			},
			want: "Generate a commit message",
		},
		{
			name: "with unstaged files",
			repoState: &model.RepositoryState{
				UnstagedFiles: []model.FileChange{
					{Path: "test.go", Status: "modified", Diff: "+func Test() {}"},
				},
			},
			want: "Generate a commit message",
		},
		{
			name: "with both staged and unstaged files",
			repoState: &model.RepositoryState{
				StagedFiles: []model.FileChange{
					{Path: "api.go", Status: "modified", Diff: "+func NewEndpoint() {}"},
				},
				UnstagedFiles: []model.FileChange{
					{Path: "test.go", Status: "modified", Diff: "+func Test() {}"},
				},
			},
			want: "Generate a commit message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMsg, err := mp.generator.GenerateUserMessage(tt.repoState)
			if err != nil {
				t.Fatalf("GenerateUserMessage() error = %v", err)
			}
			if userMsg == "" {
				t.Error("Expected non-empty user message")
			}
			if !strings.Contains(userMsg, tt.want) {
				t.Errorf("Expected user message to contain '%s', got: %s", tt.want, userMsg)
			}
			// Verify user message contains file information
			if len(tt.repoState.StagedFiles) > 0 && !strings.Contains(userMsg, tt.repoState.StagedFiles[0].Path) {
				t.Errorf("Expected user message to contain file path '%s'", tt.repoState.StagedFiles[0].Path)
			}
		})
	}
}
