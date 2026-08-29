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
		{
			name:   "nil config does not panic",
			config: nil,
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

// TestOpenRouterProvider_GenerateCommitMessage_Timeout tests that a configured
// timeout surfaces a deadline error instead of hanging on a slow provider
func TestOpenRouterProvider_GenerateCommitMessage_Timeout(t *testing.T) {
	// Server slower than the configured client timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second)
	}))
	defer server.Close()

	config := &model.AIProviderConfig{
		Name:     "openrouter",
		APIKey:   "sk-or-test",
		Model:    "openrouter/auto",
		Endpoint: server.URL,
		Timeout:  100 * time.Millisecond,
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
		t.Fatal("Expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected deadline/timeout error, got: %v", err)
	}
}

// TestOpenRouterProvider_ContextCancellation tests that cancelling the context
// mid-request propagates through the SDK and preserves its identity
func TestOpenRouterProvider_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping context cancellation test in short mode")
	}

	// Server responds slowly so we can cancel while the request is in-flight
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
	}))
	defer server.Close()

	config := &model.AIProviderConfig{
		Name:     "openrouter",
		APIKey:   "sk-or-test-12345",
		Model:    "openrouter/auto",
		Endpoint: server.URL,
		Timeout:  30 * time.Second,
	}

	provider := NewOpenRouterProvider(config)

	state := &model.RepositoryState{
		StagedFiles: []model.FileChange{
			{Path: "test.go", Status: "modified", Diff: "func Test() {}"},
		},
	}

	// Cancel the context while the request is in-flight
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := provider.GenerateCommitMessage(ctx, state)
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
	// Error must preserve its context identity
	if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context-related error, got: %v", err)
	}
}

// TestOpenRouterProvider_GenerateCommitMessage_Success tests the happy path:
// request shape (auth header, model) and content extraction
func TestOpenRouterProvider_GenerateCommitMessage_Success(t *testing.T) {
	tests := []struct {
		name      string
		model     string
		wantModel string
	}{
		{
			name:      "explicit model",
			model:     "openrouter/auto",
			wantModel: "openrouter/auto",
		},
		{
			name:      "default model fallback",
			model:     "",
			wantModel: "openrouter/auto",
		},
	}

	const (
		wantContent = "feat(api): add endpoint"
		wantAPIKey  = "sk-or-test"
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				gotAuth  string
				gotModel string
				gotMsgCt int
			)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotAuth = r.Header.Get("Authorization")

				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}
				gotModel, _ = body["model"].(string)
				if msgs, ok := body["messages"].([]interface{}); ok {
					gotMsgCt = len(msgs)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":    "chatcmpl-test",
					"model": tt.wantModel,
					"choices": []map[string]interface{}{
						{
							"index":         0,
							"message":       map[string]interface{}{"role": "assistant", "content": wantContent},
							"finish_reason": "stop",
						},
					},
				})
			}))
			defer server.Close()

			config := &model.AIProviderConfig{
				Name:     "openrouter",
				APIKey:   wantAPIKey,
				Model:    tt.model,
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

			got, err := provider.GenerateCommitMessage(ctx, state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != wantContent {
				t.Errorf("got %q, want %q", got, wantContent)
			}
			if gotAuth != "Bearer "+wantAPIKey {
				t.Errorf("got Authorization %q, want %q", gotAuth, "Bearer "+wantAPIKey)
			}
			if gotModel != tt.wantModel {
				t.Errorf("got model %q, want %q", gotModel, tt.wantModel)
			}
			if gotMsgCt != 2 {
				t.Errorf("got %d messages, want 2 (system + user)", gotMsgCt)
			}
		})
	}
}

// TestOpenRouterProvider_GenerateCommitMessage_EmptyResponse tests empty-response
// handling (no choices, and present-but-empty content)
func TestOpenRouterProvider_GenerateCommitMessage_EmptyResponse(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantSubstr string
	}{
		{
			name:       "no choices",
			body:       `{"id":"chatcmpl-test","model":"openrouter/auto","choices":[]}`,
			wantSubstr: "no choices",
		},
		{
			name:       "present but empty content",
			body:       `{"id":"chatcmpl-test","model":"openrouter/auto","choices":[{"index":0,"message":{"role":"assistant","content":""},"finish_reason":"stop"}]}`,
			wantSubstr: "empty response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
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
				t.Fatal("Expected error for empty response")
			}
			if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
				t.Errorf("Expected ErrAIProviderUnavailable, got: %v", err)
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Errorf("Expected error containing %q, got: %v", tt.wantSubstr, err)
			}
		})
	}
}

// TestOpenRouterProvider_GenerateCommitMessage_Refusal tests that a content
// refusal is surfaced instead of a generic empty-response error
func TestOpenRouterProvider_GenerateCommitMessage_Refusal(t *testing.T) {
	const wantRefusal = "I cannot help with that"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "chatcmpl-test",
			"model": "openrouter/auto",
			"choices": []map[string]interface{}{
				{
					"index":         0,
					"message":       map[string]interface{}{"role": "assistant", "content": "", "refusal": wantRefusal},
					"finish_reason": "content_filter",
				},
			},
		})
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
		t.Fatal("Expected error for refused response")
	}
	if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
		t.Errorf("Expected ErrAIProviderUnavailable, got: %v", err)
	}
	if !strings.Contains(err.Error(), "refused") || !strings.Contains(err.Error(), wantRefusal) {
		t.Errorf("Expected refusal error mentioning %q, got: %v", wantRefusal, err)
	}
}

// TestOpenRouterProvider_GenerateCommitMessage_Truncated tests that a
// finish_reason of "length" is surfaced as truncation rather than empty response
func TestOpenRouterProvider_GenerateCommitMessage_Truncated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    "chatcmpl-test",
			"model": "openrouter/auto",
			"choices": []map[string]interface{}{
				{
					"index":         0,
					"message":       map[string]interface{}{"role": "assistant", "content": ""},
					"finish_reason": "length",
				},
			},
		})
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
		t.Fatal("Expected error for truncated response")
	}
	if !utils.IsError(err, utils.ErrAIProviderUnavailable) {
		t.Errorf("Expected ErrAIProviderUnavailable, got: %v", err)
	}
	if !strings.Contains(err.Error(), "truncated") {
		t.Errorf("Expected truncation error, got: %v", err)
	}
}
