# Contract: LoadConfig with Environment Variable Placeholder Substitution

**Feature**: 014-env-var-placeholders
**Date**: 2025-01-27
**Component**: `internal/config`

## Function Signature

```go
func LoadConfig(configPath string) (*Config, error)
```

## Preconditions

- `configPath` is a valid file path (or empty string for default path)
- If `configPath` is empty, home directory is accessible
- Config file exists (or will be created by previous feature 015)
- Config file is readable
- Environment variables referenced in placeholders may or may not be set

## Postconditions

### Success Case

- Config file content is read successfully
- All placeholders in config file are identified
- All placeholder syntax is validated (alphanumeric/underscores only, no nested, no multiline)
- All referenced environment variables are found
- All placeholders are replaced with environment variable values
- Substituted content is valid YAML
- Config is loaded successfully
- Returns `*Config` with all values populated
- Returns `nil` error

### Failure Cases

#### Invalid Placeholder Syntax

- **Condition**: Config file contains invalid placeholder syntax (spaces, nested placeholders, multiline, invalid characters)
- **Behavior**: Function returns `nil` config and error describing invalid syntax
- **Error Message**: Must identify the invalid placeholder pattern
- **Exit**: Application exits immediately (handled by caller)

#### Missing Environment Variable

- **Condition**: Config file contains placeholder for environment variable that is not set
- **Behavior**: Function returns `nil` config and error listing all missing variables
- **Error Message**: Must list all missing environment variable names
- **Exit**: Application exits immediately (handled by caller)

#### File Read Error

- **Condition**: Config file cannot be read
- **Behavior**: Function returns `nil` config and wrapped error
- **Error Message**: Must include file path and original error

#### YAML Parse Error

- **Condition**: Substituted content is invalid YAML
- **Behavior**: Function returns `nil` config and viper error
- **Error Message**: Standard viper YAML parsing error

## Invariants

- Config file on disk is never modified (substitution happens in memory)
- YAML comments are preserved unchanged
- Non-placeholder content is preserved exactly as written
- Environment variable values are substituted verbatim (no escaping or transformation)
- Empty string values from environment variables are valid substitutions

## Side Effects

- None (pure function with respect to file system)
- Environment variable access (read-only, no modification)

## Test Requirements

### Unit Tests

1. **Test placeholder identification**: Verify regex correctly identifies valid placeholders
2. **Test syntax validation**: Verify invalid syntax is rejected (spaces, nested, multiline)
3. **Test variable lookup**: Verify `os.LookupEnv()` correctly distinguishes unset from empty
4. **Test substitution**: Verify placeholders are replaced with correct values
5. **Test multiple occurrences**: Verify same placeholder appearing multiple times is replaced consistently
6. **Test comment handling**: Verify placeholders in comments are ignored
7. **Test error messages**: Verify error messages clearly identify missing variables

### Integration Tests

1. **Test full LoadConfig flow**: Verify end-to-end placeholder substitution works
2. **Test missing variable**: Verify application exits with clear error when variable missing
3. **Test invalid syntax**: Verify application exits with clear error when syntax invalid
4. **Test backward compatibility**: Verify config files without placeholders work unchanged
5. **Test nested YAML structures**: Verify placeholders work in nested YAML values

## Performance Requirements

- Placeholder substitution completes in under 10ms for typical config files
- No noticeable delay in LoadConfig execution
- Efficient regex matching (compile once, reuse for multiple matches)

## Security Considerations

- Environment variable values are not logged
- Error messages identify variable names but not values
- No file system modifications (substitution in memory only)
- Config file permissions remain unchanged (handled by feature 015)
