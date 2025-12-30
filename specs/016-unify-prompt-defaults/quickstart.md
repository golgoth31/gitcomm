# Quickstart: Unify Prompt Functions to Use Default Variants

**Date**: 2025-01-27
**Feature**: 016-unify-prompt-defaults

## Overview

This refactoring simplifies the `promptCommitMessage` function by always using prompt functions with default value support, eliminating conditional logic. It also removes unused non-default prompt functions.

## What Changed

### Before

The `promptCommitMessage` function used conditional logic to select between regular and "WithDefault" prompt functions:

```go
if prefilled != nil && prefilled.Type != "" {
    commitType, err = ui.PromptCommitTypeWithPreselection(s.reader, prefilled.Type)
} else {
    commitType, err = ui.PromptCommitType(s.reader)
}
```

### After

The function always uses "WithDefault" variants, passing empty strings when no pre-filled data exists:

```go
defaultType := ""
if prefilled != nil && prefilled.Type != "" {
    defaultType = prefilled.Type
}
commitType, err := ui.PromptCommitTypeWithPreselection(s.reader, defaultType)
```

## Key Concepts

### Unified Prompt Pattern

All prompts now follow the same pattern:
1. Determine default value (empty string if no pre-filled data)
2. Call "WithDefault" variant with the default value
3. Function handles empty defaults correctly (validation, optional fields, etc.)

### Removed Functions

The following functions are removed from `internal/ui/prompts.go`:
- `PromptScope`
- `PromptSubject`
- `PromptBody`
- `PromptFooter`
- `PromptCommitType`

These were only used in `commit_service.go` and are replaced by their "WithDefault" equivalents.

## Developer Guide

### Understanding the Refactoring

1. **Read the contract**: See `contracts/prompt-functions-contract.md` for detailed function behavior
2. **Review research**: See `research.md` for behavior verification of "WithDefault" functions with empty strings
3. **Check data model**: See `data-model.md` - no data model changes

### Testing the Changes

1. **Manual Testing**:
   - Run commit creation without pre-filled data (should behave identically)
   - Run commit creation with pre-filled data (should behave identically)
   - Verify all prompts work correctly

2. **Automated Testing**:
   - Update tests that reference removed functions
   - Add tests verifying empty string default behavior
   - Run existing integration tests

### Implementation Steps

1. **Update `promptCommitMessage` function**:
   - Replace conditional logic with unified pattern
   - Use "WithDefault" variants for all prompts
   - Pass empty strings when no pre-filled data

2. **Remove unused functions**:
   - Delete `PromptScope`, `PromptSubject`, `PromptBody`, `PromptFooter`, `PromptCommitType` from `prompts.go`
   - Verify no other references exist

3. **Update tests**:
   - Update any tests referencing removed functions
   - Verify all tests pass

## Common Patterns

### Pattern: Getting Default Value

```go
defaultValue := ""
if prefilled != nil && prefilled.Field != "" {
    defaultValue = prefilled.Field
}
result, err := ui.PromptFieldWithDefault(s.reader, defaultValue)
```

### Pattern: Optional Fields (Scope, Body, Footer)

```go
defaultValue := ""
if prefilled != nil {
    defaultValue = prefilled.Field
}
result, err := ui.PromptFieldWithDefault(s.reader, defaultValue)
```

### Pattern: Required Fields (Subject)

```go
defaultValue := ""
if prefilled != nil && prefilled.Subject != "" {
    defaultValue = prefilled.Subject
}
result, err := ui.PromptSubjectWithDefault(s.reader, defaultValue)
```

## Troubleshooting

### Issue: Function not found

**Problem**: Compiler error about missing function (e.g., `PromptScope`)

**Solution**: Replace with `PromptScopeWithDefault` and pass empty string as default if no pre-filled data exists.

### Issue: Validation failing with empty defaults

**Problem**: Subject prompt fails validation when empty string is passed as default

**Solution**: This is expected behavior - `PromptSubjectWithDefault` validates that either user input or default value is non-empty. The function correctly requires user input when default is empty.

### Issue: Tests failing

**Problem**: Tests reference removed functions

**Solution**: Update tests to use "WithDefault" variants with appropriate default values.

## Next Steps

1. Review the implementation plan: `plan.md`
2. Check task breakdown: `tasks.md` (after `/speckit.tasks`)
3. Implement following TDD approach: Write/update tests → Verify failures → Implement → Verify passes
