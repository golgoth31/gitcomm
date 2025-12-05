# Research: Unify AI Provider Prompts with Validation Rules

**Date**: 2025-01-27
**Feature**: 011-unify-ai-prompts

## Research Questions & Findings

### 1. How to Extract Validation Rules from MessageValidator Programmatically?

**Decision**: Extend MessageValidator interface with methods to expose validation rules.

**Rationale**:
- Current `MessageValidator` interface only provides `Validate()` method
- Need to extract: valid types list, subject length limit (72), body length limit (320), scope format rules
- Best practice: Add getter methods to the interface to maintain encapsulation while exposing needed data
- Alternative: Reflection-based extraction was rejected as it breaks encapsulation and is fragile

**Alternatives Considered**:
- **Reflection-based extraction**: Rejected - breaks encapsulation, fragile to refactoring
- **Separate validation rules configuration**: Rejected - violates DRY principle, creates sync issues
- **Extend interface with getter methods**: ✅ Chosen - maintains encapsulation, explicit contract

**Implementation Approach**:
- Add methods to `MessageValidator` interface:
  - `GetValidTypes() []string` - returns list of valid commit types
  - `GetSubjectMaxLength() int` - returns subject length limit (72)
  - `GetBodyMaxLength() int` - returns body length limit (320)
  - `GetScopeFormatDescription() string` - returns scope format description
- Implement these methods in `Validator` struct
- Use these methods in prompt generator

### 2. Where Should the Unified Prompt Generator Be Located?

**Decision**: Place in `pkg/ai/prompt/` as a shared utility package.

**Rationale**:
- Follows Clean Architecture: shared utilities in `pkg/`
- Can be used by all providers in `internal/ai/`
- Maintains dependency direction: `internal/` depends on `pkg/`, not vice versa
- Aligns with existing structure: `pkg/conventional/` for validation, `pkg/tokenization/` for tokenization

**Alternatives Considered**:
- **Place in `internal/ai/prompt/`**: Rejected - would be private to ai package, but prompt generation is a shared concern
- **Place in `pkg/conventional/`**: Rejected - prompt generation is AI-specific, not validation-specific
- **Place in `pkg/ai/prompt/`**: ✅ Chosen - clear separation, shared utility

**Implementation Approach**:
- Create `pkg/ai/prompt/generator.go`
- Define `PromptGenerator` interface
- Implement `UnifiedPromptGenerator` struct
- Inject `MessageValidator` via constructor

### 3. How to Structure the Unified Prompt (System vs User Message)?

**Decision**: Separate system message (validation rules) and user message (diff content).

**Rationale**:
- System message contains static instructions and validation rules
- User message contains dynamic repository state (diffs)
- Separation allows providers to adapt to their API structure (e.g., Anthropic prepends system to user)
- Maintains semantic clarity: system = instructions, user = data

**Alternatives Considered**:
- **Single combined message**: Rejected - loses semantic separation, harder for providers to adapt
- **Separate system and user messages**: ✅ Chosen - clear separation, flexible adaptation

**Implementation Approach**:
- `PromptGenerator` interface provides:
  - `GenerateSystemMessage(validator MessageValidator) (string, error)`
  - `GenerateUserMessage(repoState *model.RepositoryState) (string, error)`
- Providers use both methods and adapt to their API structure

### 4. How to Format Validation Rules in the Prompt?

**Decision**: Structured bullet points with explicit constraints.

**Rationale**:
- AI models respond better to structured, explicit instructions
- Bullet points are easy to parse and understand
- Explicit constraints reduce ambiguity
- Format: "• Type must be one of: feat, fix, docs, style, refactor, test, chore, version"
- Format: "• Subject must be ≤72 characters"
- Format: "• Body must be ≤320 characters (if provided)"
- Format: "• Scope must be a valid identifier (alphanumeric, hyphens, underscores only)"

**Alternatives Considered**:
- **Natural language prose**: Rejected - less structured, harder for AI to parse constraints
- **JSON/YAML format**: Rejected - too verbose, not natural for AI prompts
- **Structured bullet points**: ✅ Chosen - clear, explicit, easy to parse

**Implementation Approach**:
- Build system message as structured text with bullet points
- Include all validation rules from MessageValidator
- Format consistently across all rules

### 5. How to Handle Anthropic's System Message Limitation?

**Decision**: Prepend system message content to user message for Anthropic.

**Rationale**:
- Anthropic API doesn't support separate system messages
- Need to maintain unified content across all providers
- Prepending system content to user message preserves semantic meaning
- Other providers (OpenAI, Mistral, local) can use separate system/user messages

**Alternatives Considered**:
- **Different prompt for Anthropic**: Rejected - violates requirement for unified prompts
- **Embed validation rules in user message only**: Rejected - loses semantic separation, harder to maintain
- **Prepend system to user for Anthropic**: ✅ Chosen - maintains unified content, adapts to API constraints

**Implementation Approach**:
- `AnthropicProvider` calls both `GenerateSystemMessage()` and `GenerateUserMessage()`
- Combines them: `systemMessage + "\n\n" + userMessage`
- Passes combined message as single user message to Anthropic API

### 6. How to Ensure Prompt Consistency When Adding New Providers?

**Decision**: All providers must use the unified prompt generator.

**Rationale**:
- Centralized prompt generation ensures consistency
- New providers automatically get unified prompts
- Changes to validation rules propagate to all providers
- Reduces maintenance burden

**Implementation Approach**:
- Document requirement in code comments
- Add integration tests that verify prompt consistency
- Include in provider implementation guidelines

## Technical Decisions Summary

| Decision | Approach | Rationale |
|----------|----------|-----------|
| Rule Extraction | Extend MessageValidator interface with getter methods | Maintains encapsulation, explicit contract |
| Prompt Generator Location | `pkg/ai/prompt/` | Shared utility, follows Clean Architecture |
| Message Structure | Separate system and user messages | Semantic clarity, flexible adaptation |
| Validation Rules Format | Structured bullet points | Clear, explicit, easy for AI to parse |
| Anthropic Handling | Prepend system to user message | Maintains unified content, adapts to API |
| Consistency | All providers use unified generator | Centralized, reduces maintenance |

## Dependencies & Integration Points

- **MessageValidator** (`pkg/conventional/validator.go`): Must be extended with getter methods
- **AIProvider Interface** (`internal/ai/provider.go`): Remains unchanged (backward compatible)
- **Provider Implementations** (`internal/ai/*_provider.go`): Must be updated to use unified generator
- **RepositoryState** (`internal/model/repository_state.go`): Used for user message generation (unchanged)

## Open Questions Resolved

All research questions have been resolved. No outstanding ambiguities.
