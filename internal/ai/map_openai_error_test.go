package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/openai/openai-go/v3"
)

// mustOpenAIError builds a minimal *openai.Error for unit tests.
// The SDK populates Request/Response for real API calls; we provide them here
// so Error() never dereferences a nil pointer.
func mustOpenAIError(statusCode int, message string) *openai.Error {
	req := httptest.NewRequest(http.MethodPost, openRouterDefaultBaseURL+"/chat/completions", nil)
	return &openai.Error{
		StatusCode: statusCode,
		Message:    message,
		Request:    req,
		Response:   &http.Response{StatusCode: statusCode, Status: http.StatusText(statusCode), Request: req},
	}
}

// TestMapOpenAIError tests the shared OpenAI SDK error mapping contract
func TestMapOpenAIError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantContextOK bool
		isUnavailable bool
		wantSubstr    string
		notWantSubstr string
		wantUnwrapped bool
	}{
		{
			name:          "preserves context canceled identity",
			err:           context.Canceled,
			wantContextOK: true,
		},
		{
			name:          "preserves wrapped context canceled identity",
			err:           fmt.Errorf("network error: %w", context.Canceled),
			wantContextOK: true,
		},
		{
			name:          "401 maps to API key invalid and preserves original",
			err:           mustOpenAIError(http.StatusUnauthorized, "Invalid API key provided"),
			isUnavailable: true,
			wantSubstr:    "API key invalid",
			wantUnwrapped: true,
		},
		{
			name:          "403 maps to API key invalid",
			err:           mustOpenAIError(http.StatusForbidden, "Forbidden"),
			isUnavailable: true,
			wantSubstr:    "API key invalid",
			wantUnwrapped: true,
		},
		{
			name:          "429 maps to rate limit exceeded",
			err:           mustOpenAIError(http.StatusTooManyRequests, "Rate limit reached"),
			isUnavailable: true,
			wantSubstr:    "rate limit exceeded",
			wantUnwrapped: true,
		},
		{
			name:          "404 maps to model not found",
			err:           mustOpenAIError(http.StatusNotFound, `The model "x" does not exist or you do not have access`),
			isUnavailable: true,
			wantSubstr:    "model not found",
			wantUnwrapped: true,
		},
		{
			name:          "500 falls back to generic unavailable but preserves original",
			err:           mustOpenAIError(http.StatusInternalServerError, "Internal server error"),
			isUnavailable: true,
			wantSubstr:    "AI provider unavailable",
			wantUnwrapped: true,
		},
		{
			name:          "invalid model string error is not mislabeled as API key invalid",
			err:           errors.New("invalid model 'openrouter/foo'"),
			isUnavailable: true,
			wantSubstr:    "AI provider unavailable",
			notWantSubstr: "API key invalid",
		},
		{
			name:          "authentication string error maps to API key invalid",
			err:           errors.New("authentication failed for this request"),
			isUnavailable: true,
			wantSubstr:    "API key invalid",
		},
		{
			name:          "generic string error maps to unavailable",
			err:           errors.New("connection reset by peer"),
			isUnavailable: true,
			wantSubstr:    "AI provider unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOpenAIError(tt.err)

			if tt.wantContextOK {
				if !errors.Is(got, context.Canceled) && !errors.Is(got, context.DeadlineExceeded) {
					t.Errorf("expected context identity preserved, got: %v", got)
				}
				return
			}

			if tt.isUnavailable && !utils.IsError(got, utils.ErrAIProviderUnavailable) {
				t.Errorf("expected ErrAIProviderUnavailable, got: %v", got)
			}
			if tt.wantSubstr != "" && !strings.Contains(got.Error(), tt.wantSubstr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantSubstr, got)
			}
			if tt.notWantSubstr != "" && strings.Contains(got.Error(), tt.notWantSubstr) {
				t.Errorf("expected error NOT containing %q, got: %v", tt.notWantSubstr, got)
			}
			if tt.wantUnwrapped {
				var apiErr *openai.Error
				if !errors.As(got, &apiErr) {
					t.Errorf("expected original *openai.Error preserved in chain, got: %v", got)
				}
			}
		})
	}
}
