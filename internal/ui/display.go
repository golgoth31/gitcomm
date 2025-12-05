package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/golgoth31/gitcomm/internal/model"
)

// DisplayCommitMessage formats and displays a commit message for user review
func DisplayCommitMessage(message *model.CommitMessage) string {
	var lines []string

	// Format header
	header := message.Type
	if message.Scope != "" {
		header = fmt.Sprintf("%s(%s)", header, message.Scope)
	}
	header = fmt.Sprintf("%s: %s", header, message.Subject)
	lines = append(lines, header)

	// Add body if present
	if message.Body != "" {
		lines = append(lines, "")
		// Wrap body at 72 characters
		wrappedBody := wrapText(message.Body, 72)
		lines = append(lines, wrappedBody)
	}

	// Add footer if present
	if message.Footer != "" {
		lines = append(lines, "")
		lines = append(lines, message.Footer)
	}

	// Add signoff indicator if enabled
	if message.Signoff {
		lines = append(lines, "")
		lines = append(lines, "(Signed-off-by will be added)")
	}

	return strings.Join(lines, "\n")
}

// GetVisualIndicator returns the visual indicator character for the given prompt state
// with appropriate lipgloss styling applied
func GetVisualIndicator(state PromptState) string {
	style := lipgloss.NewStyle()

	switch state {
	case StatePending, StateActive:
		// Blue '?' for pending/active state
		return style.Foreground(lipgloss.Color("4")).Render("?")
	case StateCompleted:
		// Green '✓' for completed state
		return style.Foreground(lipgloss.Color("2")).Render("✓")
	case StateCancelled:
		// Red '✗' for cancelled state
		return style.Foreground(lipgloss.Color("1")).Render("✗")
	case StateError:
		// Yellow '⚠' for error state
		return style.Foreground(lipgloss.Color("3")).Render("⚠")
	default:
		// Default to blue '?' for unknown states
		return style.Foreground(lipgloss.Color("4")).Render("?")
	}
}

// wrapText wraps text at the specified width
func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+len(word)+1 <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// FormatPostValidationSummary formats a post-validation summary line with green checkmark
// Format: "✓ <title>: <value>"
// Exported for testing purposes
func FormatPostValidationSummary(title string, value interface{}) string {
	// Convert value to string representation
	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
		// For multiline strings, truncate to first line if too long
		if strings.Contains(valueStr, "\n") {
			lines := strings.Split(valueStr, "\n")
			firstLine := lines[0]
			if len(firstLine) > 50 {
				valueStr = firstLine[:50] + "..."
			} else if len(lines) > 1 {
				valueStr = firstLine + "..."
			} else {
				valueStr = firstLine
			}
		} else if len(valueStr) > 100 {
			// Truncate very long single-line values
			valueStr = valueStr[:100] + "..."
		}
	case bool:
		if v {
			valueStr = "Yes"
		} else {
			valueStr = "No"
		}
	case int:
		valueStr = fmt.Sprintf("%d", v)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	// Use green checkmark (✓) with the formatted string
	return fmt.Sprintf("✓ %s: %s", title, valueStr)
}

// printPostValidationSummary prints a post-validation summary line with green checkmark
// Format: "✓ <title>: <value>"
func printPostValidationSummary(title string, value interface{}) {
	fmt.Println(FormatPostValidationSummary(title, value))
}
