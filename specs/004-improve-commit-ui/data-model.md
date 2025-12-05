# Data Model: Improved Commit Message UI

**Feature**: 004-improve-commit-ui
**Date**: 2025-01-27

## Domain Entities

### CommitTypeItem

Represents a single option in the commit type select list.

**Fields**:
- `Type string` - The commit type identifier (e.g., "feat", "fix", "docs")
- `Description string` - Human-readable description (e.g., "[new feature]", "[bug fix]")

**Methods**:
- `FilterValue() string` - Returns the value used for filtering/letter navigation (typically the Type)

**Validation Rules**:
- `Type` must be one of: "feat", "fix", "docs", "style", "refactor", "test", "chore", "version"
- `Description` must not be empty

**Lifecycle**:
1. Created at UI initialization with predefined list
2. Used in select list model for display and selection
3. Discarded after selection is made

---

### SelectListModel

Represents the state of the interactive commit type select list.

**Fields**:
- `Items []CommitTypeItem` - List of available commit types
- `SelectedIndex int` - Index of currently selected item (0-based)
- `Filter string` - Current filter value for letter-based navigation
- `Width int` - Terminal width for rendering
- `Height int` - Terminal height for rendering

**Methods**:
- `GetSelectedType() string` - Returns the type of the currently selected item
- `IsSelected(index int) bool` - Returns true if the given index is selected
- `MoveUp()` - Moves selection up (decrements SelectedIndex)
- `MoveDown()` - Moves selection down (increments SelectedIndex)
- `JumpToLetter(letter string)` - Jumps to first item starting with letter

**State Transitions**:
- Initial: SelectedIndex = 0 (first item pre-selected)
- Arrow Up: SelectedIndex = max(0, SelectedIndex - 1)
- Arrow Down: SelectedIndex = min(len(Items) - 1, SelectedIndex + 1)
- Letter typed: Filter updated, jump to first matching item
- Enter pressed: Selection confirmed, return selected type

**Validation Rules**:
- `SelectedIndex` must be in range [0, len(Items) - 1]
- `Filter` must be single character for letter navigation

---

### MultilineInputModel

Represents the state of a multiline input field (body or footer).

**Fields**:
- `Value string` - Current input value (may contain line breaks)
- `EmptyLineCount int` - Count of consecutive empty lines (for completion detection)
- `IsComplete bool` - Whether input is complete (two consecutive empty lines detected)
- `Width int` - Terminal width for rendering
- `Height int` - Terminal height for rendering

**Methods**:
- `GetValue() string` - Returns the input value (trimmed if whitespace-only)
- `IsWhitespaceOnly() bool` - Returns true if value contains only whitespace
- `AddLine(line string)` - Adds a line to the value
- `GetCurrentLine() string` - Returns the current line being edited
- `HandleEnter()` - Handles Enter key press (adds newline or completes if two empty lines)

**State Transitions**:
- Initial: Value = "", EmptyLineCount = 0, IsComplete = false
- Enter on non-empty line: Value += "\n", EmptyLineCount = 0
- Enter on empty line: EmptyLineCount++, if EmptyLineCount >= 2 then IsComplete = true
- Character typed: Value updated, EmptyLineCount = 0 (if was > 0)

**Validation Rules**:
- `EmptyLineCount` must be in range [0, 2]
- `Value` must preserve line breaks
- Whitespace-only values are treated as empty

---

## Relationships

- `SelectListModel` contains multiple `CommitTypeItem` instances
- `MultilineInputModel` is used independently for body and footer input
- Both models are part of the UI layer and interact with the service layer via interfaces

---

## State Transitions

### Commit Type Selection Workflow

```
[UI Initialization]
  ↓
[Create SelectListModel with 8 CommitTypeItems]
  ↓
[Pre-select first item (feat) - SelectedIndex = 0]
  ↓
[Display select list]
  ↓
[User navigates with arrow keys or types letter]
  ↓
[SelectedIndex updated, visual indicator updated]
  ↓
[User presses Enter]
  ↓
[Return selected CommitTypeItem.Type]
  ↓
[Selection complete]
```

### Multiline Input Workflow

```
[UI Initialization]
  ↓
[Create MultilineInputModel]
  ↓
[Display input field]
  ↓
[User types content]
  ↓
[Value updated, EmptyLineCount = 0]
  ↓
[User presses Enter on non-empty line]
  ↓
[Newline added, Value += "\n"]
  ↓
[User presses Enter on empty line]
  ↓
[EmptyLineCount = 1]
  ↓
[User presses Enter again on empty line]
  ↓
[EmptyLineCount = 2, IsComplete = true]
  ↓
[Return Value (trimmed if whitespace-only)]
  ↓
[Input complete]
```

---

## Data Flow

1. **Commit Type Selection**:
   - User interaction → `SelectListModel` state update → Selected type returned → Service layer

2. **Multiline Input**:
   - User keystrokes → `MultilineInputModel` state update → Value accumulated → Completion detected → Value returned → Service layer

---

## Persistence

- **No Persistent Storage**: All UI state is in-memory only during the prompt session
- **No Configuration File**: UI behavior is hardcoded (no user preferences needed)
- **No Database**: No persistent storage required (stateless UI)

---

## Error Types

- `ErrSelectionCancelled` - User pressed Escape during selection
- `ErrInputCancelled` - User cancelled multiline input
- `ErrInvalidSelection` - Selected index out of range (should not occur in normal flow)

---

## Integration with Existing Models

### CommitMessage (existing)
- No changes needed (UI layer provides values, model remains unchanged)
- `Type` field receives value from `SelectListModel.GetSelectedType()`
- `Body` field receives value from `MultilineInputModel.GetValue()`
- `Footer` field receives value from `MultilineInputModel.GetValue()`

### UI Prompts (existing)
- `PromptCommitType` will use `SelectListModel` instead of text input
- `PromptBody` will use `MultilineInputModel` instead of simple scanner
- `PromptFooter` will use `MultilineInputModel` instead of single-line input
