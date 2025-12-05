# Contract: Prompt Functions API

**Feature**: 013-huh-prompts
**Date**: 2025-01-27
**Type**: Function API Contract

## Overview

This contract defines the public API for prompt functions in `internal/ui/prompts.go`. These functions must maintain backward compatibility while migrating internal implementation to use the `huh` library.

## Contract Guarantees

### Backward Compatibility

- All function signatures remain unchanged
- Return types and error handling behavior preserved
- Function behavior (validation, defaults, pre-selection) preserved
- Callers require no changes

### Implementation Requirements

- Internal implementation uses `huh` library
- Prompts render inline (no alt screen)
- Post-validation summary lines displayed: `✓ <title>: <value>`
- Validation errors displayed inline within prompt UI

## Function Contracts

### Text Input Prompts

#### `PromptScope(reader *bufio.Reader) (string, error)`

**Purpose**: Prompt user for optional commit scope.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused, kept for compatibility)

**Returns**:
- `string`: User-entered scope value (empty string if skipped)
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays prompt: "Scope (optional)"
- User can enter value or press Enter to skip
- If user cancels (Escape/Ctrl+C), returns error
- After validation, displays: `✓ Scope (optional): <value>` (or empty if skipped)
- Returns empty string if user skips

**Preconditions**: None

**Postconditions**:
- If successful: Summary line displayed, value returned
- If cancelled: Error returned, no summary line

---

#### `PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Purpose**: Prompt user for commit scope with pre-filled default value.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `defaultValue string`: Default value to pre-fill

**Returns**:
- `string`: User-entered scope value or default if empty input
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays prompt with default value pre-filled
- User can accept default (Enter) or modify value
- If user enters empty string, returns `defaultValue`
- After validation, displays: `✓ Scope: <final-value>`
- Returns default value if user accepts without modification

**Preconditions**: None

**Postconditions**: Same as `PromptScope`

---

#### `PromptSubject(reader *bufio.Reader) (string, error)`

**Purpose**: Prompt user for required commit subject.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `string`: User-entered subject (trimmed)
- `error`: Non-nil if prompt was cancelled, validation failed, or error occurred

**Behavior**:
- Displays prompt: "Subject (required)"
- Validates: Subject cannot be empty or whitespace-only
- Shows validation error inline if empty
- After validation, displays: `✓ Subject (required): <value>`
- Returns trimmed subject value

**Preconditions**: None

**Postconditions**:
- If successful: Summary line displayed, non-empty trimmed value returned
- If validation fails: Error shown inline, user can correct
- If cancelled: Error returned, no summary line

---

#### `PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Purpose**: Prompt user for commit subject with pre-filled default value.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `defaultValue string`: Default value to pre-fill

**Returns**:
- `string`: User-entered subject or default if empty input
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Similar to `PromptSubject` but with default value
- If user enters empty string and default exists, returns `defaultValue`
- If both empty, validation error shown

**Preconditions**: None

**Postconditions**: Same as `PromptSubject`

---

### Multiline Input Prompts

#### `PromptBody(reader *bufio.Reader) (string, error)`

**Purpose**: Prompt user for optional commit body (multiline).

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `string`: User-entered body text (multiline)
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays multiline input prompt: "Body"
- User can enter multiple lines
- User validates with double Enter on empty line (or `huh` default)
- If body > 320 characters, shows warning (handled by existing logic)
- After validation, displays: `✓ Body: <first-line>` (truncated if long)
- Returns full body text

**Preconditions**: None

**Postconditions**:
- If successful: Summary line displayed (may truncate), full value returned
- If cancelled: Error returned, no summary line

---

#### `PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Purpose**: Prompt user for commit body with pre-filled default value.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `defaultValue string`: Default multiline text to pre-fill

**Returns**:
- `string`: User-entered body text or default
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**: Similar to `PromptBody` with default value pre-filled

---

#### `PromptFooter(reader *bufio.Reader) (string, error)`

**Purpose**: Prompt user for optional commit footer (multiline).

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `string`: User-entered footer text (multiline)
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**: Similar to `PromptBody` but for footer field

---

#### `PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Purpose**: Prompt user for commit footer with pre-filled default value.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `defaultValue string`: Default footer text to pre-fill

**Returns**:
- `string`: User-entered footer text or default
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**: Similar to `PromptFooter` with default value

---

### Selection Prompts

#### `PromptCommitType(reader *bufio.Reader) (string, error)`

**Purpose**: Prompt user to select commit type from predefined list.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `string`: Selected commit type (e.g., "feat", "fix", "docs")
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays select list: "Choose a type"
- Shows predefined commit types (feat, fix, docs, style, refactor, perf, test, chore, etc.)
- User navigates and selects with Enter
- After validation, displays: `✓ Choose a type: <selected-type>`
- Returns selected type string

**Preconditions**: None

**Postconditions**:
- If successful: Summary line displayed, selected type returned
- If cancelled: Error returned, no summary line

---

#### `PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error)`

**Purpose**: Prompt user to select commit type with pre-selected value.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `preselectedType string`: Pre-selected commit type

**Returns**:
- `string`: Selected commit type
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**: Similar to `PromptCommitType` but with pre-selected value highlighted

---

### Confirmation Prompts

#### `PromptEmptyCommit(reader *bufio.Reader) (bool, error)`

**Purpose**: Prompt user to confirm creating an empty commit.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `bool`: `true` if user confirms, `false` if declines
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays confirmation: "No changes detected. Create an empty commit?"
- User selects Yes/No
- After validation, displays: `✓ No changes detected. Create an empty commit?: <Yes/No>`
- Returns boolean result

---

#### `PromptConfirm(reader *bufio.Reader, message string) (bool, error)`

**Purpose**: General confirmation prompt with custom message.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `message string`: Confirmation message to display

**Returns**:
- `bool`: `true` if user confirms, `false` if declines
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays confirmation with custom message
- After validation, displays: `✓ <message>: <Yes/No>`
- Returns boolean result

---

#### `PromptAIUsage(reader *bufio.Reader, tokenCount int) (bool, error)`

**Purpose**: Prompt user to choose whether to use AI for commit message generation.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `tokenCount int`: Estimated token count to display

**Returns**:
- `bool`: `true` if user wants to use AI, `false` otherwise
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays: "Estimated tokens: <count>\nUse AI to generate commit message?"
- Defaults to "Yes" (true)
- After validation, displays: `✓ Use AI to generate commit message?: <Yes/No>`
- Returns boolean result

---

#### `PromptAIMessageEdit(reader *bufio.Reader, errors []string) (bool, error)`

**Purpose**: Prompt user to edit AI message when validation errors are found.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `errors []string`: List of validation error messages

**Returns**:
- `bool`: `true` if user wants to edit, `false` to use as-is
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays validation errors and asks: "Edit the message? (y=edit, n=use as-is)"
- Defaults to "Yes" (edit) when errors present
- After validation, displays summary line
- Returns boolean result

---

#### `PromptRejectChoice(reader *bufio.Reader) (bool, error)`

**Purpose**: Prompt user to choose between generating new AI message or manual input.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `bool`: `true` to generate new AI message, `false` for manual input
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays: "Generate new AI message? (y=new AI, n=manual input)"
- Defaults to "Yes" (generate new AI)
- After validation, displays summary line
- Returns boolean result

---

### Multi-Choice Prompts

#### `PromptAIMessageAcceptanceOptions(reader *bufio.Reader, message string) (AIMessageAcceptance, error)`

**Purpose**: Prompt user to choose action for AI-generated commit message.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)
- `message string`: AI-generated commit message to display

**Returns**:
- `AIMessageAcceptance`: User's choice (AcceptAndCommit, AcceptAndEdit, or Reject)
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays AI message and options:
  - 1. Accept and commit directly
  - 2. Accept and edit
  - 3. Reject
- User selects option
- After validation, displays: `✓ <prompt-title>: <selected-option>`
- Returns enum value

---

#### `PromptCommitFailureChoice(reader *bufio.Reader) (CommitFailureChoice, error)`

**Purpose**: Prompt user to choose action when commit fails.

**Parameters**:
- `reader *bufio.Reader`: Input reader (may be unused)

**Returns**:
- `CommitFailureChoice`: User's choice (RetryCommit, EditMessage, or CancelCommit)
- `error`: Non-nil if prompt was cancelled or error occurred

**Behavior**:
- Displays options:
  - 1. Retry commit
  - 2. Edit message
  - 3. Cancel
- User selects option
- After validation, displays summary line
- Returns enum value

---

## Error Handling

### Error Types

- **Cancellation Error**: Returned when user cancels prompt (Escape/Ctrl+C)
  - Format: `fmt.Errorf("<prompt-name> input cancelled")`
- **Validation Error**: Returned when validation fails (for required fields)
  - Format: `fmt.Errorf("<validation-message>")`
- **Execution Error**: Returned when form execution fails
  - Format: `fmt.Errorf("failed to run <prompt-name>: %w", err)`

### Error Behavior

- Errors are wrapped with context using `fmt.Errorf` with `%w` verb
- Callers can use `errors.Is()` or `errors.As()` to check error types
- Cancellation errors should be handled gracefully by callers

## Post-Validation Display Format

All successfully validated prompts must display a summary line:

```
✓ <prompt-title>: <validated-value>
```

Where:
- `✓` is a green checkmark character
- `<prompt-title>` is the prompt title/question
- `<validated-value>` is the user's validated input/selection

For multiline values, the summary may truncate to first line with indication of more content.

## Testing Requirements

- All functions must have unit tests
- Tests must verify function signatures and return types
- Tests must verify post-validation display format
- Tests must verify error handling (cancellation, validation)
- Tests must verify default values and pre-selection
- Integration tests must verify end-to-end prompt flow
