package integration

import (
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/ui"
)

// TestVisualIndicatorInTextInputPrompt tests that visual indicators appear correctly
// in text input prompts throughout the interaction lifecycle
func TestVisualIndicatorInTextInputPrompt(t *testing.T) {
	// This test verifies that when a text input prompt is created and used,
	// it displays the correct visual indicator at each stage:
	// - Blue '?' when active/pending
	// - Green '✓' when completed

	// Note: This is a placeholder test that will be expanded when TextInputModel is implemented
	// For now, we verify the visual indicator function works correctly

	// Test that active state shows blue question mark
	activeIndicator := ui.GetVisualIndicator(ui.StateActive)
	if !strings.Contains(activeIndicator, "?") {
		t.Errorf("Active state indicator should contain '?', got %v", activeIndicator)
	}

	// Test that completed state shows green checkmark
	completedIndicator := ui.GetVisualIndicator(ui.StateCompleted)
	if !strings.Contains(completedIndicator, "✓") {
		t.Errorf("Completed state indicator should contain '✓', got %v", completedIndicator)
	}

	// Test that the indicators are different
	if activeIndicator == completedIndicator {
		t.Error("Active and completed indicators should be different")
	}
}
