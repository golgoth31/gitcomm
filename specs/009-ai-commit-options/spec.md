# Feature Specification: AI Commit Message Acceptance Options

**Feature Branch**: `009-ai-commit-options`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "when the AI message is accepted, commit directly @gitcomm/internal/service/commit_service.go:437-438"

## User Scenarios & Testing

### User Story 1 - Accept and Commit Directly (Priority: P1)

When a user receives an AI-generated commit message that they are satisfied with, they should be able to accept it and have the commit created immediately without any additional prompts or manual editing.

**Why this priority**: This is the most common use case - users who are happy with the AI-generated message want the fastest path to commit. This reduces friction and improves the user experience for the majority of AI-generated messages.

**Independent Test**: User runs gitcomm with AI enabled, receives an AI-generated message, selects "accept and commit", and the commit is created immediately with the AI message. This can be fully tested independently and delivers immediate value by reducing steps for satisfied users.

**Acceptance Scenarios**:

1. **Given** an AI-generated commit message is displayed, **When** user selects "accept and commit directly", **Then** the commit is created immediately with the AI message and the workflow completes
2. **Given** an AI-generated commit message that passes validation, **When** user selects "accept and commit directly", **Then** the commit is created without showing validation warnings or edit prompts
3. **Given** an AI-generated commit message, **When** user selects "accept and commit directly", **Then** the commit uses the same author and signing configuration as manual commits

---

### User Story 2 - Accept and Edit (Priority: P2)

When a user receives an AI-generated commit message that is mostly correct but needs minor adjustments, they should be able to accept it and then edit specific fields (type, scope, subject, body, footer) with the AI values pre-filled.

**Why this priority**: Many users will want to make small tweaks to AI-generated messages. Pre-filling the fields with AI values saves time compared to starting from scratch while still allowing customization.

**Independent Test**: User runs gitcomm with AI enabled, receives an AI-generated message, selects "accept and edit", and is presented with the manual commit message prompts where all fields are pre-filled with values from the AI message. User can then modify any field and proceed to commit. This can be fully tested independently and delivers value by streamlining the editing workflow.

**Acceptance Scenarios**:

1. **Given** an AI-generated commit message is displayed, **When** user selects "accept and edit", **Then** the manual commit message prompts are shown with all fields (type, scope, subject, body, footer) pre-filled with values from the AI message, with the commit type automatically selected in the interactive list
2. **Given** user is editing a pre-filled commit message, **When** user modifies any field, **Then** the modified value is used in the final commit
3. **Given** user is editing a pre-filled commit message, **When** user leaves a field unchanged, **Then** the original AI value is used in the final commit
4. **Given** an AI message has an empty scope, **When** user selects "accept and edit", **Then** the scope field is shown as empty (not pre-filled) and user can optionally add one

---

### User Story 3 - Reject and Start Over (Priority: P3)

When a user receives an AI-generated commit message that is not suitable, they should be able to reject it and start the commit message creation process from scratch (either with a new AI generation or manual input).

**Why this priority**: This maintains the existing behavior for users who are not satisfied with the AI message, ensuring backward compatibility and providing an escape hatch when AI generates poor messages.

**Independent Test**: User runs gitcomm with AI enabled, receives an AI-generated message, selects "reject", and is presented with the option to generate a new AI message or proceed with manual input. This can be fully tested independently and maintains existing functionality.

**Acceptance Scenarios**:

1. **Given** an AI-generated commit message is displayed, **When** user selects "reject", **Then** the user is prompted to choose between generating a new AI message or proceeding with manual input
2. **Given** user rejected an AI message, **When** user chooses to generate a new AI message, **Then** a new AI message is generated and the acceptance options are shown again
3. **Given** user rejected an AI message, **When** user chooses manual input, **Then** the manual commit message prompts are shown with empty fields (no pre-filling)

---

### Edge Cases

- What happens when AI message has validation errors and user selects "accept and commit directly"? (Should be prevented - only valid messages can be accepted directly)
- What happens when AI message parsing fails partially (e.g., has type and subject but body parsing fails)? (Should handle gracefully - use available fields, leave others empty)
- What happens when user cancels during "accept and edit" flow? (Staging state is restored if user cancels at any point during field editing - Ctrl+C, Escape, or explicit cancel)
- What happens when AI message has very long fields that exceed recommended lengths? (Should still pre-fill but show warnings as per existing validation)
- What happens when AI message has special characters or formatting that needs escaping? (Should preserve the values as-is, let validation handle issues)
- What happens when commit fails after user selects "accept and commit directly"? (Staging state is restored, error message displayed, user can retry commit or edit message)
- What happens when AI-generated commit type doesn't match any option in the selection list? (Type selection list shown with no pre-selection, user must choose from available options)
- What happens if AI generation fails when user chooses to generate a new message after rejection? (Fall back to manual input with empty fields, display error message explaining AI generation failed)

## Requirements

### Functional Requirements

- **FR-001**: System MUST provide three distinct options when displaying an AI-generated commit message: "accept and commit directly", "accept and edit", and "reject"
- **FR-002**: System MUST only allow "accept and commit directly" option when the AI message passes all validation checks
- **FR-003**: System MUST create a commit immediately when user selects "accept and commit directly" option
- **FR-012**: System MUST restore staging state and display error message if commit fails after "accept and commit directly", allowing user to retry commit or edit message
- **FR-004**: System MUST pre-fill all commit message fields (type, scope, subject, body, footer) with AI message values when user selects "accept and edit" option
- **FR-013**: System MUST automatically select the matching commit type in the interactive selection list when pre-filling during "accept and edit" flow, if the AI type matches an available option
- **FR-014**: System MUST show the type selection list with no pre-selection if the AI-generated commit type doesn't match any available option, requiring user to choose
- **FR-005**: System MUST allow user to modify any pre-filled field during the "accept and edit" flow
- **FR-006**: System MUST preserve unchanged pre-filled values in the final commit when user selects "accept and edit"
- **FR-007**: System MUST prompt user to choose between new AI generation or manual input when user selects "reject" option
- **FR-015**: System MUST fall back to manual input with empty fields and display error message if AI generation fails when user chooses to generate a new message after rejection
- **FR-008**: System MUST show empty fields (no pre-filling) when user proceeds with manual input after rejecting an AI message
- **FR-009**: System MUST handle partial AI message parsing gracefully (use available fields, leave missing fields empty)
- **FR-010**: System MUST maintain existing validation behavior for pre-filled fields during "accept and edit" flow
- **FR-011**: System MUST restore staging state if user cancels at any point during field editing in "accept and edit" flow (Ctrl+C, Escape, or explicit cancel)

### Key Entities

- **AI Message Acceptance Response**: Represents the user's choice when presented with an AI-generated commit message. Has three possible values: accept-and-commit, accept-and-edit, reject
- **Pre-filled Commit Message**: Represents a commit message structure where fields are populated with values from an AI-generated message, ready for user editing

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can commit with an AI-generated message in under 5 seconds when selecting "accept and commit directly" (measured from message display to commit completion)
- **SC-002**: Users can edit and commit a pre-filled AI message in under 30 seconds when selecting "accept and edit" (measured from selection to commit completion, including editing time)
- **SC-003**: 80% of users who accept AI messages choose "accept and commit directly" over "accept and edit" (indicating AI message quality is sufficient for direct use)
- **SC-004**: Users who select "accept and edit" modify on average fewer than 2 fields per message (indicating pre-filling is effective)
- **SC-005**: Zero increase in commit errors or validation failures compared to current manual input flow (maintaining quality while improving speed)

## Assumptions

- AI-generated messages that pass validation are suitable for direct commit in most cases
- Users who want to edit will typically make minor adjustments rather than major rewrites
- Pre-filling fields with AI values will reduce user input time compared to manual entry
- The existing commit message validation and formatting logic will work correctly with pre-filled values
- Users understand the difference between the three acceptance options

## Dependencies

- Existing AI message generation functionality
- Existing commit message validation system
- Existing manual commit message prompt system
- Existing commit creation workflow

## Clarifications

### Session 2025-01-27

- Q: How should the system handle commit failures after the user has selected "accept and commit directly"? → A: Restore staging state, display error message, allow user to retry commit or edit message
- Q: How should the commit type be pre-filled when the user selects "accept and edit"? → A: Automatically select the matching type in the interactive list, showing it as pre-selected when the list appears
- Q: What should happen when the AI-generated commit type doesn't match any option in the interactive selection list? → A: Show type selection list with no pre-selection, user must choose from available options
- Q: When should cancellation during "accept and edit" trigger staging state restoration? → A: Restore staging state if user cancels at any point during field editing (Ctrl+C, Escape, or explicit cancel)
- Q: What should happen if AI message generation fails when the user chooses to generate a new message after rejecting the first one? → A: Fall back to manual input with empty fields, display error message explaining AI generation failed

## Out of Scope

- Changing the AI message generation algorithm or prompts
- Modifying commit message validation rules
- Adding new commit message fields
- Changing the commit creation process itself (author, signing, etc.)
- Batch processing of multiple AI messages
