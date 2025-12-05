# Quickstart: Rewrite All CLI Prompts

**Feature**: 001-rewrite-cli-prompts
**Date**: 2025-01-27

## Overview

This feature rewrites all CLI prompts to use bubbletea consistently with unified visual design. All prompts display blue '?' initially, green '✓' on completion, red '✗' on cancellation, and yellow '⚠' on validation errors. Prompts render inline without alt screen mode.

## Key Concepts

### Visual Indicators

- **Blue '?'**: Prompt is active/pending
- **Green '✓'**: Prompt completed successfully
- **Red '✗'**: Prompt cancelled
- **Yellow '⚠'**: Validation error (user can correct)

### State Management

All prompts use `PromptState` enum:
- `StatePending`: Initial state
- `StateActive`: User interacting
- `StateCompleted`: Successfully completed
- `StateCancelled`: User cancelled
- `StateError`: Validation error

### Prompt Types

1. **Text Input**: Single-line text with optional validation
2. **Multiline Input**: Multi-line text with completion on double Enter
3. **Select List**: Interactive list selection
4. **Yes/No Choice**: Binary choice (y/n)

## Implementation Steps

### Step 1: Create Text Input Model

Create `internal/ui/text_input.go`:

```go
type TextInputModel struct {
    value        string
    defaultValue string
    prompt       string
    state        PromptState
    errorMessage string
    validator    func(string) error
    width        int
    focused      bool
}

func NewTextInputModel(prompt, defaultValue string, validator func(string) error) TextInputModel {
    // Initialize with StatePending
}

func (m TextInputModel) Init() tea.Cmd {
    // Transition to StateActive
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle keyboard input, validation, cancellation
}

func (m TextInputModel) View() string {
    // Render prompt title with visual indicator, input field, error message
}
```

### Step 2: Create Yes/No Choice Model

Create `internal/ui/yes_no_choice.go`:

```go
type YesNoChoiceModel struct {
    message string
    selected bool
    state   PromptState
    width   int
    focused bool
}

func NewYesNoChoiceModel(message string) YesNoChoiceModel {
    // Initialize with StatePending
}

func (m YesNoChoiceModel) Init() tea.Cmd {
    // Transition to StateActive
}

func (m YesNoChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle 'y'/'n' keys, Enter, Escape
}

func (m YesNoChoiceModel) View() string {
    // Render prompt with visual indicator and current selection
}
```

### Step 3: Update Select List Model

Update `internal/ui/select_list.go`:

- Add `state PromptState` field
- Add `confirmedValue string` field
- Remove `tea.WithAltScreen()` from program initialization
- Update `View()` to show selected value after title only after confirmation
- Update visual indicator rendering based on state

### Step 4: Update Multiline Input Model

Update `internal/ui/multiline_input.go`:

- Add `state PromptState` field
- Add `defaultValue string` field
- Remove `tea.WithAltScreen()` from program initialization
- Update `View()` to show result under title with wrapping after completion
- Update visual indicator rendering based on state

### Step 5: Refactor Prompt Functions

Update `internal/ui/prompts.go`:

For each prompt function:
1. Create appropriate bubbletea model
2. Initialize program without alt screen: `tea.NewProgram(model)`
3. Run program: `finalModel, err := p.Run()`
4. Extract result from model
5. Check state for cancellation/error
6. Return result or error

Example for `PromptScope`:

```go
func PromptScope(reader *bufio.Reader) (string, error) {
    model := NewTextInputModel("Scope (optional, press Enter to skip)", "", nil)
    p := tea.NewProgram(model) // No alt screen
    finalModel, err := p.Run()
    if err != nil {
        return "", fmt.Errorf("failed to run scope input: %w", err)
    }
    textModel, ok := finalModel.(TextInputModel)
    if !ok {
        return "", fmt.Errorf("unexpected model type: %T", finalModel)
    }
    if textModel.state == StateCancelled {
        return "", fmt.Errorf("scope input cancelled")
    }
    return textModel.value, nil
}
```

### Step 6: Add Visual Indicator Helper

Create helper function in `internal/ui/display.go`:

```go
func GetVisualIndicator(state PromptState) string {
    style := lipgloss.NewStyle()
    switch state {
    case StatePending, StateActive:
        return style.Foreground(lipgloss.Color("4")).Render("?")
    case StateCompleted:
        return style.Foreground(lipgloss.Color("2")).Render("✓")
    case StateCancelled:
        return style.Foreground(lipgloss.Color("1")).Render("✗")
    case StateError:
        return style.Foreground(lipgloss.Color("3")).Render("⚠")
    default:
        return "?"
    }
}
```

## Testing Strategy

### Unit Tests

For each model:
- Test state transitions
- Test keyboard handling (Enter, Escape, input)
- Test visual indicator rendering
- Test validation (if applicable)
- Test default values (if applicable)

### Integration Tests

- Test full prompt workflows
- Test cancellation scenarios
- Test validation error scenarios
- Test backward compatibility (function signatures)

## Common Patterns

### Pattern 1: Text Input with Validation

```go
validator := func(value string) error {
    if value == "" {
        return fmt.Errorf("value cannot be empty")
    }
    if len(value) > 72 {
        return fmt.Errorf("value too long (max 72 characters)")
    }
    return nil
}

model := NewTextInputModel("Prompt:", "", validator)
```

### Pattern 2: Multiline Input with Default

```go
model := NewMultilineInputModelWithValue("Body:", defaultValue)
```

### Pattern 3: Select List with Preselection

```go
model := NewSelectListModelWithPreselection("feat")
```

## Visual Indicator Colors

- Blue: `lipgloss.Color("4")`
- Green: `lipgloss.Color("2")`
- Red: `lipgloss.Color("1")`
- Yellow: `lipgloss.Color("3")`

## Key Reminders

1. **No Alt Screen**: Never use `tea.WithAltScreen()` - all prompts render inline
2. **Visual Indicators**: Always show appropriate indicator based on state
3. **Backward Compatibility**: Maintain all function signatures
4. **State Management**: Use `PromptState` enum consistently
5. **Error Handling**: Show yellow '⚠' for validation errors, allow correction
6. **Result Display**: Select lists show value after title, multiline shows result below with wrapping

## Next Steps

1. Write tests for new models (TDD)
2. Implement TextInputModel
3. Implement YesNoChoiceModel
4. Update SelectListModel
5. Update MultilineInputModel
6. Refactor all prompt functions
7. Run integration tests
8. Verify backward compatibility
