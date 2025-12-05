# Research: Upgrade OpenAI Provider to SDK v3

**Feature**: 009-openai-sdk-v3
**Date**: 2025-01-27

## Technology Decisions

### 1. OpenAI Go SDK v3 Migration

**Decision**: Upgrade from `github.com/openai/openai-go` v1.12.0 to `github.com/openai/openai-go/v3`.

**Rationale**:
- Official SDK v3 provides latest features, bug fixes, and improvements
- Maintains backward compatibility with existing API patterns
- Follows Go module versioning best practices (v3 as major version)
- Reduces technical debt by staying current with SDK releases
- Improves maintainability and access to latest OpenAI API features

**Alternatives Considered**:
- Continue using SDK v1: Accumulates technical debt, misses bug fixes and improvements
- Use third-party OpenAI Go libraries: Less reliable, may not be maintained, no official support

**SDK v3 Features to Use**:
- Chat completions API (same core functionality as v1)
- Context support for cancellation/timeout (same as v1)
- Automatic retries (if available, improves reliability)
- Structured error types (mapped to existing error handling)

**SDK v3 Features NOT to Use**:
- Streaming responses (out of scope, not needed for commit message generation)
- Fine-tuning APIs (not relevant)
- Embeddings APIs (not relevant)
- Other OpenAI APIs not relevant to commit message generation

**Migration Pattern**:
```go
// Current (SDK v1)
import "github.com/openai/openai-go"
import "github.com/openai/openai-go/option"
import "github.com/openai/openai-go/shared"

client := openai.NewClient(
    option.WithAPIKey(config.APIKey),
)

req := openai.ChatCompletionNewParams{
    Model: shared.ChatModel(modelName),
    Messages: []openai.ChatCompletionMessageParamUnion{
        openai.SystemMessage(systemPrompt),
        openai.UserMessage(userPrompt),
    },
    MaxCompletionTokens: openai.Int(int64(maxTokens)),
}

resp, err := client.Chat.Completions.New(ctx, req)

// SDK v3 (expected pattern - to be verified)
import "github.com/openai/openai-go/v3"
import "github.com/openai/openai-go/v3/option"
import "github.com/openai/openai-go/v3/shared"

// Similar initialization and API call patterns expected
// May have minor API changes that need to be adapted
```

**API Documentation Reference**: https://github.com/openai/openai-go

### 2. Error Handling Strategy

**Decision**: Map SDK v3 error types to existing error handling patterns. For unmappable errors, wrap generically with `ErrAIProviderUnavailable` while preserving the original SDK error message for debugging.

**Rationale**:
- Maintains consistent user-facing error behavior
- Preserves debugging information for troubleshooting
- Avoids introducing new error types in the infrastructure
- Follows existing error handling patterns established in the codebase

**Error Mapping Patterns**:
- Authentication errors (401, invalid API key) → `ErrAIProviderUnavailable: API key invalid`
- Rate limit errors (429) → `ErrAIProviderUnavailable: rate limit exceeded`
- Timeout errors → `ErrAIProviderUnavailable: timeout`
- Unmappable errors → `ErrAIProviderUnavailable: <original SDK error message>`

**Implementation**:
```go
func (p *OpenAIProvider) mapSDKError(err error) error {
    errStr := err.Error()

    // Map known error patterns
    if strings.Contains(strings.ToLower(errStr), "authentication") ||
       strings.Contains(strings.ToLower(errStr), "invalid") ||
       strings.Contains(errStr, "401") {
        return fmt.Errorf("%w: API key invalid", utils.ErrAIProviderUnavailable)
    }
    // ... other mappings ...

    // Generic fallback for unmappable errors
    return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
}
```

### 3. Backward Compatibility Strategy

**Decision**: Maintain 100% backward compatibility with existing configuration, interfaces, and behavior.

**Rationale**:
- No breaking changes for users
- No configuration file updates required
- No changes to other providers
- Seamless upgrade experience

**Compatibility Guarantees**:
- `AIProvider` interface unchanged
- `AIProviderConfig` structure unchanged
- Configuration file format unchanged
- Error handling behavior identical
- Timeout and context cancellation behavior identical
- Prompt building and response parsing logic unchanged

### 4. Testing Strategy

**Decision**: Update existing tests minimally to match SDK v3 API, maintain same test coverage and behavior.

**Rationale**:
- Ensures functional parity with SDK v1
- Validates backward compatibility
- Maintains test quality and coverage

**Test Updates Required**:
- Update SDK client initialization in test setup (if needed)
- Update API call mocks to match SDK v3 structure (if changed)
- Verify error handling tests still pass with SDK v3 errors
- Ensure integration tests work with SDK v3

**Test Coverage**:
- Unit tests for SDK v3 client initialization
- Unit tests for API calls with SDK v3
- Unit tests for error mapping with SDK v3 errors
- Integration tests for end-to-end commit message generation

### 5. Dependency Management

**Decision**: Update `go.mod` to use `github.com/openai/openai-go/v3` instead of `github.com/openai/openai-go`.

**Rationale**:
- Go module versioning allows v1 and v3 to coexist if needed
- Clean dependency management
- Follows Go best practices for major version upgrades

**Migration Steps**:
1. Update import paths in `internal/ai/openai_provider.go`
2. Run `go get github.com/openai/openai-go/v3@latest`
3. Run `go mod tidy` to clean up dependencies
4. Verify no conflicts with other dependencies

## Implementation Notes

### API Changes to Watch For

While SDK v3 is expected to maintain similar API structure, potential changes to verify:

1. **Client Initialization**: Verify `NewClient` signature and options
2. **Chat Completions API**: Verify `ChatCompletionNewParams` structure and field names
3. **Message Types**: Verify `ChatCompletionMessageParamUnion` and helper functions (`SystemMessage`, `UserMessage`)
4. **Response Structure**: Verify response structure and content extraction
5. **Error Types**: Verify error types and error handling patterns
6. **Context Support**: Verify context cancellation and timeout handling

### Migration Checklist

- [ ] Review SDK v3 documentation and release notes
- [ ] Update import paths to use `/v3`
- [ ] Update `go.mod` dependency
- [ ] Verify client initialization works
- [ ] Verify API calls work with SDK v3
- [ ] Verify error handling maps correctly
- [ ] Update unit tests
- [ ] Update integration tests
- [ ] Verify all tests pass
- [ ] Verify backward compatibility

## References

- OpenAI Go SDK: https://github.com/openai/openai-go
- Go Module Versioning: https://go.dev/doc/modules/version-numbers
- Current Implementation: `internal/ai/openai_provider.go`
