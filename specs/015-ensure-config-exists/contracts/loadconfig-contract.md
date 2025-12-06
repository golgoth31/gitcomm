# LoadConfig Function Contract: Ensure Config File Exists

**Feature**: 015-ensure-config-exists
**Date**: 2025-01-27
**Function**: `LoadConfig` in `internal/config/config.go`

## Function Signature

```go
func LoadConfig(configPath string) (*Config, error)
```

## Input Contract

### Parameters

- **`configPath string`**
  - **Type**: `string`
  - **Required**: No (empty string uses default path)
  - **Default**: `~/.gitcomm/config.yaml` (when empty)
  - **Validation**:
    - Must not point to an existing directory (FR-005)
    - Path must be valid for the operating system
  - **Behavior**:
    - If empty, resolves to `$HOME/.gitcomm/config.yaml`
    - If provided, used as-is (no resolution)

## Output Contract

### Return Values

- **`*Config`**
  - **Type**: `*Config` (pointer to Config struct)
  - **Description**: Application configuration loaded from file or environment variables
  - **Nil Handling**: Never nil on success (always returns valid Config)
  - **Content**:
    - Loaded from config file if exists and readable
    - Loaded from environment variables (GITCOMM_* prefix)
    - Defaults used for missing values

- **`error`**
  - **Type**: `error`
  - **Nil Handling**: Nil on success, non-nil on failure
  - **Error Types**:
    - `fmt.Errorf("failed to get home directory: %w", err)` - Home directory resolution failure
    - `fmt.Errorf("config path is a directory, not a file: %s", configPath)` - Path is directory (FR-005)
    - `fmt.Errorf("failed to create config directory %s: %w", dir, err)` - Directory creation failure (FR-004)
    - `fmt.Errorf("failed to create config file: %w", err)` - File creation failure (FR-002)
    - `fmt.Errorf("failed to set config file permissions: %w", err)` - Permission setting failure (FR-010)
    - Viper read errors (wrapped) - Config file read failure (FR-006)

## Behavior Contract

### Pre-Conditions

1. **File Existence Check** (FR-001):
   - System checks if config file exists at `configPath` (or default path)
   - Uses `os.Stat()` to check existence

2. **Directory Validation** (FR-005):
   - If path exists and is a directory, return error immediately
   - Check performed before file creation

### Main Behavior

1. **Path Resolution**:
   - If `configPath` is empty, resolve to `$HOME/.gitcomm/config.yaml`
   - If `configPath` is provided, use as-is

2. **File Creation** (if missing, FR-002):
   - Create parent directories recursively with 0755 permissions (FR-004)
   - Create empty file (0 bytes) with 0600 permissions (FR-010)
   - Log creation at debug level (FR-011)
   - Handle race conditions gracefully (FR-007)

3. **Config Loading**:
   - Read config file using viper (existing behavior)
   - Load from environment variables (GITCOMM_* prefix)
   - Merge file and environment config
   - Return Config struct with loaded values

### Post-Conditions

1. **File State**:
   - Config file exists at expected path (SC-001)
   - File has 0600 permissions (FR-010)
   - File is empty (0 bytes) if newly created (FR-002)

2. **Return Value**:
   - `*Config` is non-nil on success
   - `error` is nil on success
   - `error` contains context on failure

## Error Handling Contract

### Error Categories

1. **Home Directory Errors** (FR-009):
   - **Condition**: `os.UserHomeDir()` fails
   - **Error**: `fmt.Errorf("failed to get home directory: %w", err)`
   - **Recovery**: Cannot recover, return error immediately

2. **Path Validation Errors** (FR-005):
   - **Condition**: Config path points to existing directory
   - **Error**: `fmt.Errorf("config path is a directory, not a file: %s", configPath)`
   - **Recovery**: Cannot recover, return error immediately

3. **Directory Creation Errors** (FR-004):
   - **Condition**: `os.MkdirAll()` fails (permissions, disk space, etc.)
   - **Error**: `fmt.Errorf("failed to create config directory %s: %w", dir, err)`
   - **Recovery**: Cannot recover, return error immediately

4. **File Creation Errors** (FR-002):
   - **Condition**: `os.Create()` or `os.OpenFile()` fails
   - **Error**: `fmt.Errorf("failed to create config file: %w", err)`
   - **Recovery**: Cannot recover, return error immediately
   - **Special Case**: If `os.IsExist(err)`, treat as success (race condition, FR-007)

5. **Permission Setting Errors** (FR-010):
   - **Condition**: `os.Chmod()` fails
   - **Error**: `fmt.Errorf("failed to set config file permissions: %w", err)`
   - **Recovery**: Close file handle, return error

6. **File Read Errors** (FR-006):
   - **Condition**: Viper `ReadInConfig()` fails (but file exists)
   - **Error**: Viper error (wrapped)
   - **Recovery**: Return error, but file remains created

## Thread Safety Contract

- **Concurrent Calls**: Safe for concurrent calls
- **Race Conditions**: Handled gracefully (FR-007)
  - If file created between check and creation, treat as success
  - Use `os.O_EXCL` flag or check `os.IsExist()` after creation
- **File Operations**: Atomic at OS level
- **No Shared State**: No shared mutable state introduced

## Logging Contract

- **File Creation Logging** (FR-011):
  - **Level**: Debug/Info
  - **Message**: "Created config file"
  - **Fields**: `path` (config file path)
  - **Format**: `utils.Logger.Debug().Str("path", configPath).Msg("Created config file")`
  - **Condition**: Only logged when file is actually created (not when it already exists)

## Performance Contract

- **File Creation Time**: < 100ms (SC-002)
- **Existence Check**: < 1ms (minimal overhead)
- **Overall Impact**: No noticeable delay in LoadConfig execution

## Backward Compatibility

- **Function Signature**: Unchanged
- **Return Types**: Unchanged
- **Behavior**: Enhanced (creates file if missing, but maintains existing read behavior)
- **Error Types**: New errors possible, but existing error handling still works

## Testing Contract

### Unit Tests Required

1. File creation when missing
2. File existence check
3. Directory creation (recursive)
4. Permission setting (0600)
5. Error handling (permissions, disk space, invalid path)
6. Race condition handling
7. Path validation (directory check)

### Integration Tests Required

1. LoadConfig with missing file (creates file)
2. LoadConfig with existing file (no modification)
3. LoadConfig with custom path
4. LoadConfig error scenarios (read-only directory, etc.)
