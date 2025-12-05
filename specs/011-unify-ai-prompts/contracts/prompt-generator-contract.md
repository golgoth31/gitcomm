# Contract: PromptGenerator Interface

**Package**: `pkg/ai/prompt`
**Date**: 2025-01-27
**Feature**: 011-unify-ai-prompts

## Interface Definition

```go
package prompt

import (
    "github.com/golgoth31/gitcomm/internal/model"
    "github.com/golgoth31/gitcomm/pkg/conventional"
)

// PromptGenerator defines the interface for generating unified AI prompts
type PromptGenerator interface {
    // GenerateSystemMessage generates the system message with validation rules
    // extracted from the MessageValidator
    GenerateSystemMessage(validator conventional.MessageValidator) (string, error)

    // GenerateUserMessage generates the user message with repository state
    // formatted for AI consumption
    GenerateUserMessage(repoState *model.RepositoryState) (string, error)
}
```

## Contract Specifications

### GenerateSystemMessage

**Signature**: `GenerateSystemMessage(validator conventional.MessageValidator) (string, error)`

**Preconditions**:
- `validator` must not be `nil`
- `validator` must implement the extended `MessageValidator` interface (with getter methods)

**Postconditions**:
- Returns a non-empty string containing validation rules formatted as structured bullet points
- Returns error if validator is `nil` or if rule extraction fails
- System message must include:
  - Conventional Commits format specification
  - Valid commit types (from `validator.GetValidTypes()`)
  - Subject length constraint (from `validator.GetSubjectMaxLength()`)
  - Body length constraint (from `validator.GetBodyMaxLength()`)
  - Scope format constraint (from `validator.GetScopeFormatDescription()`)

**Error Conditions**:
- `ErrNilValidator`: validator is `nil`
- `ErrRuleExtractionFailed`: failed to extract validation rules

**Thread Safety**: Must be thread-safe (stateless implementation)

**Example Output**:
```
You are a git commit message generator. When receiving a git diff, you will ONLY generate commit messages following the Conventional Commits specification.

Format: type(scope): subject

body

footer

Validation Rules:
• Type must be one of: feat, fix, docs, style, refactor, test, chore, version
• Subject must be ≤72 characters
• Body must be ≤320 characters (if provided)
• Scope must be a valid identifier (alphanumeric, hyphens, underscores only)
```

---

### GenerateUserMessage

**Signature**: `GenerateUserMessage(repoState *model.RepositoryState) (string, error)`

**Preconditions**:
- `repoState` must not be `nil`

**Postconditions**:
- Returns a string containing formatted repository state
- Returns error if `repoState` is `nil`
- User message must include:
  - Header: "Generate a commit message for the following changes:\n\n"
  - Staged files section (if any)
  - Unstaged files section (if any)
  - Each file entry: "- {path} ({status})\n{diff}\n"

**Error Conditions**:
- `ErrNilRepositoryState`: repoState is `nil`

**Thread Safety**: Must be thread-safe (stateless implementation)

**Example Output**:
```
Generate a commit message for the following changes:

Staged files:
- internal/ai/openai_provider.go (modified)
diff --git a/internal/ai/openai_provider.go b/internal/ai/openai_provider.go
...

Unstaged files:
- test/integration/prompt_test.go (added)
...
```

---

## Implementation Contract

### UnifiedPromptGenerator

**Type**: `struct` (implements `PromptGenerator`)

**Constructor**: `NewUnifiedPromptGenerator() PromptGenerator`

**Fields**: None (stateless)

**Methods**:
- `GenerateSystemMessage(validator conventional.MessageValidator) (string, error)`
- `GenerateUserMessage(repoState *model.RepositoryState) (string, error)`

**Behavior**:
- Must extract validation rules programmatically from validator
- Must format rules as structured bullet points
- Must handle empty repository state gracefully
- Must be thread-safe (no shared mutable state)

---

## Provider Usage Contract

### OpenAI/Mistral/Local Providers

**Usage Pattern**:
```go
systemMsg, err := generator.GenerateSystemMessage(validator)
userMsg, err := generator.GenerateUserMessage(repoState)
// Use as separate system and user messages in API call
```

### Anthropic Provider

**Usage Pattern**:
```go
systemMsg, err := generator.GenerateSystemMessage(validator)
userMsg, err := generator.GenerateUserMessage(repoState)
combinedMsg := systemMsg + "\n\n" + userMsg
// Use combined message as single user message in API call
```

---

## Testing Contract

**Unit Tests Required**:
- Test `GenerateSystemMessage` with valid validator
- Test `GenerateSystemMessage` with `nil` validator (error case)
- Test `GenerateUserMessage` with valid repository state
- Test `GenerateUserMessage` with empty repository state
- Test `GenerateUserMessage` with `nil` repository state (error case)
- Test prompt consistency across multiple calls

**Integration Tests Required**:
- Test all providers use identical prompts (system + user)
- Test Anthropic prepends system to user correctly
- Test prompt content matches validation rules

---

## Breaking Changes

**None**: This is a new interface, no breaking changes to existing code.

---

## Version History

- **v1.0.0** (2025-01-27): Initial contract definition
