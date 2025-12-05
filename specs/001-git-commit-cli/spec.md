# Feature Specification: Git Commit Message Automation CLI

**Feature Branch**: `001-git-commit-cli`
**Created**: 2025-01-27
**Status**: Draft
**Input**: User description: "this project is a simple cli that implements git message template, utomating user commit message workflow with some rules. this project is inspired by https://github.com/studyzy/gitcomm . The resulting commit message MUST follow the https://www.conventionalcommits.org/en/v1.0.0/#specification  . The cli should have the following steps:"

## Clarifications

### Session 2025-01-27

- Q: Which AI provider(s) should the CLI support, and how should credentials be configured? → A: Multiple AI providers (OpenAI, Anthropic, local models) with provider selection via CLI flag or config
- Q: When an AI provider call fails (network error, API error, timeout), what should the CLI do? → A: Immediately fallback to manual input with clear error message displayed to user
- Q: How should the CLI calculate token estimates for different AI providers? → A: Provider-specific tokenization libraries (tiktoken for OpenAI, custom for Anthropic) with character-based fallback for unknown providers
- Q: What should the CLI do when there are no changes to commit (no staged or unstaged files)? → A: Prompt user to confirm they want to create an empty commit, then proceed with manual input
- Q: When the AI provider returns a commit message that doesn't conform to Conventional Commits specification, what should the CLI do? → A: Validate and reject with option to edit the AI message to fix format, or use as-is (with warning)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Manual Commit Message Creation (Priority: P1)

A developer wants to create a commit message following the Conventional Commits specification without using AI assistance. The CLI guides them through entering the required components (scope, subject, body, footer) and validates the final message before committing.

**Why this priority**: This is the core functionality that must work independently. Even without AI integration, users can create properly formatted commit messages. This delivers immediate value and is the foundation for all other features.

**Independent Test**: Can be fully tested by running the CLI in a git repository, declining AI assistance, and manually entering all commit message components. The test verifies that the resulting commit message follows Conventional Commits format and the commit is successfully created.

**Acceptance Scenarios**:

1. **Given** a developer is in a git repository with staged or unstaged changes, **When** they run the CLI and decline AI assistance, **Then** the CLI prompts for scope, subject, body, and footer in sequence
2. **Given** a developer enters a scope, **When** they proceed to the subject prompt, **Then** the CLI accepts the scope (which may be empty)
3. **Given** a developer is prompted for the subject, **When** they enter an empty subject, **Then** the CLI rejects it and prompts again until a non-empty subject is provided
4. **Given** a developer enters all commit message components, **When** they review the formatted message, **Then** the CLI displays the complete message in Conventional Commits format for validation
5. **Given** a developer validates the commit message, **When** they confirm, **Then** the CLI creates the git commit with the formatted message
6. **Given** a developer reviews the commit message, **When** they reject it, **Then** the CLI allows them to edit or restart the process
7. **Given** there are no changes to commit (no staged or unstaged files), **When** the developer runs the CLI, **Then** the CLI prompts to confirm creating an empty commit and proceeds with manual input if confirmed

---

### User Story 2 - AI-Assisted Commit Message Generation (Priority: P2)

A developer wants to use AI to automatically generate a commit message based on the repository state. The CLI analyzes the git changes, calculates token usage, and optionally calls an AI provider to suggest a commit message that the user can accept, modify, or reject.

**Why this priority**: This adds significant value by automating the most time-consuming part of commit message creation, but the core manual workflow (P1) must work first. This feature enhances the user experience but is not required for basic functionality.

**Independent Test**: Can be fully tested by running the CLI in a git repository, accepting AI assistance, and verifying that the AI-generated message follows Conventional Commits format. The test verifies token calculation, AI provider integration, and the ability to accept or reject the AI suggestion.

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI in a git repository, **When** the CLI analyzes the repository state using git porcelain commands, **Then** it calculates and displays the estimated AI token count
2. **Given** the CLI displays token count, **When** the developer chooses to use AI, **Then** the CLI calls the AI provider with repository state information
3. **Given** the AI provider returns a commit message, **When** the developer reviews it, **Then** the CLI validates the message against Conventional Commits format and displays it with acceptance/rejection prompt
4. **Given** the AI-generated message does not conform to Conventional Commits format, **When** the developer reviews it, **Then** the CLI rejects it and provides options to edit the message to fix format issues or use as-is with a warning
5. **Given** the developer accepts the AI-generated message (valid or with warning), **When** they validate it, **Then** the CLI proceeds to create the commit with that message
6. **Given** the developer rejects the AI-generated message, **When** they choose to proceed manually, **Then** the CLI falls back to the manual input workflow (User Story 1)
7. **Given** the developer chooses not to use AI initially, **When** they proceed, **Then** the CLI skips AI generation and goes directly to manual input
8. **Given** the AI provider call fails (network error, API error, timeout), **When** the error occurs, **Then** the CLI immediately falls back to manual input with a clear error message displayed to the user

---

### User Story 3 - CLI Options and Configuration (Priority: P3)

A developer wants to configure the CLI behavior using command-line options, specifically to automatically stage files and control commit signoff behavior.

**Why this priority**: These are convenience features that enhance usability but are not essential for core functionality. The CLI must work without these options, and they can be added incrementally.

**Independent Test**: Can be fully tested by running the CLI with different option combinations (-a, -s) and verifying that files are automatically staged when requested and commit signoff is controlled appropriately.

**Acceptance Scenarios**:

1. **Given** a developer runs the CLI with the `-a` option, **When** the CLI starts, **Then** it automatically stages all unstaged files using `git add -A` before proceeding
2. **Given** a developer runs the CLI without the `-a` option, **When** the CLI starts, **Then** it uses the current git repository state without modifying the staging area
3. **Given** a developer runs the CLI with the `-s` option, **When** the commit is created, **Then** the commit does not include a signoff line
4. **Given** a developer runs the CLI without the `-s` option, **When** the commit is created, **Then** the commit includes a signoff line by default

---

### Edge Cases

- What happens when the CLI is run outside a git repository?
- What happens when there are no changes to commit (no staged or unstaged files)? → The CLI prompts the user to confirm they want to create an empty commit, then proceeds with manual input if confirmed
- What happens when the selected AI provider is unavailable or returns an error? → The CLI immediately falls back to manual input with a clear error message displayed to the user
- What happens when no AI provider is configured or credentials are missing?
- What happens when the user selects an AI provider that is not configured?
- What happens when the AI-generated message does not follow Conventional Commits format? → The CLI validates the message, rejects it if invalid, and provides options to edit the AI message to fix format issues or use as-is with a warning
- What happens when the user interrupts the CLI during input (Ctrl+C)?
- What happens when git operations fail (e.g., repository is locked, no write permissions)?
- What happens when the repository state changes between token calculation and commit creation?
- How does the CLI handle very long commit messages (subject > 72 chars, body > 320 chars per Conventional Commits)?
- What happens when the user provides invalid scope or subject characters that conflict with Conventional Commits format?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST use git porcelain commands to retrieve current repository state (staged and unstaged changes)
- **FR-002**: System MUST calculate and display estimated AI token count before prompting user for AI usage decision using provider-specific tokenization libraries (tiktoken for OpenAI, custom for Anthropic) with character-based fallback for unknown providers
- **FR-003**: System MUST prompt user to choose whether to use AI provider for commit message generation
- **FR-004**: System MUST support multiple AI providers (OpenAI, Anthropic, local models) with provider selection via CLI flag or configuration file
- **FR-005**: System MUST call the selected AI provider with repository state information when user chooses AI assistance
- **FR-006**: System MUST validate AI-generated commit message against Conventional Commits specification and display it with acceptance/rejection prompt
- **FR-022**: System MUST reject AI-generated messages that do not conform to Conventional Commits format and provide options to edit the message to fix format issues or use as-is with a warning
- **FR-020**: System MUST immediately fallback to manual input with a clear error message when AI provider calls fail (network error, API error, timeout)
- **FR-007**: System MUST allow user to proceed with manual input if they decline AI or reject AI-generated message
- **FR-008**: System MUST prompt user for commit scope (optional, can be empty)
- **FR-009**: System MUST prompt user for commit subject and reject empty subjects until a non-empty value is provided
- **FR-010**: System MUST prompt user for commit body (optional, can be empty)
- **FR-011**: System MUST prompt user for commit footer (optional, can be empty)
- **FR-012**: System MUST format the complete commit message according to Conventional Commits specification
- **FR-013**: System MUST display the formatted commit message for user validation before creating the commit
- **FR-014**: System MUST create the git commit with the validated message when user confirms
- **FR-015**: System MUST support `-a` option to automatically stage all unstaged files using `git add -A`
- **FR-016**: System MUST support `-s` option to disable commit signoff (omit `Signed-off-by` line)
- **FR-017**: System MUST validate that the final commit message conforms to Conventional Commits specification before committing
- **FR-018**: System MUST handle errors gracefully and provide clear error messages to the user
- **FR-019**: System MUST detect if running outside a git repository and display an appropriate error message
- **FR-021**: System MUST detect when there are no changes to commit (no staged or unstaged files) and prompt user to confirm they want to create an empty commit before proceeding with manual input

### Key Entities *(include if feature involves data)*

- **Repository State**: Represents the current state of the git repository, including staged changes, unstaged changes, and repository metadata. Used to calculate AI tokens and generate commit messages.

- **Commit Message Components**: Represents the structured parts of a commit message (scope, subject, body, footer) that together form a Conventional Commits compliant message.

- **AI Token Calculation**: Represents the estimated token count needed for AI provider API calls, calculated from repository state information. Token calculation uses provider-specific tokenization libraries (tiktoken for OpenAI, custom implementations for Anthropic) with character-based fallback for unknown or local model providers. This ensures accurate token estimates that account for different tokenization methods used by different providers.

- **AI Provider Configuration**: Represents the configuration for multiple supported AI providers (OpenAI, Anthropic, local models). Provider selection can be specified via CLI flag or configuration file, with credentials stored securely in the configuration file.

- **Formatted Commit Message**: Represents the final commit message string that conforms to Conventional Commits specification, ready for git commit.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can create a Conventional Commits compliant commit message in under 2 minutes when using manual input
- **SC-002**: Users can create a commit message using AI assistance in under 30 seconds when AI suggestion is accepted
- **SC-003**: 100% of generated commit messages conform to Conventional Commits specification format
- **SC-004**: Users can successfully create commits in 95% of attempts without encountering blocking errors
- **SC-005**: CLI responds to user input within 100ms for all interactive prompts
- **SC-006**: AI token calculation accuracy is within 10% of actual token usage for 90% of repository states
