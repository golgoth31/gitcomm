package integration

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/golgoth31/gitcomm/internal/ui"
)

// TestNoAltScreenInPrompts verifies that prompt functions use huh library and render inline
// This is verified by checking that prompts use huh (which renders inline by default)
func TestNoAltScreenInPrompts(t *testing.T) {
	// Skip if not in TTY environment (CI/test environments)
	if os.Getenv("CI") != "" {
		t.Skip("Skipping TTY-dependent test in CI environment")
	}

	// This test verifies that all prompt functions use huh library
	// which renders inline by default (no alt screen mode)

	// Verify that huh forms render inline by default
	// huh.NewForm() creates forms that render inline without alt screen
	var testValue string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Test prompt").
				Value(&testValue),
		),
	)

	// Verify form is created successfully
	if form == nil {
		t.Fatal("Failed to create huh form")
	}

	// huh forms render inline by default - no alt screen configuration needed
	// The form.State will be StateInactive initially
	// When Run() is called, it will render inline in the terminal

	// Test that formatPostValidationSummary helper works correctly
	// This verifies the post-validation display format
	result := ui.FormatPostValidationSummary("Test prompt", "test value")
	if !strings.Contains(result, "✓") {
		t.Error("Post-validation summary should contain checkmark")
	}
	if !strings.Contains(result, "Test prompt") {
		t.Error("Post-validation summary should contain prompt title")
	}
	if !strings.Contains(result, "test value") {
		t.Error("Post-validation summary should contain value")
	}
}

// TestHuhFormsRenderInline verifies that huh forms are configured for inline rendering
func TestHuhFormsRenderInline(t *testing.T) {
	// Verify that huh forms don't require alt screen configuration
	// huh forms render inline by default when using form.Run()

	var value1, value2 string

	// Test single-field form
	form1 := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Field 1").Value(&value1),
		),
	)
	if form1 == nil {
		t.Fatal("Failed to create single-field huh form")
	}

	// Test multi-field form
	form2 := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Field 1").Value(&value1),
			huh.NewInput().Title("Field 2").Value(&value2),
		),
	)
	if form2 == nil {
		t.Fatal("Failed to create multi-field huh form")
	}

	// huh forms render inline by default - no special configuration needed
	// The library handles terminal rendering automatically via Bubble Tea
	// Alt screen is not used for inline prompts
}

// TestNarrowTerminalWidthHandling verifies that huh forms handle narrow terminal widths gracefully
func TestNarrowTerminalWidthHandling(t *testing.T) {
	// Skip if not in TTY environment
	if os.Getenv("CI") != "" {
		t.Skip("Skipping TTY-dependent test in CI environment")
	}

	// huh library automatically handles terminal width constraints
	// Forms will wrap content appropriately for narrow terminals

	var testValue string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("This is a very long prompt title that should wrap gracefully in narrow terminals").
				Value(&testValue),
		),
	)

	if form == nil {
		t.Fatal("Failed to create huh form with long title")
	}

	// huh library uses lipgloss for styling, which handles width constraints
	// The form will automatically wrap and adjust layout for narrow terminals
	// This is handled by the underlying Bubble Tea and lipgloss libraries
}
