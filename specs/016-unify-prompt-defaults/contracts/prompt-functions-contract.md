# Prompt Functions Contract: Unify Prompt Functions to Use Default Variants

**Date**: 2025-01-27
**Feature**: 016-unify-prompt-defaults
**Type**: Internal Function Contract

## Overview

This contract documents the changes to prompt function usage in `commit_service.go` and the removal of non-default prompt functions from `prompts.go`.

## Functions to Remove

The following functions will be removed from `internal/ui/prompts.go`:

1. `PromptScope(reader *bufio.Reader) (string, error)`
2. `PromptSubject(reader *bufio.Reader) (string, error)`
3. `PromptBody(reader *bufio.Reader) (string, error)`
4. `PromptFooter(reader *bufio.Reader) (string, error)`
5. `PromptCommitType(reader *bufio.Reader) (string, error)`

**Rationale**: These functions are only used in `commit_service.go` and will be replaced with "WithDefault" variants.

## Functions to Use (Unchanged)

The following functions remain in `internal/ui/prompts.go` and will be used exclusively:

1. `PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`
2. `PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`
3. `PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`
4. `PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`
5. `PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error)`

## Function Behavior Contract

### PromptScopeWithDefault

**Signature**: `func PromptScopeWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Behavior**:
- Initializes form with `defaultValue` as initial value
- Allows empty input (scope is optional)
- If user enters empty value and `defaultValue != ""`, returns `defaultValue`
- Returns user input or empty string

**Usage in refactored code**:
- When `prefilled == nil` or `prefilled.Scope == ""`: Pass `""` as `defaultValue`
- When `prefilled != nil` and `prefilled.Scope != ""`: Pass `prefilled.Scope` as `defaultValue`

### PromptSubjectWithDefault

**Signature**: `func PromptSubjectWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Behavior**:
- Initializes form with `defaultValue` as initial value
- Validates: If both user input (trimmed) and `defaultValue` are empty, returns error "subject cannot be empty"
- If user enters empty value and `defaultValue != ""`, returns `defaultValue`
- Returns trimmed user input

**Usage in refactored code**:
- When `prefilled == nil` or `prefilled.Subject == ""`: Pass `""` as `defaultValue`
- When `prefilled != nil` and `prefilled.Subject != ""`: Pass `prefilled.Subject` as `defaultValue`

### PromptBodyWithDefault

**Signature**: `func PromptBodyWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Behavior**:
- Initializes form with `defaultValue` as initial value
- Allows empty input (body is optional)
- Warns if body exceeds 320 characters
- Returns user input or empty string

**Usage in refactored code**:
- When `prefilled == nil` or `prefilled.Body == ""`: Pass `""` as `defaultValue`
- When `prefilled != nil` and `prefilled.Body != ""`: Pass `prefilled.Body` as `defaultValue`

### PromptFooterWithDefault

**Signature**: `func PromptFooterWithDefault(reader *bufio.Reader, defaultValue string) (string, error)`

**Behavior**:
- Initializes form with `defaultValue` as initial value
- Allows empty input (footer is optional)
- Returns user input or empty string

**Usage in refactored code**:
- When `prefilled == nil` or `prefilled.Footer == ""`: Pass `""` as `defaultValue`
- When `prefilled != nil` and `prefilled.Footer != ""`: Pass `prefilled.Footer` as `defaultValue`

### PromptCommitTypeWithPreselection

**Signature**: `func PromptCommitTypeWithPreselection(reader *bufio.Reader, preselectedType string) (string, error)`

**Behavior**:
- Initializes form with `preselectedType` as initial value
- If `preselectedType` matches an option, that option is pre-selected
- If `preselectedType` is empty or doesn't match, no option is pre-selected
- User must select a type from available options
- Returns selected type

**Usage in refactored code**:
- When `prefilled == nil` or `prefilled.Type == ""`: Pass `""` as `preselectedType`
- When `prefilled != nil` and `prefilled.Type != ""`: Pass `prefilled.Type` as `preselectedType`

## Refactored Function Contract

### promptCommitMessage (Refactored)

**Location**: `internal/service/commit_service.go`

**Signature**: `func (s *CommitService) promptCommitMessage(prefilled *ui.PrefilledCommitMessage) (*model.CommitMessage, error)`

**Behavior Changes**:
- **Before**: Used conditional logic to select between regular and "WithDefault" prompt functions
- **After**: Always uses "WithDefault" variants, passing empty strings when no pre-filled data exists

**Implementation Pattern**:
```go
// Type
defaultType := ""
if prefilled != nil && prefilled.Type != "" {
    defaultType = prefilled.Type
}
commitType, err := ui.PromptCommitTypeWithPreselection(s.reader, defaultType)

// Scope
defaultScope := ""
if prefilled != nil {
    defaultScope = prefilled.Scope
}
scope, err := ui.PromptScopeWithDefault(s.reader, defaultScope)

// Subject
defaultSubject := ""
if prefilled != nil && prefilled.Subject != "" {
    defaultSubject = prefilled.Subject
}
subject, err := ui.PromptSubjectWithDefault(s.reader, defaultSubject)

// Body
defaultBody := ""
if prefilled != nil {
    defaultBody = prefilled.Body
}
body, err := ui.PromptBodyWithDefault(s.reader, defaultBody)

// Footer
defaultFooter := ""
if prefilled != nil {
    defaultFooter = prefilled.Footer
}
footer, err := ui.PromptFooterWithDefault(s.reader, defaultFooter)
```

**Error Handling**: Unchanged - errors are wrapped and returned as before.

## Validation

- All "WithDefault" functions maintain existing validation behavior when empty strings are passed
- Functionality is preserved - users experience identical behavior
- No breaking changes to external interfaces
