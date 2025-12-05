# Feature Specification: Display Commit Type Selection Confirmation

**Feature Branch**: `006-type-confirmation-display`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "when the commit type is choosen, a line MUST be displayed with the format: ✔ Choose a type(<scope>): <choosen scope>"

## Clarifications

### Session 2025-01-27

- Q: How should the confirmation line be displayed relative to the original prompt and the next prompt? → A: Display on a new line after alt-screen clears, with standard formatting (e.g., `fmt.Printf("✔ Choose a type(<scope>): %s\n", chosenType)`)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Display Commit Type Confirmation (Priority: P1)

A developer selects a commit type from the interactive list and expects to see a confirmation line displaying their selection in a clear, formatted way before proceeding to the next prompt.

**Why this priority**: This provides immediate visual feedback that the selection was successful and shows the chosen type clearly, improving user confidence and reducing confusion about what was selected. This is a simple but important UX enhancement that makes the workflow more transparent.

**Independent Test**: Can be fully tested by running the CLI, selecting a commit type from the interactive list, and verifying that a confirmation line appears with the format "✔ Choose a type(<scope>): <chosen type>" where <chosen type> is the actual selected type (e.g., "feat", "fix").

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI and reaches the commit type selection screen, **When** they select a commit type (e.g., "feat") and press Enter, **Then** a confirmation line is displayed with the format "✔ Choose a type(<scope>): feat" before the next prompt appears
2. **Given** a developer selects a different commit type (e.g., "fix"), **When** they confirm their selection, **Then** the confirmation line displays "✔ Choose a type(<scope>): fix"
3. **Given** a developer selects any valid commit type, **When** the confirmation line is displayed, **Then** the chosen type value matches exactly what was selected from the list

---

### Edge Cases

- What happens if the user cancels the selection (presses Escape)? → No confirmation line should be displayed, and the workflow should exit/cancel as normal
- How is the confirmation line displayed if the terminal is in alt-screen mode (bubbletea)? → The confirmation line should be displayed on a new line after the select list screen (alt-screen) is cleared, using standard terminal output formatting (e.g., `fmt.Printf` with newline)
- What if the selected type contains special characters? → The confirmation line should display the type exactly as it appears in the list (types are predefined and safe: feat, fix, docs, style, refactor, test, chore, version)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display a confirmation line immediately after a commit type is selected and confirmed
- **FR-002**: The confirmation line MUST use the exact format: "✔ Choose a type(<scope>): <chosen type>" where <chosen type> is the selected commit type value
- **FR-003**: The confirmation line MUST be displayed before the next prompt (scope prompt) appears, on a new line after the alt-screen clears with standard formatting
- **FR-004**: The confirmation line MUST only be displayed when a valid selection is made (not when cancelled)
- **FR-005**: The chosen type value in the confirmation line MUST match exactly the type that was selected from the interactive list
- **FR-006**: The confirmation line MUST use a checkmark symbol (✔) at the beginning
- **FR-007**: The confirmation line MUST preserve the original prompt text "Choose a type(<scope>):" format

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of successful commit type selections display the confirmation line with correct format
- **SC-002**: The confirmation line appears within 100ms of selection confirmation (perceived as immediate by users)
- **SC-003**: The chosen type value in the confirmation line matches the selected type in 100% of test cases
- **SC-004**: No confirmation line is displayed when selection is cancelled (0% false positives)

## Assumptions

- The confirmation line is displayed in standard terminal output (not within the bubbletea alt-screen)
- The checkmark symbol (✔) is supported by the user's terminal/encoding
- The format text "Choose a type(<scope>):" is preserved exactly as shown in the current prompt
- The confirmation is a one-time display after selection, not a persistent status indicator
- The scope part "(<scope>)" in the format is literal text, not a variable to be filled (scope is collected separately in the next prompt)

## Dependencies

- Requires existing commit type selection functionality (interactive select list from feature 004-improve-commit-ui)
- Requires `PromptCommitType` function in `internal/ui/prompts.go`
- Requires `SelectListModel` and commit type selection UI components

## Notes

- The user's original description used "choosen" but the implementation should use "chosen" (standard spelling)
- The format preserves the original prompt structure to maintain consistency with the user experience
- This is a display-only enhancement and does not change the selection logic or data flow
