# Research: Rewrite CLI Prompts with Huh Library

**Feature**: 013-huh-prompts
**Date**: 2025-01-27
**Status**: Complete

## Research Questions

### 1. `huh` Library Capabilities and API

**Question**: What are the capabilities of the `huh` library? Does it support inline rendering, validation, default values, and all required prompt types?

**Research Task**: Research `huh` library API, capabilities, and best practices for inline rendering.

**Findings**:
- **Decision**: Use `github.com/charmbracelet/huh` v0.8.0 or latest stable version
- **Rationale**:
  - `huh` is built on Bubble Tea (which we already use) and provides form-building capabilities
  - Supports all required field types: `NewInput()`, `NewText()`, `NewSelect()`, `NewConfirm()`, `NewMultiSelect()`
  - Built-in validation support via `.Validate()` method
  - Supports default values via `.Value()` method
  - Supports pre-selection via `.Selected()` for options
  - Can render inline (not in alt screen) by default or via configuration
  - Forms can be run with `.Run()` which handles the Bubble Tea program lifecycle
- **Alternatives Considered**:
  - Continue using custom Bubble Tea models: Rejected - violates requirement to use `huh` library
  - Use other form libraries: Rejected - `huh` is explicitly required and is well-maintained by Charm team
- **API Patterns**:
  ```go
  // Single field form
  form := huh.NewForm(
    huh.NewGroup(
      huh.NewInput().
        Title("Prompt title").
        Value(&value).
        Validate(validator),
    ),
  )
  err := form.Run()

  // Multi-field form
  form := huh.NewForm(
    huh.NewGroup(
      huh.NewInput().Title("Field 1").Value(&field1),
      huh.NewInput().Title("Field 2").Value(&field2),
    ),
  )
  err := form.Run()
  ```

### 2. Inline Rendering Configuration

**Question**: How to configure `huh` forms to render inline without alt screen?

**Research Task**: Research `huh` configuration options for inline rendering.

**Findings**:
- **Decision**: Use `huh` forms without alt screen configuration (default behavior) or use `form.WithAccessible(false)` if needed
- **Rationale**:
  - `huh` forms render inline by default when run with `.Run()`
  - Alt screen is typically used for full-screen TUI applications, not inline prompts
  - The library handles terminal rendering automatically via Bubble Tea
  - No special configuration needed for inline rendering
- **Alternatives Considered**:
  - Custom terminal mode switching: Rejected - `huh` handles this automatically
  - Force alt screen off: Not needed - default behavior is inline

### 3. Post-Validation Display Implementation

**Question**: How to implement the post-validation display format (green checkmark + summary line) after form completion?

**Research Task**: Research `huh` form completion callbacks and custom display logic.

**Findings**:
- **Decision**: Implement custom display logic after form completion using `form.State` to detect completion, then manually print summary lines
- **Rationale**:
  - `huh` forms have a `State` field that indicates completion (`huh.StateCompleted`)
  - After `form.Run()` completes successfully, we can access field values and print custom summary lines
  - For progressive display (multi-field forms), we may need to use form events or implement custom rendering
  - The library doesn't provide built-in post-validation display hooks, so custom implementation is required
- **Alternatives Considered**:
  - Use `huh`'s built-in completion display: Rejected - doesn't match required format (green checkmark + summary)
  - Modify `huh` library: Rejected - better to work with library as-is and add custom display logic
- **Implementation Approach**:
  ```go
  form := huh.NewForm(...)
  err := form.Run()
  if err != nil {
    return err
  }

  // After form completion, print summary lines
  for _, field := range form.Fields {
    if field.State == huh.StateCompleted {
      fmt.Printf("✓ %s: %v\n", field.Title, field.GetValue())
    }
  }
  ```

### 4. Progressive Field Display in Multi-Field Forms

**Question**: How to show summary lines progressively as each field is completed in a multi-field form?

**Research Task**: Research `huh` form field state transitions and event handling.

**Findings**:
- **Decision**: Use `huh` form field state monitoring or implement custom Bubble Tea integration to detect field completion and print summary lines progressively
- **Rationale**:
  - `huh` forms track field state (`huh.StateCompleted` per field)
  - We can monitor form state during execution or use form events
  - Alternative: Run forms field-by-field (single-field forms) for simpler progressive display
  - For combined forms, we may need to use `huh`'s Bubble Tea integration to hook into field completion events
- **Alternatives Considered**:
  - Run separate single-field forms: Rejected - requirement is to combine related prompts
  - Wait until entire form completes: Rejected - requirement is progressive display per field
- **Implementation Approach**:
  - Option A: Use `huh.Form` as `tea.Model` and integrate with Bubble Tea to detect field completion
  - Option B: After form completion, iterate through fields and print summaries in order (simpler but not truly progressive)
  - **Recommended**: Option B for initial implementation (simpler), Option A for future enhancement if needed

### 5. Validation Error Display

**Question**: How does `huh` display validation errors inline?

**Research Task**: Research `huh` validation error display mechanism.

**Findings**:
- **Decision**: Use `huh`'s built-in validation error display via `.Validate()` method
- **Rationale**:
  - `huh` automatically displays validation errors inline below the input field
  - Errors are shown when validation fails and user attempts to proceed
  - Error display is handled by the library automatically
  - No custom error display logic needed
- **Alternatives Considered**:
  - Custom error display: Rejected - `huh` provides this automatically
  - External error messages: Rejected - requirement is inline display

### 6. Backward Compatibility Strategy

**Question**: How to maintain existing prompt function signatures while using `huh` internally?

**Research Task**: Research adapter pattern for wrapping `huh` forms in existing function signatures.

**Findings**:
- **Decision**: Create adapter functions that wrap `huh` form creation and execution, maintaining existing function signatures
- **Rationale**:
  - Existing functions like `PromptScope(reader *bufio.Reader) (string, error)` can be kept
  - Internal implementation will use `huh` instead of custom Bubble Tea models
  - `bufio.Reader` parameter can be ignored (not needed with `huh`)
  - Return types remain the same: `(string, error)`, `(bool, error)`, etc.
- **Alternatives Considered**:
  - Change function signatures: Rejected - violates backward compatibility requirement
  - Create new functions: Rejected - would break existing callers
- **Implementation Approach**:
  ```go
  // Existing signature maintained
  func PromptScope(reader *bufio.Reader) (string, error) {
    var scope string
    form := huh.NewForm(
      huh.NewGroup(
        huh.NewInput().
          Title("Scope (optional)").
          Value(&scope),
      ),
    )
    if err := form.Run(); err != nil {
      return "", err
    }
    // Print summary line
    fmt.Printf("✓ Scope (optional): %s\n", scope)
    return scope, nil
  }
  ```

### 7. Combined Forms Implementation

**Question**: How to combine related prompts (e.g., commit message fields) into a single `huh.Form` while maintaining individual function signatures?

**Research Task**: Research form composition patterns and adapter strategies.

**Findings**:
- **Decision**: Create a shared form builder function that creates combined forms, but wrap it in individual prompt functions that extract only their specific field value
- **Rationale**:
  - Can create a combined form with all commit message fields (type, scope, subject, body, footer)
  - Individual functions can call the shared form builder and extract only their field
  - Alternative: Create a new combined prompt function that returns all values, but keep individual functions for backward compatibility
  - For progressive display, we'll need to handle field-by-field completion
- **Alternatives Considered**:
  - Always use individual forms: Rejected - requirement is to combine related prompts
  - Change all callers to use combined form: Rejected - violates backward compatibility
- **Implementation Approach**:
  ```go
  // Shared form builder
  func buildCommitMessageForm(commitMsg *CommitMessage) *huh.Form {
    return huh.NewForm(
      huh.NewGroup(
        huh.NewSelect[string]().Title("Type").Options(...).Value(&commitMsg.Type),
        huh.NewInput().Title("Scope").Value(&commitMsg.Scope),
        huh.NewInput().Title("Subject").Value(&commitMsg.Subject),
        huh.NewText().Title("Body").Value(&commitMsg.Body),
        huh.NewText().Title("Footer").Value(&commitMsg.Footer),
      ),
    )
  }

  // Individual functions can use shared form or create their own
  func PromptScope(reader *bufio.Reader) (string, error) {
    var scope string
    form := huh.NewForm(huh.NewGroup(huh.NewInput().Title("Scope").Value(&scope)))
    // ... or use shared form and extract scope
  }
  ```

## Dependencies

### New Dependencies

- `github.com/charmbracelet/huh` - Form building library (version to be determined - latest stable)

### Existing Dependencies (Used)

- `github.com/charmbracelet/bubbletea` - Already in use, `huh` depends on it
- `github.com/charmbracelet/lipgloss` - Already in use, for styling

## Integration Points

1. **Prompt Functions** (`internal/ui/prompts.go`): Replace custom Bubble Tea models with `huh` forms
2. **Display Utilities** (`internal/ui/display.go`): May need updates for post-validation format helper
3. **Commit Service** (`internal/service/commit_service.go`): No changes needed (backward compatible)

## Open Questions Resolved

- ✅ `huh` library supports all required field types
- ✅ Inline rendering is default behavior
- ✅ Validation errors display inline automatically
- ✅ Default values and pre-selection supported
- ✅ Backward compatibility achievable via adapter pattern
- ✅ Combined forms can be implemented with shared form builders

## Remaining Implementation Details

- Exact `huh` version to use (check latest stable release)
- Progressive display implementation details (may need Bubble Tea integration)
- Post-validation summary line formatting helper function
- Testing strategy for `huh` integration (mocking may be challenging)
