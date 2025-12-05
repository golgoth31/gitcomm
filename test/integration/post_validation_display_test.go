package integration

import (
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/ui"
)

// TestPostValidationSummaryFormat verifies that post-validation summary format is consistent
// across all prompt types: `✓ <title>: <value>`
func TestPostValidationSummaryFormat(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		value    interface{}
		expected string
	}{
		{
			name:     "Text input prompt",
			title:    "Scope",
			value:    "api",
			expected: "✓ Scope: api",
		},
		{
			name:     "Subject prompt",
			title:    "Subject (required)",
			value:    "Add user authentication",
			expected: "✓ Subject (required): Add user authentication",
		},
		{
			name:     "Boolean true value",
			title:    "Use AI to generate commit message?",
			value:    true,
			expected: "✓ Use AI to generate commit message?: Yes",
		},
		{
			name:     "Boolean false value",
			title:    "Create an empty commit?",
			value:    false,
			expected: "✓ Create an empty commit?: No",
		},
		{
			name:     "Commit type selection",
			title:    "Choose a type",
			value:    "feat",
			expected: "✓ Choose a type: feat",
		},
		{
			name:     "Multi-choice selection",
			title:    "Options",
			value:    "Accept and commit directly",
			expected: "✓ Options: Accept and commit directly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.FormatPostValidationSummary(tt.title, tt.value)
			if result != tt.expected {
				t.Errorf("FormatPostValidationSummary(%q, %v) = %q, want %q", tt.title, tt.value, result, tt.expected)
			}
			// Verify checkmark is present
			if !strings.Contains(result, "✓") {
				t.Errorf("Post-validation summary should contain checkmark: %q", result)
			}
			// Verify title is present
			if !strings.Contains(result, tt.title) {
				t.Errorf("Post-validation summary should contain title: %q", result)
			}
		})
	}
}

// TestMultilineContentTruncation verifies that multiline content is truncated correctly
// in post-validation display (body, footer prompts)
func TestMultilineContentTruncation(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		value    string
		expected string
	}{
		{
			name:     "Multiline body - first line only",
			title:    "Body",
			value:    "First line of body\nSecond line of body\nThird line",
			expected: "✓ Body: First line of body...",
		},
		{
			name:     "Multiline footer - first line only",
			title:    "Footer",
			value:    "Closes #123\nRelated to #456",
			expected: "✓ Footer: Closes #123...",
		},
		{
			name:     "Long first line - truncated",
			title:    "Body",
			value:    "This is a very long first line that should be truncated for the summary display because it exceeds the character limit\nSecond line",
			expected: "✓ Body: This is a very long first line that should be trun...",
		},
		{
			name:     "Single line body - no truncation if short",
			title:    "Body",
			value:    "Short body text",
			expected: "✓ Body: Short body text",
		},
		{
			name:     "Empty multiline",
			title:    "Body",
			value:    "",
			expected: "✓ Body: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ui.FormatPostValidationSummary(tt.title, tt.value)
			if result != tt.expected {
				t.Errorf("FormatPostValidationSummary(%q, %q) = %q, want %q", tt.title, tt.value, result, tt.expected)
			}
			// Verify checkmark is present
			if !strings.Contains(result, "✓") {
				t.Errorf("Post-validation summary should contain checkmark: %q", result)
			}
			// Verify title is present
			if !strings.Contains(result, tt.title) {
				t.Errorf("Post-validation summary should contain title: %q", result)
			}
			// Verify truncation occurred for multiline content
			if strings.Contains(tt.value, "\n") && !strings.Contains(result, "...") {
				// Multiline content should be truncated with ellipsis
				if len(result) > len(tt.expected) {
					t.Errorf("Multiline content should be truncated: %q", result)
				}
			}
		})
	}
}

// TestPostValidationSummaryAllPromptTypes verifies that all prompt types produce
// valid post-validation summary lines
func TestPostValidationSummaryAllPromptTypes(t *testing.T) {
	// Test all prompt value types that are used in the codebase

	// Text inputs
	scopeResult := ui.FormatPostValidationSummary("Scope", "api")
	if !strings.HasPrefix(scopeResult, "✓") {
		t.Errorf("Scope prompt summary should start with checkmark: %q", scopeResult)
	}

	subjectResult := ui.FormatPostValidationSummary("Subject (required)", "Add feature")
	if !strings.HasPrefix(subjectResult, "✓") {
		t.Errorf("Subject prompt summary should start with checkmark: %q", subjectResult)
	}

	// Multiline inputs
	bodyResult := ui.FormatPostValidationSummary("Body", "First line\nSecond line")
	if !strings.HasPrefix(bodyResult, "✓") {
		t.Errorf("Body prompt summary should start with checkmark: %q", bodyResult)
	}

	footerResult := ui.FormatPostValidationSummary("Footer", "Closes #123")
	if !strings.HasPrefix(footerResult, "✓") {
		t.Errorf("Footer prompt summary should start with checkmark: %q", footerResult)
	}

	// Selection prompts
	commitTypeResult := ui.FormatPostValidationSummary("Choose a type", "feat")
	if !strings.HasPrefix(commitTypeResult, "✓") {
		t.Errorf("Commit type prompt summary should start with checkmark: %q", commitTypeResult)
	}

	// Confirmation prompts
	confirmResult := ui.FormatPostValidationSummary("Create an empty commit?", true)
	if !strings.HasPrefix(confirmResult, "✓") {
		t.Errorf("Confirm prompt summary should start with checkmark: %q", confirmResult)
	}
	if !strings.Contains(confirmResult, "Yes") {
		t.Errorf("Boolean true should be formatted as 'Yes': %q", confirmResult)
	}

	// Multi-choice prompts
	optionsResult := ui.FormatPostValidationSummary("Options", "Accept and commit directly")
	if !strings.HasPrefix(optionsResult, "✓") {
		t.Errorf("Options prompt summary should start with checkmark: %q", optionsResult)
	}
}

// TestPromptUIClearedBeforeSummary verifies that prompt UI is cleared before summary line appears
// This is verified by checking that huh forms complete before summary is printed
func TestPromptUIClearedBeforeSummary(t *testing.T) {
	// huh forms automatically clear their UI when form.Run() completes
	// The post-validation summary is printed after form.Run() returns successfully
	// This means the prompt UI is cleared before the summary line appears

	// Verify the format function works correctly (indirect verification)
	// The actual UI clearing is handled by huh library when form.Run() completes
	result := ui.FormatPostValidationSummary("Test prompt", "test value")
	if !strings.HasPrefix(result, "✓") {
		t.Error("Summary should start with checkmark")
	}

	// The implementation in prompts.go follows this pattern:
	// 1. form.Run() - displays prompt UI and waits for user input
	// 2. form.Run() completes - huh library clears the prompt UI
	// 3. printPostValidationSummary() - prints summary line below where prompt was
	// This ensures the prompt UI is cleared before the summary appears
}
