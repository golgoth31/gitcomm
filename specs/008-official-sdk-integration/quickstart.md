# Quick Start: Use Official SDKs for AI Providers

**Feature**: 008-official-sdk-integration
**Date**: 2025-01-27

## Overview

This feature replaces HTTP client implementations with official Go SDKs for OpenAI, Anthropic, and Mistral AI providers. The change is transparent to users - existing configurations continue to work without modification, and the behavior remains identical while benefiting from official SDK features like automatic retries and better error handling.

## What Changed

### For Users

**No changes required!** Your existing configuration files continue to work exactly as before. The CLI behavior, error messages, and output format remain identical.

### For Developers

- Provider implementations now use official SDKs instead of raw HTTP clients
- SDK initialization failures are handled gracefully with clear error messages
- Automatic retries are used if available (improves reliability)
- Error handling remains the same (SDK errors mapped to existing error types)

## Installation

### Dependencies

The feature adds three new Go module dependencies:

```bash
go get github.com/openai/openai-go
go get github.com/anthropics/anthropic-sdk-go
go get github.com/Gage-Technologies/mistral-go
```

These are automatically added when you build the project.

### Build

```bash
# Build from source
go build ./cmd/gitcomm

# Or install globally
go install ./cmd/gitcomm
```

## Configuration

### Existing Configuration (Unchanged)

Your existing `~/.gitcomm/config.yaml` continues to work without modification:

```yaml
ai:
  default_provider: openai
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
      timeout: 30s
    anthropic:
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-opus
      timeout: 30s
    mistral:
      api_key: ${MISTRAL_API_KEY}
      model: mistral-large-latest
      timeout: 30s
```

### Optional SDK-Specific Configuration

If SDKs require additional configuration options in the future, they will be added as optional fields. Existing configs without these fields will continue to work (backward compatible).

## Usage

### Basic Usage (Unchanged)

```bash
# Use default provider (OpenAI)
gitcomm

# Use specific provider
gitcomm --provider anthropic
gitcomm --provider mistral

# Skip AI and use manual input
gitcomm --skip-ai
```

### Behavior

- **Commit message generation**: Works exactly as before
- **Error handling**: Same error messages and fallback behavior
- **Timeout handling**: Same timeout behavior (30s default)
- **Performance**: Same or better (SDK retries improve reliability)

## Troubleshooting

### SDK Initialization Failures

If an SDK fails to initialize (e.g., missing dependency, incompatible version), the CLI will:

1. Display a clear error message indicating the SDK issue
2. Automatically fall back to manual input
3. Allow you to create the commit message manually

**Example error**:
```
Error: failed to initialize OpenAI SDK: [error details]
Falling back to manual input...
```

### Build Errors

If you encounter build errors related to SDK dependencies:

1. Ensure you have the latest Go version (1.25.0+)
2. Run `go mod tidy` to update dependencies
3. Check that SDK versions are compatible with your Go version

### API Errors

API errors are handled the same way as before:

- Invalid API key → Clear error message, fallback to manual input
- Rate limiting → Clear error message, fallback to manual input
- Timeout → Clear error message, fallback to manual input

All errors are wrapped with `ErrAIProviderUnavailable` for consistent handling.

## Migration Guide

### For Existing Users

**No migration needed!** Your existing configuration and usage patterns continue to work without any changes.

### For Developers

If you're extending the provider implementations:

1. **Use SDK clients instead of HTTP clients**: Initialize SDK clients in provider constructors
2. **Map SDK errors**: Map SDK-specific errors to existing error types
3. **Respect context**: Pass `context.Context` to all SDK API calls
4. **Maintain interface**: Keep the `AIProvider` interface contract unchanged

## Testing

### Unit Tests

Provider unit tests are updated to mock SDK clients instead of HTTP clients:

```go
// Mock SDK client
mockClient := &MockOpenAIClient{}
provider := &OpenAIProvider{
    client: mockClient,
    config: config,
}

// Test GenerateCommitMessage
message, err := provider.GenerateCommitMessage(ctx, repoState)
```

### Integration Tests

Integration tests verify SDK integration:

```go
// Test with real SDK (requires API key)
config := &model.AIProviderConfig{
    APIKey: os.Getenv("OPENAI_API_KEY"),
    Model:  "gpt-4",
}
provider := ai.NewOpenAIProvider(config)
message, err := provider.GenerateCommitMessage(ctx, repoState)
```

## Benefits

### Reliability

- **Automatic retries**: SDKs handle transient errors automatically
- **Better error handling**: SDKs provide structured error types
- **API compatibility**: Official SDKs ensure compatibility with API updates

### Maintainability

- **Less code**: No manual HTTP request/response handling
- **Official support**: SDKs maintained by provider teams
- **Reduced bugs**: SDKs handle edge cases and error scenarios

### User Experience

- **Same behavior**: No changes to CLI behavior or output
- **Better reliability**: Automatic retries improve success rate
- **Clear errors**: SDK initialization failures provide clear error messages

## Limitations

- **Local provider**: Local provider implementation remains unchanged (no SDK replacement)
- **Streaming**: Streaming responses not used (out of scope)
- **SDK features**: Only chat completion APIs used, other SDK features not utilized

## Support

If you encounter issues:

1. Check error messages for SDK initialization failures
2. Verify API keys are correctly configured
3. Ensure SDK dependencies are installed (`go mod tidy`)
4. Check that Go version is 1.25.0 or later

For SDK-specific issues, refer to:
- OpenAI SDK: https://github.com/openai/openai-go
- Anthropic SDK: https://github.com/anthropics/anthropic-sdk-go
- Mistral SDK: https://github.com/Gage-Technologies/mistral-go
