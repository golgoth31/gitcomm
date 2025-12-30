# Feature Specification: Unify Prompt Functions to Use Default Variants

**Feature Branch**: `016-unify-prompt-defaults`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "fix commit_service to only use prompts with default from prompts.go"

## Clarifications

### Session 2025-01-27

- Q: What should happen to the non-default prompt functions (`PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType`) after this refactoring? → A: Remove the non-default functions entirely (assumes no other callers)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Refactor Commit Service to Use Unified Prompt Functions (Priority: P1)

When a developer uses the commit service to create a commit message, the system should consistently use prompt functions that support default values, regardless of whether pre-filled data is available. This simplifies the codebase by removing conditional logic and ensures consistent behavior across all prompt interactions.

**Why this priority**: This is a code refactoring task that improves maintainability and reduces code complexity. It's the primary and only user story for this feature.

**Independent Test**: Can be fully tested by examining the `promptCommitMessage` function in `commit_service.go` and verifying that all prompt calls use the "WithDefault" or "WithPreselection" variants, even when no pre-filled data is provided (passing empty strings or appropriate defaults).

**Acceptance Scenarios**:

1. **Given** a commit service instance with no pre-filled commit message data, **When** the system prompts for commit message components, **Then** all prompts use the "WithDefault" or "WithPreselection" variants with empty string defaults
2. **Given** a commit service instance with pre-filled commit message data, **When** the system prompts for commit message components, **Then** all prompts use the "WithDefault" or "WithPreselection" variants with the pre-filled values as defaults
3. **Given** the refactored code, **When** a developer reviews the `promptCommitMessage` function, **Then** there are no conditional branches selecting between regular and "WithDefault" prompt variants

---

### Edge Cases

- What happens when an empty string is passed as default to a prompt that requires validation (e.g., subject)? The prompt should still validate and require user input
- How does the system handle the commit type prompt when no pre-filled type is available? The system should use `PromptCommitTypeWithPreselection` with an empty string as default (the form will allow user to select any type)
- What if a prompt function with default doesn't exist for a particular field? This should be identified during implementation and addressed
- What happens when non-default prompt functions are removed? All references to these functions must be updated to use the "WithDefault" variants, and the function definitions must be deleted from `prompts.go`

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The `promptCommitMessage` function MUST use `PromptCommitTypeWithPreselection` for commit type selection, passing an empty string or appropriate default when no pre-filled type is available
- **FR-002**: The `promptCommitMessage` function MUST use `PromptScopeWithDefault` for scope input, passing an empty string when no pre-filled scope is available
- **FR-003**: The `promptCommitMessage` function MUST use `PromptSubjectWithDefault` for subject input, passing an empty string when no pre-filled subject is available
- **FR-004**: The `promptCommitMessage` function MUST use `PromptBodyWithDefault` for body input, passing an empty string when no pre-filled body is available
- **FR-005**: The `promptCommitMessage` function MUST use `PromptFooterWithDefault` for footer input, passing an empty string when no pre-filled footer is available
- **FR-006**: The system MUST remove all conditional logic that selects between regular prompt functions and "WithDefault" variants based on the presence of pre-filled data
- **FR-007**: The system MUST maintain existing functionality - prompts should behave identically to current behavior when no defaults are provided (empty strings)
- **FR-008**: The system MUST remove the non-default prompt functions (`PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType`) from `prompts.go` after confirming they are not used elsewhere in the codebase

### Key Entities

- **CommitMessage**: Represents the structured commit message with type, scope, subject, body, and footer components
- **PrefilledCommitMessage**: Represents pre-populated commit message data that can be used as defaults in prompts

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Code complexity is reduced - the `promptCommitMessage` function contains zero conditional branches for selecting between regular and "WithDefault" prompt variants
- **SC-002**: Code maintainability improves - all prompt calls in `promptCommitMessage` use a consistent pattern (always using "WithDefault" variants)
- **SC-003**: Functionality is preserved - users experience identical behavior when creating commit messages with or without pre-filled data
- **SC-004**: Code review confirms no regression - all existing tests pass and commit message creation works correctly in both manual and AI-assisted workflows
