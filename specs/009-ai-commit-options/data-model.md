# Data Model: AI Commit Message Acceptance Options

**Feature**: 009-ai-commit-options
**Date**: 2025-01-27

## Entities

### AIMessageAcceptance

Represents the user's choice when presented with an AI-generated commit message.

**Type**: Custom enum-like type (`AIMessageAcceptance`)

**Values**:
- `AcceptAndCommit` (0): User wants to commit immediately with the AI message
- `AcceptAndEdit` (1): User wants to edit the AI message before committing
- `Reject` (2): User wants to reject the AI message and start over

**Relationships**:
- Returned by `PromptAIMessageAcceptanceOptions()` function
- Used by `CommitService.generateWithAI()` to determine workflow path

**Validation Rules**:
- Must be one of the three valid values
- Invalid input should prompt user again (with error message)

**State Transitions**:
- Created: When user selects an option
- Used: When workflow branches based on selection
- No state changes (immutable value)

---

### PrefilledCommitMessage

Represents a commit message structure where fields are populated with values from an AI-generated message, ready for user editing.

**Fields**:
- `Type` (string): Pre-filled commit type from AI message
- `Scope` (string): Pre-filled scope from AI message (may be empty)
- `Subject` (string): Pre-filled subject from AI message
- `Body` (string): Pre-filled body from AI message (may be empty)
- `Footer` (string): Pre-filled footer from AI message (may be empty)

**Relationships**:
- Created from parsed AI message when user selects "accept and edit"
- Used by `promptCommitMessage()` to pre-fill interactive prompts
- Converted to `model.CommitMessage` after user editing

**Validation Rules**:
- All fields are optional (may be empty)
- Type must match one of the available commit types if provided (for pre-selection)
- Subject should be non-empty for meaningful commits
- Values are used as-is (validation happens after user edits)

**State Transitions**:
- Created: When AI message is parsed and user selects "accept and edit"
- Pre-filled: When values are displayed in prompts
- Modified: When user edits any field
- Committed: When converted to `CommitMessage` and used for commit

---

## Data Flow

### Accept and Commit Directly Flow

```
1. AI message generated
   ↓
2. PromptAIMessageAcceptanceOptions() called
   ↓
3. User selects "AcceptAndCommit"
   ↓
4. AIMessageAcceptance value returned
   ↓
5. CommitService creates commit immediately with AI message
   ↓
6. CommitMessage created from parsed AI message
   ↓
7. Commit created via GitRepository.CreateCommit()
```

### Accept and Edit Flow

```
1. AI message generated
   ↓
2. PromptAIMessageAcceptanceOptions() called
   ↓
3. User selects "AcceptAndEdit"
   ↓
4. AIMessageAcceptance value returned
   ↓
5. AI message parsed into PrefilledCommitMessage
   ↓
6. promptCommitMessage() called with PrefilledCommitMessage
   ↓
7. Interactive prompts shown with pre-filled values:
   - Commit type: Pre-selected in list (if matches)
   - Scope: Default value shown (if provided)
   - Subject: Default value shown
   - Body: Pre-populated in multiline editor (if provided)
   - Footer: Pre-populated in multiline editor (if provided)
   ↓
8. User edits fields (or accepts defaults)
   ↓
9. CommitMessage created from user input
   ↓
10. Commit created via GitRepository.CreateCommit()
```

### Reject Flow

```
1. AI message generated
   ↓
2. PromptAIMessageAcceptanceOptions() called
   ↓
3. User selects "Reject"
   ↓
4. AIMessageAcceptance value returned
   ↓
5. User prompted: "Generate new AI message or manual input?"
   ↓
6a. If "new AI message":
    - Generate new AI message
    - Return to step 2
6b. If "manual input":
    - promptCommitMessage() called with empty PrefilledCommitMessage
    - User enters all fields manually
    - CommitMessage created
    - Commit created
```

---

## Default Values

When pre-filling fields from AI message:

- **Type**: Empty if AI type doesn't match available options (no pre-selection)
- **Scope**: Empty if AI message has no scope
- **Subject**: Empty if AI message parsing fails (should not happen for valid messages)
- **Body**: Empty if AI message has no body
- **Footer**: Empty if AI message has no footer

When user rejects and chooses manual input:

- All fields start empty (no pre-filling)

---

## Error Handling

**Invalid Acceptance Response**:
- User enters invalid option (not 1/2/3): Re-prompt with error message
- User cancels (Ctrl+C): Return error, restore staging state

**Commit Failure After Accept and Commit**:
- Error returned from `CreateCommit()`: Restore staging state, display error, prompt for retry/edit/cancel
- User chooses retry: Attempt commit again with same message
- User chooses edit: Fall back to "accept and edit" flow
- User chooses cancel: Restore staging state, exit workflow

**AI Generation Failure After Reject**:
- Error from AI provider: Display error message, fall back to manual input with empty fields
- User proceeds with manual input

**Cancellation During Edit**:
- User cancels at any point (Ctrl+C, Escape): Restore staging state, return error
- Workflow exits gracefully

---

## Relationships Diagram

```
AIMessageAcceptance
    │
    ├──> Returned by
    │    └──> PromptAIMessageAcceptanceOptions()
    │
    └──> Used by
         └──> CommitService.generateWithAI()
              │
              ├──> If AcceptAndCommit
              │    └──> Create commit immediately
              │
              ├──> If AcceptAndEdit
              │    └──> Create PrefilledCommitMessage
              │         └──> Used by
              │              └──> promptCommitMessage()
              │                   └──> Pre-fills interactive prompts
              │                        └──> User edits
              │                             └──> Creates CommitMessage
              │
              └──> If Reject
                   └──> Prompt for new AI or manual
                        ├──> New AI: Return to acceptance options
                        └──> Manual: Empty PrefilledCommitMessage
                             └──> promptCommitMessage()
                                  └──> Creates CommitMessage

PrefilledCommitMessage
    │
    ├──> Created from
    │    └──> Parsed AI message
    │
    └──> Converted to
         └──> model.CommitMessage
              └──> Used by
                   └──> GitRepository.CreateCommit()
```

---

## Integration with Existing Models

**CommitMessage** (existing):
- Used as the final structure for commit creation
- Created from `PrefilledCommitMessage` after user editing
- No changes to existing `CommitMessage` structure

**SelectListModel** (existing):
- Extended with `NewSelectListModelWithPreselection()` constructor
- `SelectedIndex` set based on pre-filled type
- No changes to existing model structure

**MultilineInputModel** (existing):
- Extended with `NewMultilineInputModelWithValue()` constructor
- Initial value set from pre-filled body/footer
- No changes to existing model structure
