# Data Model: Upgrade OpenAI Provider to SDK v3

**Feature**: 009-openai-sdk-v3
**Date**: 2025-01-27

## Overview

This feature does not introduce new domain entities or modify existing data structures. The upgrade is purely an internal implementation change that maintains 100% backward compatibility with existing data models.

## Existing Entities (Unchanged)

### AIProviderConfig

**Location**: `internal/model/config.go`

**Structure**: Unchanged - no modifications required

```go
type AIProviderConfig struct {
    Name     string        // Provider name (unchanged)
    APIKey   string        // API key (unchanged)
    Model    string        // Model identifier (unchanged)
    Endpoint string        // Custom endpoint (unchanged, not used for OpenAI)
    Timeout  time.Duration // Request timeout (unchanged)
    MaxTokens int          // Maximum tokens (unchanged)
}
```

**Validation Rules**: Unchanged
- `Name` must be "openai" for OpenAI provider
- `APIKey` validated at runtime (empty key returns error)
- `Model` defaults to "gpt-4" if empty
- `Timeout` defaults to 30 seconds if zero
- `MaxTokens` defaults to 500 if zero

### RepositoryState

**Location**: `internal/model/repository_state.go`

**Structure**: Unchanged - no modifications required

**Usage**: Input to `GenerateCommitMessage` method (unchanged)

### AIProvider Interface

**Location**: `internal/ai/provider.go`

**Structure**: Unchanged - no modifications required

```go
type AIProvider interface {
    GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)
}
```

## Internal Implementation Changes

### OpenAIProvider Struct

**Location**: `internal/ai/openai_provider.go`

**Changes**: Internal implementation only - public interface unchanged

**Before (SDK v1)**:
```go
type OpenAIProvider struct {
    config *model.AIProviderConfig
    client openai.Client  // SDK v1 client
}
```

**After (SDK v3)**:
```go
type OpenAIProvider struct {
    config *model.AIProviderConfig
    client openai.Client  // SDK v3 client (different import path)
}
```

**Note**: The struct fields remain the same, but the `openai.Client` type comes from a different package (`github.com/openai/openai-go/v3` instead of `github.com/openai/openai-go`).

## Data Flow (Unchanged)

1. **Input**: `RepositoryState` → `GenerateCommitMessage`
2. **Processing**: Build prompt from `RepositoryState`, call SDK API
3. **Output**: Commit message string or error

**No changes to data flow** - only internal SDK implementation changes.

## Validation Rules (Unchanged)

- API key validation: Empty key returns `ErrAIProviderUnavailable`
- Model validation: Empty model defaults to "gpt-4"
- Response validation: Empty response returns `ErrAIProviderUnavailable`
- Error validation: All errors mapped to existing error types

## State Transitions (N/A)

No state transitions - stateless API calls only.

## Relationships (Unchanged)

- `OpenAIProvider` implements `AIProvider` interface (unchanged)
- `OpenAIProvider` uses `AIProviderConfig` for configuration (unchanged)
- `OpenAIProvider` receives `RepositoryState` as input (unchanged)

## Summary

**No new entities, no modified entities, no deleted entities.**

This upgrade is a pure implementation change that maintains complete backward compatibility with existing data models and interfaces.
