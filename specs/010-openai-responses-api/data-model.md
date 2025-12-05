# Data Model: Migrate OpenAI Provider to Responses API

**Feature**: 010-openai-responses-api
**Date**: 2025-01-27

## Overview

No new domain entities are introduced. The migration only changes the internal API implementation while maintaining the same external interfaces and data structures.

## Existing Entities (Unchanged)

### AIProviderConfig

**Location**: `internal/model/config.go`

**Structure**: Unchanged
- `Name string` - Provider name ("openai")
- `APIKey string` - API key for authentication
- `Model string` - Model identifier (e.g., "gpt-4-1")
- `Endpoint string` - Not used for OpenAI (reserved for local providers)
- `Timeout time.Duration` - Request timeout
- `MaxTokens int` - Maximum tokens for response

**Validation Rules**: Unchanged
- API key must be provided (validated at runtime)
- Model defaults to "gpt-4-1" if empty
- MaxTokens defaults to 500 if 0

### RepositoryState

**Location**: `internal/model/commit_message.go` (or similar)

**Structure**: Unchanged
- `StagedFiles []FileChange` - List of staged file changes
- `UnstagedFiles []FileChange` - List of unstaged file changes

**Usage**: Unchanged - used to build prompts for commit message generation

### FileChange

**Location**: `internal/model/commit_message.go` (or similar)

**Structure**: Unchanged
- `Path string` - File path
- `Status string` - File status (modified, added, deleted, etc.)
- `Diff string` - File diff content

## Internal Request/Response Structures

### Responses API Request (Internal)

**Location**: `internal/ai/openai_provider.go` (implementation detail)

**Structure**:
```go
// If using SDK:
type ResponseNewParams struct {
    Model    string
    Input    []MessageObject  // Array of {role, content} objects
    MaxCompletionTokens int64
    Store    bool             // false for stateless mode
}

// If using custom HTTP client:
type ResponseRequest struct {
    Model    string           `json:"model"`
    Input    []MessageObject  `json:"input"`
    MaxCompletionTokens int64 `json:"max_completion_tokens"`
    Store    bool             `json:"store"` // false for stateless
}

type MessageObject struct {
    Role    string `json:"role"`    // "system" or "user"
    Content string `json:"content"`
}
```

**Mapping from Current Implementation**:
- `Messages []ChatCompletionMessageParamUnion` → `Input []MessageObject`
- Structure is identical (role + content), only parameter name changes

### Responses API Response (Internal)

**Location**: `internal/ai/openai_provider.go` (implementation detail)

**Structure** (expected, to be verified during implementation):
```go
// If using SDK:
type Response struct {
    Content string  // or Text string
    // ... other fields
}

// If using custom HTTP client:
type ResponseResponse struct {
    Content string `json:"content"`  // or Text string `json:"text"`
    // May be nested similar to Chat Completions structure
}
```

**Extraction Logic**:
- Extract `content` or `text` field (similar to `choices[0].message.content`)
- Handle empty responses
- Map to string return value

## State Transitions

No state transitions - stateless API calls. Each `GenerateCommitMessage` call is independent.

## Validation Rules

**Request Validation**:
- API key must be non-empty (checked before API call)
- Model must be valid (defaults to "gpt-4-1" if empty)
- Input array must contain at least system and user messages
- MaxCompletionTokens must be positive (defaults to 500)

**Response Validation**:
- Response must contain non-empty content/text field
- Handle empty response gracefully with error

## Relationships

- `OpenAIProvider` uses `AIProviderConfig` for configuration
- `OpenAIProvider.GenerateCommitMessage` takes `RepositoryState` as input
- `OpenAIProvider` implements `AIProvider` interface (unchanged)
- No new relationships introduced

## Notes

- All external interfaces remain unchanged
- Internal request/response structures are implementation details
- Data model changes are minimal (parameter name changes only)
- No database or persistent storage involved
