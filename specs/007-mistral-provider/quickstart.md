# Quick Start: Add Mistral as AI Provider

**Feature**: 007-mistral-provider
**Date**: 2025-01-27

## Overview

This feature adds Mistral AI as a new provider option for generating commit messages. Mistral joins OpenAI and Anthropic as a first-class provider choice, offering users more flexibility in selecting their preferred AI service.

## Prerequisites

- Mistral AI API key (obtain from https://mistral.ai/)
- Go 1.25.0+ installed
- Existing gitcomm installation or ability to build from source

## Setup

### 1. Obtain Mistral API Key

1. Sign up for Mistral AI account at https://mistral.ai/
2. Navigate to API keys section
3. Create a new API key
4. Copy the API key (keep it secure)

### 2. Configure Mistral Provider

Edit your `~/.gitcomm/config.yaml` file (or create it if it doesn't exist):

```yaml
ai:
  default_provider: mistral  # Optional: set Mistral as default
  providers:
    mistral:
      api_key: ${MISTRAL_API_KEY}  # Use environment variable (recommended)
      model: mistral-large-latest  # Optional, default: mistral-large-latest
      timeout: 30s                 # Optional, default: 30s
```

**Security Best Practice**: Use environment variables for API keys:

```bash
export MISTRAL_API_KEY="your-api-key-here"
```

### 3. Available Mistral Models

- `mistral-tiny`: Fastest, least capable
- `mistral-small`: Good balance of speed and quality
- `mistral-medium`: Better quality
- `mistral-large-latest`: Best quality (recommended for commit messages)

## Usage

### Basic Usage (Default Provider)

If Mistral is set as default provider:

```bash
gitcomm
# Mistral will be used automatically
```

### Explicit Provider Selection

Select Mistral via CLI flag:

```bash
gitcomm --provider mistral
```

### Example Workflow

```bash
# 1. Make some changes
echo "New feature" > feature.txt
git add feature.txt

# 2. Run gitcomm with Mistral
gitcomm --provider mistral

# 3. Interactive session:
#    - AI usage? (y/n): y
#    - Review Mistral-generated message
#    - Accept, edit, or reject
```

## Configuration Examples

### Minimal Configuration

```yaml
ai:
  providers:
    mistral:
      api_key: ${MISTRAL_API_KEY}
```

### Full Configuration

```yaml
ai:
  default_provider: mistral
  providers:
    mistral:
      api_key: ${MISTRAL_API_KEY}
      model: mistral-large-latest
      timeout: 30s
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
    anthropic:
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-opus
```

### Multiple Providers

You can configure multiple providers and switch between them:

```bash
# Use Mistral
gitcomm --provider mistral

# Use OpenAI
gitcomm --provider openai

# Use Anthropic
gitcomm --provider anthropic
```

## Troubleshooting

### API Key Issues

**Error**: "Mistral API key not configured"

**Solution**:
- Verify `MISTRAL_API_KEY` environment variable is set
- Or check `api_key` field in config file
- Ensure API key is valid and not expired

### API Errors

**Error**: "API returned status 401"

**Solution**: API key is invalid or expired. Generate a new key from Mistral AI dashboard.

**Error**: "API returned status 429"

**Solution**: Rate limit exceeded. Wait a few moments and try again, or upgrade your Mistral plan.

**Error**: "API returned status 500"

**Solution**: Mistral API is experiencing issues. Try again later or use a different provider.

### Timeout Issues

**Error**: "context deadline exceeded"

**Solution**:
- Increase timeout in config: `timeout: 60s`
- Check network connectivity
- Mistral API may be slow, try again

### Model Not Found

**Error**: "API returned status 404"

**Solution**: Model name is incorrect. Use one of: `mistral-tiny`, `mistral-small`, `mistral-medium`, `mistral-large-latest`

## Integration

This feature integrates seamlessly with:
- Existing provider selection mechanism
- Token calculation infrastructure (uses character-based fallback)
- Configuration management
- Error handling and fallback mechanisms
- Commit message validation

## Comparison with Other Providers

| Feature | OpenAI | Anthropic | Mistral |
|---------|--------|----------|---------|
| Default Model | gpt-4 | claude-3-opus | mistral-large-latest |
| Token Calculation | tiktoken | Custom | Character-based |
| API Endpoint | api.openai.com | api.anthropic.com | api.mistral.ai |
| Authentication | Bearer token | x-api-key header | Bearer token |
| Request Format | Chat completions | Messages API | Chat completions (OpenAI-compatible) |

## Best Practices

1. **Use Environment Variables**: Never hardcode API keys in config files
2. **Set Appropriate Timeout**: 30s default is usually sufficient, increase for large repositories
3. **Choose Right Model**: Use `mistral-large-latest` for best commit message quality
4. **Monitor API Usage**: Check Mistral dashboard for usage and rate limits
5. **Have Fallback**: Configure multiple providers so you can switch if one is unavailable

## Next Steps

After configuring Mistral:
1. Test with a simple commit to verify setup
2. Compare generated messages with other providers
3. Adjust model or timeout if needed
4. Set as default provider if preferred
