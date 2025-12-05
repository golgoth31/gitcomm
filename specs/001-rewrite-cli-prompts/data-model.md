# Data Model: Rewrite All CLI Prompts

**Feature**: 001-rewrite-cli-prompts
**Date**: 2025-01-27

## Entities

### PromptState

Represents the current state of a prompt during its lifecycle.

**Fields**:
- `Pending` (constant): Initial state, prompt not yet displayed
- `Active` (constant): Prompt is displayed and accepting input
- `Completed` (constant): User confirmed input/selection successfully
- `Cancelled` (constant): User cancelled (Escape key)
- `Error` (constant): Validation error occurred

**Type**: Enum-like constants (using `iota` in Go)

**Validation Rules**:
- State transitions: Pending → Active → (Completed | Cancelled | Error)
- Error state can transition back to Active when user corrects input
- Completed and Cancelled are terminal states

**Usage**: Used by all prompt models to track current state for rendering visual indicators.

---

### VisualIndicator

Represents the visual indicator character displayed at the start of prompt titles.

**Fields**:
- `Pending` (string): Blue '?' character
- `Completed` (string): Green '✓' character
- `Cancelled` (string): Red '✗' or 'X' character
- `Error` (string): Yellow '⚠' character

**Type**: Constants with lipgloss styling applied

**Validation Rules**:
- Must match current PromptState
- Colors: Blue (4), Green (2), Red (1), Yellow (3)

**Usage**: Rendered in prompt title based on current state.

---

### TextInputModel

Represents a single-line text input prompt using bubbletea.

**Fields**:
- `value` (string): Current input value
- `defaultValue` (string): Optional default value to pre-fill
- `prompt` (string): Prompt title text
- `state` (PromptState): Current state of the prompt
- `errorMessage` (string): Validation error message (empty if no error)
- `validator` (func(string) error): Validation function (optional)
- `width` (int): Terminal width for rendering
- `focused` (bool): Whether prompt is currently focused

**Relationships**:
- Implements `tea.Model` interface
- Uses `VisualIndicator` for state-based rendering

**Validation Rules**:
- `value` can be empty if prompt is optional
- `errorMessage` only set when `state == Error`
- `validator` is called on Enter key press
- If validation fails, state transitions to Error, user can continue editing

**State Transitions**:
- Pending → Active: When prompt is displayed
- Active → Completed: User presses Enter with valid input
- Active → Error: User presses Enter with invalid input
- Error → Active: User continues typing after error
- Active → Cancelled: User presses Escape

---

### YesNoChoiceModel

Represents a yes/no binary choice prompt using bubbletea.

**Fields**:
- `message` (string): Prompt question text
- `selected` (bool): Current selection (true=yes, false=no)
- `state` (PromptState): Current state of the prompt
- `width` (int): Terminal width for rendering
- `focused` (bool): Whether prompt is currently focused

**Relationships**:
- Implements `tea.Model` interface
- Uses `VisualIndicator` for state-based rendering

**Validation Rules**:
- `selected` must be set before completion
- No validation errors (binary choice is always valid)

**State Transitions**:
- Pending → Active: When prompt is displayed
- Active → Completed: User presses Enter
- Active → Cancelled: User presses Escape

---

### SelectListModel

Represents a select list prompt using bubbletea (updated from existing).

**Fields**:
- `items` ([]CommitTypeItem): List of selectable items
- `selectedIndex` (int): Currently highlighted item index
- `confirmedValue` (string): Selected value after confirmation (empty until confirmed)
- `state` (PromptState): Current state of the prompt
- `width` (int): Terminal width for rendering
- `height` (int): Terminal height for rendering
- `cancelled` (bool): Whether user cancelled (legacy field, use state instead)

**Relationships**:
- Implements `tea.Model` interface
- Uses `VisualIndicator` for state-based rendering
- `CommitTypeItem` has `Type` (string) and `Description` (string)

**Validation Rules**:
- `selectedIndex` must be within bounds of `items`
- `confirmedValue` only set when `state == Completed`
- Title shows selected value only after confirmation

**State Transitions**:
- Pending → Active: When prompt is displayed
- Active → Completed: User presses Enter
- Active → Cancelled: User presses Escape

**Changes from Existing**:
- Added `state` field (replaces `cancelled` boolean)
- Added `confirmedValue` field for post-confirmation display
- Removed alt screen usage

---

### MultilineInputModel

Represents a multiline text input prompt using bubbletea (updated from existing).

**Fields**:
- `lines` ([]string): Completed lines of input
- `currentLine` (string): Currently being edited line
- `defaultValue` (string): Optional default value to pre-fill
- `emptyLineCount` (int): Count of consecutive empty lines (for completion detection)
- `isComplete` (bool): Whether input is complete (two empty lines)
- `state` (PromptState): Current state of the prompt
- `width` (int): Terminal width for rendering and wrapping
- `height` (int): Terminal height for rendering
- `prompt` (string): Prompt title text
- `focused` (bool): Whether prompt is currently focused
- `cancelled` (bool): Whether user cancelled (legacy field, use state instead)

**Relationships**:
- Implements `tea.Model` interface
- Uses `VisualIndicator` for state-based rendering

**Validation Rules**:
- `isComplete` is true when `emptyLineCount >= 2`
- Result text is all `lines` joined with newlines plus `currentLine`
- Result display wraps to `width` using lipgloss

**State Transitions**:
- Pending → Active: When prompt is displayed
- Active → Completed: User presses Enter twice on empty line
- Active → Cancelled: User presses Escape

**Changes from Existing**:
- Added `state` field (replaces `cancelled` boolean)
- Added `defaultValue` field for pre-filling
- Removed alt screen usage
- Added result display with wrapping after completion

---

## State Machine

```
PromptState Transitions:

[Pending] → [Active] → [Completed] ✓
                ↓
            [Error] → [Active] (user corrects)
                ↓
            [Cancelled] ✗
```

**Notes**:
- All prompts start in Pending state
- Active is the main interaction state
- Completed and Cancelled are terminal states
- Error can transition back to Active for correction
- Visual indicators match state: Pending/Active='?', Completed='✓', Cancelled='✗', Error='⚠'

---

## Data Flow

1. **Prompt Initialization**:
   - Caller invokes prompt function (e.g., `PromptScope(reader)`)
   - Function creates appropriate model (TextInputModel, etc.)
   - Model initialized with Pending state
   - Bubbletea program created (without alt screen)

2. **User Interaction**:
   - Program runs, model transitions to Active state
   - User input updates model fields (value, selectedIndex, etc.)
   - Visual indicator shows blue '?' during Active state

3. **Validation/Completion**:
   - User confirms (Enter key)
   - If validation passes: state → Completed, indicator → green '✓'
   - If validation fails: state → Error, indicator → yellow '⚠', error message shown
   - If user cancels (Escape): state → Cancelled, indicator → red '✗'

4. **Result Return**:
   - Model extracts result value
   - Function returns result and error (if cancelled/error)
   - Caller receives result as before (backward compatible)

---

## Default Values

- **PromptState**: Pending (initial state)
- **VisualIndicator**: Blue '?' (pending state)
- **TextInputModel.value**: Empty string (or defaultValue if provided)
- **YesNoChoiceModel.selected**: false (no selected by default)
- **SelectListModel.selectedIndex**: 0 (first item) or provided preselection
- **MultilineInputModel.lines**: Empty slice (or defaultValue split into lines if provided)

---

## Error Handling

- **Validation Errors**: Stored in `errorMessage` field, state set to Error, user can correct
- **Cancellation**: State set to Cancelled, function returns error with "cancelled" message
- **Ctrl+C Interruption**: Handled by bubbletea, state set to Cancelled, function returns error
- **Terminal Size Issues**: Models handle window resize via `tea.WindowSizeMsg`, update width/height

---

## Relationships Summary

```
PromptState (enum)
    ↓
VisualIndicator (constants)
    ↓
TextInputModel ──implements──> tea.Model
YesNoChoiceModel ──implements──> tea.Model
SelectListModel ──implements──> tea.Model
MultilineInputModel ──implements──> tea.Model
```

All models use PromptState and VisualIndicator for consistent behavior and rendering.
