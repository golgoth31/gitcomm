# Research: Display Commit Type Selection Confirmation

**Feature**: 006-type-confirmation-display
**Date**: 2025-01-27

## Technology Decisions

### 1. Display Method for Confirmation Line

**Decision**: Use `fmt.Printf` with standard formatting to display the confirmation line after bubbletea alt-screen exits.

**Rationale**:
- `fmt.Printf` is the idiomatic Go way to format and print to standard output
- Simple and straightforward - no additional dependencies needed
- Works reliably across all platforms (Linux, macOS, Windows)
- Standard library solution, no external dependencies
- Matches existing codebase patterns for terminal output (e.g., `PromptScope`, `PromptSubject` use `fmt.Print`)

**Alternatives Considered**:
- Using bubbletea's output formatting: More complex, requires modifying the model's View() method, unnecessary for simple confirmation
- Using a separate display utility: Over-engineering for a simple printf statement
- Using structured logging: Not appropriate for user-facing confirmation messages

**Implementation Pattern**:
```go
// After bubbletea program exits and before returning
selectedType, err := selectModel.GetSelectedType()
if err != nil {
    return "", fmt.Errorf("failed to get selected type: %w", err)
}

// Display confirmation line
fmt.Printf("✔ Choose a type(<scope>): %s\n", selectedType)

return selectedType, nil
```

### 2. Timing of Confirmation Display

**Decision**: Display confirmation line immediately after bubbletea alt-screen exits, before returning from `PromptCommitType`.

**Rationale**:
- Alt-screen is already cleared when bubbletea program exits
- Standard terminal output is immediately available
- No race conditions or timing issues
- User sees confirmation before next prompt appears
- Matches clarification: "on a new line after the alt-screen clears"

**Alternatives Considered**:
- Displaying within bubbletea View(): Would require keeping alt-screen active, conflicts with requirement to show in standard output
- Displaying in calling code: Breaks encapsulation, confirmation logic should be with selection logic
- Delayed display: Unnecessary complexity, immediate feedback is better UX

### 3. Checkmark Symbol Support

**Decision**: Use Unicode checkmark symbol (✔) as specified in requirements, with assumption that modern terminals support it.

**Rationale**:
- Specified in requirements (FR-006)
- Modern terminals (Linux, macOS, Windows with modern terminal emulators) support Unicode
- Provides clear visual feedback
- Matches common CLI confirmation patterns

**Alternatives Considered**:
- ASCII fallback (e.g., "[OK]"): Less visually appealing, but could be added if encoding issues arise
- No symbol: Less clear visual indication of success

**Fallback Strategy**: If encoding issues are discovered during testing, can add terminal encoding detection and fallback to ASCII "[OK]" or similar.

### 4. Format String Consistency

**Decision**: Preserve exact format "Choose a type(<scope>):" from original prompt, only adding checkmark prefix and chosen type suffix.

**Rationale**:
- Maintains consistency with existing UI
- Users recognize the format from the original prompt
- Clear connection between prompt and confirmation
- Specified in requirements (FR-007)

**Implementation**:
- Format string: `"✔ Choose a type(<scope>): %s\n"`
- Literal text "(<scope>)" preserved as-is (not a variable)
- Chosen type inserted via `%s` format specifier
