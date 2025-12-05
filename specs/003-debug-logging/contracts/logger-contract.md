# Logger Interface Contract: Debug Logging Configuration

**Feature**: 003-debug-logging
**Date**: 2025-01-27
**Interface**: `InitLogger` function

## Extended Interface

### InitLogger

Initializes the global logger with configuration based on debug and verbose flags.

```go
InitLogger(verbose bool, debug bool)
```

**Input**:
- `verbose bool` - Existing verbose flag (becomes no-op when debug is true)
- `debug bool` - New debug flag (enables debug logging)

**Output**:
- None (configures global `Logger` variable)

**Behavior**:
- If `debug` is true: Configure logger for debug mode (raw text format, no timestamps, DEBUG level enabled)
- If `debug` is false and `verbose` is true: Configure logger for silent mode (verbose is no-op)
- If both are false: Configure logger for silent mode (no output)
- Logger configuration is applied to global `Logger` variable

**Error Cases**:
- None (logger initialization is simple, uses default zerolog behavior)

---

## Logger Usage Contract

### Debug Method

All logging statements MUST use `Logger.Debug()` method exclusively.

```go
// Correct usage
utils.Logger.Debug().Msg("message")
utils.Logger.Debug().Err(err).Msg("error occurred")
utils.Logger.Debug().Str("key", "value").Msg("message with field")

// Incorrect usage (must be changed)
utils.Logger.Info().Msg("message")  // ❌ Must use Debug()
utils.Logger.Warn().Msg("warning")  // ❌ Must use Debug()
utils.Logger.Error().Msg("error")   // ❌ Must use Debug()
```

**Behavior**:
- When debug mode is enabled: Messages are output in raw text format
- When debug mode is disabled: Messages are suppressed (no output)
- Format: `[DEBUG] message key=value key2=value2` (no timestamp)

---

## Implementation Contract

### Thread Safety

- Logger is thread-safe (zerolog handles this internally)
- No additional synchronization needed
- Safe to call from multiple goroutines

### Error Handling

- Logger initialization never fails (uses default zerolog behavior)
- No error return needed
- Logging calls never fail (output is best-effort)

### Performance

- Logger initialization: <10ms
- Log message formatting: <1ms per message
- No performance impact when debug disabled

### Resource Management

- Logger uses stderr for output
- No file handles or network connections
- No cleanup needed

---

## Testing Contract

### Unit Tests

- Test InitLogger with debug=true, verbose=false
- Test InitLogger with debug=false, verbose=true
- Test InitLogger with debug=true, verbose=true (debug takes precedence)
- Test InitLogger with debug=false, verbose=false (silent mode)
- Test log output format (raw text, no timestamp)
- Test log suppression when debug disabled

### Integration Tests

- Test CLI flag parsing and logger initialization
- Test actual log output format
- Test that verbose flag is no-op when debug is set

### Mock Requirements

- No mocks needed (logger is global, can be tested directly)
