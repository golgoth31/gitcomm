# Data Model: Add Mistral as AI Provider

**Feature**: 007-mistral-provider
**Date**: 2025-01-27

## Domain Entities

### No New Entities Required

This feature extends existing data structures without introducing new domain entities. The feature uses existing data structures:

- **AIProviderConfig**: Existing configuration structure used by all providers (Name, APIKey, Model, Endpoint, Timeout, MaxTokens)
- **AIProvider Interface**: Existing interface that MistralProvider will implement
- **RepositoryState**: Existing structure passed to all providers for commit message generation
- **CommitMessage**: Existing structure for parsed and validated commit messages

## Data Flow

### Provider Selection to Message Generation Flow

```
[User selects Mistral provider via config or CLI flag]
  ↓
[CommitService.generateWithAI receives providerName="mistral"]
  ↓
[GetProviderConfig("mistral") returns AIProviderConfig]
  ↓
[NewMistralProvider(providerConfig) creates MistralProvider instance]
  ↓
[MistralProvider.GenerateCommitMessage(ctx, repoState) called]
  ↓
[MistralProvider builds prompt from repoState]
  ↓
[MistralProvider makes HTTP POST to Mistral API]
  ↓
[MistralProvider parses API response]
  ↓
[Generated commit message string returned]
  ↓
[CommitService parses and validates message]
  ↓
[Validated CommitMessage returned to user]
```

## Validation Rules

- **API Key**: Must be non-empty string when provider is selected
- **Model**: Optional, defaults to "mistral-large-latest" if not specified
- **Endpoint**: Optional, defaults to "https://api.mistral.ai/v1/chat/completions" if not specified
- **Timeout**: Optional, defaults to 30 seconds if not specified
- **MaxTokens**: Optional, defaults to 500 if not specified
- **Provider Name**: Must be "mistral" (case-sensitive in switch statement)

## State Transitions

### Provider Initialization State Machine

```
[Provider Config Loaded]
  ↓
[API Key Present?]
  ├─ No → [Provider Unavailable Error]
  └─ Yes → [MistralProvider Created]
            ↓
            [HTTP Client Initialized with Timeout]
            ↓
            [Provider Ready]
```

### API Call State Machine

```
[Provider Ready]
  ↓
[GenerateCommitMessage Called]
  ↓
[Build Prompt from RepositoryState]
  ↓
[Create HTTP Request with Context]
  ↓
[Execute Request]
  ↓
[Response Received?]
  ├─ No (Error) → [Error Handling → Fallback to Manual]
  └─ Yes → [Parse Response]
            ↓
            [Extract Message Content]
            ↓
            [Return Generated Message]
```

## Relationships

- **Uses**: `AIProviderConfig` (existing) - provides configuration for Mistral API
- **Uses**: `RepositoryState` (existing) - provides context for commit message generation
- **Implements**: `AIProvider` interface (existing) - defines provider contract
- **Used By**: `CommitService` (modified) - calls MistralProvider for message generation
- **Uses**: `TokenCalculator` (modified) - calculates token estimates for Mistral requests
- **No Dependencies**: No new entities depend on this feature

## Configuration Structure

### AIProviderConfig Fields (Existing, Used by Mistral)

- `Name`: "mistral" (provider identifier)
- `APIKey`: Mistral API key (required)
- `Model`: Mistral model name (optional, default: "mistral-large-latest")
- `Endpoint`: API endpoint URL (optional, default: "https://api.mistral.ai/v1/chat/completions")
- `Timeout`: Request timeout (optional, default: 30 seconds)
- `MaxTokens`: Maximum tokens in response (optional, default: 500)

## API Request/Response Structures

### Request Structure (Mistral API)

```go
type MistralRequest struct {
    Model      string    `json:"model"`
    Messages   []Message `json:"messages"`
    MaxTokens  int       `json:"max_tokens"`
}

type Message struct {
    Role    string `json:"role"`    // "system" or "user"
    Content string `json:"content"`
}
```

### Response Structure (Mistral API)

```go
type MistralResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}
```

**Note**: These structures are internal to MistralProvider implementation, not exposed as domain entities.
