# AI Provider Contract: Mistral Provider

**Feature**: 007-mistral-provider
**Date**: 2025-01-27

## Interface Contract

### AIProvider Interface (Existing)

**Signature**: `type AIProvider interface { GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) }`

**Implementation**: `MistralProvider` implements this interface

## MistralProvider Contract

### NewMistralProvider

**Signature**: `func NewMistralProvider(config *model.AIProviderConfig) AIProvider`

**Preconditions**:
- `config` is not nil
- `config.APIKey` may be empty (will be checked in GenerateCommitMessage)
- `config.Model` defaults to "mistral-large-latest" if empty
- `config.Endpoint` defaults to "https://api.mistral.ai/v1/chat/completions" if empty
- `config.Timeout` defaults to 30 seconds if zero

**Postconditions**:
- Returns a `MistralProvider` instance that implements `AIProvider`
- HTTP client is initialized with configured timeout
- Provider is ready to generate commit messages

**Behavior**:
1. Create HTTP client with timeout from config (default 30s)
2. Store config for later use
3. Return provider instance

**Error Cases**:
- None (constructor always succeeds, API key validation happens in GenerateCommitMessage)

### GenerateCommitMessage

**Signature**: `func (p *MistralProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)`

**Preconditions**:
- `ctx` is not nil (may be cancelled for timeout)
- `repoState` is not nil
- MistralProvider is properly initialized
- API key is configured (checked at runtime)

**Postconditions**:
- If successful: Returns generated commit message string
- If error: Returns error wrapped with `utils.ErrAIProviderUnavailable`
- Context cancellation is respected (request aborted if context cancelled)

**Behavior**:
1. Validate API key is present
2. Build prompt from repository state (same pattern as OpenAI/Anthropic)
3. Create HTTP POST request to Mistral API endpoint
4. Set Authorization header with Bearer token
5. Set Content-Type header to application/json
6. Execute request with context (respects cancellation/timeout)
7. Parse JSON response
8. Extract message content from response
9. Return generated message string

**Request Format**:
```json
{
  "model": "mistral-large-latest",
  "messages": [
    {
      "role": "system",
      "content": "You are a git commit message generator. Generate commit messages following the Conventional Commits specification..."
    },
    {
      "role": "user",
      "content": "Generate a commit message for the following changes:\n\n..."
    }
  ],
  "max_tokens": 500
}
```

**Response Format**:
```json
{
  "choices": [
    {
      "message": {
        "content": "feat(api): add new endpoint\n\n..."
      }
    }
  ]
}
```

**Error Cases**:
- Missing API key → Returns error: "Mistral API key not configured"
- HTTP request failure → Returns wrapped error with context
- HTTP status != 200 → Returns error with status code and response body
- JSON parse failure → Returns error: "failed to decode response"
- Empty response → Returns error: "no response from API"
- Context cancellation → Returns context error
- Timeout → Returns context deadline exceeded error

**Side Effects**:
- Makes HTTP request to external Mistral API
- No state changes, no file I/O, no logging of sensitive data (API keys)

## Integration Points

### Provider Registration

**Location**: `internal/service/commit_service.go` → `generateWithAI` method

**Contract**:
- Add "mistral" case to provider switch statement
- Call `ai.NewMistralProvider(providerConfig)` when providerName == "mistral"
- Return error if provider not found (existing behavior)

**Example**:
```go
case "mistral":
    aiProvider = ai.NewMistralProvider(providerConfig)
```

### Token Calculation

**Location**: `pkg/tokenization/token_calculator.go` → `NewTokenCalculator` function

**Contract**:
- Add "mistral" case to switch statement
- Return `NewFallbackTokenCalculator()` for Mistral (character-based estimation)
- Maintains existing interface contract

**Example**:
```go
case "mistral":
    return NewFallbackTokenCalculator()
```

### Configuration

**Location**: `configs/config.yaml.example`

**Contract**:
- Add mistral provider configuration example
- Follow same format as openai/anthropic providers
- Include api_key, model, timeout fields

**Example**:
```yaml
mistral:
  api_key: ${MISTRAL_API_KEY}
  model: mistral-large-latest
  timeout: 30s
```

## Error Handling Contract

### Error Types

- **Missing API Key**: `fmt.Errorf("%w: Mistral API key not configured", utils.ErrAIProviderUnavailable)`
- **HTTP Errors**: `fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)` (wrapped)
- **API Response Errors**: `fmt.Errorf("%w: API returned status %d: %s", utils.ErrAIProviderUnavailable, statusCode, body)`
- **Parse Errors**: `fmt.Errorf("failed to decode response: %w", err)`
- **Empty Response**: `fmt.Errorf("%w: no response from API", utils.ErrAIProviderUnavailable)`

### Error Propagation

- All errors are wrapped with `utils.ErrAIProviderUnavailable` for consistent handling
- `CommitService` catches all provider errors and falls back to manual input
- Error messages are user-friendly and do not expose API keys or sensitive data

## Performance Contract

### Timeout Handling

- Default timeout: 30 seconds (configurable)
- Context cancellation respected
- HTTP client timeout matches config timeout
- No retries on failure (single attempt per request)

### Resource Management

- HTTP client created per provider instance (reused for multiple calls)
- Response body always closed (defer pattern)
- No connection pooling requirements (simple HTTP client sufficient)
