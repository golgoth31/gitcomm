# Research: Add Mistral as AI Provider

**Feature**: 007-mistral-provider
**Date**: 2025-01-27

## Technology Decisions

### 1. Mistral API Endpoint and Authentication

**Decision**: Use Mistral AI's standard API endpoint `https://api.mistral.ai/v1/chat/completions` with Bearer token authentication via `Authorization` header.

**Rationale**:
- Mistral AI follows OpenAI-compatible API patterns (similar to OpenAI's chat completions endpoint)
- Standard REST API with HTTP POST requests
- Bearer token authentication is standard and secure
- Compatible with existing HTTP client patterns used for OpenAI/Anthropic

**Alternatives Considered**:
- Custom endpoint: Not needed, standard endpoint is well-documented and stable
- Different authentication method: Bearer token is standard and matches OpenAI pattern

**Implementation Pattern**:
```go
// Create request
req, err := http.NewRequestWithContext(ctx, "POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
```

**API Documentation Reference**: Mistral AI uses OpenAI-compatible chat completions API at `https://api.mistral.ai/v1/chat/completions`

### 2. Default Mistral Model Selection

**Decision**: Use `mistral-large-latest` as the default model for commit message generation.

**Rationale**:
- `mistral-large-latest` is Mistral's most capable model, suitable for structured text generation
- "latest" suffix ensures users get the most recent version automatically
- Large models provide better adherence to format specifications (Conventional Commits)
- Matches the pattern of using high-quality models for commit message generation (similar to gpt-4, claude-3-opus)

**Alternatives Considered**:
- `mistral-small-latest`: Less capable, may not follow format as well
- `mistral-medium-latest`: Good balance but large is better for structured output
- Specific versioned models (e.g., `mistral-large-2402`): Less flexible, requires updates

**Model Options**:
- `mistral-tiny`: Fastest, least capable
- `mistral-small`: Good balance
- `mistral-medium`: Better quality
- `mistral-large-latest`: Best quality (default choice)

### 3. Mistral API Request Format

**Decision**: Use OpenAI-compatible chat completions request format with messages array.

**Rationale**:
- Mistral API is OpenAI-compatible, using the same request/response structure
- Matches existing OpenAI provider implementation pattern
- Simplifies implementation by reusing similar code structure
- Well-documented and stable API format

**Request Structure**:
```json
{
  "model": "mistral-large-latest",
  "messages": [
    {
      "role": "system",
      "content": "You are a git commit message generator..."
    },
    {
      "role": "user",
      "content": "Generate a commit message for..."
    }
  ],
  "max_tokens": 500
}
```

**Response Structure**:
```json
{
  "choices": [
    {
      "message": {
        "content": "feat(api): add new endpoint..."
      }
    }
  ]
}
```

**Implementation Pattern**: Follow OpenAI provider structure with Mistral-specific endpoint and model.

### 4. Token Calculation for Mistral

**Decision**: Use character-based fallback token calculator for Mistral (no Mistral-specific tokenization library available).

**Rationale**:
- No widely-available Go library for Mistral tokenization (unlike tiktoken for OpenAI)
- Character-based estimation is acceptable for token counting (within 10% accuracy requirement)
- Matches fallback strategy used for other providers without specific tokenization
- Simple and maintainable solution
- Mistral uses similar tokenization to other models, so character-based estimation is reasonably accurate

**Alternatives Considered**:
- Implement Mistral tokenization from scratch: Too complex, not worth the effort for estimation
- Use OpenAI tiktoken as approximation: Inaccurate, Mistral uses different tokenization
- Skip token calculation: Violates requirement FR-010

**Implementation**: Add "mistral" case to `NewTokenCalculator` that returns `NewFallbackTokenCalculator()`.

### 5. Error Handling and API Response Parsing

**Decision**: Follow existing provider error handling patterns with Mistral-specific error messages.

**Rationale**:
- Consistent error handling across all providers
- Clear error messages help users diagnose issues
- Graceful fallback to manual input maintains user workflow
- HTTP status code handling matches existing patterns

**Error Scenarios**:
- Missing API key → Clear error message, fallback to manual
- Invalid API key → HTTP 401, display error, fallback
- Rate limiting → HTTP 429, display error, fallback
- Timeout → Context deadline exceeded, display error, fallback
- Invalid response format → Parse error, display error, fallback

**Implementation Pattern**: Use same error handling structure as OpenAI/Anthropic providers with Mistral-specific error context.

### 6. Provider Registration and Selection

**Decision**: Add "mistral" case to provider switch in `CommitService.generateWithAI` method.

**Rationale**:
- Follows existing pattern (openai, anthropic, local cases)
- Simple switch statement extension
- No architectural changes needed
- Maintains consistency with existing code

**Implementation Pattern**:
```go
switch providerName {
case "openai":
    aiProvider = ai.NewOpenAIProvider(providerConfig)
case "anthropic":
    aiProvider = ai.NewAnthropicProvider(providerConfig)
case "mistral":
    aiProvider = ai.NewMistralProvider(providerConfig)  // NEW
case "local":
    aiProvider = ai.NewLocalProvider(providerConfig)
default:
    return nil, fmt.Errorf("%w: unknown provider %s", utils.ErrAIProviderUnavailable, providerName)
}
```

### 7. Configuration File Updates

**Decision**: Add Mistral provider example to `config.yaml.example` following existing provider format.

**Rationale**:
- Provides clear documentation for users
- Shows expected configuration structure
- Maintains consistency with existing provider examples
- Helps users configure Mistral quickly

**Configuration Example**:
```yaml
mistral:
  api_key: ${MISTRAL_API_KEY}
  model: mistral-large-latest  # Optional, default: mistral-large-latest
  timeout: 30s                 # Optional, default: 30s
```
