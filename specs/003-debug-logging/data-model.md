# Data Model: Debug Logging Configuration

**Feature**: 003-debug-logging
**Date**: 2025-01-27

## Domain Entities

### LogConfiguration

Represents the logging configuration state for the CLI.

**Fields**:
- `DebugEnabled bool` - Whether debug logging is currently enabled
- `Format string` - Output format ("raw" or "json")
- `TimestampEnabled bool` - Whether timestamps are included in output
- `Level string` - Log level ("debug", "info", "warn", "error", "disabled")

**Methods**:
- `IsDebugMode() bool` - Returns true if debug logging is enabled
- `ShouldOutput() bool` - Returns true if any log output should be produced
- `GetFormat() string` - Returns the output format

**Validation Rules**:
- `Format` must be one of: "raw", "json"
- `Level` must be one of: "debug", "info", "warn", "error", "disabled"
- When `DebugEnabled` is false, `Level` must be "disabled"
- When `DebugEnabled` is true, `Level` must be "debug" and `Format` must be "raw"

**Lifecycle**:
1. Created at CLI startup based on debug flag
2. Used throughout CLI execution to determine logging behavior
3. Discarded at CLI exit

---

### DebugFlagState

Represents the state of the debug flag from command-line arguments.

**Fields**:
- `Enabled bool` - Whether debug flag was provided
- `Source string` - How flag was provided ("--debug", "-d", or "none")

**Methods**:
- `IsSet() bool` - Returns true if debug flag was provided

**Validation Rules**:
- `Source` must be one of: "--debug", "-d", "none"
- If `Enabled` is true, `Source` must not be "none"

**Lifecycle**:
1. Created when parsing command-line arguments
2. Passed to logger initialization
3. Discarded after logger initialization

---

## Relationships

- `DebugFlagState` is used to create `LogConfiguration`
- `LogConfiguration` controls behavior of `zerolog.Logger` instance
- All logging statements use `Logger.Debug()` method (single log level)

---

## State Transitions

### Logger Initialization Workflow

```
[CLI Startup]
  ↓
[Parse Command-Line Flags]
  ↓
[Debug Flag Set?] → No → [Create LogConfiguration: DebugEnabled=false, Level=disabled]
  ↓ Yes
[Create LogConfiguration: DebugEnabled=true, Format=raw, TimestampEnabled=false, Level=debug]
  ↓
[Initialize Logger with Configuration]
  ↓
[Logger Ready]
```

### Logging Output Workflow

```
[Log Statement Executed]
  ↓
[Check LogConfiguration.ShouldOutput()]
  ↓
[Should Output?] → No → [Skip Output]
  ↓ Yes
[Format Message in Raw Text Format]
  ↓
[Output to stderr]
```

---

## Data Flow

1. **Flag Parsing**:
   - Command-line arguments → `DebugFlagState`

2. **Logger Initialization**:
   - `DebugFlagState` → `LogConfiguration` → `zerolog.Logger` configuration

3. **Logging Execution**:
   - Code calls `Logger.Debug()` → `LogConfiguration` checked → Message formatted → Output

---

## Persistence

- **No Persistent Storage**: Configuration is in-memory only, created at CLI startup
- **No Configuration File**: Debug flag is command-line only (no persistent config needed)
- **No Database**: No persistent storage required (stateless CLI)

---

## Error Types

- No new error types needed (logger initialization is simple, uses existing error handling)

---

## Integration with Existing Models

### Logger (existing)
- Global `zerolog.Logger` instance in `internal/utils/logger.go`
- Extended with new configuration options
- All methods remain the same, only configuration changes

### CommitOptions (existing)
- No changes needed (debug flag is separate from commit options)
