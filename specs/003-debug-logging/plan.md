# Implementation Plan: Debug Logging Configuration

**Branch**: `003-debug-logging` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-debug-logging/spec.md`

## Summary

This feature changes the CLI logging behavior to only output log messages when debug mode is enabled via a debug flag. When enabled, logs are displayed in human-readable structured text format (not JSON) without timestamps. The implementation will:

1. **Add debug flag** (`--debug` or `-d`) to enable debug logging
2. **Modify logger initialization** to support raw text format without timestamps
3. **Change default behavior** to suppress all logging unless debug flag is set
4. **Update all logging calls** throughout codebase to use DEBUG level only
5. **Handle verbose flag** by making it a no-op when debug flag is present

The technical approach extends the existing `zerolog` logger configuration to support raw text output format and conditional logging based on debug flag. All existing logging statements must be converted to use DEBUG level.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `github.com/rs/zerolog` - Structured logging (existing, needs configuration changes)
- `github.com/spf13/cobra` - CLI framework (existing, needs flag addition)
- `os` - Standard library for output handling

**Storage**: N/A (logging configuration is in-memory only)

**Testing**:
- Standard Go testing framework (`testing` package)
- `github.com/onsi/ginkgo/v2` and `github.com/onsi/gomega` for BDD-style tests (existing)
- Unit tests for logger configuration
- Integration tests for CLI flag behavior

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Logger initialization completes within 10ms
- Log message output has minimal performance impact (<1ms per message)
- No noticeable delay in CLI startup

**Constraints**:
- Must maintain backward compatibility (existing verbose flag behavior changes)
- Must not break existing error message display
- All logging statements must use DEBUG level only
- No global state changes (logger is already global, but configuration is acceptable)

**Scale/Scope**:
- Single CLI binary
- All logging statements in codebase need review/update
- Typical log volume: low (only when debug enabled)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Extends existing `internal/utils` layer with logger configuration
  - No new layers required, fits existing structure
  - Logger configuration is utility-level concern

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - Logger interface already exists (zerolog.Logger)
  - No new interfaces needed
  - Configuration passed via function parameters

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: Write tests for logger configuration first
  - Unit tests for InitLogger with debug flag
  - Integration tests for CLI flag behavior
  - Table-driven tests for format variations

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Uses standard library where possible
  - No panics in library code

- **Error Handling**: ✅ **COMPLIANT**
  - Logger initialization errors handled explicitly
  - No error types needed (logger setup is simple)

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - Logger is thread-safe (zerolog handles this)
  - No shared mutable state introduced

- **Technical Constraints**: ✅ **COMPLIANT**
  - Logger is global (acceptable for logging utility)
  - No resource cleanup needed
  - Configuration is stateless

- **Operational Constraints**: ✅ **COMPLIANT**
  - Logging strategy defined (raw text, no timestamps, DEBUG only)
  - No secrets in logs (existing constraint maintained)

**Violations**: None. All principles are satisfied.

### Post-Design Constitution Check

After Phase 1 design completion, all principles remain satisfied:

- **Clean Architecture**: ✅ Logger configuration remains in `internal/utils` layer
- **Interface-Driven Development**: ✅ No new interfaces needed, uses existing logger interface
- **Test-First Development**: ✅ Test strategy defined for logger configuration
- **Idiomatic Go**: ✅ Uses standard zerolog patterns and Go conventions
- **Error Handling**: ✅ No new error types needed, existing error handling sufficient
- **Context & Thread Safety**: ✅ Logger is thread-safe, no new concurrency concerns
- **Technical Constraints**: ✅ No global state changes (logger is already global, acceptable)
- **Operational Constraints**: ✅ Logging strategy fully defined and documented

## Project Structure

### Documentation (this feature)

```text
specs/003-debug-logging/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
internal/
├── utils/
│   └── logger.go         # Modify InitLogger to support debug flag and raw text format
└── (no other changes needed)

cmd/gitcomm/
└── main.go                # Add debug flag, pass to InitLogger

(all other files)
└── (review and update logging calls to use DEBUG level only)
```

**Structure Decision**: Minimal changes to existing structure. Primary modifications:
- `internal/utils/logger.go` - Logger initialization and configuration
- `cmd/gitcomm/main.go` - Debug flag addition
- All files with logging statements - Convert to DEBUG level

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all principles satisfied.
