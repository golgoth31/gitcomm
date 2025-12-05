# Research: AI Commit Message Acceptance Options

**Feature**: 009-ai-commit-options
**Date**: 2025-01-27
**Purpose**: Resolve technical unknowns and document design decisions

## Research Questions

### RQ1: UI Pattern for Three-Option Selection

**Question**: What is the best UI pattern for presenting three options (accept-and-commit, accept-and-edit, reject) to the user in a CLI context?

**Research**:
- Current pattern: Simple Y/n prompt for binary choice
- Existing patterns in codebase:
  - `PromptAIMessageEdit`: Uses numbered options (1/2) for edit vs use-as-is
  - `PromptCommitType`: Uses interactive bubbletea list for multiple options
  - `PromptConfirm`: Uses Y/n for binary choice

**Decision**: Use numbered options pattern (1/2/3) similar to `PromptAIMessageEdit`

**Rationale**:
- Consistent with existing UI patterns in codebase
- Simple and clear for three options
- No need for interactive list (only 3 options, not 8+ like commit types)
- Fast to implement and test
- Works well in CLI context

**Alternatives Considered**:
- Interactive bubbletea list: Overkill for 3 options, adds complexity
- Letter-based selection (a/e/r): Less discoverable, requires explanation
- Separate prompts: Too many steps, poor UX

---

### RQ2: Pre-filling Commit Type in Interactive List

**Question**: How should the commit type be pre-selected in the interactive bubbletea list when pre-filling during "accept and edit"?

**Research**:
- `SelectListModel` has `SelectedIndex` field that controls which item is highlighted
- `NewSelectListModel()` currently defaults to `SelectedIndex: 0` (first item)
- `PromptCommitType()` creates a new model and runs bubbletea program
- Need to find matching type in items list and set `SelectedIndex` accordingly

**Decision**: Create a new constructor `NewSelectListModelWithPreselection(type string)` that finds the matching type and sets `SelectedIndex` to that index

**Rationale**:
- Maintains existing `NewSelectListModel()` for backward compatibility
- Clear separation of concerns
- Easy to test (can verify SelectedIndex is set correctly)
- If type doesn't match, SelectedIndex remains 0 (first item) as fallback

**Alternatives Considered**:
- Modify `NewSelectListModel()` to accept optional type parameter: Breaks existing callers, requires refactoring
- Set SelectedIndex after model creation: Works but less clean, harder to test
- Skip type selection entirely: Violates requirement to pre-fill all fields

---

### RQ3: Pre-filling Text Fields (Scope, Subject, Body, Footer)

**Question**: How should text input fields be pre-filled when user selects "accept and edit"?

**Research**:
- `PromptScope`, `PromptSubject`: Use simple `bufio.Reader.ReadString('\n')` with default empty
- `PromptBody`, `PromptFooter`: Use `MultilineInputModel` with bubbletea
- `MultilineInputModel` has `GetValue()` method and `Cancelled` field
- Need to set initial value in the model before running bubbletea program

**Decision**:
- For simple prompts (scope, subject): Display default value in prompt text, allow user to press Enter to accept or type to replace
- For multiline prompts (body, footer): Create `NewMultilineInputModelWithValue(fieldName, initialValue)` that pre-populates the model's value

**Rationale**:
- Simple prompts: Minimal change, clear UX (shows default, user can accept or edit)
- Multiline prompts: Need to pre-populate model state, requires new constructor similar to commit type approach
- Maintains backward compatibility with existing constructors

**Alternatives Considered**:
- Always show empty and require user to type: Violates requirement to pre-fill
- Auto-accept pre-filled values without showing: Poor UX, user can't see what they're accepting
- Separate "accept default" vs "edit" prompt: Too many steps, violates requirement

---

### RQ4: Return Type for Acceptance Options

**Question**: What should `PromptAIMessageAcceptance` return to represent three options instead of a boolean?

**Research**:
- Current: Returns `(bool, error)` where true = accept, false = reject
- Go best practices: Use custom type for enum-like values
- Options: string constant, int constant, custom type with methods

**Decision**: Create custom type `AIMessageAcceptance` with three constants and String() method

**Rationale**:
- Type-safe (compiler catches invalid values)
- Self-documenting (clear what each value means)
- Extensible (easy to add more options later if needed)
- Follows Go idioms for enum-like values

**Implementation**:
```go
type AIMessageAcceptance int

const (
    AcceptAndCommit AIMessageAcceptance = iota
    AcceptAndEdit
    Reject
)

func (a AIMessageAcceptance) String() string { ... }
```

**Alternatives Considered**:
- String constants: Less type-safe, can have typos
- Int constants: Less readable, magic numbers
- Separate boolean flags: Doesn't scale, awkward API

---

### RQ5: Error Recovery After Commit Failure

**Question**: What user experience should be provided when commit fails after "accept and commit directly"?

**Research**:
- Existing error handling: `CreateCommit` returns error, service handles it
- Current flow: Error triggers staging state restoration via defer
- User needs: Clear error message, option to retry or edit

**Decision**: Display error message, restore staging state, then prompt user to choose: retry commit with same message, edit message, or cancel

**Rationale**:
- Maintains data safety (staging restored)
- Gives user control (retry, edit, or abort)
- Consistent with existing error handling patterns
- Clear recovery path

**Alternatives Considered**:
- Auto-retry: Could loop indefinitely, poor UX
- Auto-fallback to edit: Assumes user wants to edit, may not be true
- Just show error and exit: Poor UX, user loses work

---

## Summary of Decisions

1. **UI Pattern**: Use numbered options (1/2/3) for three-option selection
2. **Pre-filling Commit Type**: Create `NewSelectListModelWithPreselection()` constructor
3. **Pre-filling Text Fields**: Display defaults for simple prompts, pre-populate model for multiline
4. **Return Type**: Custom `AIMessageAcceptance` type with three constants
5. **Error Recovery**: Restore staging, show error, prompt for retry/edit/cancel

## Open Questions Resolved

All technical unknowns have been resolved. No blocking issues remain.

## Dependencies

- **Existing**: `github.com/charmbracelet/bubbletea` v1.3.10 - for interactive UI
- **Existing**: `github.com/charmbracelet/lipgloss` v1.1.0 - for UI styling
- **No new dependencies required**
