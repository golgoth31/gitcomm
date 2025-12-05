# Feature Specification: Rewrite All CLI Prompts

**Feature Branch**: `001-rewrite-cli-prompts`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "rewrite all cli prompts: - you MUST use bubletea for each prompts - you MUST NOT use altscreen - all prompts must start with a bleue '?' - when a prompt has been answered, the '?' is replaced by a green checkmark - for select list, the selected value will be printed right after the select list title - for multilines, once validated, the result is writen under the prompt"

## User Scenarios & Testing

### User Story 1 - Consistent Prompt Visual Design (Priority: P1)

When a user interacts with any CLI prompt, they should see a consistent visual design where all prompts start with a blue question mark, and once answered, the question mark is replaced with a green checkmark.

**Why this priority**: This establishes the visual foundation for all prompts. Without consistent visual design, the user experience will be fragmented and confusing. This is the most critical aspect as it affects every user interaction.

**Independent Test**: User runs gitcomm and interacts with any prompt (scope, subject, type selection, etc.). All prompts display a blue '?' at the start, and after answering, the '?' is replaced with a green checkmark. This can be fully tested independently and delivers immediate visual consistency.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit scope, **When** the prompt is displayed, **Then** it starts with a blue '?' character
2. **Given** user enters a scope value and confirms (presses Enter), **When** the prompt completes, **Then** the blue '?' is replaced with a green checkmark
3. **Given** user is prompted for commit type selection, **When** the prompt is displayed, **Then** it starts with a blue '?' character
4. **Given** user selects a commit type and confirms (presses Enter), **When** the selection completes, **Then** the blue '?' is replaced with a green checkmark and the selected value is displayed after the prompt title

---

### User Story 2 - Bubbletea for All Prompts (Priority: P1)

When a user interacts with any CLI prompt, the prompt should use bubbletea for rendering and interaction, ensuring consistent behavior and visual feedback across all prompt types.

**Why this priority**: Using bubbletea for all prompts ensures consistent interaction patterns, better visual feedback, and unified keyboard handling. This is critical as it affects the core user experience for every prompt type.

**Independent Test**: User runs gitcomm and interacts with all prompt types (text input, select lists, multiline input, yes/no choices). All prompts use bubbletea for rendering and interaction. This can be fully tested independently and delivers consistent interaction patterns.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit scope (text input), **When** the prompt is displayed, **Then** it uses bubbletea for rendering and input handling
2. **Given** user is prompted for commit subject (text input with validation), **When** the prompt is displayed, **Then** it uses bubbletea for rendering and input handling
3. **Given** user is prompted for commit type (select list), **When** the prompt is displayed, **Then** it uses bubbletea for rendering and selection
4. **Given** user is prompted for commit body (multiline input), **When** the prompt is displayed, **Then** it uses bubbletea for rendering and multiline input handling
5. **Given** user is prompted for yes/no choice (AI usage, empty commit, etc.), **When** the prompt is displayed, **Then** it uses bubbletea for rendering and choice selection

---

### User Story 3 - No AltScreen Usage (Priority: P2)

When a user interacts with CLI prompts, the prompts should render inline without using alt screen mode, keeping the terminal history visible and maintaining context.

**Why this priority**: Not using alt screen keeps the terminal history visible, which helps users maintain context of previous prompts and system output. This improves usability, especially when users need to reference earlier information.

**Independent Test**: User runs gitcomm and interacts with prompts. All prompts render inline without switching to alt screen mode. Terminal history remains visible throughout the interaction. This can be fully tested independently and delivers better context preservation.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit type (select list), **When** the prompt is displayed, **Then** it renders inline without using alt screen mode
2. **Given** user is prompted for commit body (multiline input), **When** the prompt is displayed, **Then** it renders inline without using alt screen mode
3. **Given** user interacts with any prompt, **When** the prompt is active, **Then** previous terminal output remains visible above the prompt

---

### User Story 4 - Select List Value Display (Priority: P2)

When a user selects a value from a select list prompt, the selected value should be displayed immediately after the prompt title, providing clear visual confirmation of the selection.

**Why this priority**: Displaying the selected value after the title provides immediate visual feedback and confirmation, improving user confidence in their selection. This is important for select lists which are used for commit type selection.

**Independent Test**: User runs gitcomm and selects a commit type from the select list. After selection, the selected value (e.g., "feat") is displayed immediately after the prompt title. This can be fully tested independently and delivers clear selection feedback.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit type with title "? Choose a type(<scope>):", **When** user selects "feat" and confirms (presses Enter), **Then** the prompt displays "✓ Choose a type(<scope>): feat" with the selected value after the title
2. **Given** user is prompted for commit type with preselection, **When** the prompt is displayed, **Then** the preselected value is shown after the title (before confirmation)
3. **Given** user navigates through options in the list, **When** user has not yet confirmed, **Then** the value after the title shows the currently highlighted option (if preselection exists) or remains empty until confirmation

---

### User Story 5 - Multiline Input Result Display (Priority: P2)

When a user completes a multiline input prompt (body or footer), the entered text should be displayed under the prompt title after validation, providing clear confirmation of what was entered.

**Why this priority**: Displaying the multiline result under the prompt provides visual confirmation of the entered text, especially important for longer multiline inputs where users may want to verify what they entered.

**Independent Test**: User runs gitcomm and enters text in the commit body multiline prompt. After validation (double Enter on empty line), the entered text is displayed under the prompt title. This can be fully tested independently and delivers clear input confirmation.

**Acceptance Scenarios**:

1. **Given** user is prompted for commit body with title "? Body:", **When** user enters multiline text and validates, **Then** the entered text is displayed under "? Body:" (or "✓ Body:" after completion)
2. **Given** user is prompted for commit footer with title "? Footer:", **When** user enters multiline text and validates, **Then** the entered text is displayed under "? Footer:" (or "✓ Footer:" after completion)
3. **Given** user enters empty multiline input, **When** user validates, **Then** the prompt shows completion with empty result indication

---

### Edge Cases

- What happens when a prompt is cancelled (Escape key)? (The prompt should show red 'X' cancellation indicator in place of the blue '?')
- What happens when terminal width is very narrow? (Prompts should handle narrow terminals gracefully, wrapping text appropriately)
- What happens when multiline input is very long? (Result display should show full text with line wrapping to respect terminal width)
- What happens when select list has many items and terminal is short? (List should scroll appropriately without alt screen)
- What happens when user presses Ctrl+C during a prompt? (Prompt should handle interruption gracefully, showing cancellation state)
- What happens when a prompt has a default value? (Default value should be displayed appropriately with the blue '?' indicator)
- What happens when validation fails in a prompt? (Prompt should show yellow/orange warning indicator with error message displayed, allowing user to correct input)

## Requirements

### Functional Requirements

- **FR-001**: System MUST use bubbletea for rendering and interaction for all CLI prompts (text input, select lists, multiline input, yes/no choices)
- **FR-002**: System MUST NOT use alt screen mode for any prompt (all prompts render inline)
- **FR-003**: System MUST display a blue '?' character at the start of every prompt title
- **FR-004**: System MUST replace the blue '?' with a green checkmark when a prompt is successfully answered/completed (after user confirms with Enter/confirmation action)
- **FR-005**: System MUST display the selected value after the prompt title for select list prompts (after user confirms selection with Enter)
- **FR-006**: System MUST display the entered text under the prompt title for multiline input prompts after validation, with line wrapping for long lines to respect terminal width
- **FR-007**: System MUST maintain terminal history visibility throughout all prompt interactions
- **FR-008**: System MUST handle prompt cancellation (Escape key) gracefully, showing red 'X' or similar cancellation indicator in place of the blue '?'
- **FR-009**: System MUST handle Ctrl+C interruption during prompts gracefully
- **FR-010**: System MUST support all existing prompt types: scope, subject, body, footer, commit type, AI usage, empty commit confirmation, AI message acceptance, AI message edit, commit failure choice, reject choice
- **FR-011**: System MUST maintain backward compatibility with existing prompt function signatures (callers should not need changes)
- **FR-012**: System MUST preserve all existing prompt functionality (validation, default values, pre-selection, etc.)
- **FR-013**: System MUST show yellow/orange warning indicator with error message when validation fails in a prompt

### Key Entities

- **Prompt State**: Represents the current state of a prompt (pending, active, completed, cancelled)
- **Visual Indicator**: Represents the prompt indicator character ('?' for pending, '✓' for completed, '✗' or 'X' for cancelled, '⚠' or similar for validation error)
- **Prompt Result**: Represents the result of a prompt interaction (value selected/entered, cancellation, error)

## Success Criteria

### Measurable Outcomes

- **SC-001**: All prompts display blue '?' at start and green checkmark on completion (100% consistency across all prompt types)
- **SC-002**: All prompts use bubbletea for rendering (100% of prompts migrated from bufio.Reader to bubbletea)
- **SC-003**: No prompts use alt screen mode (0 prompts using alt screen)
- **SC-004**: Select list prompts display selected value after title (100% of select list prompts)
- **SC-005**: Multiline input prompts display result under title after validation (100% of multiline prompts)
- **SC-006**: Terminal history remains visible during all prompt interactions (user can scroll up to see previous output)
- **SC-007**: All existing prompt functionality preserved (100% backward compatibility with existing callers)

## Assumptions

- Users prefer inline prompts over alt screen mode for better context preservation
- Blue '?' and green checkmark provide sufficient visual feedback for prompt state
- Bubbletea can handle all prompt types (text input, select lists, multiline, yes/no) without alt screen
- Existing prompt function signatures can be maintained while changing internal implementation to bubbletea
- Terminal width/height constraints can be handled gracefully by bubbletea without alt screen
- All existing validation and business logic can be preserved when migrating to bubbletea

## Dependencies

- Existing bubbletea library (already in use for select lists and multiline input)
- Existing prompt function signatures and callers
- Existing validation logic for prompts
- Existing default value and pre-selection functionality

## Clarifications

### Session 2025-01-27

- Q: When should the blue '?' change to a green checkmark? → A: Checkmark appears only after user confirms (presses Enter/confirms selection)
- Q: What visual indicator should show when a prompt is cancelled (Escape)? → A: Show red 'X' or similar cancellation indicator
- Q: When should the selected value appear after the prompt title for select lists? → A: Show selected value only after user confirms selection (presses Enter)
- Q: What visual indicator should show when validation fails in a prompt? → A: Show yellow/orange warning indicator with error message
- Q: How should multiline input result be displayed under the prompt title? → A: Display full text with line wrapping for long lines (respect terminal width)

## Out of Scope

- Changing prompt function signatures or return types
- Modifying validation logic or business rules
- Changing prompt content or wording
- Adding new prompt types
- Changing prompt behavior beyond visual design and rendering method
