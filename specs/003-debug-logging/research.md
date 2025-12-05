# Research: Debug Logging Configuration

**Feature**: 003-debug-logging
**Date**: 2025-01-27

## Technology Decisions

### 1. Zerolog Raw Text Format Configuration

**Decision**: Use `zerolog.ConsoleWriter` for human-readable structured text output format.

**Rationale**:
- `zerolog` provides `ConsoleWriter` which outputs human-readable text format
- Supports structured fields in key=value format
- Can be configured to disable timestamps
- Already a dependency, no new libraries needed
- Thread-safe and performant

**Alternatives Considered**:
- Custom formatter: More complex, unnecessary when ConsoleWriter exists
- Standard library `log` package: Less structured, doesn't support fields well
- `logrus`: Would require new dependency, zerolog is already used

**Implementation Pattern**:
```go
// For raw text format without timestamps
writer := zerolog.ConsoleWriter{
    Out: os.Stderr,
    NoColor: false, // Optional: disable colors
    TimeFormat: "", // Empty string disables timestamp
}
Logger = zerolog.New(writer).Level(zerolog.DebugLevel)
```

### 2. Debug Flag Implementation

**Decision**: Add `--debug` and `-d` flags using Cobra's flag system.

**Rationale**:
- Cobra is already used for CLI flags
- Standard pattern for CLI tools
- Supports both long and short forms
- Integrates with existing flag infrastructure

**Alternatives Considered**:
- Environment variable only: Less discoverable, flag is more user-friendly
- Config file only: Too complex for simple debug toggle
- Single flag only: Supporting both `--debug` and `-d` is more flexible

**Implementation Pattern**:
```go
rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
```

### 3. Logger Initialization Strategy

**Decision**: Modify `InitLogger` to accept debug flag parameter and configure logger based on flag state.

**Rationale**:
- Maintains existing function signature pattern
- Clear separation of concerns
- Easy to test
- Backward compatible (can add parameter without breaking)

**Alternatives Considered**:
- Separate InitDebugLogger function: Duplicates code, harder to maintain
- Global debug flag variable: Violates dependency injection principle
- Environment variable check: Less explicit, harder to test

**Implementation Pattern**:
```go
func InitLogger(verbose bool, debug bool) {
    if debug {
        // Configure for debug mode: raw text, no timestamp, DEBUG level
    } else {
        // Configure for silent mode: no output
    }
}
```

### 4. Log Level Migration Strategy

**Decision**: Update all logging calls throughout codebase to use `Logger.Debug()` instead of `Logger.Info()`, `Logger.Warn()`, `Logger.Error()`.

**Rationale**:
- FR-008 requires DEBUG-only log levels in codebase
- Simplifies logging configuration (only one level to manage)
- Consistent with feature requirement
- Makes codebase easier to maintain

**Alternatives Considered**:
- Keep multiple levels but filter: More complex, violates requirement
- Conditional level selection: Unnecessary complexity
- Separate logger instances: Over-engineered for this use case

**Implementation Pattern**:
```go
// Before
utils.Logger.Info().Msg("message")
utils.Logger.Warn().Err(err).Msg("warning")
utils.Logger.Error().Err(err).Msg("error")

// After
utils.Logger.Debug().Msg("message")
utils.Logger.Debug().Err(err).Msg("warning")
utils.Logger.Debug().Err(err).Msg("error")
```

### 5. Verbose Flag Handling

**Decision**: Make verbose flag a no-op when debug flag is present (debug flag takes precedence).

**Rationale**:
- Aligns with FR-007 clarification
- Simplifies logic (debug replaces verbose)
- Clear precedence rules
- Avoids confusion about which flag controls what

**Alternatives Considered**:
- Both flags work together: More complex, unclear behavior
- Verbose overrides debug: Contradicts requirement
- Error on both flags: Too strict, user-unfriendly

**Implementation Pattern**:
```go
// In InitLogger
if debug {
    // Enable debug logging
} else if verbose {
    // Do nothing (verbose is no-op without debug)
} else {
    // Silent mode (default)
}
```

## Best Practices

### Zerolog Configuration
- Use `ConsoleWriter` for human-readable output
- Set `TimeFormat` to empty string to disable timestamps
- Configure log level to `DebugLevel` when debug enabled
- Use `Disabled` level when debug disabled

### Flag Design
- Support both long (`--debug`) and short (`-d`) forms
- Place debug flag with other logging-related flags
- Use descriptive help text
- Follow existing flag naming conventions

### Code Migration
- Search for all `Logger.Info()`, `Logger.Warn()`, `Logger.Error()` calls
- Replace with `Logger.Debug()`
- Maintain existing log message content
- Update any level-specific logic

## Integration Points

### Existing Code
- `InitLogger` function: Extend signature to accept debug flag
- `main.go`: Add debug flag, pass to InitLogger
- All files using `utils.Logger`: Update to use Debug() method
- Error message display: Ensure separate from logging (already handled via fmt.Printf)

### Dependencies
- No new external dependencies required
- Uses existing `zerolog` capabilities
- Uses existing `cobra` flag system

## Performance Considerations

- Logger initialization: <10ms (negligible)
- Log message formatting: <1ms per message (ConsoleWriter is efficient)
- No performance impact when debug disabled (logger is disabled)
- Minimal memory overhead (logger instance is reused)

## Security Considerations

- No secrets in logs (existing constraint maintained)
- Debug logs may contain sensitive information (user responsibility)
- No new attack vectors introduced
- Error messages remain separate (not affected by logging changes)

## Testing Strategy

- Unit tests: Logger configuration with different flag combinations
- Integration tests: CLI flag behavior, output format verification
- Code review: Verify all logging calls use DEBUG level
- Manual testing: Verify output format and timestamp absence
