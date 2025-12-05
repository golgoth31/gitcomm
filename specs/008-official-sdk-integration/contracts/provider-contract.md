# Provider Contract: SDK-Based AI Providers

**Feature**: 008-official-sdk-integration
**Date**: 2025-01-27

## Interface Contract

### AIProvider Interface (Existing, Unchanged)

**Signature**: `type AIProvider interface { GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) }`

**Implementation**: All providers (OpenAIProvider, AnthropicProvider, MistralProvider) implement this interface using official SDKs instead of HTTP clients.

## Provider Constructor Contracts

### NewOpenAIProvider

**Signature**: `func NewOpenAIProvider(config *model.AIProviderConfig) AIProvider`

**Preconditions**:
- `config` is not nil
- `config.APIKey` may be empty (will be checked in GenerateCommitMessage)
- `config.Model` defaults to "gpt-4" if empty
- `config.Timeout` defaults to 30 seconds if zero
- `config.MaxTokens` defaults to 500 if zero
- OpenAI SDK (`github.com/openai/openai-go`) is available

**Postconditions**:
- Returns an `OpenAIProvider` instance that implements `AIProvider`
- OpenAI SDK client is initialized with API key from config
- Provider is ready to generate commit messages
- If SDK initialization fails, constructor returns error (fail fast)

**Behavior**:
1. Validate config is not nil
2. Initialize OpenAI SDK client with API key
3. Store config for later use
4. Return provider instance

**Error Cases**:
- SDK initialization failure → Returns error with clear message, provider not created
- Missing SDK dependency → Build error (fail at build time)

### NewAnthropicProvider

**Signature**: `func NewAnthropicProvider(config *model.AIProviderConfig) AIProvider`

**Preconditions**:
- `config` is not nil
- `config.APIKey` may be empty (will be checked in GenerateCommitMessage)
- `config.Model` defaults to "claude-3-opus-20240229" if empty
- `config.Timeout` defaults to 30 seconds if zero
- `config.MaxTokens` defaults to 500 if zero
- Anthropic SDK (`github.com/anthropics/anthropic-sdk-go`) is available

**Postconditions**:
- Returns an `AnthropicProvider` instance that implements `AIProvider`
- Anthropic SDK client is initialized with API key from config
- Provider is ready to generate commit messages
- If SDK initialization fails, constructor returns error (fail fast)

**Behavior**:
1. Validate config is not nil
2. Initialize Anthropic SDK client with API key
3. Store config for later use
4. Return provider instance

**Error Cases**:
- SDK initialization failure → Returns error with clear message, provider not created
- Missing SDK dependency → Build error (fail at build time)

### NewMistralProvider

**Signature**: `func NewMistralProvider(config *model.AIProviderConfig) AIProvider`

**Preconditions**:
- `config` is not nil
- `config.APIKey` may be empty (will be checked in GenerateCommitMessage)
- `config.Model` defaults to "mistral-large-latest" if empty
- `config.Endpoint` defaults to "https://api.mistral.ai/v1/chat/completions" if empty
- `config.Timeout` defaults to 30 seconds if zero
- `config.MaxTokens` defaults to 500 if zero
- Mistral SDK (`github.com/Gage-Technologies/mistral-go`) is available

**Postconditions**:
- Returns a `MistralProvider` instance that implements `AIProvider`
- Mistral SDK client is initialized with API key from config
- Provider is ready to generate commit messages
- If SDK initialization fails, constructor returns error (fail fast)

**Behavior**:
1. Validate config is not nil
2. Initialize Mistral SDK client with API key
3. Store config for later use
4. Return provider instance

**Error Cases**:
- SDK initialization failure → Returns error with clear message, provider not created
- Missing SDK dependency → Build error (fail at build time)

## GenerateCommitMessage Contract

**Signature**: `func (p *Provider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)`

**Preconditions**:
- `ctx` is not nil (may be cancelled for timeout)
- `repoState` is not nil
- Provider is properly initialized (SDK client ready)
- API key is configured (checked at runtime)

**Postconditions**:
- If successful: Returns generated commit message string
- If error: Returns error wrapped with `utils.ErrAIProviderUnavailable`
- Context cancellation is respected (request aborted if context cancelled)
- SDK errors are mapped to existing error types

**Behavior**:
1. Validate API key is present
2. Build prompt from repository state (same logic as current implementation)
3. Create SDK request with model, messages, max tokens
4. Execute SDK API call with context (respects cancellation/timeout)
5. SDK handles request (with automatic retries if available)
6. Extract message content from SDK response
7. Return generated message string

**Error Cases**:
- Missing API key → Returns error: "API key not configured"
- SDK API call failure → Returns wrapped error with context
- SDK authentication error → Returns error: "API key invalid"
- SDK rate limit error → Returns error: "rate limit exceeded"
- SDK timeout error → Returns error: "timeout"
- Empty response → Returns error: "no response from API"
- Context cancellation → Returns context error
- All errors wrapped with `utils.ErrAIProviderUnavailable` for consistent handling

**Error Mapping**:
- SDK authentication errors → `ErrAIProviderUnavailable` with "API key invalid" message
- SDK rate limit errors → `ErrAIProviderUnavailable` with "rate limit exceeded" message
- SDK timeout errors → `ErrAIProviderUnavailable` with "timeout" message
- SDK network errors → `ErrAIProviderUnavailable` with generic error message
- All errors preserve user-facing messages, no SDK-specific error types exposed

**Side Effects**:
- Makes API call to external AI provider via SDK
- No state changes, no file I/O, no logging of sensitive data (API keys)
- SDK may perform automatic retries (transparent to caller)

## Integration Points

### Provider Registration

**Location**: `internal/service/commit_service.go` → `generateWithAI` method

**Contract**:
- Provider selection mechanism remains unchanged
- Switch statement calls provider constructors (unchanged)
- Provider constructors now initialize SDK clients instead of HTTP clients
- Error handling remains the same (provider errors caught and fallback to manual input)

**Example**:
```go
switch providerName {
case "openai":
    aiProvider = ai.NewOpenAIProvider(providerConfig)  // Now uses SDK
case "anthropic":
    aiProvider = ai.NewAnthropicProvider(providerConfig)  // Now uses SDK
case "mistral":
    aiProvider = ai.NewMistralProvider(providerConfig)  // Now uses SDK
case "local":
    aiProvider = ai.NewLocalProvider(providerConfig)  // Unchanged
}
```

### Configuration Loading

**Location**: `internal/config/config.go` → `GetProviderConfig` method

**Contract**:
- Configuration loading remains unchanged
- AIProviderConfig structure may be extended with optional SDK-specific fields
- Existing configs without new fields continue to work (backward compatible)
- Config parsing handles missing optional fields gracefully

## Performance Contract

### Timeout Handling

- Default timeout: 30 seconds (configurable via config)
- Context cancellation respected by SDKs
- SDK timeout matches config timeout
- SDK automatic retries do not exceed timeout (SDK handles this)

### Resource Management

- SDK clients created per provider instance (reused for multiple calls)
- SDK clients properly closed/cleaned up (SDK handles this)
- No connection pooling requirements (SDK manages connections)
- Context cancellation properly propagated to SDKs

## Backward Compatibility Contract

### Interface Compatibility

- `AIProvider` interface remains unchanged
- All providers implement the same interface contract
- Method signatures unchanged
- Return types unchanged

### Configuration Compatibility

- Existing `AIProviderConfig` fields remain unchanged
- New optional fields only added if SDKs require them
- Existing configs without new fields continue to work
- Config file format unchanged (YAML structure preserved)

### Error Handling Compatibility

- Error types remain the same (`utils.ErrAIProviderUnavailable`)
- User-facing error messages remain the same
- Error handling behavior unchanged (fallback to manual input)
- No SDK-specific error types exposed
