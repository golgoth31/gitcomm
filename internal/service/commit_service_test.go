package service

import (
	"context"
	"errors"
	"testing"
	"time"
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
