# Feature Specification: Upgrade OpenAI Provider to SDK v3

**Feature Branch**: `009-openai-sdk-v3`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "modify openai provider to use the sdk v3, github.com/openai/openai-go/v3"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Upgrade OpenAI Provider to SDK v3 (Priority: P1)

As a developer using gitcomm, I want the OpenAI provider to use the latest SDK version (v3) so that I benefit from the latest features, bug fixes, and improvements in the official OpenAI Go SDK.

**Why this priority**: This is a technical upgrade that improves maintainability, reliability, and access to the latest SDK features. It's a single provider upgrade that should be straightforward.

**Independent Test**: Configure gitcomm with OpenAI API credentials, run gitcomm with OpenAI provider, and verify that commit messages are generated successfully using SDK v3. The behavior should be identical to the current SDK v1 implementation, with no breaking changes to the user experience.

**Acceptance Scenarios**:

1. **Given** gitcomm is configured with OpenAI API credentials, **When** a user runs gitcomm with OpenAI provider, **Then** commit messages are generated successfully using SDK v3
2. **Given** gitcomm is configured with OpenAI API credentials, **When** a user runs gitcomm with OpenAI provider, **Then** all existing functionality (error handling, timeout, context cancellation) works identically to SDK v1
3. **Given** gitcomm is configured with invalid OpenAI API credentials, **When** a user runs gitcomm with OpenAI provider, **Then** error messages are user-friendly and maintain the same error handling behavior as SDK v1

### Edge Cases

- What happens when SDK v3 has breaking API changes compared to v1? (Must maintain backward compatibility with existing configuration and behavior)
- How does the system handle SDK v3 initialization failures? (Should fail gracefully with clear error messages)
- What if SDK v3 introduces new error types? (Must map to existing error handling patterns; unmappable errors should be wrapped generically with `ErrAIProviderUnavailable` while preserving original SDK error message for debugging)
- How does the system handle SDK v3 response format changes? (Must extract commit messages correctly)
- What if SDK v3 has different timeout/context handling? (Must respect existing timeout configurations)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST use OpenAI Go SDK v3 (`github.com/openai/openai-go/v3`) for all OpenAI provider interactions
- **FR-002**: System MUST maintain 100% backward compatibility with existing OpenAI provider configuration (no changes to `AIProviderConfig` structure or configuration file format)
- **FR-003**: System MUST maintain identical error handling behavior (same error types, same user-facing error messages)
- **FR-004**: System MUST maintain identical timeout and context cancellation behavior
- **FR-005**: System MUST maintain identical prompt building and response parsing logic
- **FR-006**: System MUST handle SDK v3 API changes gracefully (map new API calls to existing functionality)
- **FR-007**: System MUST map SDK v3 error types to existing error handling patterns. For unmappable errors, wrap generically with `ErrAIProviderUnavailable` while preserving the original SDK error message for debugging
- **FR-008**: System MUST preserve all existing functionality (token calculation, commit message generation, format validation)

### Key Entities

- **OpenAIProvider**: The AI provider implementation that uses OpenAI SDK v3 to generate commit messages
- **AIProviderConfig**: Configuration structure for OpenAI provider (unchanged, maintains backward compatibility)
- **RepositoryState**: Input to commit message generation (unchanged)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: OpenAI provider successfully generates commit messages using SDK v3 with 100% functional parity to SDK v1 implementation
- **SC-002**: All existing unit tests for OpenAI provider pass without modification (or with minimal updates to match SDK v3 API)
- **SC-003**: All existing integration tests for OpenAI provider pass without modification
- **SC-004**: Error handling behavior remains identical (same error types, same user-facing messages, same fallback behavior)
- **SC-005**: Configuration file format remains unchanged (no breaking changes to user configuration)
- **SC-006**: Timeout and context cancellation behavior works identically to SDK v1 implementation
- **SC-007**: SDK v3 upgrade completes without requiring changes to other providers (Anthropic, Mistral, local)

## Assumptions

- SDK v3 maintains similar API structure to v1 (chat completions, message format, error handling)
- SDK v3 is backward compatible with v1 in terms of core functionality (or changes are minimal and well-documented)
- SDK v3 supports the same authentication mechanism (API key via configuration)
- SDK v3 supports context cancellation and timeout handling
- SDK v3 response format is compatible with existing response parsing logic (or changes are minimal)

## Dependencies

- Existing OpenAI provider implementation (SDK v1)
- Existing `AIProvider` interface (unchanged)
- Existing `AIProviderConfig` structure (unchanged)
- Existing error handling infrastructure
- Existing test suite for OpenAI provider

## Constraints

- Must maintain 100% backward compatibility with existing configuration
- Must not break existing functionality or user experience
- Must maintain identical error handling behavior
- Must not require changes to other AI providers
- Must not require changes to the `AIProvider` interface

## Clarifications

### Session 2025-01-27

- Q: If SDK v3 introduces error types that cannot be cleanly mapped to existing error handling patterns, what should the system do? → A: Wrap unmappable errors generically with `ErrAIProviderUnavailable`, preserving the original SDK error message in the wrapped error for debugging
