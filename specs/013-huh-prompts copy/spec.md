# Feature Specification: Rewrite CLI Prompts with Huh Library

**Feature Branch**: `013-huh-prompts`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "fully rewrite the cli prompts using https://github.com/charmbracelet/huh . all prompts must use this library. they must be inline and NOT in altscreen. once a prompt is validated by the user, it must be written to be replaced in the out put by a line with the format: <green checkmark> <prompt title or question>: <prompt validated answer>"

## User Scenarios & Testing

### User Story 1 - Migrate All Prompts to Huh Library (Priority: P1)

When a user interacts with any CLI prompt, the prompt should use the `huh` library for rendering and interaction, replacing all existing custom Bubble Tea implementations.

**Why this priority**: This is the core requirement - migrating all prompts to use the `huh` library. Without this, the feature cannot be considered complete. This affects every user interaction with the CLI.

**Independent Test**: User runs gitcomm and interacts with all prompt types (text input, select lists, multiline input, yes/no choices, confirmations). All prompts use the `huh` library for rendering and interaction. This can be fully tested independently and delivers consistent interaction patterns using a well-maintained library.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit scope (text input), **When** the prompt is displayed, **Then** it uses `huh.NewInput()` for rendering and input handling
2. **Given** user is prompted for commit subject (text input with validation), **When** the prompt is displayed, **Then** it uses `huh.NewInput()` with validation
3. **Given** user is prompted for commit type (select list), **When** the prompt is displayed, **Then** it uses `huh.NewSelect()` for rendering and selection
4. **Given** user is prompted for commit body (multiline input), **When** the prompt is displayed, **Then** it uses `huh.NewText()` for multiline input handling
5. **Given** user is prompted for yes/no choice (AI usage, empty commit, etc.), **When** the prompt is displayed, **Then** it uses `huh.NewConfirm()` for rendering and choice selection
6. **Given** user is prompted for commit footer (multiline input), **When** the prompt is displayed, **Then** it uses `huh.NewText()` for multiline input handling

---

### User Story 2 - Inline Rendering Without Alt Screen (Priority: P1)

When a user interacts with CLI prompts, the prompts should render inline without using alt screen mode, keeping the terminal history visible and maintaining context.

**Why this priority**: Inline rendering is a core requirement. Using alt screen would hide terminal history and break the user's context. This is critical for user experience.

**Independent Test**: User runs gitcomm and interacts with prompts. All prompts render inline without switching to alt screen mode. Terminal history remains visible throughout the interaction. This can be fully tested independently and delivers better context preservation.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit type (select list), **When** the prompt is displayed, **Then** it renders inline without using alt screen mode
2. **Given** user is prompted for commit body (multiline input), **When** the prompt is displayed, **Then** it renders inline without using alt screen mode
3. **Given** user interacts with any prompt, **When** the prompt is active, **Then** previous terminal output remains visible above the prompt
4. **Given** user runs gitcomm with multiple prompts, **When** each prompt completes, **Then** all previous prompts and their answers remain visible in the terminal

---

### User Story 3 - Post-Validation Display Format (Priority: P1)

When a user validates a prompt (confirms their input/selection), the prompt UI should be cleared and replaced in the output with a single line below where the prompt was, showing a green checkmark, the prompt title/question, and the validated answer. Previous terminal output must remain visible.

**Why this priority**: This is a core visual requirement. The post-validation display format provides clear confirmation of user input and maintains a clean, readable output. This affects every prompt interaction.

**Independent Test**: User runs gitcomm and completes any prompt (text input, select, multiline, confirmation). After validation, the prompt is replaced with a line showing: `<green checkmark> <prompt title>: <validated answer>`. This can be fully tested independently and delivers clear visual confirmation.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit scope with title "Scope (optional)", **When** user enters "api" and validates, **Then** the output shows "✓ Scope (optional): api" with a green checkmark
2. **Given** user is prompted for commit subject with title "Subject (required)", **When** user enters "Add user authentication" and validates, **Then** the output shows "✓ Subject (required): Add user authentication" with a green checkmark
3. **Given** user is prompted for commit type with title "Choose a type", **When** user selects "feat" and validates, **Then** the output shows "✓ Choose a type: feat" with a green checkmark
4. **Given** user is prompted for commit body with title "Body" within a combined form, **When** user enters multiline text and completes the field, **Then** the body field UI is cleared and the output shows "✓ Body: <first line of body text>" with a green checkmark, then the form moves to the next field (multiline content may be truncated or shown on subsequent lines)
5. **Given** user is prompted for AI usage with title "Use AI to generate commit message?", **When** user selects "Yes" and validates, **Then** the output shows "✓ Use AI to generate commit message?: Yes" with a green checkmark
6. **Given** user is prompted for empty commit confirmation with title "No changes detected. Create an empty commit?", **When** user selects "No" and validates, **Then** the output shows "✓ No changes detected. Create an empty commit?: No" with a green checkmark

---

### Edge Cases

- What happens when a prompt is cancelled (Escape key or Ctrl+C)? **Answer**: The prompt should handle cancellation gracefully. The cancelled prompt may show a cancellation indicator or be removed from output without showing the validated answer format
- What happens when terminal width is very narrow? **Answer**: Prompts should handle narrow terminals gracefully, wrapping text appropriately. The post-validation line should also wrap if needed
- What happens when multiline input is very long? **Answer**: The post-validation display should show the validated answer. For very long multiline content, the system may truncate with ellipsis or show the first line followed by indication of more content
- What happens when select list has many items and terminal is short? **Answer**: List should scroll appropriately without alt screen, and the post-validation line should show the selected value
- What happens when validation fails in a prompt? **Answer**: Prompt should show validation error inline within the prompt UI (below the input field, using `huh`'s built-in validation display) and allow user to correct input. The post-validation format should only appear after successful validation
- What happens when a prompt has a default value? **Answer**: Default value should be displayed appropriately in the prompt, and the post-validation line should show the final value (default or user-entered)
- What happens when user presses Ctrl+C during a prompt? **Answer**: Prompt should handle interruption gracefully. The workflow may exit or show cancellation state, but should not show the post-validation format for cancelled prompts
- What happens with AI message acceptance prompts that have multiple options? **Answer**: These should use appropriate `huh` field type (e.g., `huh.NewSelect()` for multiple choice options) and show the selected option in the post-validation format

## Requirements

### Functional Requirements

- **FR-001**: System MUST use the `huh` library (github.com/charmbracelet/huh) for all CLI prompts
- **FR-002**: System MUST NOT use alt screen mode for any prompt (all prompts render inline)
- **FR-003**: System MUST replace all existing custom Bubble Tea prompt implementations (TextInputModel, MultilineInputModel, YesNoChoiceModel, SelectListModel) with `huh` equivalents
- **FR-004**: System MUST render prompts inline, keeping terminal history visible throughout all prompt interactions
- **FR-005**: System MUST display post-validation format for all successfully validated prompts: `<green checkmark> <prompt title or question>: <prompt validated answer>`. The prompt UI MUST be cleared and the summary line MUST be displayed below where the prompt was, keeping previous terminal output visible. For combined forms with multiple fields, summary lines MUST be shown progressively: clear each field's UI as it's completed and show its summary line, then move to the next field
- **FR-006**: System MUST use `huh.NewInput()` for single-line text input prompts (scope, subject)
- **FR-007**: System MUST use `huh.NewText()` for multiline text input prompts (body, footer)
- **FR-008**: System MUST use `huh.NewSelect()` for single-selection prompts (commit type selection)
- **FR-009**: System MUST use `huh.NewConfirm()` for yes/no confirmation prompts (AI usage, empty commit, general confirmations)
- **FR-010**: System MUST use `huh.NewMultiSelect()` or `huh.NewSelect()` for prompts with multiple choice options (AI message acceptance options, commit failure choices)
- **FR-011**: System MUST preserve all existing prompt functionality (validation, default values, pre-selection, etc.) when migrating to `huh`
- **FR-018**: System MUST display validation errors inline within the prompt UI (below the input field) using `huh`'s built-in validation display mechanism
- **FR-012**: System MUST maintain backward compatibility with existing prompt function signatures (callers should not need changes). When prompts are combined into a single form, individual prompt functions may internally use a shared form but must still return values compatible with existing callers
- **FR-013**: System MUST handle prompt cancellation (Escape key, Ctrl+C) gracefully
- **FR-014**: System MUST configure `huh` forms to render inline (not in alt screen) using appropriate `huh` configuration options
- **FR-019**: System MUST combine related prompts into single `huh.Form` instances with multiple fields where appropriate (e.g., commit message fields: type, scope, subject, body, footer in one form)
- **FR-015**: System MUST display green checkmark character (✓) in the post-validation format
- **FR-016**: System MUST format multiline answers appropriately in the post-validation display (may truncate or show first line with indication of more content)
- **FR-017**: System MUST support all existing prompt types: scope, subject, body, footer, commit type, AI usage, empty commit confirmation, AI message acceptance, AI message edit, commit failure choice, reject choice

### Key Entities

- **Huh Form**: Represents a `huh.Form` instance used to render prompts
- **Huh Field**: Represents individual `huh` field types (Input, Text, Select, Confirm, MultiSelect)
- **Post-Validation Display**: Represents the formatted output line shown after prompt validation: `<green checkmark> <title>: <answer>`
- **Prompt Result**: Represents the result of a prompt interaction (value selected/entered, cancellation, error)

## Success Criteria

### Measurable Outcomes

- **SC-001**: 100% of prompts use the `huh` library for rendering and interaction (all custom Bubble Tea models replaced)
- **SC-002**: 0 prompts use alt screen mode (all prompts render inline)
- **SC-003**: 100% of successfully validated prompts display the post-validation format: `<green checkmark> <title>: <answer>`
- **SC-004**: All existing prompt functionality preserved (100% backward compatibility with existing callers)
- **SC-005**: Terminal history remains visible during all prompt interactions (user can scroll up to see previous output)
- **SC-006**: All prompt types successfully migrated: text input (scope, subject), multiline input (body, footer), select lists (commit type), confirmations (AI usage, empty commit, etc.), multi-choice prompts (AI message acceptance, commit failure)
- **SC-007**: All existing validation rules work correctly with `huh` prompts (required fields, length limits, format validation)
- **SC-008**: Default values and pre-selections work correctly with `huh` prompts

## Assumptions

- The `huh` library supports inline rendering without alt screen mode
- The `huh` library can handle all required prompt types (Input, Text, Select, Confirm, MultiSelect)
- The `huh` library supports validation, default values, and pre-selection functionality
- Existing prompt function signatures can be maintained while changing internal implementation to `huh`
- Terminal width/height constraints can be handled gracefully by `huh` without alt screen
- All existing validation and business logic can be preserved when migrating to `huh`
- The `huh` library provides mechanisms to customize post-validation display or we can implement custom display logic after form completion
- Green checkmark character (✓) is supported in all user terminals (no fallback needed)
- For multiline answers, showing the first line or a truncated version in the post-validation format is acceptable

## Dependencies

- `huh` library (github.com/charmbracelet/huh) - needs to be added as a dependency
- Existing prompt function signatures and callers
- Existing validation logic for prompts
- Existing default value and pre-selection functionality
- Existing commit service and workflow that calls prompt functions

## Clarifications

### Session 2025-01-27

- Q: How should multiline answers be displayed in the post-validation format? → A: Show the first line or a truncated version with indication of more content if needed (reasonable default)
- Q: What should happen when a prompt is cancelled (Escape or Ctrl+C)? → A: Handle cancellation gracefully - do not show post-validation format for cancelled prompts (reasonable default)
- Q: Should the post-validation format appear immediately after validation or after the entire form completes? → A: Appear immediately after each individual prompt is validated (user requirement: "once a prompt is validated by the user")
- Q: How should the prompt be replaced in the output after validation? → A: Clear the prompt UI and show the summary line below where the prompt was, keeping previous terminal output visible
- Q: How should the system handle terminal compatibility for the checkmark character? → A: Always use the checkmark character (✓) - assume all terminals support it
- Q: How should validation errors be displayed during prompt interaction? → A: Show validation errors inline within the prompt UI (below the input field, using `huh`'s built-in validation display)
- Q: Should each prompt function use its own individual `huh.Form` or combine related prompts into a single form? → A: Combine related prompts into a single `huh.Form` with multiple fields (e.g., all commit message fields in one form)
- Q: How should post-validation summary lines appear when prompts are combined into a single form with multiple fields? → A: Show summary lines progressively: clear each field's UI as it's completed and show its summary line, then move to next field

## Out of Scope

- Changing prompt function signatures or return types (unless required by `huh` library)
- Modifying validation logic or business rules (unless required for `huh` integration)
- Changing prompt content or wording
- Adding new prompt types beyond what `huh` supports
- Changing prompt behavior beyond visual design and rendering method
- Custom theming or styling beyond what `huh` provides (unless required for inline rendering)
- Accessibility features beyond what `huh` provides natively
