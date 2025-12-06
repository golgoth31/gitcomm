# Research: Ensure Config File Exists Before Reading

**Feature**: 015-ensure-config-exists
**Date**: 2025-01-27

## Technology Decisions

### 1. File Existence Check Pattern

**Decision**: Use `os.Stat()` to check file existence before attempting to read, then create if missing.

**Rationale**:
- `os.Stat()` is the standard Go way to check file existence
- Returns `os.ErrNotExist` when file doesn't exist, which can be checked with `os.IsNotExist()`
- More efficient than attempting to read and catching errors
- Clear separation between existence check and read operation

**Alternatives Considered**:
- Try to read first, catch error: Less efficient, mixes concerns
- Use `os.Open()` with O_CREATE flag: Creates file even if exists, not what we want
- Use `filepath.Walk()`: Overkill for single file check

**Implementation Pattern**:
```go
if _, err := os.Stat(configPath); os.IsNotExist(err) {
    // File doesn't exist, create it
}
```

### 2. Empty File Creation Pattern

**Decision**: Create empty file using `os.Create()` or `os.OpenFile()` with `O_CREATE|O_WRONLY|O_TRUNC` flags, then immediately close.

**Rationale**:
- `os.Create()` creates file if it doesn't exist, truncates if it exists (but we check first)
- Creates file with 0 bytes (empty) as required
- Simple and idiomatic Go
- File handle must be closed immediately after creation

**Alternatives Considered**:
- `ioutil.WriteFile()` with empty content: Deprecated in Go 1.16+, use `os.WriteFile()` instead
- `os.WriteFile()` with empty slice: Works but less explicit than Create
- `os.OpenFile()` with explicit flags: More verbose, same result

**Implementation Pattern**:
```go
file, err := os.Create(configPath)
if err != nil {
    return fmt.Errorf("failed to create config file: %w", err)
}
file.Close() // Close immediately, file is empty
```

### 3. File Permissions (0600) Pattern

**Decision**: Use `os.Chmod()` after file creation to set 0600 permissions, or use `os.OpenFile()` with explicit permission parameter.

**Rationale**:
- `os.Create()` uses default permissions (typically 0666), need to change to 0600
- `os.Chmod()` is explicit and clear
- Can also use `os.OpenFile()` with `os.FileMode(0600)` parameter
- Must set permissions before closing file handle

**Alternatives Considered**:
- Rely on umask: Not portable, doesn't guarantee 0600
- Use `os.WriteFile()` with mode: Works but file is empty, no write needed
- Set permissions in separate step: More explicit, easier to test

**Implementation Pattern**:
```go
// Option 1: Create then chmod
file, err := os.Create(configPath)
if err != nil {
    return fmt.Errorf("failed to create config file: %w", err)
}
if err := os.Chmod(configPath, 0600); err != nil {
    file.Close()
    return fmt.Errorf("failed to set config file permissions: %w", err)
}
file.Close()

// Option 2: OpenFile with explicit mode
file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
if err != nil {
    return fmt.Errorf("failed to create config file: %w", err)
}
file.Close()
```

### 4. Parent Directory Creation Pattern

**Decision**: Use `os.MkdirAll()` to create parent directories recursively with 0755 permissions.

**Rationale**:
- `os.MkdirAll()` creates all necessary parent directories in one call
- Handles existing directories gracefully (no error if already exists)
- Can specify permissions directly (0755)
- Standard Go pattern for directory creation
- Must be called before file creation

**Alternatives Considered**:
- `os.Mkdir()` in loop: More complex, error-prone
- Check each directory level: Unnecessary, `MkdirAll` handles it
- Use `filepath.Dir()` to get parent, then create: `MkdirAll` does this internally

**Implementation Pattern**:
```go
configDir := filepath.Dir(configPath)
if err := os.MkdirAll(configDir, 0755); err != nil {
    return fmt.Errorf("failed to create config directory: %w", err)
}
```

### 5. Error Handling for File Operations

**Decision**: Wrap all file operation errors with context using `fmt.Errorf()` with `%w` verb, providing clear error messages.

**Rationale**:
- Wrapped errors preserve original error for debugging
- Context helps identify which operation failed
- Clear error messages help users understand what went wrong
- Follows Go error handling best practices
- Enables error inspection with `errors.Is()` and `errors.As()`

**Alternatives Considered**:
- Return raw errors: Less helpful for users
- Custom error types: Unnecessary complexity for this feature
- Panic on errors: Violates constitution (no panics in library code)

**Implementation Pattern**:
```go
if err := os.MkdirAll(configDir, 0755); err != nil {
    return fmt.Errorf("failed to create config directory %s: %w", configDir, err)
}
```

### 6. Race Condition Handling

**Decision**: Rely on OS-level atomicity of file creation operations. Handle `os.ErrExist` errors gracefully if file is created between check and creation.

**Rationale**:
- File creation is atomic at OS level
- If two processes try to create simultaneously, one succeeds, other gets `os.ErrExist`
- Check for `os.IsExist()` after creation attempt
- If file exists after creation attempt, treat as success (another process created it)
- No need for file locking (out of scope per spec)

**Alternatives Considered**:
- File locking: Out of scope, adds complexity
- Retry logic: Unnecessary, OS handles atomicity
- Check-then-create pattern: Race condition possible, but handled gracefully

**Implementation Pattern**:
```go
file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
if err != nil {
    if os.IsExist(err) {
        // File was created by another process, treat as success
        return nil
    }
    return fmt.Errorf("failed to create config file: %w", err)
}
file.Close()
```

### 7. Logging File Creation

**Decision**: Use `utils.Logger.Debug()` to log file creation events with file path.

**Rationale**:
- Follows existing logging pattern in codebase
- Debug level appropriate for operational visibility
- Logs only when debug mode enabled (per existing logger configuration)
- Provides troubleshooting information without cluttering normal output

**Alternatives Considered**:
- No logging: Reduces visibility for troubleshooting
- Info level: Too verbose for routine operation
- Separate logger: Unnecessary, existing logger sufficient

**Implementation Pattern**:
```go
utils.Logger.Debug().Str("path", configPath).Msg("Created config file")
```

### 8. Path Validation

**Decision**: Validate that config path is not a directory before attempting file operations using `os.Stat()` and `os.ModeType` check.

**Rationale**:
- FR-005 requires error if path points to existing directory
- `os.Stat()` returns file info including mode
- Check `info.Mode().IsDir()` to detect directory
- Must check before file creation to provide clear error

**Alternatives Considered**:
- Try to create and catch error: Less clear error message
- Use `filepath.Dir()` and check: Doesn't detect if path itself is directory
- Separate validation function: More explicit, easier to test

**Implementation Pattern**:
```go
info, err := os.Stat(configPath)
if err == nil && info.Mode().IsDir() {
    return fmt.Errorf("config path is a directory, not a file: %s", configPath)
}
```

## Best Practices Summary

1. **Check before create**: Use `os.Stat()` to check existence, handle `os.IsNotExist()`
2. **Create directories first**: Use `os.MkdirAll()` for parent directories before file creation
3. **Set permissions explicitly**: Use `os.Chmod()` or `os.OpenFile()` with mode parameter
4. **Close file handles**: Always close file handles immediately after creation
5. **Handle race conditions**: Check for `os.IsExist()` after creation, treat as success
6. **Wrap errors with context**: Use `fmt.Errorf()` with `%w` for error wrapping
7. **Validate paths**: Check if path is directory before file operations
8. **Log operations**: Use debug-level logging for operational visibility
