# Quickstart: Environment Variable Placeholder Substitution in Config Files

**Feature**: 014-env-var-placeholders
**Date**: 2025-01-27

## Overview

This feature enables environment variable placeholder substitution in config files. You can use `${ENV_VAR_NAME}` syntax in your config file, and the application will automatically replace these placeholders with values from your environment when loading configuration.

## Key Changes

1. **Placeholder Syntax**: Use `${ENV_VAR_NAME}` in config file values
2. **Automatic Substitution**: Placeholders are replaced with environment variable values during config loading
3. **Validation**: Invalid placeholder syntax or missing environment variables cause immediate exit with clear error messages
4. **Backward Compatibility**: Config files without placeholders continue to work unchanged

## Usage Examples

### Basic Usage

1. **Set environment variable**:
   ```bash
   export OPENAI_API_KEY="sk-12345"
   ```

2. **Use placeholder in config file** (`~/.gitcomm/config.yaml`):
   ```yaml
   ai:
     default_provider: openai
     providers:
       openai:
         api_key: ${OPENAI_API_KEY}
         model: gpt-4
   ```

3. **Run application**: Placeholder is automatically replaced with `sk-12345` when config is loaded

### Multiple Placeholders

```yaml
ai:
  default_provider: openai
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
    anthropic:
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-opus
```

**Environment variables required**:
- `OPENAI_API_KEY`
- `ANTHROPIC_API_KEY`

### Nested YAML Values

Placeholders work in nested YAML structures:

```yaml
ai:
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}
      endpoint: ${OPENAI_ENDPOINT:-https://api.openai.com}
```

### Empty String Values

Empty string values are treated as valid:

```bash
export OPTIONAL_VAR=""
```

```yaml
optional_setting: ${OPTIONAL_VAR}  # Substitutes to empty string
```

## Error Handling

### Missing Environment Variable

If a required environment variable is not set, the application exits immediately:

```bash
$ gitcomm
Error: missing environment variables: OPENAI_API_KEY, ANTHROPIC_API_KEY
```

### Invalid Placeholder Syntax

Invalid placeholder syntax causes immediate exit:

```yaml
# Invalid: spaces in variable name
api_key: ${VAR NAME}

# Invalid: nested placeholder
api_key: ${${NESTED}}

# Invalid: multiline placeholder
api_key: ${VAR
NAME}
```

**Error message**:
```
Error: invalid placeholder syntax: ${VAR NAME}
```

## Best Practices

1. **Use descriptive variable names**: `OPENAI_API_KEY` is clearer than `KEY1`
2. **Document required variables**: List required environment variables in README
3. **Use consistent naming**: Follow your project's environment variable naming convention
4. **Keep secrets out of config files**: Use placeholders for sensitive values like API keys
5. **Test with missing variables**: Verify error messages are clear and helpful

## Migration Guide

### Existing Config Files

Existing config files without placeholders continue to work unchanged. No migration required.

### Adding Placeholders

To add placeholders to existing config files:

1. **Identify values to externalize**: API keys, endpoints, tokens, etc.
2. **Replace values with placeholders**: Change `api_key: sk-12345` to `api_key: ${OPENAI_API_KEY}`
3. **Set environment variables**: Export required variables before running application
4. **Test**: Verify config loads correctly with placeholders

### Example Migration

**Before**:
```yaml
ai:
  providers:
    openai:
      api_key: sk-12345abc  # Hardcoded secret
```

**After**:
```yaml
ai:
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}  # From environment
```

**Environment setup**:
```bash
export OPENAI_API_KEY="sk-12345abc"
```

## Testing

### Manual Testing

1. **Test valid substitution**:
   ```bash
   export TEST_VAR="test-value"
   # Config: test: ${TEST_VAR}
   # Expected: test: test-value
   ```

2. **Test missing variable**:
   ```bash
   unset MISSING_VAR
   # Config: test: ${MISSING_VAR}
   # Expected: Error exit with "missing environment variables: MISSING_VAR"
   ```

3. **Test invalid syntax**:
   ```bash
   # Config: test: ${VAR NAME}
   # Expected: Error exit with "invalid placeholder syntax: ${VAR NAME}"
   ```

### Integration Testing

Run the application with various config file scenarios:

```bash
# Valid config with placeholders
export OPENAI_API_KEY="sk-test"
gitcomm  # Should work

# Missing variable
unset OPENAI_API_KEY
gitcomm  # Should exit with error

# Invalid syntax in config
# Edit config file to include invalid placeholder
gitcomm  # Should exit with error
```

## Troubleshooting

### Placeholder Not Substituted

**Possible causes**:
- Environment variable not set (check with `echo $VAR_NAME`)
- Invalid placeholder syntax (check for spaces, nested placeholders)
- Placeholder in comment (comments are ignored)

**Solution**: Verify environment variable is set and placeholder syntax is correct

### Application Exits Immediately

**Possible causes**:
- Missing required environment variable
- Invalid placeholder syntax

**Solution**: Check error message for specific variable name or syntax issue

### Empty String Substitution

**Behavior**: Empty string values are valid substitutions. If you need to distinguish unset from empty, use a different approach (e.g., default values in config).

## Key Implementation Details

1. **Processing order**: Placeholders are processed before YAML parsing
2. **Comment handling**: Placeholders in YAML comments are ignored
3. **Validation**: All placeholders are validated before any substitution
4. **Error reporting**: All missing variables are reported in a single error message
5. **Performance**: Substitution completes in under 10ms for typical config files
