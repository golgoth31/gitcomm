# Research: Improved Commit Message UI

**Feature**: 004-improve-commit-ui
**Date**: 2025-01-27

## Technology Decisions

### 1. Interactive Select List Implementation

**Decision**: Use `bubbletea` list component with custom styling via `lipgloss` for commit type selection.

**Rationale**:
- `bubbletea` is already planned for use in the project (from 001-git-commit-cli research)
- Provides built-in list component with arrow key navigation
- Supports letter-based navigation (typing to jump to options)
- Customizable styling with `lipgloss` for checkmarks and highlighting
- Thread-safe, follows Go idioms
- Well-documented and actively maintained

**Alternatives Considered**:
- Custom implementation with `termbox-go`: More control but significantly more code, reinventing the wheel
- `promptui`: Simpler but less flexible, doesn't support letter-based navigation well
- Plain `fmt.Scan` with manual rendering: Too low-level, requires terminal control sequences

**Implementation Pattern**:
```go
// Use bubbletea list component
import "github.com/charmbracelet/bubbletea/list"

type commitTypeItem struct {
    Type        string
    Description string
}

func (i commitTypeItem) FilterValue() string {
    return i.Type
}

// Customize with lipgloss for checkmarks
func (m *selectListModel) View() string {
    // Render with checkmark for selected item
    // Use lipgloss for styling
}
```

### 2. Multiline Input Implementation

**Decision**: Use `bubbletea` textarea component with custom completion logic (double Enter detection).

**Rationale**:
- `bubbletea` provides `textarea` component for multiline input
- Supports line breaks and cursor movement
- Can detect consecutive empty lines for completion
- Handles terminal resizing gracefully
- Integrates with existing bubbletea workflow

**Alternatives Considered**:
- Custom `bufio.Scanner` implementation: More control but requires manual line break handling and cursor management
- `readline` library: Good for single-line, multiline support is limited
- Plain `bufio.Reader`: Too low-level, requires manual state management

**Implementation Pattern**:
```go
import "github.com/charmbracelet/bubbletea/textarea"

type multilineInputModel struct {
    textarea textarea.Model
    emptyLineCount int
}

func (m *multilineInputModel) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            if m.textarea.Value() == "" {
                m.emptyLineCount++
                if m.emptyLineCount >= 2 {
                    // Complete input
                }
            } else {
                m.emptyLineCount = 0
                // Add newline
            }
        }
    }
}
```

### 3. Letter-Based Navigation

**Decision**: Use bubbletea list's built-in filtering/navigation with custom key handler.

**Rationale**:
- Bubbletea list component supports filtering by typing
- Can be configured to jump to first matching option
- Standard TUI pattern, familiar to users
- No additional dependencies needed

**Alternatives Considered**:
- Custom key handler: More control but duplicates existing functionality
- Separate search field: More complex UI, not needed for small list (8 items)

**Implementation Pattern**:
```go
// Bubbletea list supports filtering by default
// Configure to jump to first match on letter press
list.NewModel(items, itemDelegate, width, height)
list.SetFilteringEnabled(true)
// Custom key handler for letter navigation
```

### 4. Visual Indicators (Checkmarks and Highlighting)

**Decision**: Use `lipgloss` for styling with Unicode checkmarks (✓) and color highlighting.

**Rationale**:
- `lipgloss` is already a dependency (used with bubbletea)
- Provides consistent styling and color support
- Unicode checkmarks work across platforms
- Graceful degradation for terminals without color support
- No additional dependencies

**Alternatives Considered**:
- ASCII-only indicators: Less visually appealing, harder to distinguish
- External icon library: Unnecessary dependency for simple checkmarks
- Terminal-specific escape sequences: Platform-dependent, harder to maintain

**Implementation Pattern**:
```go
import "github.com/charmbracelet/lipgloss"

var (
    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("2")). // Green
        Bold(true)
    checkmark = "✓"
)

func renderItem(item commitTypeItem, isSelected bool) string {
    if isSelected {
        return checkmark + " " + selectedStyle.Render(item.Type + " " + item.Description)
    }
    return "  " + item.Type + " " + item.Description
}
```

### 5. Whitespace Handling

**Decision**: Trim whitespace-only input using `strings.TrimSpace()` before validation.

**Rationale**:
- Standard Go library function, no dependencies
- Consistent with user expectation (whitespace-only = empty)
- Prevents accidental commits with meaningless whitespace
- Simple to implement and test

**Alternatives Considered**:
- Preserve all whitespace: Could lead to commits with only whitespace, confusing
- Trim only leading/trailing: More complex, doesn't solve the core issue

**Implementation Pattern**:
```go
func (m *multilineInputModel) GetValue() string {
    value := m.textarea.Value()
    trimmed := strings.TrimSpace(value)
    if trimmed == "" {
        return "" // Treat as empty
    }
    return value // Preserve original with line breaks
}
```

### 6. Double Enter Completion Detection

**Decision**: Track consecutive empty lines in multiline input model, complete on second empty line.

**Rationale**:
- Simple state tracking (emptyLineCount)
- Clear user intent (two empty lines = done)
- Allows single blank lines within content
- No special key combinations needed

**Alternatives Considered**:
- Ctrl+D (EOF): Standard but requires special key, less discoverable
- Esc then Enter: More complex, requires two-step action
- Timeout-based: Unreliable, could complete unintentionally

**Implementation Pattern**:
```go
type multilineInputModel struct {
    textarea      textarea.Model
    emptyLineCount int
}

func (m *multilineInputModel) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            currentLine := m.getCurrentLine()
            if strings.TrimSpace(currentLine) == "" {
                m.emptyLineCount++
                if m.emptyLineCount >= 2 {
                    return tea.Quit // Complete
                }
            } else {
                m.emptyLineCount = 0
                // Add newline normally
            }
        }
    }
}
```

## Best Practices

### Bubbletea Component Design
- Keep models small and focused (one model per UI component)
- Use composition for complex UIs (combine list + textarea models)
- Handle all key messages explicitly
- Return appropriate commands for state transitions

### Multiline Input
- Track empty line count for completion detection
- Preserve line breaks in final value
- Handle terminal resizing gracefully
- Provide clear visual feedback for completion state

### Select List
- Pre-select first option by default
- Support both arrow keys and letter navigation
- Provide clear visual indicators (checkmark, highlight)
- Handle Escape key for cancellation

## Integration Points

### Existing Code
- `internal/ui/prompts.go`: Refactor `PromptCommitType` to use new select list component
- `internal/ui/prompts.go`: Refactor `PromptBody` and `PromptFooter` to use multiline input component
- `internal/service/commit_service.go`: No changes needed (uses UI layer interfaces)

### Dependencies
- `github.com/charmbracelet/bubbletea`: Already planned, needs to be added to go.mod
- `github.com/charmbracelet/lipgloss`: Already planned, needs to be added to go.mod
- Standard library: `strings`, `bufio` (already used)

## Performance Considerations

- Select list: <50ms response time for arrow key navigation (bubbletea handles efficiently)
- Multiline input: <10ms per keystroke (bubbletea textarea is optimized)
- No performance impact on commit workflow (UI is interactive, not blocking)

## Security Considerations

- No secrets in UI (user input only)
- Input validation prevents injection (already handled in validation layer)
- No new attack vectors introduced
- Terminal escape sequences handled safely by bubbletea

## Testing Strategy

- Unit tests: Select list model, multiline input model (state transitions)
- Integration tests: Full interactive workflow (type selection, multiline input)
- Manual testing: Visual verification of checkmarks, highlighting, completion
- Edge case tests: Whitespace handling, blank lines, cancellation
