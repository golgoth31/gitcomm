# Research: Rewrite All CLI Prompts

**Feature**: 001-rewrite-cli-prompts
**Date**: 2025-01-27
**Purpose**: Resolve technical unknowns and establish implementation patterns

## Research Questions

### Q1: Bubbletea Text Input Model Pattern

**Question**: How to implement a bubbletea text input model that supports validation, default values, and visual indicators without alt screen?

**Research**:
- Bubbletea supports inline rendering by omitting `tea.WithAltScreen()` option
- Text input can be implemented using a simple model with string state
- Validation can be done in `Update()` method before allowing completion
- Visual indicators (blue '?', green '✓', red 'X', yellow '⚠') can be rendered using lipgloss styling

**Decision**: Create a `TextInputModel` struct implementing `tea.Model` interface with:
- `value` string field for current input
- `defaultValue` string for pre-filling
- `validator` function for validation
- `state` enum (pending, active, completed, cancelled, error)
- `errorMessage` string for validation errors

**Rationale**: Standard bubbletea pattern, supports all required features, no alt screen needed.

**Alternatives Considered**:
- Using existing bubbletea textinput component: Rejected - too complex, need custom validation and visual indicators
- Keeping bufio.Reader: Rejected - doesn't support visual indicators and inline rendering

---

### Q2: Bubbletea Yes/No Choice Model Pattern

**Question**: How to implement a yes/no choice prompt using bubbletea without alt screen?

**Research**:
- Yes/no choices can be implemented as a simple select list with 2 items
- Or as a dedicated model with boolean state
- Visual indicators can be shown in the prompt title

**Decision**: Create a `YesNoChoiceModel` struct implementing `tea.Model` interface with:
- `message` string for the prompt text
- `selected` bool for current selection (true=yes, false=no)
- `state` enum (pending, active, completed, cancelled)
- Simple keyboard navigation (y/n keys or arrow keys)

**Rationale**: Dedicated model is simpler than reusing select list, provides better UX for binary choices.

**Alternatives Considered**:
- Reusing SelectListModel: Rejected - overkill for binary choice, less intuitive
- Using text input with validation: Rejected - less user-friendly than dedicated choice model

---

### Q3: Visual Indicator Rendering Without Alt Screen

**Question**: How to render visual indicators (blue '?', green '✓', red 'X', yellow '⚠') in inline prompts?

**Research**:
- lipgloss supports color styling: `lipgloss.NewStyle().Foreground(lipgloss.Color("4"))` for blue
- Colors: blue=4, green=2, red=1, yellow=3
- Indicators can be part of the prompt title string
- State transitions happen in `Update()` method when user confirms/cancels

**Decision**: Use lipgloss styling for indicators:
- Blue '?' (color 4): `lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render("?")`
- Green '✓' (color 2): `lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("✓")`
- Red 'X' (color 1): `lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("✗")`
- Yellow '⚠' (color 3): `lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("⚠")`

**Rationale**: Standard lipgloss approach, works inline, provides clear visual feedback.

**Alternatives Considered**:
- Using ANSI escape codes directly: Rejected - less maintainable, lipgloss is already a dependency
- Using emoji without colors: Rejected - less visible, colors provide better contrast

---

### Q4: Select List Value Display After Title

**Question**: How to display selected value after prompt title in select list without alt screen?

**Research**:
- Current `SelectListModel.View()` renders the list
- Title can be modified to include selected value after confirmation
- During navigation, title shows blue '?' without value
- After confirmation, title shows green '✓' with selected value

**Decision**: Update `SelectListModel.View()` to:
- Show `"? Choose a type(<scope>):"` during navigation
- Show `"✓ Choose a type(<scope>): feat"` after confirmation (with selected value)
- Store `confirmedValue` in model state
- Update title rendering based on state

**Rationale**: Simple state-based rendering, maintains inline display, clear visual feedback.

**Alternatives Considered**:
- Showing value during navigation: Rejected - spec requires value only after confirmation
- Separate confirmation step: Rejected - adds unnecessary complexity

---

### Q5: Multiline Input Result Display with Wrapping

**Question**: How to display multiline input result under prompt title with line wrapping?

**Research**:
- lipgloss supports text wrapping: `lipgloss.NewStyle().Width(terminalWidth).Wrap()`
- Result can be displayed after prompt completes
- Need to respect terminal width for wrapping
- Full text should be displayed (no truncation per spec)

**Decision**: Update `MultilineInputModel.View()` to:
- Show prompt title with indicator during input
- After completion, show title with checkmark and full result below
- Use `lipgloss.NewStyle().Width(m.Width).Wrap(result)` for wrapping
- Store `isComplete` state to switch between input and result display

**Rationale**: lipgloss wrapping handles terminal width automatically, preserves full text, inline display.

**Alternatives Considered**:
- Truncating long text: Rejected - spec requires full text display
- Using external text wrapping library: Rejected - lipgloss already provides wrapping

---

### Q6: Removing Alt Screen from Existing Models

**Question**: How to remove alt screen from existing select list and multiline input models?

**Research**:
- Alt screen is enabled via `tea.WithAltScreen()` option
- Removing it allows inline rendering
- Terminal history remains visible
- Models need no changes, only program initialization

**Decision**: Remove `tea.WithAltScreen()` from:
- `PromptCommitType()` and `PromptCommitTypeWithPreselection()`
- `PromptBody()`, `PromptBodyWithDefault()`, `PromptFooter()`, `PromptFooterWithDefault()`

Change from:
```go
p := tea.NewProgram(model, tea.WithAltScreen())
```

To:
```go
p := tea.NewProgram(model)
```

**Rationale**: Simple change, no model modifications needed, achieves inline rendering requirement.

**Alternatives Considered**:
- Conditional alt screen: Rejected - spec requires no alt screen for any prompt
- Keeping alt screen for some prompts: Rejected - violates spec requirement

---

### Q7: Validation Error Display in Text Input

**Question**: How to show validation errors with yellow/orange warning indicator in text input prompts?

**Research**:
- Validation can be checked in `Update()` method when user confirms
- If validation fails, set state to `error` and store error message
- `View()` method can show warning indicator and error message below input
- User can continue editing after error

**Decision**: Implement validation in `TextInputModel`:
- Add `state` field with `error` state
- Add `errorMessage` string field
- In `Update()`, validate on Enter key press
- If invalid, set state to `error`, store message, don't quit
- In `View()`, show yellow '⚠' indicator and error message when state is `error`
- User can continue typing to fix error

**Rationale**: Clear error feedback, allows correction, maintains inline display.

**Alternatives Considered**:
- Blocking input on error: Rejected - less user-friendly, prevents correction
- Separate error prompt: Rejected - adds complexity, breaks inline flow

---

### Q8: Backward Compatibility with Function Signatures

**Question**: How to maintain backward compatibility while changing internal implementation?

**Research**:
- Current functions: `func PromptXxx(reader *bufio.Reader) (result, error)`
- New implementation uses bubbletea internally
- Function signatures can remain the same
- `reader` parameter can be ignored (not used in bubbletea)

**Decision**: Keep all function signatures unchanged:
- `func PromptScope(reader *bufio.Reader) (string, error)`
- `func PromptSubject(reader *bufio.Reader) (string, error)`
- `func PromptBody(reader *bufio.Reader) (string, error)`
- etc.

Internal implementation changes from bufio.Reader to bubbletea, but callers see no difference.

**Rationale**: Maintains backward compatibility, no breaking changes, gradual migration possible.

**Alternatives Considered**:
- Changing signatures to remove `reader`: Rejected - breaking change, violates spec
- Adding new functions: Rejected - spec requires maintaining existing signatures

---

## Summary

All technical unknowns resolved. Implementation approach:
1. Create new bubbletea models for text input and yes/no choices
2. Update existing select list and multiline input models (remove alt screen, add visual indicators)
3. Refactor all prompt functions to use bubbletea internally
4. Maintain backward-compatible function signatures
5. Use lipgloss for visual indicators and text wrapping
6. Implement state-based rendering for all prompt types
