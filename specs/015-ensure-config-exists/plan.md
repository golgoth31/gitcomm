# Implementation Plan: Ensure Config File Exists Before Reading

**Branch**: `015-ensure-config-exists` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/015-ensure-config-exists/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature modifies the `LoadConfig` function to ensure the config file exists before attempting to read it. If the file doesn't exist, the system will create an empty file (0 bytes) with restrictive permissions (0600) and create any necessary parent directories with default permissions (0755). The implementation will:

1. **Check file existence** before reading in `LoadConfig` function
2. **Create empty config file** (0 bytes) if missing, with 0600 permissions
3. **Create parent directories** recursively with 0755 permissions if needed
4. **Handle errors gracefully** with clear error messages for permission/disk issues
5. **Log file creation** at debug/info level for operational visibility

The technical approach uses Go's standard library (`os`, `path/filepath`) for file system operations, ensuring thread-safe file creation and proper error handling. The implementation maintains backward compatibility while adding the auto-creation behavior.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- `os` - Standard library for file system operations (file creation, directory creation, permission setting)
- `path/filepath` - Standard library for path manipulation
- `github.com/spf13/viper` - Configuration management (existing, used for reading config)
- `github.com/rs/zerolog` - Structured logging (existing, used for logging file creation)

**Storage**: File system (YAML config file at `~/.gitcomm/config.yaml` or custom path)

**Testing**:
- Standard Go testing framework (`testing` package)
- Unit tests for file creation logic
- Integration tests for LoadConfig behavior with missing/existing files
- Table-driven tests for error scenarios

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Config file creation completes in under 100ms (SC-002)
- No noticeable delay in LoadConfig execution
- File existence check has minimal overhead (<1ms)

**Constraints**:
- Must maintain backward compatibility (existing LoadConfig behavior preserved)
- Must handle concurrent access gracefully (race conditions)
- Must work across all supported platforms (Linux, macOS, Windows)
- Must respect file permissions (0600 for file, 0755 for directories)
- Must not expose secrets in logs or error messages

**Scale/Scope**:
- Single config file per user
- Single LoadConfig call per application run
- No concurrent LoadConfig calls expected (but must handle gracefully)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **PASS** - Feature modifies existing `internal/config/config.go` which follows layer separation. No new layers needed. Repository Pattern not applicable (file system operations, not data access).

- **Interface-Driven Development**: ✅ **PASS** - No new interfaces needed. File system operations use standard library (`os`, `filepath`) which are already abstracted. No dependency injection changes required.

- **Test-First Development**: ✅ **PASS** - TDD approach defined: tests will be written before implementation. Unit tests for file creation logic, integration tests for LoadConfig behavior. Table-driven tests for error scenarios.

- **Idiomatic Go**: ✅ **PASS** - Design follows Go conventions: standard library for file operations, explicit error handling, small focused functions. Naming follows Go conventions (PascalCase for exported, camelCase for unexported).

- **Error Handling**: ✅ **PASS** - Error handling strategy defined: wrapped errors with context (`fmt.Errorf("context: %w", err)`), clear error messages for permission/disk issues. No custom error types needed (standard errors sufficient).

- **Context & Thread Safety**: ✅ **PASS** - No context.Context needed (synchronous file operations). Thread safety handled: file creation operations are atomic at OS level, race conditions handled gracefully (FR-007). No shared state introduced.

- **Technical Constraints**: ✅ **PASS** - No global state introduced. No graceful shutdown needed (synchronous operation). Resource cleanup: file handles closed automatically by Go's defer/GC. No panics, all errors returned.

- **Operational Constraints**: ✅ **PASS** - Logging strategy defined: use `utils.Logger.Debug()` for file creation events (FR-011). Secrets management: config file contains API keys, protected with 0600 permissions (FR-010). No secrets exposed in logs.

**Violations**: None. All principles satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/015-ensure-config-exists/
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
└── config/
    ├── config.go        # Modified: Add file existence check and creation logic
    └── config_test.go   # New: Tests for file creation behavior

test/
└── integration/
    └── config_test.go   # Modified: Integration tests for LoadConfig with file creation
```

**Structure Decision**: Single project structure. Feature modifies existing `internal/config/config.go` file. No new packages or modules needed. Tests added alongside source file (unit tests) and in `test/integration/` (integration tests).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. All constitution principles satisfied.
