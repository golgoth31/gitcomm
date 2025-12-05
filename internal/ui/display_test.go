package ui

import (
	"strings"
	"testing"
)

func TestGetVisualIndicator(t *testing.T) {
	tests := []struct {
		name           string
		state          PromptState
		expectedSymbol string
		expectedColor  string
	}{
		{
			name:           "StatePending returns blue question mark",
			state:          StatePending,
			expectedSymbol: "?",
			expectedColor:  "4", // Blue
		},
		{
			name:           "StateActive returns blue question mark",
			state:          StateActive,
			expectedSymbol: "?",
			expectedColor:  "4", // Blue
		},
		{
			name:           "StateCompleted returns green checkmark",
			state:          StateCompleted,
			expectedSymbol: "✓",
			expectedColor:  "2", // Green
		},
		{
			name:           "StateCancelled returns red X",
			state:          StateCancelled,
			expectedSymbol: "✗",
			expectedColor:  "1", // Red
		},
		{
			name:           "StateError returns yellow warning",
			state:          StateError,
			expectedSymbol: "⚠",
			expectedColor:  "3", // Yellow
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetVisualIndicator(tt.state)

			// Check that the result contains the expected symbol
			// Note: lipgloss may return plain text when not in a TTY, so we just verify the symbol is present
			if !strings.Contains(result, tt.expectedSymbol) {
				t.Errorf("GetVisualIndicator(%v) = %v, expected to contain %v", tt.state, result, tt.expectedSymbol)
			}

			// Verify the result is not empty
			if result == "" {
				t.Errorf("GetVisualIndicator(%v) returned empty string", tt.state)
			}
		})
	}
}

func TestGetVisualIndicator_Rendering(t *testing.T) {
	// Test that visual indicators render correctly for prompt titles
	indicator := GetVisualIndicator(StateActive)
	if indicator == "" {
		t.Error("GetVisualIndicator(StateActive) returned empty string")
	}

	// Test that completed state shows checkmark
	completedIndicator := GetVisualIndicator(StateCompleted)
	if !strings.Contains(completedIndicator, "✓") {
		t.Errorf("GetVisualIndicator(StateCompleted) = %v, expected to contain ✓", completedIndicator)
	}
}

// TestVisualIndicatorInPromptTitle tests that visual indicators are correctly used in prompt titles
func TestVisualIndicatorInPromptTitle(t *testing.T) {
	tests := []struct {
		name     string
		state    PromptState
		prompt   string
		expected string
	}{
		{
			name:     "Active state shows blue question mark",
			state:    StateActive,
			prompt:   "Choose a type",
			expected: "?",
		},
		{
			name:     "Completed state shows green checkmark",
			state:    StateCompleted,
			prompt:   "Choose a type",
			expected: "✓",
		},
		{
			name:     "Cancelled state shows red X",
			state:    StateCancelled,
			prompt:   "Choose a type",
			expected: "✗",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indicator := GetVisualIndicator(tt.state)
			title := indicator + " " + tt.prompt + ":"

			// Verify the title contains the expected indicator symbol
			if !strings.Contains(title, tt.expected) {
				t.Errorf("Prompt title = %v, expected to contain %v", title, tt.expected)
			}

			// Verify the title contains the prompt text
			if !strings.Contains(title, tt.prompt) {
				t.Errorf("Prompt title = %v, expected to contain %v", title, tt.prompt)
			}
		})
	}
}

func TestFormatPostValidationSummary(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		value    interface{}
		expected string
	}{
		{
			name:     "Simple text value",
			title:    "Scope (optional)",
			value:    "api",
			expected: "✓ Scope (optional): api",
		},
		{
			name:     "Boolean true value",
			title:    "Use AI?",
			value:    true,
			expected: "✓ Use AI?: Yes",
		},
		{
			name:     "Boolean false value",
			title:    "Create empty commit?",
			value:    false,
			expected: "✓ Create empty commit?: No",
		},
		{
			name:     "Integer value",
			title:    "Token count",
			value:    1234,
			expected: "✓ Token count: 1234",
		},
		{
			name:     "Multiline string - first line only",
			title:    "Body",
			value:    "First line\nSecond line\nThird line",
			expected: "✓ Body: First line...",
		},
		{
			name:     "Long single-line string - truncated",
			title:    "Subject",
			value:    strings.Repeat("a", 150),
			expected: "✓ Subject: " + strings.Repeat("a", 100) + "...",
		},
		{
			name:     "Long multiline first line - truncated",
			title:    "Body",
			value:    strings.Repeat("a", 100) + "\nSecond line",
			expected: "✓ Body: " + strings.Repeat("a", 50) + "...",
		},
		{
			name:     "Empty string",
			title:    "Scope",
			value:    "",
			expected: "✓ Scope: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPostValidationSummary(tt.title, tt.value)
			if result != tt.expected {
				t.Errorf("formatPostValidationSummary(%q, %v) = %q, expected %q", tt.title, tt.value, result, tt.expected)
			}
		})
	}
}
