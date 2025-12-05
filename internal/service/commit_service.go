package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golgoth31/gitcomm/internal/ai"
	"github.com/golgoth31/gitcomm/internal/config"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/repository"
	"github.com/golgoth31/gitcomm/internal/ui"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/golgoth31/gitcomm/pkg/tokenization"
)

// CommitService orchestrates the commit message creation workflow
type CommitService struct {
	gitRepo     repository.GitRepository
	formatter   *FormattingService
	validator   *ValidationService
	reader      *bufio.Reader
	options     *model.CommitOptions
	config      *config.Config
	restoreDone chan struct{} // Channel to signal restoration completion (optional)
}

// NewCommitService creates a new commit service
func NewCommitService(gitRepo repository.GitRepository, options *model.CommitOptions, cfg *config.Config) *CommitService {
	return &CommitService{
		gitRepo:     gitRepo,
		formatter:   NewFormattingService(),
		validator:   NewValidationService(),
		reader:      bufio.NewReader(os.Stdin),
		options:     options,
		config:      cfg,
		restoreDone: nil, // Will be set if needed
	}
}

// SetRestoreDoneChannel sets the channel to signal restoration completion
func (s *CommitService) SetRestoreDoneChannel(ch chan struct{}) {
	s.restoreDone = ch
}

// CreateCommit orchestrates the complete commit creation workflow
func (s *CommitService) CreateCommit(ctx context.Context) error {
	utils.Logger.Debug().Msg("Starting commit creation workflow")

	// Capture pre-CLI staging state for restoration
	preCLIState, err := s.gitRepo.CaptureStagingState(ctx)
	if err != nil {
		return fmt.Errorf("failed to capture staging state: %w", err)
	}

	// Set up deferred restoration on cancellation/error
	// Use pointer so we can modify it and defer will see the updated value
	restoreOnExit := true
	defer func() {
		// Closure captures restoreOnExit by reference, so it sees current value
		if restoreOnExit && preCLIState != nil {
			// Check if context was cancelled (signal interrupt)
			isInterrupted := ctx.Err() == context.Canceled
			if isInterrupted {
				utils.Logger.Debug().Msg("Context cancelled - restoring staging state with timeout")
			}

			// Create appropriate context for restoration
			var restoreCtx context.Context
			var restoreCancel context.CancelFunc
			if isInterrupted {
				// Use timeout context (3 seconds) when interrupted by Ctrl+C
				restoreCtx, restoreCancel = context.WithTimeout(context.Background(), 3*time.Second)
				defer restoreCancel()
			} else {
				// Use background context for normal restoration (backward compatibility)
				restoreCtx = context.Background()
			}

			// Restore state on any exit (unless commit succeeded)
			if err := s.restoreStagingState(restoreCtx, preCLIState); err != nil {
				// Check if error is due to timeout
				if errors.Is(err, context.DeadlineExceeded) {
					utils.Logger.Debug().Err(err).Msg("Restoration timed out")
					fmt.Printf("Warning: Restoration timed out. Repository may be in unexpected state.\n")
					fmt.Printf("Please check git status and manually restore if needed.\n")
				} else {
					utils.Logger.Debug().Err(err).Msg("Failed to restore staging state in defer")
				}
			} else {
				utils.Logger.Debug().Msg("Staging state restored")
			}

			// Signal restoration completion if channel is set
			if s.restoreDone != nil {
				close(s.restoreDone)
			}
		}
	}()

	// Auto-stage modified files (always, before any prompts)
	utils.Logger.Debug().Msg("Auto-staging modified files")
	var stagingResult *model.AutoStagingResult
	useAllFiles := s.options != nil && s.options.AutoStage

	if useAllFiles {
		// Stage all files including untracked when -a flag is used
		stagingResult, err = s.gitRepo.StageAllFilesIncludingUntracked(ctx)
	} else {
		// Stage only modified files
		stagingResult, err = s.gitRepo.StageModifiedFiles(ctx)
	}

	if err != nil {
		// Staging failed - restore state and exit
		utils.Logger.Debug().Err(err).Msg("Auto-staging failed")
		if restoreErr := s.restoreStagingState(ctx, preCLIState); restoreErr != nil {
			utils.Logger.Debug().Err(restoreErr).Msg("Failed to restore staging state after staging failure")
		}
		return fmt.Errorf("failed to stage files: %w", err)
	}

	if stagingResult.HasFailures() {
		// Partial failure - abort and restore
		utils.Logger.Debug().Msg("Partial staging failure - aborting")
		if restoreErr := s.restoreStagingState(ctx, preCLIState); restoreErr != nil {
			utils.Logger.Debug().Err(restoreErr).Msg("Failed to restore staging state after partial failure")
		}
		failedFiles := stagingResult.GetFailedFilePaths()
		return fmt.Errorf("%w: failed to stage files: %v", utils.ErrStagingFailed, failedFiles)
	}

	utils.Logger.Debug().Int("staged_count", len(stagingResult.StagedFiles)).Msg("Files auto-staged successfully")

	// Set context value for repository filtering based on addAll flag
	// This ensures GetRepositoryState respects the addAll flag when filtering new files
	ctx = context.WithValue(ctx, repository.IncludeNewFilesKey, useAllFiles)

	// Get repository state after staging
	state, err := s.gitRepo.GetRepositoryState(ctx)
	if err != nil {
		// Error getting state - restore and exit
		if restoreErr := s.restoreStagingState(ctx, preCLIState); restoreErr != nil {
			utils.Logger.Debug().Err(restoreErr).Msg("Failed to restore staging state after state retrieval failure")
		}
		return fmt.Errorf("failed to get repository state: %w", err)
	}

	// Handle empty repository state
	if state.IsEmpty() {
		confirm, err := ui.PromptEmptyCommit(s.reader)
		if err != nil {
			// User cancelled - restore state (defer will handle it)
			return fmt.Errorf("failed to prompt for empty commit: %w", err)
		}
		if !confirm {
			// User declined empty commit - restore state (defer will handle it)
			return utils.ErrNoChanges
		}
	}

	// Determine if AI should be used
	useAI := false
	if s.options == nil || !s.options.SkipAI {
		// Calculate token count
		tokenCalc := tokenization.NewTokenCalculator("openai") // Default provider for calculation
		tokenCount, err := tokenCalc.CalculateForRepositoryState(state)
		if err != nil {
			utils.Logger.Debug().Err(err).Msg("Failed to calculate tokens")
		}
		// Prompt for AI usage
		useAI, err = ui.PromptAIUsage(s.reader, tokenCount)
		if err != nil {
			// User cancelled - restore state (defer will handle it)
			return fmt.Errorf("failed to prompt for AI usage: %w", err)
		}
	}

	var message *model.CommitMessage
	if useAI {
		// Try AI generation
		message, err = s.generateWithAI(ctx, state)
		if err != nil {
			// Check if commit was already created (AcceptAndCommit path)
			if errors.Is(err, utils.ErrCommitAlreadyCreated) {
				// Commit was already created - disable restoration and return success
				restoreOnExit = false
				return nil
			}
			utils.Logger.Debug().Err(err).Msg("AI generation failed, falling back to manual input")
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Falling back to manual input...")
			// Fall through to manual input
			useAI = false
		}
	}

	if !useAI {
		// Prompt for commit message components manually
		message, err = s.promptCommitMessage(nil)
		if err != nil {
			// User cancelled - restore state (defer will handle it)
			return fmt.Errorf("failed to prompt for commit message: %w", err)
		}
	}

	// Validate message
	valid, errors := s.validator.Validate(message)
	if !valid {
		fmt.Println("\nValidation errors:")
		for _, e := range errors {
			fmt.Printf("  - %s: %s\n", e.Field, e.Message)
		}
		confirm, err := ui.PromptConfirm(s.reader, "Continue anyway?", false)
		if err != nil || !confirm {
			// User declined - restore state (defer will handle it)
			return utils.ErrInvalidFormat
		}
	}

	// Display formatted message for review
	formatted := ui.DisplayCommitMessage(message)
	fmt.Println("\n--- Commit Message ---")
	fmt.Println(formatted)
	fmt.Println("---")

	// Confirm before committing
	confirm, err := ui.PromptConfirm(s.reader, "Create commit with this message?", true)
	if err != nil {
		// User cancelled - restore state (defer will handle it)
		return fmt.Errorf("failed to prompt for confirmation: %w", err)
	}
	if !confirm {
		// User cancelled - restore state (defer will handle it)
		return fmt.Errorf("commit cancelled by user")
	}

	// Set signoff based on options
	if s.options != nil {
		message.Signoff = !s.options.NoSignoff
	} else {
		message.Signoff = true // Default to signoff
	}

	// Create commit
	if err := s.gitRepo.CreateCommit(ctx, message); err != nil {
		// Commit failed - restore state (defer will handle it)
		return fmt.Errorf("failed to create commit: %w", err)
	}

	// Commit succeeded - do NOT restore state
	// Disable restoration since commit succeeded (defer captures by value, so we need to set before return)
	restoreOnExit = false
	utils.Logger.Debug().Msg("Commit created successfully")
	fmt.Println("✓ Commit created successfully")
	return nil
}

// restoreStagingState restores the staging state to pre-CLI state
func (s *CommitService) restoreStagingState(ctx context.Context, preCLIState *model.StagingState) error {
	if preCLIState == nil {
		return nil
	}

	// Get current staging state
	currentState, err := s.gitRepo.CaptureStagingState(ctx)
	if err != nil {
		utils.Logger.Debug().Err(err).Msg("Failed to capture current state for restoration")
		// Continue with best-effort restoration
	}

	// Create restoration plan
	plan := &model.RestorationPlan{
		PreCLIState:  preCLIState,
		CurrentState: currentState,
	}

	// Calculate files to unstage (files staged by CLI)
	if currentState != nil {
		plan.FilesToUnstage = currentState.Diff(preCLIState)
	} else {
		// If we can't get current state, unstage all files that were in pre-CLI state
		// This is a fallback - ideally we track which files were staged by CLI
		utils.Logger.Debug().Msg("Cannot determine current state - skipping restoration")
		return nil
	}

	if plan.IsEmpty() {
		// No restoration needed
		return nil
	}

	// Validate plan
	if err := plan.Validate(); err != nil {
		utils.Logger.Debug().Err(err).Msg("Restoration plan validation failed")
		// Continue with best-effort
	}

	// Execute restoration
	utils.Logger.Debug().Int("files_to_unstage", len(plan.FilesToUnstage)).Msg("Restoring staging state")
	if err := s.gitRepo.UnstageFiles(ctx, plan.GetFilesToUnstage()); err != nil {
		// Check if error is due to timeout
		if errors.Is(err, context.DeadlineExceeded) {
			utils.Logger.Debug().Err(err).Msg("Restoration timed out")
			fmt.Printf("Warning: Restoration timed out. Repository may be in unexpected state.\n")
			fmt.Printf("Please check git status and manually restore if needed.\n")
			return fmt.Errorf("%w: %v", utils.ErrRestorationFailed, err)
		}
		utils.Logger.Debug().Err(err).Msg("Failed to restore staging state")
		fmt.Printf("Warning: failed to restore staging state. Repository may be in unexpected state.\n")
		fmt.Printf("Please check git status and manually restore if needed.\n")
		return fmt.Errorf("%w: %v", utils.ErrRestorationFailed, err)
	}

	utils.Logger.Debug().Msg("Staging state restored successfully")
	return nil
}

// promptCommitMessage prompts the user for all commit message components
// If prefilled is not nil, the fields will be pre-filled with values from prefilled
func (s *CommitService) promptCommitMessage(prefilled *ui.PrefilledCommitMessage) (*model.CommitMessage, error) {
	message := &model.CommitMessage{}

	// Prompt for type
	var commitType string
	var err error
	if prefilled != nil && prefilled.Type != "" {
		commitType, err = ui.PromptCommitTypeWithPreselection(s.reader, prefilled.Type)
	} else {
		commitType, err = ui.PromptCommitType(s.reader)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to prompt for type: %w", err)
	}
	message.Type = commitType

	// Prompt for scope
	var scope string
	if prefilled != nil {
		scope, err = ui.PromptScopeWithDefault(s.reader, prefilled.Scope)
	} else {
		scope, err = ui.PromptScope(s.reader)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to prompt for scope: %w", err)
	}
	message.Scope = scope

	// Prompt for subject (required, with validation)
	var subject string
	if prefilled != nil && prefilled.Subject != "" {
		subject, err = ui.PromptSubjectWithDefault(s.reader, prefilled.Subject)
	} else {
		subject, err = ui.PromptSubject(s.reader)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to prompt for subject: %w", err)
	}
	message.Subject = subject

	// Prompt for body
	var body string
	if prefilled != nil {
		body, err = ui.PromptBodyWithDefault(s.reader, prefilled.Body)
	} else {
		body, err = ui.PromptBody(s.reader)
	}
	if err != nil {
		// Body is optional, so we can continue if user cancels
		utils.Logger.Debug().Err(err).Msg("Body input cancelled or failed")
		message.Body = ""
	} else {
		message.Body = body
	}

	// Prompt for footer
	var footer string
	if prefilled != nil {
		footer, err = ui.PromptFooterWithDefault(s.reader, prefilled.Footer)
	} else {
		footer, err = ui.PromptFooter(s.reader)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to prompt for footer: %w", err)
	}
	message.Footer = footer

	return message, nil
}

// generateWithAI generates a commit message using AI
// This is the public entry point that calls the internal implementation with retry limit
func (s *CommitService) generateWithAI(ctx context.Context, repoState *model.RepositoryState) (*model.CommitMessage, error) {
	return s.generateWithAIWithRetry(ctx, repoState, 0)
}

// generateWithAIWithRetry generates a commit message using AI with retry limit tracking
func (s *CommitService) generateWithAIWithRetry(ctx context.Context, repoState *model.RepositoryState, retryCount int) (*model.CommitMessage, error) {
	// Prevent infinite recursion
	const maxRetries = 3
	if retryCount >= maxRetries {
		fmt.Println("Maximum retry limit reached. Falling back to manual input...")
		return s.promptCommitMessage(nil)
	}
	// Get provider configuration
	providerName := "openai"
	if s.options != nil && s.options.AIProvider != "" {
		providerName = s.options.AIProvider
	} else if s.config != nil && s.config.AI.DefaultProvider != "" {
		providerName = s.config.AI.DefaultProvider
	}

	providerConfig, err := s.config.GetProviderConfig(providerName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
	}

	// Create AI provider
	var aiProvider ai.AIProvider
	switch providerName {
	case "openai":
		aiProvider = ai.NewOpenAIProvider(providerConfig)
	case "anthropic":
		aiProvider = ai.NewAnthropicProvider(providerConfig)
	case "mistral":
		aiProvider = ai.NewMistralProvider(providerConfig)
	case "local":
		aiProvider = ai.NewLocalProvider(providerConfig)
	default:
		return nil, fmt.Errorf("%w: unknown provider %s", utils.ErrAIProviderUnavailable, providerName)
	}

	// Generate commit message
	aiMessage, err := aiProvider.GenerateCommitMessage(ctx, repoState)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
	}

	// Parse AI message into CommitMessage structure
	message, err := s.parseAIMessage(aiMessage)
	if err != nil {
		utils.Logger.Debug().Err(err).Msg("Failed to parse AI message")
		// Try to use as-is
		message = &model.CommitMessage{
			Type:    "feat",
			Subject: strings.TrimSpace(aiMessage),
		}
	}

	// Validate AI-generated message
	valid, validationErrors := s.validator.Validate(message)
	if !valid {
		// Show validation errors
		var errorMessages []string
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", ve.Field, ve.Message))
		}

		// Prompt user to edit or use with warning
		edit, err := ui.PromptAIMessageEdit(s.reader, errorMessages)
		if err != nil {
			return nil, fmt.Errorf("failed to prompt for edit: %w", err)
		}

		if edit {
			// User wants to edit - fall back to manual input
			return s.promptCommitMessage(nil)
		}
		// User wants to use as-is with warning
		fmt.Println("Warning: Using message that does not fully conform to Conventional Commits format")
	}

	// Show AI message and get user acceptance with three options
	acceptance, err := ui.PromptAIMessageAcceptanceOptions(s.reader, ui.DisplayCommitMessage(message))
	if err != nil {
		return nil, fmt.Errorf("failed to prompt for acceptance: %w", err)
	}

	switch acceptance {
	case ui.AcceptAndCommit:
		// User wants to commit immediately - create commit here
		// Set signoff based on options
		if s.options != nil {
			message.Signoff = !s.options.NoSignoff
		} else {
			message.Signoff = true // Default to signoff
		}

		// Create commit immediately
		if err := s.gitRepo.CreateCommit(ctx, message); err != nil {
			// Commit failed - handle failure with retry/edit/cancel options
			return s.handleCommitFailure(ctx, message, err)
		}

		// Commit succeeded - return sentinel error to signal commit was already created
		utils.Logger.Debug().Msg("Commit created successfully via AcceptAndCommit")
		fmt.Println("✓ Commit created successfully")
		return message, utils.ErrCommitAlreadyCreated

	case ui.AcceptAndEdit:
		// User wants to edit - parse AI message into PrefilledCommitMessage and pre-fill prompts
		prefilled := s.parseAIMessageToPrefilled(aiMessage)
		commitMsg, err := s.promptCommitMessage(&prefilled)
		if err != nil {
			// Handle cancellation (restore staging state)
			return nil, fmt.Errorf("failed to prompt for commit message: %w", err)
		}

		// Create commit with edited message
		// Set signoff based on options
		if s.options != nil {
			commitMsg.Signoff = !s.options.NoSignoff
		} else {
			commitMsg.Signoff = true // Default to signoff
		}

		// Create commit
		if err := s.gitRepo.CreateCommit(ctx, commitMsg); err != nil {
			return s.handleCommitFailure(ctx, commitMsg, err)
		}

		// Commit succeeded - return sentinel error to signal commit was already created
		utils.Logger.Debug().Msg("Commit created successfully via AcceptAndEdit")
		fmt.Println("✓ Commit created successfully")
		return commitMsg, utils.ErrCommitAlreadyCreated

	case ui.Reject:
		// User rejected - prompt for choice: new AI or manual input
		useNewAI, err := ui.PromptRejectChoice(s.reader)
		if err != nil {
			return nil, fmt.Errorf("failed to prompt for reject choice: %w", err)
		}

		if useNewAI {
			// Generate new AI message (recursive call with incremented retry count)
			newMessage, err := s.generateWithAIWithRetry(ctx, repoState, retryCount+1)
			if err != nil {
				// AI generation failed - fall back to manual input with error message
				fmt.Printf("Error generating new AI message: %v\n", err)
				fmt.Println("Falling back to manual input...")
				return s.promptCommitMessage(nil)
			}
			return newMessage, nil
		} else {
			// Fall back to manual input with empty fields
			return s.promptCommitMessage(nil)
		}

	default:
		// Should not happen, but handle gracefully
		return nil, fmt.Errorf("unknown acceptance option: %v", acceptance)
	}
}

// handleCommitFailure handles commit failure after AcceptAndCommit by prompting user for retry/edit/cancel
func (s *CommitService) handleCommitFailure(ctx context.Context, message *model.CommitMessage, commitErr error) (*model.CommitMessage, error) {
	// Display error message
	fmt.Printf("\nError creating commit: %v\n", commitErr)

	// Prompt for retry/edit/cancel
	choice, err := ui.PromptCommitFailureChoice(s.reader)
	if err != nil {
		return nil, fmt.Errorf("failed to prompt for commit failure choice: %w", err)
	}

	switch choice {
	case ui.RetryCommit:
		// Retry commit with same message
		if err := s.gitRepo.CreateCommit(ctx, message); err != nil {
			// Recursive retry (with limit to prevent infinite loop)
			// For now, just retry once more
			return s.handleCommitFailure(ctx, message, err)
		}
		utils.Logger.Debug().Msg("Commit created successfully after retry")
		fmt.Println("✓ Commit created successfully")
		return message, utils.ErrCommitAlreadyCreated

	case ui.EditMessage:
		// Fall back to accept and edit flow - convert CommitMessage to PrefilledCommitMessage
		prefilled := s.commitMessageToPrefilled(message)
		return s.promptCommitMessage(&prefilled)

	case ui.CancelCommit:
		// User cancelled - return error to trigger staging state restoration
		return nil, fmt.Errorf("commit cancelled by user after failure")

	default:
		// Should not happen
		return nil, fmt.Errorf("unknown commit failure choice: %v", choice)
	}
}

// parseAIMessageToPrefilled converts an AI-generated message string into PrefilledCommitMessage structure
func (s *CommitService) parseAIMessageToPrefilled(aiMessage string) ui.PrefilledCommitMessage {
	prefilled := ui.PrefilledCommitMessage{}

	lines := strings.Split(strings.TrimSpace(aiMessage), "\n")
	if len(lines) == 0 {
		return prefilled
	}

	// Parse header (first line): type(scope): subject
	header := lines[0]
	parts := strings.SplitN(header, ":", 2)
	if len(parts) == 2 {
		typeScope := strings.TrimSpace(parts[0])
		prefilled.Subject = strings.TrimSpace(parts[1])

		// Parse type and scope
		if strings.Contains(typeScope, "(") && strings.Contains(typeScope, ")") {
			openIdx := strings.Index(typeScope, "(")
			closeIdx := strings.Index(typeScope, ")")
			prefilled.Type = strings.TrimSpace(typeScope[:openIdx])
			prefilled.Scope = strings.TrimSpace(typeScope[openIdx+1 : closeIdx])
		} else {
			prefilled.Type = strings.TrimSpace(typeScope)
		}
	}

	// Parse body and footer (if present)
	if len(lines) > 1 {
		var bodyLines []string
		var footerLines []string
		inFooter := false

		for i := 1; i < len(lines); i++ {
			line := lines[i]
			if line == "" {
				if len(bodyLines) > 0 {
					inFooter = true
				}
				continue
			}
			if inFooter {
				footerLines = append(footerLines, line)
			} else {
				bodyLines = append(bodyLines, line)
			}
		}

		if len(bodyLines) > 0 {
			prefilled.Body = strings.Join(bodyLines, "\n")
		}
		if len(footerLines) > 0 {
			prefilled.Footer = strings.Join(footerLines, "\n")
		}
	}

	return prefilled
}

// commitMessageToPrefilled converts a CommitMessage to PrefilledCommitMessage
func (s *CommitService) commitMessageToPrefilled(msg *model.CommitMessage) ui.PrefilledCommitMessage {
	return ui.PrefilledCommitMessage{
		Type:    msg.Type,
		Scope:   msg.Scope,
		Subject: msg.Subject,
		Body:    msg.Body,
		Footer:  msg.Footer,
	}
}

// parseAIMessage attempts to parse an AI-generated message into CommitMessage structure
func (s *CommitService) parseAIMessage(aiMessage string) (*model.CommitMessage, error) {
	message := &model.CommitMessage{
		Signoff: true, // Default
	}

	lines := strings.Split(strings.TrimSpace(aiMessage), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	// Parse header (first line): type(scope): subject
	header := lines[0]
	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid header format")
	}

	typeScope := strings.TrimSpace(parts[0])
	message.Subject = strings.TrimSpace(parts[1])

	// Parse type and scope
	if strings.Contains(typeScope, "(") && strings.Contains(typeScope, ")") {
		openIdx := strings.Index(typeScope, "(")
		closeIdx := strings.Index(typeScope, ")")
		message.Type = strings.TrimSpace(typeScope[:openIdx])
		message.Scope = strings.TrimSpace(typeScope[openIdx+1 : closeIdx])
	} else {
		message.Type = strings.TrimSpace(typeScope)
	}

	// Parse body and footer (if present)
	if len(lines) > 1 {
		var bodyLines []string
		var footerLines []string
		inFooter := false

		for i := 1; i < len(lines); i++ {
			line := lines[i]
			if line == "" {
				if len(bodyLines) > 0 {
					inFooter = true
				}
				continue
			}
			if inFooter {
				footerLines = append(footerLines, line)
			} else {
				bodyLines = append(bodyLines, line)
			}
		}

		if len(bodyLines) > 0 {
			message.Body = strings.Join(bodyLines, "\n")
		}
		if len(footerLines) > 0 {
			message.Footer = strings.Join(footerLines, "\n")
		}
	}

	return message, nil
}
