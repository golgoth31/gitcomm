# Contract: UI Prompt Extension for AI Message Acceptance Options

**Feature**: 009-ai-commit-options
**Date**: 2025-01-27
**Component**: `internal/ui/prompts.go`

## Overview

This contract defines the extension to the UI prompt system to support three-option acceptance of AI-generated commit messages and pre-filling of commit message fields.

## Interface: PromptAIMessageAcceptanceOptions

### Signature

```go
func PromptAIMessageAcceptanceOptions(reader *bufio.Reader) (AIMessageAcceptance, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader for user responses

### Returns

- `AIMessageAcceptance`: User's selection (AcceptAndCommit, AcceptAndEdit, or Reject)
- `error`: Error if input reading fails or user cancels

### Behavior

1. **Display AI Message**: Shows the AI-generated commit message to the user
2. **Present Options**: Displays three numbered options:
   - `1`: Accept and commit directly
   - `2`: Accept and edit
   - `3`: Reject
3. **Read Input**: Waits for user to enter 1, 2, or 3
4. **Validate**: If input is invalid (not 1/2/3), re-prompt with error message
5. **Return**: Returns corresponding `AIMessageAcceptance` value

### Error Handling

- **Invalid Input**: Re-prompt with error message "Invalid option. Please enter 1, 2, or 3."
- **Read Error**: Return error wrapped with context
- **User Cancellation**: Return error (handled by caller to restore staging state)

### Example Usage

```go
acceptance, err := ui.PromptAIMessageAcceptanceOptions(reader)
if err != nil {
    // Handle error (restore staging state if cancelled)
    return err
}

switch acceptance {
case ui.AcceptAndCommit:
    // Commit immediately
case ui.AcceptAndEdit:
    // Pre-fill and edit
case ui.Reject:
    // Start over
}
```

---

## Type: AIMessageAcceptance

### Definition

```go
type AIMessageAcceptance int

const (
    AcceptAndCommit AIMessageAcceptance = iota
    AcceptAndEdit
    Reject
)

func (a AIMessageAcceptance) String() string
```

### Methods

- `String() string`: Returns human-readable string representation ("accept and commit", "accept and edit", "reject")

---

## Interface: PromptCommitTypeWithPreselection

### Signature

```go
func PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader (for compatibility, not used in bubbletea)
- `preselectedType` (string): Type to pre-select in the list (empty if no pre-selection)

### Returns

- `string`: Selected commit type
- `error`: Error if selection fails or user cancels

### Behavior

1. **Create Model**: Creates `SelectListModel` with pre-selection if `preselectedType` matches an available option
2. **Run Bubbletea**: Runs interactive selection UI
3. **Return Selection**: Returns selected type or error if cancelled

### Pre-selection Logic

- If `preselectedType` matches an item in the list: Set `SelectedIndex` to that item's index
- If `preselectedType` doesn't match: Set `SelectedIndex` to 0 (first item, no pre-selection)
- If `preselectedType` is empty: Set `SelectedIndex` to 0 (first item)

### Error Handling

- **Cancellation**: Returns error with message "commit type selection cancelled"
- **Invalid Model**: Returns error if model type assertion fails

---

## Interface: PromptScopeWithDefault

### Signature

```go
func PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader for user responses
- `defaultValue` (string): Default value to display (may be empty)

### Returns

- `string`: User-entered scope (or empty if user presses Enter with empty default)
- `error`: Error if input reading fails

### Behavior

1. **Display Prompt**: Shows prompt with default value in brackets: `? Scope (default: <value>):` or `? Scope:` if empty
2. **Read Input**: Waits for user input
3. **Return**: Returns user input (or empty if Enter pressed with empty default)

### Example

```go
scope, err := ui.PromptScopeWithDefault(reader, "auth")
// Prompt: "? Scope (default: auth): "
// User presses Enter → returns "auth"
// User types "api" → returns "api"
```

---

## Interface: PromptSubjectWithDefault

### Signature

```go
func PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader for user responses
- `defaultValue` (string): Default value to display (must not be empty for meaningful commits)

### Returns

- `string`: User-entered subject
- `error`: Error if input reading fails

### Behavior

1. **Display Prompt**: Shows prompt with default value: `? Subject (default: <value>):`
2. **Read Input**: Waits for user input
3. **Return**: Returns user input (or default if Enter pressed)

---

## Interface: PromptBodyWithDefault

### Signature

```go
func PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader (for compatibility, not used in bubbletea)
- `defaultValue` (string): Default value to pre-populate (may be empty)

### Returns

- `string`: User-entered body
- `error`: Error if input fails or user cancels

### Behavior

1. **Create Model**: Creates `MultilineInputModel` with `defaultValue` pre-populated
2. **Run Bubbletea**: Runs interactive multiline editor
3. **Return**: Returns edited body or error if cancelled

---

## Interface: PromptFooterWithDefault

### Signature

```go
func PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader (for compatibility, not used in bubbletea)
- `defaultValue` (string): Default value to pre-populate (may be empty)

### Returns

- `string`: User-entered footer
- `error`: Error if input fails or user cancels

### Behavior

1. **Create Model**: Creates `MultilineInputModel` with `defaultValue` pre-populated
2. **Run Bubbletea**: Runs interactive multiline editor
3. **Return**: Returns edited footer or error if cancelled

---

## Interface: PromptRejectChoice

### Signature

```go
func PromptRejectChoice(reader *bufio.Reader) (bool, error)
```

### Parameters

- `reader` (*bufio.Reader): Input reader for user responses

### Returns

- `bool`: `true` if user chooses "generate new AI message", `false` if "manual input"
- `error`: Error if input reading fails

### Behavior

1. **Display Options**: Shows two options:
   - `1`: Generate new AI message
   - `2`: Manual input
2. **Read Input**: Waits for user to enter 1 or 2
3. **Return**: Returns `true` for option 1, `false` for option 2

### Error Handling

- **Invalid Input**: Re-prompt with error message "Invalid option. Please enter 1 or 2."

---

## Testing Requirements

### Unit Tests

- `TestPromptAIMessageAcceptanceOptions_ValidInput`: Test all three valid options (1/2/3)
- `TestPromptAIMessageAcceptanceOptions_InvalidInput`: Test invalid input re-prompting
- `TestPromptAIMessageAcceptanceOptions_Cancellation`: Test Ctrl+C handling
- `TestPromptCommitTypeWithPreselection_MatchingType`: Test pre-selection with matching type
- `TestPromptCommitTypeWithPreselection_NonMatchingType`: Test pre-selection with non-matching type
- `TestPromptScopeWithDefault_WithDefault`: Test scope prompt with default value
- `TestPromptScopeWithDefault_EmptyDefault`: Test scope prompt with empty default
- `TestPromptSubjectWithDefault`: Test subject prompt with default
- `TestPromptBodyWithDefault`: Test body prompt with default (multiline)
- `TestPromptFooterWithDefault`: Test footer prompt with default (multiline)
- `TestPromptRejectChoice`: Test reject choice prompt

### Integration Tests

- `TestAIAcceptanceWorkflow_AcceptAndCommit`: Full workflow for accept and commit
- `TestAIAcceptanceWorkflow_AcceptAndEdit`: Full workflow for accept and edit
- `TestAIAcceptanceWorkflow_Reject`: Full workflow for reject

---

## Backward Compatibility

- Existing `PromptAIMessageAcceptance()` function is **deprecated** but remains for backward compatibility
- Existing `PromptCommitType()` function remains unchanged (used for manual input)
- Existing `PromptScope()`, `PromptSubject()`, `PromptBody()`, `PromptFooter()` functions remain unchanged (used for manual input)
- New functions are additive (no breaking changes)

---

## Dependencies

- `bufio.Reader`: For reading user input
- `github.com/charmbracelet/bubbletea`: For interactive UI (commit type, body, footer)
- `internal/ui/select_list.go`: `SelectListModel` for commit type selection
- `internal/ui/multiline_input.go`: `MultilineInputModel` for body/footer editing
