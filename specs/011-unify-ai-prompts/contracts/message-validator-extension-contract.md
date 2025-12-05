# Contract: MessageValidator Interface Extension

**Package**: `pkg/conventional`
**Date**: 2025-01-27
**Feature**: 011-unify-ai-prompts

## Interface Extension

The existing `MessageValidator` interface is extended with methods to extract validation rules programmatically.

### Existing Interface

```go
type MessageValidator interface {
    Validate(message *model.CommitMessage) (bool, []ValidationError)
}
```

### Extended Interface

```go
type MessageValidator interface {
    // Existing method
    Validate(message *model.CommitMessage) (bool, []ValidationError)

    // New methods for rule extraction
    GetValidTypes() []string
    GetSubjectMaxLength() int
    GetBodyMaxLength() int
    GetScopeFormatDescription() string
}
```

## Contract Specifications

### GetValidTypes

**Signature**: `GetValidTypes() []string`

**Preconditions**: None

**Postconditions**:
- Returns a slice of valid commit type strings
- Must return exactly: `["feat", "fix", "docs", "style", "refactor", "test", "chore", "version"]`
- Order is not guaranteed but should be consistent

**Error Conditions**: None (always succeeds)

**Thread Safety**: Must be thread-safe

**Example**:
```go
types := validator.GetValidTypes()
// Returns: ["feat", "fix", "docs", "style", "refactor", "test", "chore", "version"]
```

---

### GetSubjectMaxLength

**Signature**: `GetSubjectMaxLength() int`

**Preconditions**: None

**Postconditions**:
- Returns the maximum allowed length for commit message subject
- Must return exactly: `72`

**Error Conditions**: None (always succeeds)

**Thread Safety**: Must be thread-safe

**Example**:
```go
maxLen := validator.GetSubjectMaxLength()
// Returns: 72
```

---

### GetBodyMaxLength

**Signature**: `GetBodyMaxLength() int`

**Preconditions**: None

**Postconditions**:
- Returns the maximum allowed length for commit message body
- Must return exactly: `320`

**Error Conditions**: None (always succeeds)

**Thread Safety**: Must be thread-safe

**Example**:
```go
maxLen := validator.GetBodyMaxLength()
// Returns: 320
```

---

### GetScopeFormatDescription

**Signature**: `GetScopeFormatDescription() string`

**Preconditions**: None

**Postconditions**:
- Returns a human-readable description of valid scope format
- Must return exactly: `"alphanumeric, hyphens, underscores only"`

**Error Conditions**: None (always succeeds)

**Thread Safety**: Must be thread-safe

**Example**:
```go
desc := validator.GetScopeFormatDescription()
// Returns: "alphanumeric, hyphens, underscores only"
```

---

## Implementation Contract

### Validator (Extended)

**Type**: `struct` (implements `MessageValidator`)

**Existing Methods**: Unchanged

**New Methods**:
- `GetValidTypes() []string`: Returns hardcoded list matching `isValidType()` logic
- `GetSubjectMaxLength() int`: Returns constant `72`
- `GetBodyMaxLength() int`: Returns constant `320`
- `GetScopeFormatDescription() string`: Returns constant description matching `isValidScope()` logic

**Behavior**:
- All new methods must return values that match validation logic exactly
- Methods must be pure functions (no side effects)
- Methods must be thread-safe

**Implementation Notes**:
- Values are constants, extracted from existing validation logic
- No need to refactor existing validation code
- New methods are simple getters

---

## Backward Compatibility

**Status**: ✅ **FULLY BACKWARD COMPATIBLE**

**Rationale**:
- Interface extension (adding methods) is backward compatible in Go
- Existing code using `MessageValidator` continues to work
- New methods are optional to call
- No breaking changes to existing `Validate()` method

**Migration Path**:
1. Add new methods to `MessageValidator` interface
2. Implement methods in `Validator` struct
3. Existing code continues to work unchanged
4. New code (prompt generator) can use new methods

---

## Testing Contract

**Unit Tests Required**:
- Test `GetValidTypes()` returns correct list
- Test `GetSubjectMaxLength()` returns 72
- Test `GetBodyMaxLength()` returns 320
- Test `GetScopeFormatDescription()` returns correct description
- Test all methods are thread-safe (concurrent access)

**Integration Tests Required**:
- Test extracted rules match validation logic
- Test prompt generator uses extracted rules correctly

---

## Breaking Changes

**None**: Interface extension is backward compatible.

---

## Version History

- **v1.0.0** (2025-01-27): Initial extension contract definition
