# Data Model: Use Official SDKs for AI Providers

**Feature**: 008-official-sdk-integration
**Date**: 2025-01-27

## Domain Entities

### Existing Entities (Unchanged)

- **AIProvider Interface**: Existing interface that all providers must implement - remains unchanged
  - `GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)`
  - No changes to interface contract

- **RepositoryState**: Existing structure passed to all providers for commit message generation
  - No changes needed

- **CommitMessage**: Existing structure for parsed and validated commit messages
  - No changes needed

### Modified Entity: AIProviderConfig

**Current Structure** (existing):
```go
type AIProviderConfig struct {
    Name      string        // Provider name (openai, anthropic, mistral, local)
    APIKey    string        // API key or authentication token
    Model     string        // Optional model identifier
    Endpoint  string        // Optional custom API endpoint (for local models)
    Timeout   time.Duration // Optional request timeout (default: 30s)
    MaxTokens int           // Optional maximum tokens for response (default: 500)
}
```

**Extended Structure** (optional SDK-specific fields):
```go
type AIProviderConfig struct {
    // Existing fields (unchanged, required for backward compatibility)
    Name      string
    APIKey    string
    Model     string
    Endpoint  string
    Timeout   time.Duration
    MaxTokens int

    // Optional SDK-specific configuration (new, optional)
    // These fields are nil/empty if not provided (backward compatible)
    OpenAISDKConfig   *OpenAISDKConfig   // Optional OpenAI SDK-specific options
    AnthropicSDKConfig *AnthropicSDKConfig // Optional Anthropic SDK-specific options
    MistralSDKConfig  *MistralSDKConfig   // Optional Mistral SDK-specific options
}

// SDK-specific config structures (only if SDKs require additional options)
type OpenAISDKConfig struct {
    // OpenAI SDK-specific options (if any)
    // Example: RetryConfig, CustomHeaders, etc.
}

type AnthropicSDKConfig struct {
    // Anthropic SDK-specific options (if any)
}

type MistralSDKConfig struct {
    // Mistral SDK-specific options (if any)
}
```

**Backward Compatibility**:
- Existing configs without new fields continue to work (fields are optional)
- Config parsing handles missing optional fields gracefully (nil/zero values)
- No breaking changes to config file structure
- New fields only added if SDKs require additional configuration options

## Data Flow

### Provider Initialization Flow

```
[User configures provider in config.yaml]
  ↓
[AIProviderConfig loaded from config]
  ↓
[Provider constructor called with config]
  ↓
[SDK client initialized with API key and config]
  ↓
[SDK initialization success?]
  ├─ No → [Fail fast with error, fallback to manual input]
  └─ Yes → [Provider ready for commit message generation]
```

### Commit Message Generation Flow

```
[CommitService.generateWithAI receives providerName]
  ↓
[GetProviderConfig(providerName) returns AIProviderConfig]
  ↓
[NewProvider(providerConfig) creates provider instance]
  ↓
[Provider initializes SDK client with config]
  ↓
[Provider.GenerateCommitMessage(ctx, repoState) called]
  ↓
[Provider builds prompt from repoState (unchanged logic)]
  ↓
[Provider calls SDK API with context and request]
  ↓
[SDK handles request (with automatic retries if available)]
  ↓
[SDK returns response or error]
  ↓
[Provider maps SDK response/error to existing types]
  ↓
[Generated commit message string returned]
  ↓
[CommitService parses and validates message (unchanged)]
```

## Validation Rules

- **API Key**: Must be non-empty string when provider is selected (existing rule)
- **Model**: Optional, provider-specific defaults if not specified (existing rule)
- **Timeout**: Optional, defaults to 30 seconds if not specified (existing rule)
- **MaxTokens**: Optional, defaults to 500 if not specified (existing rule)
- **SDK-Specific Config**: Optional, nil/zero values use SDK defaults (new rule)
- **Provider Name**: Must be "openai", "anthropic", "mistral", or "local" (existing rule)
- **SDK Initialization**: Must succeed or fail fast with clear error (new rule)

## State Transitions

### SDK Client Initialization State Machine

```
[Provider Config Loaded]
  ↓
[API Key Present?]
  ├─ No → [Provider Unavailable Error → Fallback to Manual]
  └─ Yes → [Initialize SDK Client]
            ↓
            [SDK Initialization Success?]
            ├─ No → [Fail Fast with Error → Fallback to Manual]
            └─ Yes → [SDK Client Ready]
                      ↓
                      [Provider Ready]
```

### API Call State Machine

```
[Provider Ready]
  ↓
[GenerateCommitMessage Called]
  ↓
[Build Prompt from RepositoryState (unchanged)]
  ↓
[Create SDK Request with Context]
  ↓
[Execute SDK API Call]
  ↓
[SDK Handles Request (with retries if available)]
  ↓
[Response Received?]
  ├─ No (Error) → [Map SDK Error to Existing Error Type]
  │                 ↓
  │                 [Return Wrapped Error → Fallback to Manual]
  └─ Yes → [Extract Message Content from SDK Response]
            ↓
            [Map to Existing Response Format]
            ↓
            [Return Generated Message]
```

## Relationships

- **Uses**: `AIProviderConfig` (modified) - provides configuration for SDK clients, may include optional SDK-specific fields
- **Uses**: `RepositoryState` (existing) - provides context for commit message generation
- **Implements**: `AIProvider` interface (existing) - defines provider contract, unchanged
- **Used By**: `CommitService` (existing) - calls providers for message generation, unchanged
- **Uses**: Official SDKs (new) - OpenAI SDK, Anthropic SDK, Mistral SDK for API interactions
- **No Dependencies**: No new entities depend on this feature

## Configuration Structure

### AIProviderConfig Fields

**Existing Fields** (unchanged, required for backward compatibility):
- `Name`: Provider identifier ("openai", "anthropic", "mistral", "local")
- `APIKey`: API key (required)
- `Model`: Model name (optional, provider-specific defaults)
- `Endpoint`: API endpoint (optional, for local models)
- `Timeout`: Request timeout (optional, default: 30 seconds)
- `MaxTokens`: Maximum tokens (optional, default: 500)

**New Optional Fields** (only if SDKs require additional options):
- `OpenAISDKConfig`: Optional OpenAI SDK-specific configuration
- `AnthropicSDKConfig`: Optional Anthropic SDK-specific configuration
- `MistralSDKConfig`: Optional Mistral SDK-specific configuration

**Backward Compatibility**:
- Existing configs without new fields continue to work
- New fields are optional (nil/zero values use SDK defaults)
- Config parsing handles missing fields gracefully
- No breaking changes to YAML config file structure

## SDK Response Structures

### OpenAI SDK Response

```go
// SDK response structure (internal to provider implementation)
type OpenAIChatCompletionResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}
```

### Anthropic SDK Response

```go
// SDK response structure (internal to provider implementation)
type AnthropicMessageResponse struct {
    Content []struct {
        Text string `json:"text"`
    } `json:"content"`
}
```

### Mistral SDK Response

```go
// SDK response structure (internal to provider implementation)
type MistralChatResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}
```

**Note**: These structures are internal to provider implementations, not exposed as domain entities. Providers extract the commit message content and return it as a string, maintaining the existing `AIProvider` interface contract.
