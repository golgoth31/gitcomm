# Implementation Plan: Improve Commit Generation with Git Config

**Branch**: `012-git-config-commit` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/012-git-config-commit/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Extract git configuration (user.name, user.email, SSH signing settings) from `.git/config` and `~/.gitconfig` files before initializing git objects, and use these values to configure commit author and SSH signing. This ensures commits are properly attributed and signed according to user's git configuration, with graceful fallback to defaults when config is unavailable.

## Technical Context

**Language/Version**: Go 1.25.0+
**Primary Dependencies**:
- Existing: `github.com/go-git/go-git/v5` (v5.16.4) - git repository operations and SSH signing support
- Existing: `github.com/go-git/gcfg` (v1.5.1-0.20230307220236-3a3c6141e376) - git config file parsing (INI format)
- Existing: `github.com/rs/zerolog` (v1.34.0) - debug logging
- Existing: `github.com/spf13/cobra` (v1.10.1) - CLI flag parsing
- No new external dependencies required
**Storage**: N/A (in-memory config extraction, no persistence)
**Testing**: Go `testing` package with table-driven tests, existing test infrastructure
**Target Platform**: Linux/macOS/Windows (CLI tool)
**Project Type**: Single CLI application
**Performance Goals**: Config extraction should complete in <50ms per file (SC-003, SC-004)
**Constraints**:
- Must extract config before opening repository with go-git (FR-001, FR-002)
- Must silently ignore missing/unreadable config files (FR-009)
- Must maintain backward compatibility with existing GitRepository interface
- Must not break existing commit creation workflow
**Scale/Scope**:
- 2 config file locations (`.git/config`, `~/.gitconfig`)
- Single unified config extraction per commit operation
- SSH signing support (no GPG support)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with gitcomm Constitution principles:

- **Clean Architecture**: ✅ **COMPLIANT**
  - Config extraction will be placed in `pkg/git/config/` (shared utility)
  - Repository implementation in `internal/repository/` will use the shared extractor
  - Clear separation: config extraction (pkg) → repository usage (internal)

- **Interface-Driven Development**: ✅ **COMPLIANT**
  - Config extractor will implement an interface for testability
  - Dependencies injected via constructors (file paths, logger)
  - No global state

- **Test-First Development**: ✅ **COMPLIANT**
  - TDD approach: tests for config extraction first
  - Table-driven tests for various config scenarios
  - Integration tests for repository usage

- **Idiomatic Go**: ✅ **COMPLIANT**
  - Follows Go naming conventions
  - Small, focused functions
  - Proper error handling

- **Error Handling**: ✅ **COMPLIANT**
  - Explicit error handling for file I/O and parsing failures
  - Wrapped errors for traceability
  - No panics (silent ignore per FR-009)

- **Context & Thread Safety**: ✅ **COMPLIANT**
  - Config extraction is stateless (no shared mutable state)
  - Thread-safe file operations
  - No goroutines required

- **Technical Constraints**: ✅ **COMPLIANT**
  - No global state
  - Stateless config extraction
  - Resource cleanup not applicable (file reads are closed automatically)

- **Operational Constraints**: ✅ **COMPLIANT**
  - Debug logging via existing zerolog infrastructure
  - No secrets involved (git config is user-readable)
  - Error messages don't expose internal details

**Violations**: None. All principles are satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/012-git-config-commit/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
pkg/
└── git/
    └── config/
        ├── extractor.go          # NEW: Git config extraction interface and implementation
        ├── extractor_test.go     # NEW: Tests for config extraction
        └── errors.go             # NEW: Error types for config extraction

internal/
└── repository/
    └── git_repository_impl.go   # MODIFY: Use config extractor before opening repository, use extracted values for author and signing

cmd/
└── gitcomm/
    └── main.go                  # MODIFY: Add --no-sign CLI flag

test/
└── integration/
    └── git_config_test.go       # NEW: Integration tests for config extraction and commit signing
```

**Structure Decision**: Single project structure. The git config extractor is placed in `pkg/git/config/` as a shared utility that can be used by repository implementations. The repository implementation in `internal/repository/` is modified to use the extractor before opening the repository and to use extracted values for commit author and signing configuration.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. All principles satisfied.

---

## Design Artifacts Generated

### Phase 0: Research
- **research.md**: Technical decisions for INI parsing, SSH signing, config extraction timing, private key resolution, and CLI flag design

### Phase 1: Design & Contracts
- **data-model.md**: Entity definitions (GitConfig, CommitSigner), relationships, data flow, default values, error handling
- **contracts/config-extractor-contract.md**: Interface contract for ConfigExtractor
- **contracts/git-repository-extension-contract.md**: Implementation contract for GitRepository extension
- **contracts/cli-extension-contract.md**: Contract for --no-sign CLI flag
- **quickstart.md**: Integration scenarios, usage examples, troubleshooting guide

### Next Steps
- Run `/gitcomm/speckit.tasks` to generate task breakdown for implementation
