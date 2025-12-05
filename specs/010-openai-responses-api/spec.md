# Feature Specification: Migrate OpenAI Provider to Responses API

**Feature Branch**: `010-openai-responses-api`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "modify the openai provider to use responses api instead of chat completions api (as given here https://platform.openai.com/docs/guides/migrate-to-responses )."

## User Scenarios & Testing

### User Story 1 - Generate Commit Messages Using Responses API (Priority: P1)

Users continue to generate commit messages using the OpenAI provider, but the system now uses the Responses API instead of Chat Completions API. The user experience remains identical - users configure their OpenAI API key and model, and the system generates commit messages following the Conventional Commits specification.

**Why this priority**: This is the core functionality migration. Users must be able to generate commit messages using the new API without any changes to their workflow or configuration.

**Independent Test**: Configure gitcomm with OpenAI API credentials, run gitcomm with OpenAI provider, and verify that commit messages are generated successfully using the Responses API. The behavior should be identical to the current Chat Completions implementation, with no breaking changes to the user experience.

**Acceptance Scenarios**:

1. **Given** a user has configured OpenAI API credentials, **When** they run gitcomm with OpenAI provider, **Then** the system successfully generates a commit message using the Responses API
2. **Given** a user has staged changes in their repository, **When** they request AI-generated commit messages, **Then** the system uses the Responses API to generate a message following Conventional Commits format
3. **Given** a user has configured a custom model name, **When** they run gitcomm, **Then** the system uses the specified model with the Responses API
4. **Given** the Responses API returns a successful response, **When** the system processes it, **Then** the commit message is extracted and presented to the user correctly

---

### Edge Cases

- What happens when the Responses API returns an error (authentication, rate limit, timeout)?
- How does the system handle empty responses from the Responses API?
- What happens when the Responses API response structure differs from Chat Completions?
- How does the system handle conversation state if the Responses API manages it automatically? → System disables state management and uses stateless mode
- What happens when the model specified by the user is not available in the Responses API?
- How does the system handle multimodal inputs if the Responses API supports them but we only use text?

## Requirements

### Functional Requirements

- **FR-001**: System MUST use the Responses API endpoint (`/v1/responses`) instead of Chat Completions API (`/v1/chat/completions`)
- **FR-002**: System MUST convert the existing `messages` array structure to the `input` parameter format required by Responses API (input is an array of message objects with `role` and `content` fields, same structure as current `messages`)
- **FR-003**: System MUST extract commit message content from Responses API response structure correctly (response contains a `content` or `text` field with the message, similar to Chat Completions structure)
- **FR-004**: System MUST maintain identical error handling behavior (same error types, same user-facing messages) as the current Chat Completions implementation
- **FR-005**: System MUST preserve all existing configuration options (API key, model, timeout, max tokens)
- **FR-006**: System MUST maintain backward compatibility - existing user configurations must continue to work without modification
- **FR-007**: System MUST handle context cancellation and timeouts correctly with Responses API
- **FR-008**: System MUST map Responses API errors to existing error types (authentication, rate limit, timeout, generic unavailable)
- **FR-009**: System MUST support the same models that are currently supported (with Responses API compatibility)
- **FR-010**: System MUST generate commit messages following Conventional Commits specification using Responses API
- **FR-011**: System MUST use Responses API in stateless mode (disable conversation state management) - each commit message generation is independent with no state persistence

### Key Entities

- **OpenAIProvider**: The AI provider implementation that generates commit messages using OpenAI's Responses API
- **AIProviderConfig**: Configuration structure containing API key, model, timeout, and max tokens (unchanged)
- **RepositoryState**: Repository state containing staged and unstaged files (unchanged)
- **Responses API Request**: Request structure using `input` parameter instead of `messages` array
- **Responses API Response**: Response structure from Responses API containing generated commit message

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can successfully generate commit messages using OpenAI provider with Responses API - 100% of successful API calls result in valid commit message extraction
- **SC-002**: Error handling behavior remains identical - all error scenarios (authentication, rate limit, timeout) produce the same user-facing error messages as the current implementation
- **SC-003**: Configuration compatibility - 100% of existing OpenAI provider configurations continue to work without modification
- **SC-004**: Response time - Commit message generation completes within the same time bounds as the current Chat Completions implementation (no degradation)
- **SC-005**: Functional parity - All existing OpenAI provider unit tests pass with Responses API implementation
- **SC-006**: Integration compatibility - All existing integration tests pass with Responses API implementation

## Clarifications

### Session 2025-01-27

- Q: If SDK v3 doesn't support Responses API yet, what should we do? → A: Proceed with migration assuming SDK support will be available or can be added via custom HTTP client if needed
- Q: What is the exact structure of the `input` parameter for Responses API? → A: `input` is an array of message objects with `role` and `content` fields (same structure as current `messages` array)
- Q: How is the commit message content structured in the Responses API response? → A: Response contains a `content` or `text` field with the message (similar structure to Chat Completions `choices[0].message.content`)
- Q: How should we handle conversation state in Responses API for our stateless use case? → A: Disable conversation state management (use stateless mode) - each commit message generation is independent, no state persistence
- Q: What should be the default model if Responses API doesn't support the current default? → A: Assume all current models work with Responses API; keep current default (`gpt-4-1`)

## Assumptions

- The Responses API will be available in OpenAI SDK v3, or can be accessed via custom HTTP client implementation if SDK support is not yet available
- The Responses API supports the same models that are currently supported via Chat Completions (including the current default model `gpt-4-1`)
- The Responses API response structure contains the commit message content in a `content` or `text` field (similar to Chat Completions `choices[0].message.content` structure)
- Error codes and error message formats from Responses API are similar enough to map to existing error types
- The `input` parameter in Responses API accepts an array of message objects with `role` and `content` fields (same structure as current `messages` array)
- Conversation state management in Responses API will be disabled - system uses stateless mode where each commit message generation is independent with no state persistence
- The Responses API supports the same timeout and context cancellation mechanisms as Chat Completions

## Dependencies

- OpenAI SDK v3 should support Responses API endpoints and request/response structures, or custom HTTP client implementation will be used if SDK support is not available
- Existing AIProvider interface remains unchanged (no breaking changes to other providers)
- Existing configuration structure (AIProviderConfig) remains unchanged
- Existing error handling utilities remain unchanged

## Out of Scope

- Migration of other AI providers (Anthropic, Mistral, local) to Responses API
- Utilization of Responses API-specific features like conversation state management, built-in tools, or multimodal inputs
- Changes to the user-facing CLI interface or configuration format
- Performance optimizations specific to Responses API (unless they naturally result from the migration)
- Support for Responses API-specific features like function calling or tool integration
