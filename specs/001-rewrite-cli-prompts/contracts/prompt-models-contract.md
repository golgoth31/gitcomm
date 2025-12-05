# Contract: Prompt Models Interface

**Feature**: 001-rewrite-cli-prompts
**Date**: 2025-01-27
**Type**: Internal API Contract

## Overview

All prompt models must implement the `tea.Model` interface from bubbletea library. This contract defines the common behavior and state management patterns for all prompt types.

## Interface Requirements

### tea.Model Interface

All prompt models MUST implement:

```go
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() string
}
```

## Common State Management

### PromptState

All models MUST track state using a `state` field of type `PromptState`:

```go
type PromptState int

const (
    StatePending PromptState = iota
    StateActive
    StateCompleted
    StateCancelled
    StateError
)
```

### Visual Indicators

All models MUST render visual indicators based on state:

- `StatePending` / `StateActive`: Blue '?' (color 4)
- `StateCompleted`: Green '✓' (color 2)
- `StateCancelled`: Red '✗' or 'X' (color 1)
- `StateError`: Yellow '⚠' (color 3)

## Model-Specific Contracts

### TextInputModel

**Fields**:
- `value string`: Current input value
- `defaultValue string`: Optional default value
- `prompt string`: Prompt title text
- `state PromptState`: Current state
- `errorMessage string`: Validation error message
- `validator func(string) error`: Optional validation function
- `width int`: Terminal width

**Behavior**:
- `Init()`: Initialize with `StatePending`, transition to `StateActive`
- `Update()`: Handle keyboard input, validate on Enter, handle Escape for cancellation
- `View()`: Render prompt title with visual indicator, input field, error message if state is Error

**Validation**:
- If `validator` is provided, call on Enter key press
- If validation fails, set state to `StateError`, store error message
- User can continue editing to correct error

### YesNoChoiceModel

**Fields**:
- `message string`: Prompt question text
- `selected bool`: Current selection (true=yes, false=no)
- `state PromptState`: Current state
- `width int`: Terminal width

**Behavior**:
- `Init()`: Initialize with `StatePending`, transition to `StateActive`
- `Update()`: Handle 'y'/'n' keys or arrow keys, handle Enter for confirmation, Escape for cancellation
- `View()`: Render prompt title with visual indicator, show current selection

### SelectListModel

**Fields**:
- `items []CommitTypeItem`: List of selectable items
- `selectedIndex int`: Currently highlighted item index
- `confirmedValue string`: Selected value after confirmation
- `state PromptState`: Current state
- `width int`, `height int`: Terminal dimensions

**Behavior**:
- `Init()`: Initialize with `StatePending`, transition to `StateActive`
- `Update()`: Handle arrow keys for navigation, Enter for confirmation, Escape for cancellation
- `View()`: Render prompt title with visual indicator, show selected value after confirmation, render list items

**Display Rules**:
- During navigation: Show `"? Choose a type(<scope>):"` (no value)
- After confirmation: Show `"✓ Choose a type(<scope>): feat"` (with confirmed value)

### MultilineInputModel

**Fields**:
- `lines []string`: Completed lines
- `currentLine string`: Current line being edited
- `defaultValue string`: Optional default value
- `emptyLineCount int`: Count of consecutive empty lines
- `isComplete bool`: Whether input is complete
- `state PromptState`: Current state
- `width int`, `height int`: Terminal dimensions
- `prompt string`: Prompt title text

**Behavior**:
- `Init()`: Initialize with `StatePending`, transition to `StateActive`
- `Update()`: Handle text input, Enter for newline or completion (two empty lines), Escape for cancellation
- `View()`: During input: render prompt title with indicator and input lines. After completion: render prompt title with checkmark and full result below with wrapping

**Display Rules**:
- Result text wrapped to `width` using lipgloss
- Full text displayed (no truncation)

## Common Behavior Requirements

### State Transitions

All models MUST follow these state transitions:

```
Pending → Active → Completed ✓
              ↓
          Error → Active (user corrects)
              ↓
          Cancelled ✗
```

### Keyboard Handling

All models MUST handle:
- `Enter`: Confirm/complete (if valid)
- `Escape`: Cancel (transition to `StateCancelled`)
- `Ctrl+C`: Handled by bubbletea, transition to `StateCancelled`

### Rendering Requirements

All models MUST:
- Render inline (no alt screen)
- Show visual indicator in prompt title based on state
- Handle terminal resize via `tea.WindowSizeMsg`
- Respect terminal width/height constraints

## Error Handling

### Validation Errors

- Models with validation MUST set state to `StateError` on validation failure
- Error message MUST be stored in `errorMessage` field
- User MUST be able to continue editing to correct error
- Visual indicator MUST show yellow '⚠' when state is `StateError`

### Cancellation

- Models MUST transition to `StateCancelled` on Escape key
- Visual indicator MUST show red '✗' when state is `StateCancelled`
- Function MUST return error with "cancelled" message

## Backward Compatibility

### Function Signatures

All prompt functions MUST maintain existing signatures:

```go
func PromptScope(reader *bufio.Reader) (string, error)
func PromptSubject(reader *bufio.Reader) (string, error)
func PromptBody(reader *bufio.Reader) (string, error)
// ... etc
```

**Note**: `reader` parameter may be ignored (not used in bubbletea implementation), but signature must remain for backward compatibility.

## Testing Requirements

All models MUST have:
- Unit tests for state transitions
- Unit tests for keyboard handling
- Unit tests for visual indicator rendering
- Table-driven tests for validation scenarios
- Integration tests for full prompt workflows
