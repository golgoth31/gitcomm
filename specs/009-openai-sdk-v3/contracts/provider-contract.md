# Provider Contract: OpenAI Provider SDK v3 Upgrade

**Feature**: 009-openai-sdk-v3
**Date**: 2025-01-27

## Contract Overview

This document defines the function contracts for the OpenAI provider implementation after upgrading to SDK v3. The contracts maintain 100% backward compatibility with the existing `AIProvider` interface.

## AIProvider Interface Contract

**Location**: `internal/ai/provider.go`

**Interface**: Unchanged

```go
type AIProvider interface {
    GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)
}
```

## OpenAIProvider Constructor Contract

**Function**: `NewOpenAIProvider`

**Signature**: `func NewOpenAIProvider(config *model.AIProviderConfig) AIProvider`

### Preconditions

- `config` is not nil
- `config.Name` is "openai" (enforced by caller)
- SDK v3 package `github.com/openai/openai-go/v3` is available

### Postconditions

- Returns a non-nil `AIProvider` implementation
- Provider is ready to generate commit messages (if API key is configured)
- SDK v3 client is initialized (may fail silently if API key is empty, error returned on first API call)

### Behavior

- If `config.APIKey` is empty: Logs debug message, returns provider (error returned on first API call)
- If `config.APIKey` is set: Initializes SDK v3 client with API key
- SDK v3 client initialization does not return error (reads from environment or uses provided options)
- Returns `*OpenAIProvider` that implements `AIProvider` interface

### Error Handling

- No errors returned from constructor (SDK v3 `NewClient` doesn't return error)
- Errors deferred to `GenerateCommitMessage` method

## GenerateCommitMessage Contract

**Function**: `GenerateCommitMessage`

**Signature**: `func (p *OpenAIProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)`

### Preconditions

- `p` is not nil
- `ctx` is not nil
- `repoState` is not nil
- SDK v3 client is initialized (from constructor)

### Postconditions

- On success: Returns non-empty commit message string and nil error
- On failure: Returns empty string and non-nil error wrapped with `ErrAIProviderUnavailable`

### Behavior

1. **API Key Validation**:
   - If `p.config.APIKey` is empty → returns `ErrAIProviderUnavailable: OpenAI API key not configured`

2. **Prompt Building**:
   - Builds prompt from `repoState` (unchanged logic)
   - Includes staged and unstaged file information
   - Adds Conventional Commits format instructions

3. **Model Configuration**:
   - Uses `p.config.Model` if set, otherwise defaults to "gpt-4"
   - Uses `p.config.MaxTokens` if set, otherwise defaults to 500

4. **SDK v3 API Call**:
   - Creates `ChatCompletionNewParams` with SDK v3 types
   - Calls `p.client.Chat.Completions.New(ctx, req)`
   - Respects context cancellation and timeout

5. **Response Processing**:
   - Extracts content from first choice in response
   - Validates response is not empty
   - Returns commit message string

### Error Handling

**Error Mapping Strategy**:
- Authentication errors (401, invalid API key) → `ErrAIProviderUnavailable: API key invalid`
- Rate limit errors (429) → `ErrAIProviderUnavailable: rate limit exceeded`
- Timeout errors → `ErrAIProviderUnavailable: timeout`
- Unmappable SDK v3 errors → `ErrAIProviderUnavailable: <original SDK error message>`

**Error Wrapping**:
- All errors wrapped with `ErrAIProviderUnavailable`
- Original SDK error message preserved in wrapped error for debugging
- User-facing error messages remain consistent with SDK v1 behavior

### Context Handling

- Respects `ctx` cancellation (SDK v3 supports context)
- Respects `ctx` timeout/deadline
- Uses `p.config.Timeout` if context doesn't have deadline

### Thread Safety

- Method is safe for concurrent calls (SDK v3 client is thread-safe)
- No shared mutable state

## mapSDKError Contract

**Function**: `mapSDKError` (internal helper)

**Signature**: `func (p *OpenAIProvider) mapSDKError(err error) error`

### Preconditions

- `err` is not nil
- `err` is an SDK v3 error

### Postconditions

- Returns non-nil error wrapped with `ErrAIProviderUnavailable`
- Original error message preserved for debugging

### Behavior

- Maps known SDK v3 error patterns to specific error messages
- Falls back to generic wrapping for unmappable errors
- Preserves original SDK error in wrapped error

## buildPrompt Contract

**Function**: `buildPrompt` (internal helper)

**Signature**: `func (p *OpenAIProvider) buildPrompt(repoState *model.RepositoryState) string`

### Preconditions

- `repoState` is not nil

### Postconditions

- Returns non-empty prompt string
- Prompt includes file information and format instructions

### Behavior

- Unchanged from SDK v1 implementation
- Builds prompt from staged and unstaged files
- Adds Conventional Commits format instructions

## Testing Contracts

### Unit Test Requirements

- Test `NewOpenAIProvider` with valid and invalid configs
- Test `GenerateCommitMessage` with successful SDK v3 responses
- Test `GenerateCommitMessage` with SDK v3 errors (authentication, rate limit, timeout)
- Test `GenerateCommitMessage` with unmappable SDK v3 errors
- Test context cancellation and timeout handling
- Test error mapping for all error types

### Integration Test Requirements

- Test end-to-end commit message generation with SDK v3
- Test with real API key (if available in test environment)
- Test error handling with invalid API key
- Verify backward compatibility with existing test suite

## Backward Compatibility Guarantees

- `AIProvider` interface unchanged
- `NewOpenAIProvider` signature unchanged
- `GenerateCommitMessage` signature unchanged
- Error types unchanged (`ErrAIProviderUnavailable`)
- Error messages unchanged (user-facing)
- Configuration structure unchanged
- Behavior unchanged (functional parity)

## Version Compatibility

- **SDK v1**: Previous implementation (to be replaced)
- **SDK v3**: New implementation (target)
- **Interface**: Unchanged (backward compatible)
