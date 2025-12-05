# Quickstart: Rewrite CLI Prompts with Huh Library

**Feature**: 013-huh-prompts
**Date**: 2025-01-27

## Overview

This guide provides a quick start for implementing the migration of CLI prompts from custom Bubble Tea models to the `huh` library.

## Prerequisites

- Go 1.25.0 or later
- Understanding of existing prompt functions in `internal/ui/prompts.go`
- Familiarity with `huh` library (see [research.md](./research.md))

## Step 1: Add `huh` Dependency

```bash
go get github.com/charmbracelet/huh@latest
```

Verify the dependency was added to `go.mod`:

```go
require (
    github.com/charmbracelet/huh v0.8.0  // or latest version
)
```

## Step 2: Understand Existing Prompt Functions

Review the existing prompt functions in `internal/ui/prompts.go`:

- `PromptScope`, `PromptScopeWithDefault`
- `PromptSubject`, `PromptSubjectWithDefault`
- `PromptBody`, `PromptBodyWithDefault`
- `PromptFooter`, `PromptFooterWithDefault`
- `PromptCommitType`, `PromptCommitTypeWithPreselection`
- `PromptEmptyCommit`, `PromptConfirm`, `PromptAIUsage`, etc.

**Key Points**:
- All functions take `*bufio.Reader` (may be unused after migration)
- All functions return `(value, error)`
- Functions handle validation, defaults, and pre-selection
- Functions are called from `internal/service/commit_service.go`

## Step 3: Create Basic `huh` Form Example

Start with a simple prompt migration. Example for `PromptScope`:

```go
package ui

import (
    "bufio"
    "fmt"
    "github.com/charmbracelet/huh"
)

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
        return "", fmt.Errorf("scope input cancelled: %w", err)
    }

    // Print post-validation summary line
    fmt.Printf("✓ Scope (optional): %s\n", scope)

    return scope, nil
}
```

## Step 4: Add Validation

For prompts with validation (e.g., `PromptSubject`):

```go
func PromptSubject(reader *bufio.Reader) (string, error) {
    var subject string

    validator := func(value string) error {
        if strings.TrimSpace(value) == "" {
            return fmt.Errorf("subject cannot be empty")
        }
        return nil
    }

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Subject (required)").
                Value(&subject).
                Validate(validator),
        ),
    )

    if err := form.Run(); err != nil {
        return "", fmt.Errorf("subject input cancelled: %w", err)
    }

    subject = strings.TrimSpace(subject)
    fmt.Printf("✓ Subject (required): %s\n", subject)

    return subject, nil
}
```

## Step 5: Handle Default Values

For prompts with defaults (e.g., `PromptScopeWithDefault`):

```go
func PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error) {
    scope := defaultValue

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Scope").
                Value(&scope),
        ),
    )

    if err := form.Run(); err != nil {
        return "", fmt.Errorf("scope input cancelled: %w", err)
    }

    // If empty and default exists, return default
    if scope == "" && defaultValue != "" {
        scope = defaultValue
    }

    fmt.Printf("✓ Scope: %s\n", scope)
    return scope, nil
}
```

## Step 6: Implement Select Lists

For commit type selection:

```go
func PromptCommitType(reader *bufio.Reader) (string, error) {
    var commitType string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Choose a type").
                Options(
                    huh.NewOption("feat", "feat"),
                    huh.NewOption("fix", "fix"),
                    huh.NewOption("docs", "docs"),
                    // ... more options
                ).
                Value(&commitType),
        ),
    )

    if err := form.Run(); err != nil {
        return "", fmt.Errorf("commit type selection cancelled: %w", err)
    }

    fmt.Printf("✓ Choose a type: %s\n", commitType)
    return commitType, nil
}
```

## Step 7: Implement Confirmations

For yes/no prompts:

```go
func PromptConfirm(reader *bufio.Reader, message string) (bool, error) {
    var confirm bool

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewConfirm().
                Title(message).
                Value(&confirm),
        ),
    )

    if err := form.Run(); err != nil {
        return false, fmt.Errorf("confirm prompt cancelled: %w", err)
    }

    result := "Yes"
    if !confirm {
        result = "No"
    }
    fmt.Printf("✓ %s: %s\n", message, result)

    return confirm, nil
}
```

## Step 8: Implement Multiline Input

For body/footer prompts:

```go
func PromptBody(reader *bufio.Reader) (string, error) {
    var body string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewText().
                Title("Body").
                Value(&body),
        ),
    )

    if err := form.Run(); err != nil {
        return "", fmt.Errorf("body input cancelled: %w", err)
    }

    // Truncate for summary line if needed
    summary := body
    if len(body) > 50 {
        summary = body[:50] + "..."
    }
    fmt.Printf("✓ Body: %s\n", summary)

    return body, nil
}
```

## Step 9: Test the Migration

Write tests for each migrated function:

```go
func TestPromptScope(t *testing.T) {
    // Test with valid input
    // Test with empty input (skip)
    // Test with cancellation
    // Test post-validation display format
}
```

## Step 10: Remove Old Models

After all prompts are migrated:

1. Delete `text_input.go` and `text_input_test.go`
2. Delete `multiline_input.go` and `multiline_input_test.go`
3. Delete `yes_no_choice.go` and `yes_no_choice_test.go`
4. Delete `select_list.go` and `select_list_test.go`
5. Review `prompt_state.go` - may be removable if not used elsewhere

## Step 11: Update Integration Tests

Update integration tests in `test/integration/` to verify:
- Prompts render inline (no alt screen)
- Post-validation summary lines appear
- Validation errors display inline
- All existing functionality preserved

## Common Patterns

### Combined Forms (Future Enhancement)

For combining related prompts into a single form:

```go
func buildCommitMessageForm(msg *CommitMessage) *huh.Form {
    return huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().Title("Type").Value(&msg.Type),
            huh.NewInput().Title("Scope").Value(&msg.Scope),
            huh.NewInput().Title("Subject").Value(&msg.Subject),
            huh.NewText().Title("Body").Value(&msg.Body),
            huh.NewText().Title("Footer").Value(&msg.Footer),
        ),
    )
}
```

Note: For initial implementation, individual functions are simpler and maintain backward compatibility.

## Troubleshooting

### Issue: Form doesn't render inline

**Solution**: `huh` forms render inline by default. If alt screen appears, check form configuration.

### Issue: Post-validation summary doesn't appear

**Solution**: Ensure you print the summary line after `form.Run()` completes successfully.

### Issue: Validation errors don't show inline

**Solution**: Use `.Validate()` method on the field. `huh` displays errors automatically.

### Issue: Default values don't pre-fill

**Solution**: Set the value variable before creating the form:
```go
scope := defaultValue
form := huh.NewForm(...huh.NewInput().Value(&scope)...)
```

## Next Steps

1. Migrate one prompt function at a time
2. Write tests for each migrated function
3. Verify backward compatibility with existing callers
4. Update integration tests
5. Remove old Bubble Tea models after all prompts are migrated

## References

- [huh Library Documentation](https://github.com/charmbracelet/huh)
- [Research Findings](./research.md)
- [Data Model](./data-model.md)
- [Function Contracts](./contracts/prompt-functions-contract.md)
