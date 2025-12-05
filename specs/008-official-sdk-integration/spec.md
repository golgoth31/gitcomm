# Feature Specification: Use Official SDKs for AI Providers

**Feature Branch**: `008-official-sdk-integration`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "use the official sdk for each AI provider to interact with provider API"

## Clarifications

### Session 2025-01-27

- Q: When an official SDK fails to initialize (e.g., missing dependency, incompatible version, initialization error), what should the system do? → A: Fail fast with clear error message, fallback to manual input
- Q: How should SDK-specific error types be mapped to the existing error handling system? → A: Map SDK errors to existing error types, preserve user-facing messages
- Q: Should the system use SDK-provided automatic retries if the SDK offers them, or maintain single-attempt behavior? → A: Use SDK automatic retries if available (improves reliability)
- Q: If an official SDK requires configuration options not present in the current AIProviderConfig structure, how should the system handle this? → A: Extend AIProviderConfig to support SDK-specific options while maintaining backward compatibility
- Q: When an SDK version incompatibility is detected (e.g., during build or initialization), what should the system do? → A: Fail at build/initialization with clear error message indicating the incompatibility

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Replace OpenAI HTTP Client with Official SDK (Priority: P1)

As a developer using gitcomm with OpenAI, I want the CLI to use the official OpenAI Go SDK instead of raw HTTP calls, so that I benefit from official SDK features like automatic retries, better error handling, and API compatibility guarantees.

**Why this priority**: OpenAI is the default provider and most commonly used. Replacing it first provides immediate value and establishes the pattern for other providers.

**Independent Test**: Configure OpenAI API credentials, run gitcomm with OpenAI provider, and verify that commit messages are generated successfully using the official SDK. The behavior should be identical to the current HTTP client implementation.

**Acceptance Scenarios**:

1. **Given** OpenAI is configured as the default provider, **When** I run gitcomm with staged changes, **Then** the CLI successfully generates a commit message using the official OpenAI SDK
2. **Given** OpenAI API key is invalid, **When** I run gitcomm, **Then** the CLI displays an appropriate error message and falls back to manual input
3. **Given** OpenAI API returns a timeout error, **When** I run gitcomm, **Then** the CLI handles the timeout gracefully and falls back to manual input
4. **Given** OpenAI API is rate-limited, **When** I run gitcomm, **Then** the CLI displays an appropriate error message and falls back to manual input

---

### User Story 2 - Replace Anthropic HTTP Client with Official SDK (Priority: P2)

As a developer using gitcomm with Anthropic, I want the CLI to use the official Anthropic Go SDK instead of raw HTTP calls, so that I benefit from official SDK features and maintain compatibility with Anthropic API updates.

**Why this priority**: Anthropic is a commonly used provider. Replacing it maintains feature parity with OpenAI and follows the established pattern.

**Independent Test**: Configure Anthropic API credentials, run gitcomm with Anthropic provider, and verify that commit messages are generated successfully using the official Anthropic SDK. The behavior should be identical to the current HTTP client implementation.

**Acceptance Scenarios**:

1. **Given** Anthropic is configured as the provider, **When** I run gitcomm with staged changes, **Then** the CLI successfully generates a commit message using the official Anthropic SDK
2. **Given** Anthropic API key is invalid, **When** I run gitcomm, **Then** the CLI displays an appropriate error message and falls back to manual input
3. **Given** Anthropic API returns an error, **When** I run gitcomm, **Then** the CLI handles the error gracefully and falls back to manual input

---

### User Story 3 - Replace Mistral HTTP Client with Official SDK (Priority: P2)

As a developer using gitcomm with Mistral, I want the CLI to use the official Mistral Go SDK instead of raw HTTP calls, so that I benefit from official SDK features and maintain compatibility with Mistral API updates.

**Why this priority**: Mistral was recently added as a provider. Replacing it maintains consistency across all providers and ensures all use official SDKs.

**Independent Test**: Configure Mistral API credentials, run gitcomm with Mistral provider, and verify that commit messages are generated successfully using the official Mistral SDK. The behavior should be identical to the current HTTP client implementation.

**Acceptance Scenarios**:

1. **Given** Mistral is configured as the provider, **When** I run gitcomm with staged changes, **Then** the CLI successfully generates a commit message using the official Mistral SDK
2. **Given** Mistral API key is invalid, **When** I run gitcomm, **Then** the CLI displays an appropriate error message and falls back to manual input
3. **Given** Mistral API returns an error, **When** I run gitcomm, **Then** the CLI handles the error gracefully and falls back to manual input

---

### Edge Cases

- What happens when an official SDK is not available or fails to initialize? → System fails fast with clear error message and falls back to manual input (no HTTP client fallback)
- How does the system handle SDK version incompatibilities? → System fails at build/initialization with clear error message indicating the incompatibility, then falls back to manual input
- What happens if the official SDK has different error types than the current HTTP client implementation? → SDK errors are mapped to existing error types, preserving user-facing messages (no SDK-specific error types exposed)
- How does the system handle SDK-specific features that don't exist in the current implementation (e.g., streaming, retries)? → SDK-provided automatic retries are used if available (improves reliability), streaming is not used (out of scope)
- What happens if the official SDK requires different configuration options than the current implementation? → AIProviderConfig structure is extended with optional SDK-specific fields if needed, maintaining backward compatibility (existing configs continue to work)
- How does the system maintain backward compatibility with existing configuration files?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST use the official OpenAI Go SDK (github.com/openai/openai-go) for OpenAI provider interactions
- **FR-002**: System MUST use the official Anthropic Go SDK (github.com/anthropics/anthropic-sdk-go) for Anthropic provider interactions
- **FR-003**: System MUST use the official Mistral Go SDK (github.com/Gage-Technologies/mistral-go) for Mistral provider interactions
- **FR-004**: System MUST maintain the existing AIProvider interface contract (no breaking changes to public API)
- **FR-005**: System MUST preserve all existing functionality (commit message generation, error handling, fallback behavior) - SDK-provided automatic retries SHOULD be used if available to improve reliability without changing user-facing behavior
- **FR-006**: System MUST maintain backward compatibility with existing configuration files (no changes required to user configs) - if SDKs require additional configuration options, AIProviderConfig structure MAY be extended with optional fields to support SDK-specific options, but existing configs without these fields MUST continue to work
- **FR-007**: System MUST handle SDK initialization errors gracefully with appropriate error messages - if SDK initialization fails (including version incompatibilities), the system MUST fail at build/initialization with a clear error message and fall back to manual input (no attempt to use HTTP client as fallback)
- **FR-008**: System MUST preserve existing error handling behavior (same error types and messages for user-facing errors) - SDK-specific error types MUST be mapped to existing error types, preserving user-facing error messages
- **FR-009**: System MUST respect existing timeout configurations from user config
- **FR-010**: System MUST maintain the same prompt building logic (no changes to how prompts are constructed from repository state)
- **FR-011**: System MUST maintain the same response parsing logic (extract commit message from SDK response in the same way)
- **FR-012**: System MUST not expose SDK-specific implementation details in error messages or logs (maintain abstraction)

### Key Entities *(include if feature involves data)*

- **AIProvider Interface**: Existing interface that all providers must implement - remains unchanged
- **AIProviderConfig**: Existing configuration structure - may be extended with optional SDK-specific fields if needed, but must maintain backward compatibility (existing configs without new fields must continue to work)
- **Provider Implementations**: OpenAIProvider, AnthropicProvider, MistralProvider - internal implementation changes to use SDKs instead of HTTP clients

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All three providers (OpenAI, Anthropic, Mistral) successfully generate commit messages using their respective official SDKs with 100% functional parity to current HTTP client implementation
- **SC-002**: Existing configuration files continue to work without modification (100% backward compatibility)
- **SC-003**: Error handling behavior remains identical to current implementation (same error messages, same fallback behavior)
- **SC-004**: All existing unit tests and integration tests pass without modification (test behavior remains the same)
- **SC-005**: Commit message generation success rate remains the same or improves (no regression in functionality)
- **SC-006**: User experience remains identical (no changes to CLI behavior, prompts, or output format)
- **SC-007**: SDK initialization and API calls complete within the same time constraints as current HTTP client implementation (no performance regression)

## Assumptions

- Official SDKs are available and maintained for all three providers
- Official SDKs support the same API features needed for commit message generation (chat completions)
- Official SDKs can be configured with the same parameters (API key, model, timeout, max tokens) as current implementation
- Official SDKs support context cancellation for timeout handling
- Official SDKs provide error types that can be mapped to existing error handling patterns
- Official SDKs do not require breaking changes to existing configuration structure
- Official SDKs maintain compatibility with the same API endpoints and request/response formats
- Local provider implementation remains unchanged (no SDK replacement needed)

## Dependencies

- **External Dependencies**:
  - github.com/openai/openai-go (OpenAI official SDK)
  - github.com/anthropics/anthropic-sdk-go (Anthropic official SDK)
  - github.com/Gage-Technologies/mistral-go (Mistral official SDK)
- **Internal Dependencies**:
  - Existing AIProvider interface (must remain unchanged)
  - Existing AIProviderConfig structure (must remain unchanged)
  - Existing provider selection mechanism in CommitService (must remain unchanged)
- **No Breaking Changes**: This feature must not break any existing functionality, tests, or user configurations

## Out of Scope

- Adding new features or capabilities beyond replacing HTTP clients with SDKs
- Changing the AIProvider interface or public API
- Modifying configuration file structure or format (note: optional SDK-specific fields may be added to AIProviderConfig while maintaining backward compatibility)
- Adding support for SDK-specific features (e.g., streaming) - automatic retries provided by SDKs are used if available to improve reliability
- Changing error handling behavior or error message formats
- Modifying prompt building or response parsing logic beyond what's necessary for SDK integration
- Replacing the local provider implementation (local provider may not have an official SDK)
