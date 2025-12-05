# Quick Start: OpenAI Responses API Migration

**Feature**: 010-openai-responses-api
**Date**: 2025-01-27

## Overview

This guide provides immediate usability guidance for developers working with the OpenAI provider migration to Responses API.

## What Changed

- **API Endpoint**: Changed from `/v1/chat/completions` to `/v1/responses`
- **Request Parameter**: Changed from `messages` to `input` (same structure)
- **Response Extraction**: Adapted to Responses API response structure
- **User Experience**: No changes - everything works the same

## Configuration

### Existing Configuration (Unchanged)

Your existing configuration file (`~/.gitcomm/config.yaml`) continues to work without modification:

```yaml
ai:
  default_provider: openai
  providers:
    openai:
      api_key: sk-your-api-key-here
      model: gpt-4-1  # Optional, defaults to gpt-4-1
      max_tokens: 500  # Optional, defaults to 500
      timeout: 30s     # Optional, defaults to 30s
```

**No changes required** - the migration is transparent to users.

## Usage

### Basic Usage (Unchanged)

```bash
# Navigate to your git repository
cd /path/to/your/repo

# Stage some changes
git add .

# Run gitcomm with OpenAI provider
gitcomm

# Follow the interactive prompts
# The system will use Responses API automatically
```

### With Custom Model

```bash
# Configure custom model in config.yaml
# providers.openai.model: gpt-4o

# Run gitcomm
gitcomm
# Uses Responses API with your specified model
```

## Troubleshooting

### API Key Issues

**Error**: `AI provider unavailable: API key invalid`

**Solution**:
1. Verify your API key in `~/.gitcomm/config.yaml`
2. Ensure the API key is valid and has access to Responses API
3. Check that the API key format is correct (starts with `sk-`)

### Model Not Available

**Error**: `AI provider unavailable: model not found`

**Solution**:
1. Verify the model name is correct in your configuration
2. Check if the model is available in Responses API
3. Try using the default model (`gpt-4-1`) first

### Timeout Issues

**Error**: `AI provider unavailable: timeout`

**Solution**:
1. Increase timeout in configuration: `timeout: 60s`
2. Check your network connection
3. Verify OpenAI API is accessible

### Empty Response

**Error**: `AI provider unavailable: empty response from API`

**Solution**:
1. Check API status at https://status.openai.com
2. Verify your API key has sufficient credits
3. Try again - may be a temporary API issue

## SDK vs Custom HTTP Client

The implementation automatically handles SDK availability:

- **If SDK v3 supports Responses API**: Uses SDK client methods
- **If SDK doesn't support it yet**: Falls back to custom HTTP client

**No action required** - the implementation chooses the best method automatically.

## Testing

### Unit Tests

```bash
# Run OpenAI provider unit tests
go test ./internal/ai/... -run TestOpenAI -v
```

### Integration Tests

```bash
# Run integration tests (requires API key)
export OPENAI_API_KEY=sk-your-key
go test ./test/integration/... -v
```

## Migration Notes

### For Developers

- **No breaking changes** to external interfaces
- **Internal implementation** changed (API endpoint and request format)
- **Error handling** remains identical
- **Test structure** remains the same (update mocks to Responses API format)

### For Users

- **No configuration changes** required
- **No behavior changes** - everything works the same
- **No migration steps** needed

## API Differences (Internal)

### Request Format

**Before (Chat Completions)**:
```json
{
  "model": "gpt-4-1",
  "messages": [...],
  "max_completion_tokens": 500
}
```

**After (Responses API)**:
```json
{
  "model": "gpt-4-1",
  "input": [...],  // Same structure as messages
  "max_completion_tokens": 500,
  "store": false   // Stateless mode
}
```

### Response Format

**Before**: `choices[0].message.content`
**After**: `content` or `text` field (structure may vary)

**Note**: Extraction logic handles both formats automatically.

## Support

If you encounter issues:

1. Check the troubleshooting section above
2. Verify your configuration is correct
3. Check OpenAI API status
4. Review error messages for specific guidance

## Next Steps

After migration:

1. **Test your workflow**: Run `gitcomm` with your typical commit workflow
2. **Verify behavior**: Ensure commit messages are generated correctly
3. **Report issues**: If you find any problems, report them with error details

The migration is designed to be transparent - if everything works as before, the migration is successful!
