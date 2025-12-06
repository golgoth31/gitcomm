# Quickstart: Ensure Config File Exists Before Reading

**Feature**: 015-ensure-config-exists
**Date**: 2025-01-27

## Overview

This feature modifies the `LoadConfig` function in `internal/config/config.go` to automatically create an empty config file if it doesn't exist before attempting to read it. This eliminates the need for users to manually create the config file.

## Key Changes

1. **File Existence Check**: Check if config file exists before reading
2. **Auto-Creation**: Create empty file (0 bytes) with 0600 permissions if missing
3. **Directory Creation**: Create parent directories recursively with 0755 permissions if missing
4. **Error Handling**: Clear error messages for permission/disk issues
5. **Logging**: Log file creation at debug level

## Implementation Steps

### Step 1: Add File Existence Check

Before calling `v.ReadInConfig()`, add a check for file existence:

```go
// Check if config file exists
if _, err := os.Stat(configPath); os.IsNotExist(err) {
    // File doesn't exist, create it
}
```

### Step 2: Create Parent Directories

If file doesn't exist, create parent directories first:

```go
configDir := filepath.Dir(configPath)
if err := os.MkdirAll(configDir, 0755); err != nil {
    return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
}
```

### Step 3: Validate Path (Not Directory)

Check that the path is not a directory:

```go
info, err := os.Stat(configPath)
if err == nil && info.Mode().IsDir() {
    return nil, fmt.Errorf("config path is a directory, not a file: %s", configPath)
}
```

### Step 4: Create Empty File

Create the empty file with 0600 permissions:

```go
file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
if err != nil {
    if os.IsExist(err) {
        // File was created by another process, treat as success
    } else {
        return nil, fmt.Errorf("failed to create config file: %w", err)
    }
} else {
    file.Close()
    utils.Logger.Debug().Str("path", configPath).Msg("Created config file")
}
```

### Step 5: Complete Implementation

The complete flow in `LoadConfig`:

```go
func LoadConfig(configPath string) (*Config, error) {
    v := viper.New()

    // Set default config path
    if configPath == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil, fmt.Errorf("failed to get home directory: %w", err)
        }
        configPath = filepath.Join(homeDir, ".gitcomm", "config.yaml")
    }

    // Validate path is not a directory
    if info, err := os.Stat(configPath); err == nil && info.Mode().IsDir() {
        return nil, fmt.Errorf("config path is a directory, not a file: %s", configPath)
    }

    // Check if config file exists, create if missing
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        // Create parent directories
        configDir := filepath.Dir(configPath)
        if err := os.MkdirAll(configDir, 0755); err != nil {
            return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
        }

        // Create empty file with 0600 permissions
        file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
        if err != nil {
            if os.IsExist(err) {
                // File was created by another process, treat as success
            } else {
                return nil, fmt.Errorf("failed to create config file: %w", err)
            }
        } else {
            file.Close()
            utils.Logger.Debug().Str("path", configPath).Msg("Created config file")
        }
    }

    // Configure viper (existing code)
    v.SetConfigFile(configPath)
    v.SetConfigType("yaml")
    v.SetEnvPrefix("GITCOMM")
    v.AutomaticEnv()

    // Read config file (existing code)
    if err := v.ReadInConfig(); err != nil {
        // Config file is optional, continue with defaults
    }

    // Rest of existing LoadConfig implementation...
}
```

## Testing

### Unit Tests

Create `internal/config/config_test.go` with tests for:

1. **File Creation**:
   ```go
   func TestLoadConfig_CreatesFileIfMissing(t *testing.T) {
       // Test that file is created when missing
   }
   ```

2. **Directory Creation**:
   ```go
   func TestLoadConfig_CreatesParentDirectories(t *testing.T) {
       // Test that parent directories are created
   }
   ```

3. **Permission Setting**:
   ```go
   func TestLoadConfig_SetsFilePermissions(t *testing.T) {
       // Test that file has 0600 permissions
   }
   ```

4. **Error Handling**:
   ```go
   func TestLoadConfig_HandlesErrors(t *testing.T) {
       // Test permission errors, disk space errors, etc.
   }
   ```

5. **Race Conditions**:
   ```go
   func TestLoadConfig_HandlesRaceConditions(t *testing.T) {
       // Test concurrent file creation
   }
   ```

### Integration Tests

Add to `test/integration/config_test.go`:

```go
func TestLoadConfig_Integration(t *testing.T) {
    // Test full LoadConfig flow with file creation
}
```

## Key Implementation Details

1. **File Permissions**: Use 0600 (owner read/write only) for security
2. **Directory Permissions**: Use 0755 (default) for parent directories
3. **Empty File**: Create 0-byte file (no content)
4. **Error Wrapping**: Use `fmt.Errorf()` with `%w` for error context
5. **Logging**: Use `utils.Logger.Debug()` for file creation events
6. **Race Conditions**: Handle `os.IsExist()` errors gracefully

## Files to Modify

- `internal/config/config.go` - Add file existence check and creation logic

## Files to Create

- `internal/config/config_test.go` - Unit tests for file creation

## Files to Update

- `test/integration/config_test.go` - Integration tests (if exists)

## Dependencies

- `os` - File system operations
- `path/filepath` - Path manipulation
- `github.com/rs/zerolog` - Logging (via `utils.Logger`)
- `github.com/spf13/viper` - Config reading (existing)

## Success Criteria

- ✅ Config file created automatically when missing
- ✅ File has 0600 permissions
- ✅ Parent directories created with 0755 permissions
- ✅ Clear error messages for failures
- ✅ File creation logged at debug level
- ✅ Race conditions handled gracefully
- ✅ All tests pass
