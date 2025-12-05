# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **CLI Prompts Migration to Huh Library**: All CLI prompts have been migrated from custom Bubble Tea implementations to the `huh` library
  - All prompts now use `huh` library (`github.com/charmbracelet/huh`) for consistent, modern UI
  - Prompts render inline (no alt screen mode) - terminal history remains visible throughout interactions
  - Post-validation display: After each prompt is validated, it's replaced with a summary line: `✓ <title>: <value>` with green checkmark
  - All prompt function signatures remain unchanged (backward compatible)
  - Removed old custom Bubble Tea models: `TextInputModel`, `MultilineInputModel`, `YesNoChoiceModel`, `SelectListModel`
  - Enhanced integration tests for inline rendering and post-validation display

### Added
- **Diff Computation for Staged Files**: GetRepositoryState now computes unified diffs (patch format) for all staged files
  - Provides AI models with actual code changes for more accurate commit message generation
  - Optimized for token usage: 0 context lines, 5000 character limit per file (shows metadata for larger files)
  - Handles edge cases: binary files (empty diff), unmerged files (attempts diff, fallback to empty), empty repository (treats as new files)
  - Errors are logged but don't stop processing (graceful degradation)
  - Unstaged files always have empty diff field (only staged files include diff content)
- **AI Message Acceptance Options**: Enhanced AI message acceptance workflow with three distinct options
  - **Accept and commit directly**: Commit immediately with AI-generated message (fastest path for satisfied users)
  - **Accept and edit**: Pre-fill all commit message fields (type, scope, subject, body, footer) with AI values for quick editing
  - **Reject**: Generate a new AI message or proceed with manual input
  - Pre-filling support: Commit type automatically selected in interactive list, all text fields pre-populated
  - Error recovery: If commit fails after "accept and commit directly", user can retry, edit, or cancel
  - Retry limit: Maximum 3 retries when generating new AI messages after rejection (prevents infinite loops)
  - Graceful fallback: If AI generation fails after rejection, automatically falls back to manual input with error message
- **Git Configuration Integration**: GitComm now automatically extracts and uses git configuration from `.git/config` and `~/.gitconfig` files before initializing git objects
  - Uses `user.name` and `user.email` from git config for commit author (with local config taking precedence over global)
  - Falls back to defaults ("gitcomm <gitcomm@local>") when git config is not available
  - Silently ignores missing or unreadable config files with debug logging
- **SSH Commit Signing**: Automatic SSH commit signing when configured in git config
  - Supports SSH signing via `gpg.format = ssh` and `user.signingkey` configuration
  - Signs commits by default when SSH signing is configured (unless `commit.gpgsign = false`)
  - Handles signing failures gracefully by creating unsigned commits and logging errors
  - New `--no-sign` CLI flag to disable commit signing regardless of git config
- **Debug Logging for Config**: Debug messages are logged when config files are missing, unreadable, or when signing configuration is unavailable

### Fixed
- **CLI hang on Ctrl+C**: Fixed issue where CLI would hang indefinitely when Ctrl+C was pressed during state restoration. CLI now exits within 5 seconds with a 3-second timeout for restoration operations.

### Added
- Initial implementation of git commit message automation CLI
- Manual commit message creation with interactive prompts
- AI-assisted commit message generation (OpenAI, Anthropic, Mistral, local models)
- Token calculation for AI providers
- Conventional Commits format validation
- CLI options: auto-stage (`-a`), no-signoff (`-s`), provider selection
- Configuration file support for AI provider credentials
- Error handling with graceful fallback to manual input
- Integration with go-git for repository operations
- **Auto-staging of modified files on CLI launch**
- **State restoration on cancellation or error**
- **Signal handling for graceful interruption (Ctrl+C, SIGTERM)**
- **Debug logging with `-d`/`--debug` flag (raw text format, no timestamps)**
- **Interactive commit type selection with visual feedback (checkmarks, highlighting, arrow key navigation, letter-based navigation)**
- **Commit type selection confirmation display showing chosen type before next prompt**
- **Multiline input for body and footer fields with double-Enter completion and blank line preservation**
- **Mistral AI provider support**: Added Mistral as a new AI provider option for commit message generation
- **Official SDK integration**: Replaced custom HTTP clients with official Go SDKs for OpenAI, Anthropic, and Mistral providers, improving reliability and maintainability
- **OpenAI SDK v3 upgrade**: Upgraded OpenAI provider from SDK v1 to SDK v3 (`github.com/openai/openai-go/v3`) while maintaining 100% backward compatibility
- **OpenAI Responses API migration**: Migrated OpenAI provider from Chat Completions API to Responses API (`/v1/responses`) while maintaining 100% backward compatibility with existing configurations and user experience
- **Unified AI provider prompts**: All AI providers (OpenAI, Anthropic, Mistral, local) now use identical system and user messages with validation rules extracted dynamically from MessageValidator, ensuring consistent commit message generation and validation compliance across all providers

### Features
- Interactive prompts for commit type, scope, subject, body, footer
- AI message generation with format validation
- Edit or use-with-warning options for invalid AI messages
- Empty commit confirmation
- Auto-staging of files (modified files always, untracked with `-a` flag)
- Signoff control
- **Automatic staging state capture and restoration**
- **Graceful interruption handling with state cleanup**
- **Debug logging mode for troubleshooting (silent by default)**

## [0.1.0] - 2025-01-27

### Added
- Initial release
- Core functionality for manual and AI-assisted commit message creation
