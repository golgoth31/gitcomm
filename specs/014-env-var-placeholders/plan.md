# Implementation Plan: Environment Variable Placeholder Substitution in Config Files

**Branch**: `014-env-var-placeholders` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/014-env-var-placeholders/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature adds environment variable placeholder substitution to the config loading mechanism. Config files can contain placeholders in the form `${ENV_VAR_NAME}`, which are automatically replaced with values from the environment during config loading. If a required environment variable is missing, the application exits immediately with a clear error message. The implementation will:

1. **Parse config file content** to identify `${ENV_VAR_NAME}` placeholders
2. **Validate placeholder syntax** (alphanumeric and underscores only, no spaces, nested placeholders, or newlines)
3. **Extract environment variable names** from valid placeholders
4. **Look up environment variables** and replace placeholders with their values
5. **Handle missing variables** by exiting immediately with clear error messages
6. **Process all placeholders** before proceeding with YAML parsing

The technical approach processes the config file content as text before YAML parsing, performing substitution in memory. This ensures placeholders are replaced before viper processes the YAML structure, maintaining backward compatibility with existing config files that don't use placeholders.

## Technical Context

**Language/Version**: Go 1.25.0+

**Primary Dependencies**:
- `os` - Standard library for accessing environment variables (`os.Getenv`, `os.LookupEnv`)
- `regexp` - Standard library for pattern matching to identify placeholders
- `strings` - Standard library for string manipulation and replacement
- `github.com/spf13/viper` - Configuration management (existing, used for reading config after substitution)
- `github.com/rs/zerolog` - Structured logging (existing, via `utils.Logger`)

**Storage**: File system (YAML config file at `~/.gitcomm/config.yaml` or custom path)

**Testing**:
- Standard Go testing framework (`testing` package)
- Unit tests for placeholder parsing and substitution logic
- Integration tests for LoadConfig behavior with placeholders
- Table-driven tests for error scenarios (invalid syntax, missing variables)

**Target Platform**: Linux, macOS, Windows (CLI application)

**Project Type**: CLI tool (single binary)

**Performance Goals**:
- Placeholder substitution completes in under 10ms for typical config files (SC-002: <1 second total)
- No noticeable delay in LoadConfig execution
- Efficient regex matching and string replacement

**Constraints**:
- Must maintain backward compatibility (existing config files without placeholders work unchanged)
- Must not modify the config file on disk (substitution happens in memory)
- Must process placeholders before YAML parsing (viper reads substituted content)
- Must handle all valid YAML structures (nested objects, arrays, strings)
- Must ignore placeholders in YAML comments
- Must exit immediately on invalid syntax or missing variables (fail-fast)

**Scale/Scope**:
- Single config file per user
- Typical config files contain 1-10 placeholders
- Single LoadConfig call per application run
- No concurrent LoadConfig calls expected

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **PASS** - Feature modifies existing `internal/config/config.go` which follows layer separation. No new layers needed. Repository Pattern not applicable (in-memory string processing, not data access).

- **Interface-Driven Development**: ✅ **PASS** - No new interfaces needed. Environment variable access uses standard library (`os`) which is already abstracted. No dependency injection changes required (can add interface for testing if needed, but standard library is sufficient).

- **Test-First Development**: ✅ **PASS** - TDD approach defined: tests will be written before implementation. Unit tests for placeholder parsing and substitution logic, integration tests for LoadConfig behavior with placeholders. Table-driven tests for error scenarios.

- **Idiomatic Go**: ✅ **PASS** - Design follows Go conventions: standard library for string/regex operations, explicit error handling, small focused functions. Naming follows Go conventions (PascalCase for exported, camelCase for unexported).

- **Error Handling**: ✅ **PASS** - Error handling strategy defined: wrapped errors with context (`fmt.Errorf("context: %w", err)`), clear error messages identifying missing variables. Custom error types may be useful for distinguishing invalid syntax vs missing variables.

- **Context & Thread Safety**: ✅ **PASS** - No context.Context needed (synchronous string processing). Thread safety: environment variable access is read-only and thread-safe. No shared state introduced. String operations are immutable.

- **Technical Constraints**: ✅ **PASS** - No global state introduced. No graceful shutdown needed (synchronous operation). Resource cleanup: no resources to clean up (string processing). No panics, all errors returned.

- **Operational Constraints**: ✅ **PASS** - Logging strategy: may add debug logging for placeholder substitution (optional, not required by spec). Secrets management: environment variable values are substituted but not logged. Error messages identify variable names but not values.

**Violations**: None. All principles satisfied.

## Post-Phase 1 Design Re-check

After completing Phase 1 design (data model, contracts, quickstart), all constitution principles remain satisfied:

- **Clean Architecture**: ✅ **PASS** - Design maintains layer separation, no new layers introduced
- **Interface-Driven Development**: ✅ **PASS** - No new interfaces needed, standard library usage
- **Test-First Development**: ✅ **PASS** - Test requirements defined in contract, TDD approach maintained
- **Idiomatic Go**: ✅ **PASS** - Design uses standard library patterns, follows Go conventions
- **Error Handling**: ✅ **PASS** - Error handling strategy defined in contract, clear error messages
- **Context & Thread Safety**: ✅ **PASS** - No context needed, thread-safe environment variable access
- **Technical Constraints**: ✅ **PASS** - No global state, in-memory processing, no resource leaks
- **Operational Constraints**: ✅ **PASS** - Error messages don't expose secrets, logging optional

## Project Structure

### Documentation (this feature)

```text
specs/014-env-var-placeholders/
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
    ├── config.go        # Modified: Add placeholder substitution logic before viper.ReadInConfig()
    └── config_test.go   # Modified: Add tests for placeholder substitution

test/
└── integration/
    └── config_test.go   # Modified: Add integration tests for LoadConfig with placeholders
```

**Structure Decision**: Single project structure. Feature modifies existing `internal/config/config.go` file. No new packages or modules needed. Tests added alongside source file (unit tests) and in `test/integration/` (integration tests).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. All constitution principles satisfied.
