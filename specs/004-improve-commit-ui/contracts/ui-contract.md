# UI Interface Contract: Improved Commit Message UI

**Feature**: 004-improve-commit-ui
**Date**: 2025-01-27
**Interface**: UI prompt functions

## Extended Interface

### PromptCommitType

Prompts the user to select a commit type using an interactive select list.

```go
func PromptCommitType(reader *bufio.Reader) (string, error)
```

**Input**:
- `reader *bufio.Reader` - Input reader (may be unused if using bubbletea)

**Output**:
- `string` - Selected commit type (e.g., "feat", "fix", "docs")
- `error` - Error if selection was cancelled or failed

**Behavior**:
- Displays interactive select list with all 8 commit types
- Pre-selects first option (feat)
- Allows navigation with arrow keys (up/down)
- Supports letter-based navigation (typing a letter jumps to first matching option)
- Shows checkmark (✓) and highlight for selected option
- Returns selected type when user presses Enter
- Returns error if user presses Escape (cancellation)

**Error Cases**:
- `ErrSelectionCancelled` - User pressed Escape
- `ErrInvalidSelection` - Internal error (selected index out of range)

---

### PromptBody

Prompts the user for commit body using multiline input.

```go
func PromptBody(reader *bufio.Reader) (string, error)
```

**Input**:
- `reader *bufio.Reader` - Input reader (may be unused if using bubbletea)

**Output**:
- `string` - Body text (may contain line breaks, empty if skipped)
- `error` - Error if input was cancelled or failed

**Behavior**:
- Displays multiline input field
- Allows multiple lines of text with line breaks
- Single Enter creates new line (does not complete)
- Two consecutive empty lines (double Enter) complete the input
- Single blank lines within content are preserved
- Whitespace-only input is treated as empty (trimmed)
- Returns empty string if user skips (double Enter immediately)

**Error Cases**:
- `ErrInputCancelled` - User pressed Escape or Ctrl+C
- No error for empty input (treated as skipped)

---

### PromptFooter

Prompts the user for commit footer using multiline input.

```go
func PromptFooter(reader *bufio.Reader) (string, error)
```

**Input**:
- `reader *bufio.Reader` - Input reader (may be unused if using bubbletea)

**Output**:
- `string` - Footer text (may contain line breaks, empty if skipped)
- `error` - Error if input was cancelled or failed

**Behavior**:
- Same as `PromptBody` (multiline input with double Enter completion)
- Allows multiple footer entries (e.g., "Fixes #123\nCloses #456")
- Single blank lines preserved, two consecutive empty lines complete input
- Whitespace-only input treated as empty

**Error Cases**:
- `ErrInputCancelled` - User pressed Escape or Ctrl+C
- No error for empty input (treated as skipped)

---

## Implementation Contract

### Thread Safety

- UI components are single-threaded (bubbletea models run in main goroutine)
- No additional synchronization needed
- Safe to call from service layer (which is also single-threaded for prompts)

### Error Handling

- All prompt functions return errors for cancellation
- Empty input (skipped fields) returns empty string, no error
- Errors are wrapped with context for traceability

### Performance

- Select list navigation: <50ms response time
- Multiline input keystroke handling: <10ms per keystroke
- No blocking operations (all UI is interactive)

### Resource Management

- UI components use terminal I/O (handled by bubbletea)
- No file handles or network connections
- Cleanup handled by bubbletea lifecycle

---

## Testing Contract

### Unit Tests

- Test SelectListModel state transitions (arrow keys, letter navigation)
- Test MultilineInputModel completion detection (double Enter)
- Test whitespace handling (trimming)
- Test blank line preservation vs completion

### Integration Tests

- Test full commit type selection workflow
- Test multiline body input with various scenarios
- Test multiline footer input
- Test cancellation (Escape key)
- Test terminal resizing

### Mock Requirements

- Mock terminal for testing (bubbletea provides test utilities)
- Mock user input for automated testing
