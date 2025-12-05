# Feature Specification: Unify AI Provider Prompts with Validation Rules

**Feature Branch**: `011-unify-ai-prompts`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "ensure that all AI provider uses the same prompt. Update the prompt to follow the message validation rules as required by the MessageValidator"

## User Scenarios & Testing

### User Story 1 - Unified Prompt Across All AI Providers (Priority: P1)

Users generate commit messages using any AI provider (OpenAI, Anthropic, Mistral, or local), and all providers use the same standardized prompt that includes validation rules. This ensures consistent commit message quality and format regardless of which provider is used, and ensures generated messages pass validation.

**Why this priority**: This is the core functionality - all providers must generate commit messages that are consistent and valid. Inconsistent prompts lead to inconsistent output quality and validation failures.

**Independent Test**: Configure gitcomm with different AI providers, generate commit messages for the same repository state, and verify that all providers generate commit messages that:
1. Follow the same format and structure
2. Pass MessageValidator validation
3. Include all required validation constraints in their generation

**Acceptance Scenarios**:

1. **Given** a user has configured OpenAI provider, **When** they generate a commit message, **Then** the prompt includes all MessageValidator validation rules (type constraints, length limits, scope format)
2. **Given** a user has configured Anthropic provider, **When** they generate a commit message, **Then** the prompt is identical to OpenAI provider's prompt
3. **Given** a user has configured Mistral provider, **When** they generate a commit message, **Then** the prompt is identical to OpenAI provider's prompt
4. **Given** a user has configured local provider, **When** they generate a commit message, **Then** the prompt is identical to OpenAI provider's prompt
5. **Given** any AI provider generates a commit message, **When** the message is validated, **Then** it passes all MessageValidator checks (valid type, subject length ≤72, body length ≤320, valid scope format)

---

### Edge Cases

- What happens when the unified prompt is longer than expected? (Should be handled gracefully)
- How does the system handle prompt updates if validation rules change in the future?
- What happens if a provider-specific constraint conflicts with the unified prompt?
- How does the system ensure prompt consistency when new providers are added?

## Requirements

### Functional Requirements

- **FR-001**: System MUST use the same prompt text across all AI providers (OpenAI, Anthropic, Mistral, local) - both system message (with validation rules) and user message (with diff) must be identical. For Anthropic (which doesn't support system messages), prepend system message content to the user message to maintain unified content.
- **FR-002**: System MUST include MessageValidator validation rules in the prompt (type constraints, length limits, scope format) - rules must be presented as structured bullet points with explicit constraints
- **FR-003**: System MUST ensure the prompt instructs AI to generate commit messages that will pass MessageValidator validation
- **FR-004**: System MUST maintain backward compatibility - existing functionality must continue to work
- **FR-005**: System MUST extract the unified prompt to a shared location to ensure consistency - both system message and user message parts must be extracted. Provider implementations must adapt the unified prompt to their API structure (e.g., Anthropic prepends system content to user message)
- **FR-012**: System MUST dynamically generate the prompt from MessageValidator - validation rules must be extracted programmatically to ensure prompt stays in sync with validator implementation
- **FR-006**: System MUST update all provider implementations to use the unified prompt
- **FR-007**: System MUST include in the prompt: valid commit types (feat, fix, docs, style, refactor, test, chore, version)
- **FR-008**: System MUST include in the prompt: subject length constraint (≤72 characters)
- **FR-009**: System MUST include in the prompt: body length constraint (≤320 characters if provided)
- **FR-010**: System MUST include in the prompt: scope format constraint (alphanumeric, hyphens, underscores only)
- **FR-011**: System MUST ensure the prompt format matches Conventional Commits specification: `type(scope): subject\n\nbody\n\nfooter`

### Key Entities

- **Unified Prompt**: A single prompt template used by all AI providers that includes validation rules
- **MessageValidator**: The validation system that enforces commit message rules (type, length, format constraints)
- **AIProvider**: Interface implemented by OpenAI, Anthropic, Mistral, and local providers
- **CommitMessage**: The structured commit message model with Type, Scope, Subject, Body, Footer fields

## Success Criteria

### Measurable Outcomes

- **SC-001**: Prompt consistency - 100% of AI providers use identical prompt text (verified by code inspection)
- **SC-002**: Validation compliance - 95% of AI-generated commit messages pass MessageValidator validation on first generation
- **SC-003**: Format consistency - All providers generate commit messages in the same format structure
- **SC-004**: Type compliance - 100% of generated commit messages use valid types from the allowed list (feat, fix, docs, style, refactor, test, chore, version)
- **SC-005**: Length compliance - 90% of generated commit messages have subject ≤72 characters and body ≤320 characters
- **SC-006**: Scope compliance - 100% of generated commit messages with scope use valid identifier format (alphanumeric, hyphens, underscores)

## Clarifications

### Session 2025-01-27

- Q: Should the unified prompt include both system and user message parts, or only the user message? → A: Unify both system message (with validation rules) and user message (with diff) - all providers use same structure
- Q: How should Anthropic handle system messages since its API doesn't support separate system messages? → A: Prepend system message content (validation rules) to the user message for Anthropic - unified content, different structure
- Q: What format should validation rules use in the prompt? → A: Structured bullet points with explicit constraints (e.g., "Type must be one of: feat, fix, docs...")
- Q: Should the prompt be dynamically generated from MessageValidator or hardcoded? → A: Dynamically generate prompt from MessageValidator - extract rules programmatically to ensure sync

## Assumptions

- All AI providers can use the same prompt format and structure (both system and user message parts)
- The unified prompt can be extracted to a shared utility that dynamically generates from MessageValidator
- MessageValidator validation rules can be programmatically extracted for prompt generation
- AI models can understand and follow the validation constraints when included in the prompt
- The prompt length increase (due to validation rules) will not significantly impact API costs or response quality

## Dependencies

- MessageValidator implementation and validation rules must be stable
- All AI provider implementations must support the unified prompt structure
- Existing AIProvider interface remains unchanged (no breaking changes)

## Out of Scope

- Changes to MessageValidator validation rules
- Changes to AIProvider interface
- Provider-specific prompt optimizations
- Dynamic prompt generation based on repository context (beyond current diff-based prompts)
- Prompt versioning or migration system
