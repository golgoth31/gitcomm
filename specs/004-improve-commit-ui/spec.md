# Feature Specification: Improved Commit Message UI

**Feature Branch**: `004-improve-commit-ui`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "the \"commit type\" text and select type field MUST be replaced by a select list with the style show in the picture. The body and Footer field MUST allow multiline"

## Clarifications

### Session 2025-01-27

- Q: What is the exact key combination or method users should use to complete multiline body and footer input? → A: Double Enter (press Enter twice on empty line) - Simple, no special keys needed
- Q: How should the system handle body and footer input that contains only whitespace characters? → A: Treat as empty (trimmed) - Whitespace-only input is treated as if the field was skipped
- Q: Which commit type should be pre-selected when the interactive select list first appears? → A: First option (feat) - Always start with the first option in the list pre-selected
- Q: Should the commit type select list support letter-based navigation (typing a letter to jump to matching options)? → A: Required feature - Letter-based navigation must be implemented (typing a letter jumps to first matching option)
- Q: How should the system distinguish between a blank line that's part of the content versus a blank line that signals completion? → A: Two consecutive empty lines signal completion - Single blank lines are preserved as content, two empty lines in a row complete the input

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Interactive Commit Type Selection (Priority: P1)

A developer wants to select a commit type using an interactive select list with visual feedback (checkmarks, highlighting) instead of typing numbers or text, making the selection process more intuitive and error-resistant.

**Why this priority**: This is the core UI improvement - replacing the current numbered list with an interactive select list. This directly addresses the user's requirement and significantly improves the user experience for the most common interaction in the commit workflow.

**Independent Test**: Can be fully tested by running the CLI and verifying that the commit type selection displays as an interactive select list with visual indicators (checkmarks, highlighting) and allows navigation with arrow keys, with the selected option clearly indicated.

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI to create a commit, **When** the commit type prompt appears, **Then** an interactive select list is displayed showing all available commit types (feat, fix, docs, style, refactor, test, chore, version) with descriptions, with the first option (feat) pre-selected
2. **Given** the commit type select list is displayed, **When** the developer navigates with arrow keys, **Then** the selection highlight moves between options and the currently selected option is clearly indicated (e.g., with a checkmark or highlight)
3. **Given** the commit type select list is displayed, **When** the developer presses Enter, **Then** the selected commit type is accepted and the workflow continues
4. **Given** the commit type select list is displayed, **When** the developer types a letter matching a commit type, **Then** the selection jumps to the first matching option

---

### User Story 2 - Multiline Body Input (Priority: P1)

A developer wants to enter a commit body that spans multiple lines, allowing for detailed explanations, bullet points, or formatted text in the commit message body.

**Why this priority**: The body field is critical for detailed commit messages. While the current implementation may support multiline, it needs to be clearly documented and work reliably. This is essential for proper commit message documentation.

**Independent Test**: Can be fully tested by running the CLI and verifying that the body prompt accepts multiple lines of input, allows line breaks, and preserves the multiline format in the final commit message.

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI to create a commit, **When** the body prompt appears, **Then** the input field accepts multiple lines of text
2. **Given** the body input field is active, **When** the developer presses Enter, **Then** a new line is created (not submitted) and the cursor moves to the next line
3. **Given** the body input field is active, **When** the developer enters multiple lines of text, **Then** all lines are preserved in the commit message body
4. **Given** the body input field is active, **When** the developer presses Enter twice on consecutive empty lines, **Then** the body input is completed and the workflow continues
5. **Given** the body input field is active, **When** the developer enters a single blank line within content, **Then** the blank line is preserved as part of the body content

---

### User Story 3 - Multiline Footer Input (Priority: P1)

A developer wants to enter a commit footer that spans multiple lines, allowing for multiple footer entries (e.g., multiple "Fixes #123", "Closes #456" entries) or detailed footer information.

**Why this priority**: The footer field currently only accepts single-line input, but Conventional Commits allows multiple footer entries. This improvement enables proper use of footers for issue tracking and breaking changes documentation.

**Independent Test**: Can be fully tested by running the CLI and verifying that the footer prompt accepts multiple lines of input, allows line breaks, and preserves the multiline format in the final commit message.

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI to create a commit, **When** the footer prompt appears, **Then** the input field accepts multiple lines of text (not just single line)
2. **Given** the footer input field is active, **When** the developer presses Enter, **Then** a new line is created (not submitted) and the cursor moves to the next line
3. **Given** the footer input field is active, **When** the developer enters multiple lines of text, **Then** all lines are preserved in the commit message footer
4. **Given** the footer input field is active, **When** the developer presses Enter twice on consecutive empty lines, **Then** the footer input is completed and the workflow continues
5. **Given** the footer input field is active, **When** the developer enters a single blank line within content, **Then** the blank line is preserved as part of the footer content

---

### Edge Cases

- What happens when the user presses Escape during commit type selection? → Should cancel the commit workflow and restore staging state
- What happens when the terminal is resized during interactive selection? → UI should adapt gracefully or maintain current selection
- What happens when body/footer input exceeds recommended length? → Should show warning but allow continuation (existing behavior)
- What happens when body/footer contains only whitespace? → Should be treated as empty (trimmed) - whitespace-only input is treated as if the field was skipped
- What happens when user wants to skip body/footer entirely? → Should allow empty input (press Enter twice immediately on empty prompt to skip)
- How are blank lines within multiline content handled? → Single blank lines are preserved as content; two consecutive empty lines signal completion
- How are special characters handled in multiline body/footer? → Should preserve all characters including newlines, tabs, etc.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display commit type selection as an interactive select list with visual indicators (checkmarks, highlighting) instead of a numbered text list
- **FR-002**: System MUST allow navigation of commit type selection using arrow keys (up/down)
- **FR-003**: System MUST clearly indicate the currently selected commit type option (e.g., with checkmark, highlight, or other visual indicator), with the first option (feat) pre-selected when the list first appears
- **FR-004**: System MUST display commit type descriptions alongside type names in the select list (e.g., "feat [new feature]", "fix [bug fix]")
- **FR-005**: System MUST accept commit type selection when user presses Enter on a selected option
- **FR-014**: System MUST support letter-based navigation in commit type selection (typing a letter jumps to the first matching option)
- **FR-006**: System MUST allow body input to span multiple lines with line breaks preserved
- **FR-007**: System MUST allow footer input to span multiple lines with line breaks preserved
- **FR-008**: System MUST allow users to complete multiline body input by pressing Enter twice on consecutive empty lines
- **FR-009**: System MUST allow users to complete multiline footer input by pressing Enter twice on consecutive empty lines
- **FR-015**: System MUST preserve single blank lines within body and footer content (only two consecutive empty lines signal completion)
- **FR-010**: System MUST preserve all line breaks and whitespace in body and footer when saving to commit message
- **FR-011**: System MUST allow users to skip body input (empty body) when using multiline input
- **FR-012**: System MUST allow users to skip footer input (empty footer) when using multiline input
- **FR-013**: System MUST treat body and footer input containing only whitespace characters as empty (trimmed)

### Key Entities *(include if feature involves data)*

- **Commit Type Selection UI**: Represents the interactive select list interface for choosing commit types, with visual feedback and navigation capabilities
- **Multiline Input Field**: Represents input fields (body, footer) that accept and preserve multiple lines of text with line breaks

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can select a commit type using arrow keys in under 2 seconds (no typing required)
- **SC-002**: 100% of commit type selections use the interactive select list (no fallback to text input)
- **SC-003**: Users can enter multiline body text with line breaks preserved in 100% of cases
- **SC-004**: Users can enter multiline footer text with line breaks preserved in 100% of cases
- **SC-005**: Body and footer multiline input completion is successful on first attempt for 95% of users (clear completion mechanism)
- **SC-006**: Commit messages with multiline body/footer are correctly formatted and committed without data loss
