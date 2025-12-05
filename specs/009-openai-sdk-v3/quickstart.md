# Quick Start: Upgrade OpenAI Provider to SDK v3

**Feature**: 009-openai-sdk-v3
**Date**: 2025-01-27

## Overview

This guide provides quick reference for upgrading the OpenAI provider from SDK v1 to SDK v3. The upgrade maintains 100% backward compatibility - no changes to configuration, interfaces, or user experience.

## Prerequisites

- Go 1.25.0 or later
- Existing gitcomm installation with OpenAI provider configured
- OpenAI API key configured in `~/.gitcomm/config.yaml`

## Upgrade Steps

### 1. Update Dependencies

```bash
# Update go.mod to use SDK v3
go get github.com/openai/openai-go/v3@latest

# Clean up dependencies
go mod tidy
```

### 2. Update Import Paths

**File**: `internal/ai/openai_provider.go`

**Before (SDK v1)**:
```go
import (
    "github.com/openai/openai-go"
    "github.com/openai/openai-go/option"
    "github.com/openai/openai-go/shared"
)
```

**After (SDK v3)**:
```go
import (
    "github.com/openai/openai-go/v3"
    "github.com/openai/openai-go/v3/option"
    "github.com/openai/openai-go/v3/shared"
)
```

### 3. Verify API Compatibility

Check if SDK v3 API calls match the current implementation:

```go
// Verify client initialization
client := openai.NewClient(
    option.WithAPIKey(config.APIKey),
)

// Verify chat completion request
req := openai.ChatCompletionNewParams{
    Model: shared.ChatModel(modelName),
    Messages: []openai.ChatCompletionMessageParamUnion{
        openai.SystemMessage(systemPrompt),
        openai.UserMessage(userPrompt),
    },
    MaxCompletionTokens: openai.Int(int64(maxTokens)),
}

// Verify API call
resp, err := client.Chat.Completions.New(ctx, req)
```

### 4. Update Error Handling (if needed)

If SDK v3 introduces new error types, update the `mapSDKError` function:

```go
func (p *OpenAIProvider) mapSDKError(err error) error {
    // Map known error patterns
    // ...

    // Generic fallback for unmappable errors
    return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
}
```

### 5. Run Tests

```bash
# Run unit tests
go test ./internal/ai/... -v

# Run integration tests
go test ./test/integration/... -v
```

### 6. Verify Functionality

```bash
# Test with real API key
export OPENAI_API_KEY="your-api-key"
cd /path/to/git/repo
gitcomm

# Verify commit message is generated successfully
```

## Configuration (Unchanged)

No configuration changes required. Existing configuration continues to work:

```yaml
# ~/.gitcomm/config.yaml
ai:
  default_provider: openai
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
      timeout: 30s
```

## Troubleshooting

### Issue: Import errors after upgrade

**Solution**: Verify import paths are updated to `/v3`:
```go
import "github.com/openai/openai-go/v3"  // Correct
import "github.com/openai/openai-go"     // Wrong (old version)
```

### Issue: API call errors

**Solution**: Check if SDK v3 API structure changed. Review SDK v3 documentation and update API calls accordingly.

### Issue: Error handling not working

**Solution**: Verify error mapping function handles all SDK v3 error types. Update `mapSDKError` to include new error patterns.

### Issue: Tests failing

**Solution**: Update test mocks to match SDK v3 API structure. Verify test setup uses SDK v3 types.

## Verification Checklist

- [ ] Dependencies updated (`go.mod` shows SDK v3)
- [ ] Import paths updated to `/v3`
- [ ] Client initialization works
- [ ] API calls work with SDK v3
- [ ] Error handling maps correctly
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Real API key test succeeds
- [ ] Configuration unchanged
- [ ] No breaking changes to user experience

## Rollback (if needed)

If issues occur, rollback to SDK v1:

```bash
# Revert to SDK v1
go get github.com/openai/openai-go@v1.12.0

# Revert import paths
# (manually update internal/ai/openai_provider.go)

# Clean up
go mod tidy
```

## Next Steps

After successful upgrade:

1. Monitor error logs for any SDK v3-specific issues
2. Update documentation if SDK v3 introduces new features
3. Consider leveraging new SDK v3 features in future enhancements

## References

- OpenAI Go SDK: https://github.com/openai/openai-go
- Current Implementation: `internal/ai/openai_provider.go`
- Test Suite: `internal/ai/openai_provider_test.go`
