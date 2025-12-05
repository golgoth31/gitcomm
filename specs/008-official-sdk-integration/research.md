# Research: Use Official SDKs for AI Providers

**Feature**: 008-official-sdk-integration
**Date**: 2025-01-27

## Technology Decisions

### 1. OpenAI Go SDK Integration

**Decision**: Use `github.com/openai/openai-go` official SDK for OpenAI provider interactions.

**Rationale**:
- Official SDK maintained by OpenAI ensures API compatibility and updates
- Provides automatic retries and better error handling than raw HTTP clients
- Supports context cancellation for timeout handling
- Follows Go best practices and conventions
- Reduces maintenance burden of manual HTTP request/response handling

**Alternatives Considered**:
- Continue using raw HTTP client: More maintenance, manual error handling, no automatic retries
- Third-party OpenAI Go libraries: Less reliable, may not be maintained, no official support

**SDK Features to Use**:
- Chat completions API (replaces current HTTP POST to `/v1/chat/completions`)
- Context support for cancellation/timeout
- Automatic retries (if available, improves reliability)
- Structured error types (mapped to existing error handling)

**SDK Features NOT to Use**:
- Streaming responses (out of scope, not needed for commit message generation)
- Fine-tuning APIs (not relevant)
- Embeddings APIs (not relevant)

**Implementation Pattern**:
```go
// Initialize SDK client
client := openai.NewClient(apiKey)

// Create chat completion request
req := openai.ChatCompletionRequest{
    Model: model,
    Messages: []openai.ChatCompletionMessage{
        {Role: "system", Content: systemPrompt},
        {Role: "user", Content: userPrompt},
    },
    MaxTokens: maxTokens,
}

// Execute with context
resp, err := client.CreateChatCompletion(ctx, req)
```

**API Documentation Reference**: https://github.com/openai/openai-go

### 2. Anthropic Go SDK Integration

**Decision**: Use `github.com/anthropics/anthropic-sdk-go` official SDK for Anthropic provider interactions.

**Rationale**:
- Official SDK maintained by Anthropic ensures API compatibility
- Provides better error handling and type safety than raw HTTP clients
- Supports context cancellation for timeout handling
- Follows Go best practices
- Reduces maintenance burden

**Alternatives Considered**:
- Continue using raw HTTP client: More maintenance, manual error handling
- Third-party Anthropic Go libraries: Less reliable, may not be maintained

**SDK Features to Use**:
- Messages API (replaces current HTTP POST to `/v1/messages`)
- Context support for cancellation/timeout
- Automatic retries (if available)
- Structured error types (mapped to existing error handling)

**SDK Features NOT to Use**:
- Streaming responses (out of scope)
- Text completion APIs (not used)
- Other Anthropic APIs not relevant to commit message generation

**Implementation Pattern**:
```go
// Initialize SDK client
client, err := anthropic.NewClient(apiKey)

// Create message request
req := anthropic.MessageRequest{
    Model: model,
    Messages: []anthropic.Message{
        {Role: "user", Content: prompt},
    },
    MaxTokens: maxTokens,
}

// Execute with context
resp, err := client.CreateMessage(ctx, req)
```

**API Documentation Reference**: https://github.com/anthropics/anthropic-sdk-go

### 3. Mistral Go SDK Integration

**Decision**: Use `github.com/Gage-Technologies/mistral-go` official SDK for Mistral provider interactions.

**Rationale**:
- Official SDK maintained by Gage Technologies for Mistral AI
- Provides better error handling than raw HTTP clients
- Supports context cancellation for timeout handling
- Follows Go best practices
- Reduces maintenance burden

**Alternatives Considered**:
- Continue using raw HTTP client: More maintenance, manual error handling
- Other Mistral Go libraries: Less reliable, may not be maintained

**SDK Features to Use**:
- Chat completions API (replaces current HTTP POST to `/v1/chat/completions`)
- Context support for cancellation/timeout
- Automatic retries (if available)
- Structured error types (mapped to existing error handling)

**SDK Features NOT to Use**:
- Streaming responses (out of scope)
- Embeddings APIs (not relevant)
- Other Mistral APIs not relevant to commit message generation

**Implementation Pattern**:
```go
// Initialize SDK client
client := mistral.NewMistralClientDefault(apiKey)

// Create chat request
req := mistral.ChatRequest{
    Model: model,
    Messages: []mistral.ChatMessage{
        {Role: "system", Content: systemPrompt},
        {Role: "user", Content: userPrompt},
    },
    MaxTokens: maxTokens,
}

// Execute with context
resp, err := client.Chat(ctx, model, messages, nil)
```

**API Documentation Reference**: https://github.com/Gage-Technologies/mistral-go

### 4. Error Type Mapping Strategy

**Decision**: Map SDK-specific error types to existing error handling patterns using error wrapping.

**Rationale**:
- Maintains existing error handling behavior (FR-008)
- Preserves user-facing error messages (FR-012)
- Maintains abstraction (no SDK-specific errors exposed)
- Consistent error handling across all providers

**Implementation Pattern**:
```go
// SDK returns SDK-specific error
resp, err := sdkClient.Call(ctx, req)
if err != nil {
    // Map to existing error type
    if isAuthError(err) {
        return "", fmt.Errorf("%w: API key invalid", utils.ErrAIProviderUnavailable)
    }
    if isRateLimitError(err) {
        return "", fmt.Errorf("%w: rate limit exceeded", utils.ErrAIProviderUnavailable)
    }
    // Generic error mapping
    return "", fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
}
```

**Error Mapping Rules**:
- Authentication errors → `ErrAIProviderUnavailable` with "API key invalid" message
- Rate limit errors → `ErrAIProviderUnavailable` with "rate limit exceeded" message
- Timeout errors → `ErrAIProviderUnavailable` with "timeout" message
- Network errors → `ErrAIProviderUnavailable` with generic error message
- All errors wrapped with `utils.ErrAIProviderUnavailable` for consistent handling

### 5. SDK Initialization Failure Handling

**Decision**: Fail fast with clear error message, fallback to manual input (no HTTP client fallback).

**Rationale**:
- Prevents silent failures and confusion
- Clear error messages help users diagnose issues
- Maintains user workflow with fallback to manual input
- Aligns with FR-007 requirement

**Implementation Pattern**:
```go
// Initialize SDK client
client, err := sdk.NewClient(apiKey)
if err != nil {
    // Fail fast with clear error
    return nil, fmt.Errorf("failed to initialize SDK: %w", err)
    // Service layer will catch this and fallback to manual input
}
```

**Error Scenarios**:
- Missing SDK dependency → Build error (fail at build time)
- SDK version incompatibility → Build/initialization error with clear message
- SDK initialization error → Return error, service layer falls back to manual input
- No attempt to use HTTP client as fallback (per clarification)

### 6. Configuration Extension Strategy

**Decision**: Extend AIProviderConfig with optional SDK-specific fields if needed, maintaining backward compatibility.

**Rationale**:
- Allows SDK-specific configuration options without breaking existing configs
- Maintains backward compatibility (FR-006)
- Optional fields ensure existing configs continue to work
- Follows Go best practices for optional configuration

**Implementation Pattern**:
```go
type AIProviderConfig struct {
    // Existing fields (unchanged)
    Name     string
    APIKey   string
    Model    string
    Endpoint string
    Timeout  time.Duration
    MaxTokens int

    // Optional SDK-specific fields (new, optional)
    OpenAIConfig   *OpenAISDKConfig   // Optional OpenAI-specific config
    AnthropicConfig *AnthropicSDKConfig // Optional Anthropic-specific config
    MistralConfig  *MistralSDKConfig   // Optional Mistral-specific config
}
```

**Backward Compatibility**:
- Existing configs without new fields continue to work
- New fields are optional (nil/zero values use SDK defaults)
- No breaking changes to config file structure
- Config parsing handles missing optional fields gracefully

### 7. SDK Automatic Retries

**Decision**: Use SDK-provided automatic retries if available to improve reliability.

**Rationale**:
- Improves reliability without changing user-facing behavior
- SDK retries are well-tested and handle transient errors
- Reduces need for manual retry logic
- Aligns with FR-005 (improve functionality where possible)

**Implementation Pattern**:
```go
// SDKs with built-in retries will automatically retry on transient errors
// No additional configuration needed if SDK supports it
// If SDK doesn't support retries, maintain single-attempt behavior
```

**Retry Behavior**:
- Use SDK retries if available (transparent to user)
- No user-facing changes (same timeout, same error messages)
- Retries handled by SDK internally
- No additional configuration needed

### 8. Context and Timeout Handling

**Decision**: Use context.Context for all SDK API calls to support cancellation and timeout.

**Rationale**:
- Maintains existing timeout behavior (FR-009)
- Supports context cancellation for graceful shutdown
- All official SDKs support context.Context
- Follows Go best practices

**Implementation Pattern**:
```go
// Create context with timeout from config
ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
defer cancel()

// Pass context to SDK call
resp, err := sdkClient.Call(ctx, req)
// SDK respects context cancellation and timeout
```

**Timeout Behavior**:
- Use timeout from config (default 30s if not specified)
- Context cancellation respected by SDKs
- Timeout errors mapped to existing error types
- Same timeout behavior as current HTTP client implementation
