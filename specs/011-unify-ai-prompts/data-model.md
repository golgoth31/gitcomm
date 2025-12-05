# Data Model: Unify AI Provider Prompts with Validation Rules

**Date**: 2025-01-27
**Feature**: 011-unify-ai-prompts

## Domain Entities

### PromptGenerator

Represents the interface and implementation for generating unified AI prompts.

**Interface** (`pkg/ai/prompt/generator.go`):
```go
type PromptGenerator interface {
    // GenerateSystemMessage generates the system message with validation rules
    GenerateSystemMessage(validator conventional.MessageValidator) (string, error)

    // GenerateUserMessage generates the user message with repository state
    GenerateUserMessage(repoState *model.RepositoryState) (string, error)
}
```

**Implementation** (`UnifiedPromptGenerator`):
- **Fields**:
  - None (stateless)
- **Methods**:
  - `GenerateSystemMessage(validator MessageValidator) (string, error)`: Extracts validation rules from validator and formats as structured bullet points
  - `GenerateUserMessage(repoState *RepositoryState) (string, error)`: Formats repository state (staged/unstaged files, diffs) as user message
- **Dependencies**:
  - `conventional.MessageValidator` (injected via method parameter)
  - `model.RepositoryState` (injected via method parameter)

**State**: Stateless - all methods are pure functions

**Validation Rules**: None (generator doesn't validate, it formats)

---

### MessageValidator (Extended)

Existing validator extended with methods to extract validation rules programmatically.

**Interface** (`pkg/conventional/validator.go` - extended):
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

**Implementation** (`Validator` struct - extended):
- **New Methods**:
  - `GetValidTypes() []string`: Returns `["feat", "fix", "docs", "style", "refactor", "test", "chore", "version"]`
  - `GetSubjectMaxLength() int`: Returns `72`
  - `GetBodyMaxLength() int`: Returns `320`
  - `GetScopeFormatDescription() string`: Returns `"alphanumeric, hyphens, underscores only"`

**State**: Stateless - all methods are pure functions

**Validation Rules**: Methods return constants that match validation logic

---

### PromptParts

Represents the structured parts of a unified prompt (not a domain entity, but a conceptual model).

**Structure**:
- **SystemMessage** (string): Contains validation rules and instructions
  - Format: Structured bullet points
  - Content: Conventional Commits specification, validation constraints
- **UserMessage** (string): Contains repository state and diffs
  - Format: "Generate a commit message for the following changes:\n\n[file list and diffs]"

**Usage**: Providers adapt these parts to their API structure:
- OpenAI/Mistral/Local: Use as separate system and user messages
- Anthropic: Combine as single user message (system prepended)

---

## Relationships

```
MessageValidator (extended)
    ↑ (used by)
PromptGenerator
    ↑ (used by)
AIProvider implementations
    ├── OpenAIProvider
    ├── AnthropicProvider (adapts: prepends system to user)
    ├── MistralProvider
    └── LocalProvider
```

**Dependency Flow**:
- `pkg/ai/prompt/` depends on `pkg/conventional/` and `internal/model/`
- `internal/ai/*_provider.go` depend on `pkg/ai/prompt/`
- No circular dependencies

---

## State Transitions

**PromptGenerator**: Stateless - no state transitions

**MessageValidator**: Stateless - no state transitions

**Provider Usage Flow**:
1. Provider receives `RepositoryState`
2. Provider calls `PromptGenerator.GenerateSystemMessage(validator)`
3. Provider calls `PromptGenerator.GenerateUserMessage(repoState)`
4. Provider adapts messages to API structure
5. Provider sends to AI API

---

## Data Flow

```
RepositoryState
    ↓
PromptGenerator.GenerateUserMessage()
    ↓
User Message (string)
    ↓
Provider (adapts to API)
    ↓
AI API

MessageValidator
    ↓
PromptGenerator.GenerateSystemMessage()
    ↓
System Message (string)
    ↓
Provider (adapts to API)
    ↓
AI API
```

---

## Validation Rules (for Prompt Content)

**System Message Validation**:
- Must include all validation rules from MessageValidator
- Must be formatted as structured bullet points
- Must include Conventional Commits format specification

**User Message Validation**:
- Must include repository state information
- Must format file changes clearly
- Must be consistent across all providers

---

## Edge Cases

1. **Empty RepositoryState**: User message should handle empty state gracefully
2. **Validator Returns Empty Rules**: System message should have fallback content
3. **Provider API Changes**: Adaptation layer handles API-specific formatting
4. **Concurrent Access**: PromptGenerator is stateless, thread-safe by design

---

## Migration Notes

- **Backward Compatibility**: AIProvider interface unchanged
- **Breaking Changes**: None
- **Migration Path**:
  1. Add new methods to MessageValidator
  2. Create PromptGenerator
  3. Update providers one by one
  4. Remove old prompt building code
