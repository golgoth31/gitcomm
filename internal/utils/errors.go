package utils

import (
	"errors"
	"fmt"
)

// Domain error types for gitcomm

var (
	// ErrNotGitRepository indicates the CLI was run outside a git repository
	ErrNotGitRepository = errors.New("not a git repository: run gitcomm from within a git repository")

	// ErrNoChanges indicates there are no staged or unstaged changes to commit
	ErrNoChanges = errors.New("no changes to commit: stage some files or confirm empty commit")

	// ErrInvalidFormat indicates the commit message does not conform to Conventional Commits specification
	ErrInvalidFormat = errors.New("commit message does not conform to Conventional Commits format: see https://www.conventionalcommits.org/")

	// ErrAIProviderUnavailable indicates the AI provider is unavailable or returned an error
	ErrAIProviderUnavailable = errors.New("AI provider unavailable: check API key and network connection")

	// ErrEmptySubject indicates the commit message subject is empty
	ErrEmptySubject = errors.New("commit message subject cannot be empty: subject is required")

	// ErrTokenCalculationFailed indicates token calculation failed
	ErrTokenCalculationFailed = errors.New("token calculation failed: unable to estimate AI token usage")

	// ErrStagingFailed indicates staging operation failed (partial or complete failure)
	ErrStagingFailed = errors.New("staging operation failed: unable to stage files")

	// ErrRestorationFailed indicates state restoration operation failed
	ErrRestorationFailed = errors.New("restoration operation failed: unable to restore staging state")

	// ErrStagingStateInvalid indicates captured staging state is invalid or corrupted
	ErrStagingStateInvalid = errors.New("staging state invalid: captured state is invalid or corrupted")

	// ErrRestorationPlanInvalid indicates restoration plan is invalid (e.g., files don't exist)
	ErrRestorationPlanInvalid = errors.New("restoration plan invalid: plan is invalid or cannot be executed")

	// ErrInterruptedDuringStaging indicates CLI was interrupted while staging was in progress
	ErrInterruptedDuringStaging = errors.New("interrupted during staging: CLI was interrupted while staging was in progress. Staging state has been restored")

	// ErrCommitAlreadyCreated indicates the commit was already created (e.g., via AcceptAndCommit)
	// This is a sentinel error that should be handled by skipping further commit processing
	ErrCommitAlreadyCreated = errors.New("commit already created")
)

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsError checks if an error is of a specific type
func IsError(err error, target error) bool {
	return errors.Is(err, target)
}
