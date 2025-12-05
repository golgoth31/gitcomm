# Data Model: Rewrite CLI Prompts with Huh Library

**Feature**: 013-huh-prompts
**Date**: 2025-01-27

## Overview

This feature migrates prompt implementations from custom Bubble Tea models to the `huh` library. The data model is minimal - primarily function signatures and prompt configuration. No persistent data structures are introduced.

## Entities

### Prompt Function Signatures

These are the existing function signatures that must be maintained for backward compatibility:

#### Text Input Prompts

```go
// PromptScope prompts for commit scope (optional)
func PromptScope(reader *bufio.Reader) (string, error)

// PromptScopeWithDefault prompts for commit scope with default value
func PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error)

// PromptSubject prompts for commit subject (required)
func PromptSubject(reader *bufio.Reader) (string, error)

// PromptSubjectWithDefault prompts for commit subject with default value
func PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

#### Multiline Input Prompts

```go
// PromptBody prompts for commit body (optional, multiline)
func PromptBody(reader *bufio.Reader) (string, error)

// PromptBodyWithDefault prompts for commit body with default value
func PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error)

// PromptFooter prompts for commit footer (optional, multiline)
func PromptFooter(reader *bufio.Reader) (string, error)

// PromptFooterWithDefault prompts for commit footer with default value
func PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

#### Selection Prompts

```go
// PromptCommitType prompts for commit type using select list
func PromptCommitType(reader *bufio.Reader) (string, error)

// PromptCommitTypeWithPreselection prompts for commit type with pre-selected value
func PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error)
```

#### Confirmation Prompts

```go
// PromptEmptyCommit prompts to confirm creating an empty commit
func PromptEmptyCommit(reader *bufio.Reader) (bool, error)

// PromptConfirm prompts for general confirmation
func PromptConfirm(reader *bufio.Reader, message string) (bool, error)

// PromptAIUsage prompts to choose whether to use AI
func PromptAIUsage(reader *bufio.Reader, tokenCount int) (bool, error)

// PromptAIMessageEdit prompts to edit or use AI message with validation errors
func PromptAIMessageEdit(reader *bufio.Reader, errors []string) (bool, error)

// PromptRejectChoice prompts to choose between new AI message or manual input
func PromptRejectChoice(reader *bufio.Reader) (bool, error)
```

#### Multi-Choice Prompts

```go
// AIMessageAcceptance represents user's choice for AI message
type AIMessageAcceptance int

const (
    AcceptAndCommit AIMessageAcceptance = iota
    AcceptAndEdit
    Reject
)

// PromptAIMessageAcceptanceOptions prompts for AI message acceptance with options
func PromptAIMessageAcceptanceOptions(reader *bufio.Reader, message string) (AIMessageAcceptance, error)

// CommitFailureChoice represents user's choice when commit fails
type CommitFailureChoice int

const (
    RetryCommit CommitFailureChoice = iota
    EditMessage
    CancelCommit
)

// PromptCommitFailureChoice prompts for action when commit fails
func PromptCommitFailureChoice(reader *bufio.Reader) (CommitFailureChoice, error)
```

## Internal Data Structures

### Huh Form Configuration

Internal structures for building `huh` forms (not exported):

```go
// Form configuration for commit message fields (combined form)
type commitMessageFormConfig struct {
    Type    string
    Scope   string
    Subject string
    Body    string
    Footer  string
}

// Field validation functions (matching existing validators)
type fieldValidator func(string) error
```

## Validation Rules

### Subject Field
- **Required**: Yes
- **Max Length**: 72 characters (recommended, but not enforced)
- **Validation**: Cannot be empty or whitespace-only

### Scope Field
- **Required**: No
- **Max Length**: None (reasonable limit assumed)
- **Validation**: Optional, no validation if empty

### Body Field
- **Required**: No
- **Max Length**: 320 characters (recommended, warning shown if exceeded)
- **Validation**: Optional, multiline input

### Footer Field
- **Required**: No
- **Max Length**: None
- **Validation**: Optional, multiline input

### Commit Type Field
- **Required**: Yes
- **Options**: Predefined list (feat, fix, docs, style, refactor, perf, test, chore, etc.)
- **Validation**: Must select from valid options

## State Transitions

### Prompt Lifecycle

1. **Pending**: Prompt function called, form not yet initialized
2. **Active**: Form displayed, user interacting
3. **Validating**: User input being validated
4. **Completed**: User validated input, summary line displayed
5. **Cancelled**: User cancelled (Escape/Ctrl+C), no summary line

### Form Field States (huh internal)

- `huh.StateIdle`: Field not yet interacted with
- `huh.StateFocused`: Field currently focused
- `huh.StateCompleted`: Field completed and validated
- `huh.StateError`: Field has validation error

## Relationships

- **Prompt Functions → Huh Forms**: Each prompt function creates and runs a `huh.Form` internally
- **Combined Forms → Individual Functions**: Related prompts (commit message fields) can share a form builder, but individual functions extract their specific field value
- **Validation → Field Validators**: Each field can have a validator function that returns an error if validation fails

## Data Flow

1. **Caller** calls prompt function (e.g., `PromptScope(reader)`)
2. **Prompt Function** creates `huh.Form` with appropriate field type
3. **Huh Library** renders form inline and handles user interaction
4. **User** enters/selects value and validates
5. **Huh Library** validates input (if validator provided)
6. **Prompt Function** extracts value from form field
7. **Prompt Function** prints post-validation summary line: `✓ <title>: <value>`
8. **Prompt Function** returns value and error (if any)

## Constraints

- All function signatures must remain unchanged (backward compatibility)
- `bufio.Reader` parameter may be unused (kept for compatibility)
- Return types must match existing signatures exactly
- Error handling must preserve existing error types and messages
- Post-validation display format must be consistent: `✓ <title>: <value>`

## Migration Notes

### Removed Entities

The following custom Bubble Tea models will be removed:
- `TextInputModel` → Replaced by `huh.NewInput()`
- `MultilineInputModel` → Replaced by `huh.NewText()`
- `YesNoChoiceModel` → Replaced by `huh.NewConfirm()`
- `SelectListModel` → Replaced by `huh.NewSelect()`

### Kept Entities

- `PromptState` enum (may be kept for compatibility or removed if not needed)
- `AIMessageAcceptance` enum (kept - used in return types)
- `CommitFailureChoice` enum (kept - used in return types)
- `DisplayCommitMessage` function (kept - used for displaying formatted messages)
