# Contract: Prompt Functions Interface

**Feature**: 001-rewrite-cli-prompts
**Date**: 2025-01-27
**Type**: Public API Contract

## Overview

This contract defines the public interface for all prompt functions. All functions MUST maintain backward-compatible signatures while changing internal implementation to use bubbletea.

## Function Signatures

### Text Input Prompts

```go
// PromptScope prompts for commit scope (optional)
func PromptScope(reader *bufio.Reader) (string, error)

// PromptScopeWithDefault prompts for commit scope with default value
func PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error)

// PromptSubject prompts for commit subject (required, with validation)
func PromptSubject(reader *bufio.Reader) (string, error)

// PromptSubjectWithDefault prompts for commit subject with default value
func PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

**Behavior**:
- Display prompt with blue '?' indicator
- Accept text input
- Validate (for subject: non-empty, length warning if >72 chars)
- Show yellow '⚠' indicator with error message on validation failure
- Show green '✓' indicator on successful completion
- Return entered value or error

### Multiline Input Prompts

```go
// PromptBody prompts for commit body (optional, multiline)
func PromptBody(reader *bufio.Reader) (string, error)

// PromptBodyWithDefault prompts for commit body with default value pre-populated
func PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error)

// PromptFooter prompts for commit footer (optional, multiline)
func PromptFooter(reader *bufio.Reader) (string, error)

// PromptFooterWithDefault prompts for commit footer with default value pre-populated
func PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error)
```

**Behavior**:
- Display prompt with blue '?' indicator
- Accept multiline input (Enter for newline, double Enter on empty line to complete)
- Show green '✓' indicator on completion
- Display full entered text under prompt title with line wrapping
- Return entered text or error

### Select List Prompts

```go
// PromptCommitType prompts for commit type using interactive select list
func PromptCommitType(reader *bufio.Reader) (string, error)

// PromptCommitTypeWithPreselection prompts for commit type with pre-selected type
func PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error)
```

**Behavior**:
- Display prompt with blue '?' indicator
- Show interactive list of commit types
- Allow navigation with arrow keys
- Show selected value after prompt title only after confirmation
- Show green '✓' indicator with selected value on completion
- Return selected type or error

### Yes/No Choice Prompts

```go
// PromptEmptyCommit prompts to confirm creating an empty commit
func PromptEmptyCommit(reader *bufio.Reader) (bool, error)

// PromptConfirm prompts to confirm an action
func PromptConfirm(reader *bufio.Reader, message string) (bool, error)

// PromptAIUsage prompts to choose whether to use AI
func PromptAIUsage(reader *bufio.Reader, tokenCount int) (bool, error)

// PromptAIMessageEdit prompts to choose whether to edit AI message with validation errors
func PromptAIMessageEdit(reader *bufio.Reader, errors []string) (bool, error)

// PromptRejectChoice prompts to choose between new AI generation or manual input
func PromptRejectChoice(reader *bufio.Reader) (bool, error)
```

**Behavior**:
- Display prompt with blue '?' indicator
- Show yes/no choice (y/n keys or arrow keys)
- Show green '✓' indicator on completion
- Return true (yes) or false (no), or error

### Multi-Choice Prompts

```go
// PromptAIMessageAcceptanceOptions prompts for AI message acceptance with three options
func PromptAIMessageAcceptanceOptions(reader *bufio.Reader, message string) (AIMessageAcceptance, error)

// PromptCommitFailureChoice prompts for action when commit fails
func PromptCommitFailureChoice(reader *bufio.Reader) (CommitFailureChoice, error)
```

**Behavior**:
- Display prompt with blue '?' indicator
- Show multiple choice options (can use select list model)
- Show green '✓' indicator on completion
- Return selected choice enum value or error

## Error Handling

### Cancellation

If user cancels (Escape key), function MUST:
- Return error with message containing "cancelled"
- Show red '✗' indicator before returning

### Validation Errors

If validation fails, function MUST:
- Show yellow '⚠' indicator with error message
- Allow user to correct input
- Not return error until user cancels or provides valid input

### Ctrl+C Interruption

If user presses Ctrl+C, function MUST:
- Handle gracefully via bubbletea
- Return error with message containing "cancelled" or "interrupted"

## Visual Indicator Requirements

All functions MUST:
- Start with blue '?' indicator when prompt is displayed
- Show green '✓' indicator when prompt completes successfully
- Show red '✗' indicator when prompt is cancelled
- Show yellow '⚠' indicator when validation fails (for prompts with validation)

## Backward Compatibility Guarantees

### Function Signatures

- All function signatures MUST remain unchanged
- `reader *bufio.Reader` parameter MUST be accepted (may be ignored internally)
- Return types MUST remain unchanged

### Behavior

- All existing validation logic MUST be preserved
- Default values MUST work as before
- Pre-selection MUST work as before
- Error messages MUST be compatible (may be enhanced)

### Callers

- Existing callers MUST NOT require changes
- Function calls MUST work identically from caller perspective
- Only internal implementation changes (bufio.Reader → bubbletea)

## Performance Requirements

- Prompt rendering MUST be responsive (<100ms for state updates)
- No noticeable lag when typing or navigating
- Terminal resize MUST be handled smoothly

## Testing Requirements

All functions MUST have:
- Unit tests for all prompt types
- Tests for validation scenarios
- Tests for cancellation
- Tests for default values
- Tests for pre-selection
- Integration tests for full workflows
