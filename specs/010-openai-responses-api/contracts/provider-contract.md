# Provider Contract: OpenAI Responses API

**Feature**: 010-openai-responses-api
**Date**: 2025-01-27

## Interface Contract (Unchanged)

### AIProvider Interface

**Location**: `internal/ai/provider.go`

```go
type AIProvider interface {
    GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)
}
```

**Contract**: Unchanged - OpenAI provider continues to implement this interface.

## OpenAIProvider Contract

### Constructor

**Function**: `NewOpenAIProvider(config *model.AIProviderConfig) AIProvider`

**Pre-conditions**:
- `config` is not nil
- `config.Name` is "openai" (implicitly)

**Post-conditions**:
- Returns non-nil `AIProvider` implementation
- Provider is ready to generate commit messages
- SDK client (or HTTP client) is initialized with API key

**Behavior**:
- Initializes OpenAI SDK v3 client (or custom HTTP client if SDK doesn't support Responses API)
- Uses `config.APIKey` for authentication
- Logs debug message if API key is empty (but doesn't fail)

### GenerateCommitMessage

**Function**: `GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)`

**Pre-conditions**:
- `ctx` is not nil
- `repoState` is not nil (may have empty file lists)
- `p.config.APIKey` is non-empty (returns error if empty)

**Post-conditions**:
- On success: Returns non-empty commit message string following Conventional Commits format
- On error: Returns empty string and error wrapped with `utils.ErrAIProviderUnavailable`

**Behavior**:
1. Validates API key is configured
2. Builds prompt from `repoState` (unchanged logic)
3. Prepares model (defaults to "gpt-4-1" if empty)
4. Prepares max tokens (defaults to 500 if 0)
5. Creates Responses API request with:
   - Model name
   - Input array (converted from messages array format)
   - Max completion tokens
   - Store: false (stateless mode)
6. Executes API call with context (respects cancellation/timeout)
7. Extracts commit message content from response
8. Returns content or error

**Error Handling**:
- API key missing: Returns `ErrAIProviderUnavailable` with "OpenAI API key not configured"
- API errors: Mapped via `mapSDKError` to `ErrAIProviderUnavailable` with appropriate context
- Empty response: Returns `ErrAIProviderUnavailable` with "empty response from API"
- Context cancellation: Returns context error (wrapped)

### mapSDKError (Internal)

**Function**: `mapSDKError(err error) error`

**Pre-conditions**:
- `err` is not nil

**Post-conditions**:
- Returns error wrapped with `utils.ErrAIProviderUnavailable`
- Original error message preserved for debugging

**Behavior**:
- Maps authentication errors (401, "authentication", "invalid") → "API key invalid"
- Maps rate limit errors (429, "rate limit") → "rate limit exceeded"
- Maps timeout errors ("timeout", "deadline") → "timeout"
- Maps all other errors → generic "AI provider unavailable" with original error message

## API Request Contract

### Request Structure

**Endpoint**: `POST /v1/responses`

**Headers**:
- `Authorization: Bearer {APIKey}`
- `Content-Type: application/json`

**Body**:
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
  "max_completion_tokens": 500,
  "store": false
}
```

**Contract**:
- `model`: String, required, defaults to "gpt-4-1"
- `input`: Array of message objects, required, at least 2 messages (system + user)
- `max_completion_tokens`: Integer, optional, defaults to 500
- `store`: Boolean, required, must be false for stateless mode

## API Response Contract

### Response Structure (Expected)

**Success Response** (200 OK):
```json
{
  "content": "feat(api): add new endpoint\n\n...",
  // or
  "text": "feat(api): add new endpoint\n\n...",
  // or nested structure similar to Chat Completions
}
```

**Error Responses**:
- 401 Unauthorized: Invalid API key
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Server error
- Timeout: Context deadline exceeded

**Contract**:
- Success: Response contains non-empty `content` or `text` field
- Errors: Standard HTTP status codes with error messages

## Test Contracts

### Unit Test Requirements

**Test File**: `internal/ai/openai_provider_test.go`

**Required Tests**:
1. `TestNewOpenAIProvider_ResponsesAPIClientInitialization` - Verify client initialization
2. `TestOpenAIProvider_GenerateCommitMessage_ResponsesAPISuccess` - Verify successful API call
3. `TestOpenAIProvider_GenerateCommitMessage_ResponsesAPIErrorMapping` - Verify error mapping
4. `TestOpenAIProvider_ContextCancellation` - Verify context cancellation
5. `TestOpenAIProvider_EmptyResponse` - Verify empty response handling

**Contract**:
- All tests must pass with Responses API implementation
- Tests must mock Responses API responses (not Chat Completions)
- Error handling tests must verify same error types as current implementation

### Integration Test Requirements

**Test File**: `test/integration/ai_commit_test.go`

**Required Tests**:
- Verify end-to-end commit message generation with Responses API
- Verify error handling with real API (if credentials available)

**Contract**:
- Integration tests must pass with Responses API
- Behavior must be identical to current Chat Completions implementation

## Backward Compatibility Contract

**Guarantee**: 100% backward compatibility

**What Must Not Change**:
- `AIProvider` interface signature
- `AIProviderConfig` structure
- Error types and messages
- User-facing behavior
- Configuration format

**What Can Change** (Internal Only):
- API endpoint (`/v1/chat/completions` → `/v1/responses`)
- Request parameter names (`messages` → `input`)
- Response extraction logic (if structure differs)
- SDK client usage (if custom HTTP client needed)

## Notes

- Contract focuses on external behavior, not internal implementation
- Internal request/response structures are implementation details
- All contracts must be verified through tests
- Backward compatibility is non-negotiable
