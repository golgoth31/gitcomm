# Research: Migrate OpenAI Provider to Responses API

**Feature**: 010-openai-responses-api
**Date**: 2025-01-27
**Purpose**: Document technology decisions and rationale for migrating from Chat Completions API to Responses API

## API Endpoint

**Decision**: Use `/v1/responses` endpoint instead of `/v1/chat/completions`

**Rationale**:
- Responses API is the new unified API that combines capabilities from Chat Completions and Assistants APIs
- Provides better state management, multimodal support, and tool integration
- Migration guide recommends this endpoint for new implementations

**Alternatives Considered**:
- Continue using `/v1/chat/completions` - Rejected because spec requires migration to Responses API
- Use both APIs conditionally - Rejected because spec requires full migration

## SDK Support

**Decision**: Proceed with migration assuming SDK v3 support will be available, or implement custom HTTP client if SDK doesn't support Responses API yet

**Rationale**:
- OpenAI SDK v3 is actively maintained and likely to add Responses API support
- Custom HTTP client can be implemented using standard Go `net/http` package if needed
- Clarification from spec indicates this approach is acceptable

**Implementation Strategy**:
1. First attempt: Check if SDK v3 has Responses API support
2. If available: Use SDK client methods (e.g., `client.Responses.New()`)
3. If not available: Implement custom HTTP client with proper request/response handling
4. Maintain same interface and error handling regardless of implementation method

**Alternatives Considered**:
- Wait for SDK support - Rejected because spec allows proceeding with custom HTTP client
- Use only custom HTTP client - Rejected because SDK is preferred if available

## Request Structure

**Decision**: Convert `messages` array to `input` parameter (array of message objects with `role` and `content` fields)

**Rationale**:
- Responses API uses `input` parameter instead of `messages`
- Structure is identical: array of objects with `role` and `content` fields
- Minimal conversion needed - can directly map current `messages` array to `input`

**Request Format**:
```json
{
  "model": "gpt-4-1",
  "input": [
    {
      "role": "system",
      "content": "You are a git commit message generator..."
    },
    {
      "role": "user",
      "content": "Generate a commit message for..."
    }
  ],
  "max_completion_tokens": 500
}
```

**Alternatives Considered**:
- Single string input - Rejected because spec clarifies array format
- Different structure - Rejected because migration guide shows same message object structure

## Response Structure

**Decision**: Extract commit message from `content` or `text` field in response (similar to Chat Completions `choices[0].message.content`)

**Rationale**:
- Responses API response structure is similar to Chat Completions
- Content is accessible via `content` or `text` field
- Extraction logic can be adapted with minimal changes

**Response Format** (expected):
```json
{
  "content": "feat(api): add new endpoint\n\n...",
  // or
  "text": "feat(api): add new endpoint\n\n...",
  // or nested structure similar to Chat Completions
}
```

**Implementation Notes**:
- Need to verify exact response structure during implementation
- May need to handle different response formats (content vs text vs nested)
- Maintain same extraction logic pattern as current implementation

**Alternatives Considered**:
- Completely different structure - Rejected because spec indicates similar structure
- Nested differently - Will handle during implementation if needed

## Error Handling

**Decision**: Map Responses API errors to existing error types using same `mapSDKError` function pattern

**Rationale**:
- Maintain identical error handling behavior as current implementation
- Same error types: authentication (401), rate limit (429), timeout, generic unavailable
- User-facing error messages must remain unchanged

**Error Mapping Strategy**:
- Check error strings for common patterns (authentication, rate limit, timeout)
- Map to existing `utils.ErrAIProviderUnavailable` with appropriate context
- Preserve original SDK error message in wrapped error for debugging

**Alternatives Considered**:
- New error types - Rejected because spec requires identical error handling
- Different error messages - Rejected because spec requires same user-facing messages

## Stateless Mode

**Decision**: Disable conversation state management - use stateless mode for each commit message generation

**Rationale**:
- Current implementation is stateless (each commit is independent)
- Responses API supports stateless mode
- No need for conversation history or state persistence

**Configuration**:
- Set `store: false` or equivalent parameter to disable state management
- Each API call is independent with no state persistence
- No `previous_response_id` or conversation tracking needed

**Alternatives Considered**:
- Use state management - Rejected because current implementation is stateless
- Allow optional state - Rejected because spec requires stateless mode

## Model Compatibility

**Decision**: Assume all current models work with Responses API; keep current default (`gpt-4-1`)

**Rationale**:
- Spec clarification indicates assuming model compatibility
- Current default model (`gpt-4-1`) should work with Responses API
- If model incompatibility discovered, handle gracefully with error message

**Implementation Notes**:
- Use same model names as current implementation
- Default to `gpt-4-1` if no model specified
- If model not supported, error will be returned and mapped appropriately

**Alternatives Considered**:
- Change default model - Rejected because spec requires keeping current default
- Verify all models first - Rejected because spec allows assuming compatibility

## Timeout and Context Handling

**Decision**: Use same timeout and context cancellation mechanisms as Chat Completions

**Rationale**:
- Responses API supports standard HTTP timeout and context cancellation
- Current implementation uses `context.Context` for cancellation/timeout
- No changes needed to context handling logic

**Implementation**:
- Pass `context.Context` to API calls (SDK or HTTP client)
- Respect context cancellation and deadline
- Same timeout configuration from `AIProviderConfig`

**Alternatives Considered**:
- Different timeout mechanism - Rejected because spec requires same behavior
- Remove context support - Rejected because constitution requires context propagation

## Token Calculation

**Decision**: Use same token calculation approach (if applicable) or remove if Responses API handles it differently

**Rationale**:
- Current implementation may have token calculation logic
- Responses API may handle tokens differently
- Need to verify during implementation

**Implementation Notes**:
- Check if token calculation is needed for Responses API
- If different, adapt or remove as appropriate
- Maintain same user-facing behavior

## Testing Strategy

**Decision**: Update existing tests to work with Responses API while maintaining same test structure

**Rationale**:
- Existing tests in `openai_provider_test.go` must pass with new implementation
- Test structure remains the same (unit tests, integration tests)
- Mock Responses API responses instead of Chat Completions responses

**Test Updates Needed**:
- Update test mocks to use Responses API response structure
- Verify error handling tests still work
- Ensure integration tests pass with real API (if available)

**Alternatives Considered**:
- New test suite - Rejected because spec requires existing tests to pass
- Remove tests - Rejected because constitution requires TDD
