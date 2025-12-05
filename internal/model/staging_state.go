package model

import (
	"fmt"
	"time"

	"github.com/golgoth31/gitcomm/internal/utils"
)

// StagingState represents a snapshot of the git repository staging state at a specific point in time
type StagingState struct {
	// StagedFiles is the list of file paths that are staged
	StagedFiles []string

	// CapturedAt is the timestamp when state was captured
	CapturedAt time.Time

	// RepositoryPath is the path to the git repository root
	RepositoryPath string
}

// IsEmpty returns true if no files are staged
func (s *StagingState) IsEmpty() bool {
	return len(s.StagedFiles) == 0
}

// Contains returns true if the specified file is staged
func (s *StagingState) Contains(file string) bool {
	for _, stagedFile := range s.StagedFiles {
		if stagedFile == file {
			return true
		}
	}
	return false
}

// Diff returns files in this state but not in other (for restoration)
func (s *StagingState) Diff(other *StagingState) []string {
	if other == nil {
		return s.StagedFiles
	}

	otherSet := make(map[string]bool)
	for _, file := range other.StagedFiles {
		otherSet[file] = true
	}

	var diff []string
	for _, file := range s.StagedFiles {
		if !otherSet[file] {
			diff = append(diff, file)
		}
	}

	return diff
}

// AutoStagingResult represents the result of an automatic staging operation
type AutoStagingResult struct {
	// StagedFiles is the list of file paths successfully staged
	StagedFiles []string

	// FailedFiles is the list of files that failed to stage
	FailedFiles []StagingFailure

	// Success is the overall success status (true if all files staged)
	Success bool

	// Duration is the time taken for staging operation
	Duration time.Duration
}

// HasFailures returns true if any files failed to stage
func (r *AutoStagingResult) HasFailures() bool {
	return len(r.FailedFiles) > 0
}

// GetFailedFilePaths returns the list of failed file paths
func (r *AutoStagingResult) GetFailedFilePaths() []string {
	var paths []string
	for _, failure := range r.FailedFiles {
		paths = append(paths, failure.FilePath)
	}
	return paths
}

// StagingFailure represents a single file staging failure
type StagingFailure struct {
	// FilePath is the path to the file that failed to stage
	FilePath string

	// Error is the error that occurred during staging
	Error error

	// ErrorType is the type of error (permission, locked, conflict, etc.)
	ErrorType string
}

// RestorationPlan represents the plan for restoring staging state to pre-CLI state
type RestorationPlan struct {
	// FilesToUnstage is the list of file paths to unstage (files staged by CLI)
	FilesToUnstage []string

	// PreCLIState is the captured pre-CLI staging state
	PreCLIState *StagingState

	// CurrentState is the current staging state (for validation)
	CurrentState *StagingState
}

// IsEmpty returns true if no restoration is needed
func (p *RestorationPlan) IsEmpty() bool {
	return len(p.FilesToUnstage) == 0
}

// Validate validates that restoration plan is valid
func (p *RestorationPlan) Validate() error {
	if p.PreCLIState == nil {
		return fmt.Errorf("%w: pre-CLI state is nil", utils.ErrRestorationPlanInvalid)
	}

	if p.CurrentState == nil {
		return fmt.Errorf("%w: current state is nil", utils.ErrRestorationPlanInvalid)
	}

	// Verify FilesToUnstage is a subset of CurrentState.StagedFiles
	currentSet := make(map[string]bool)
	for _, file := range p.CurrentState.StagedFiles {
		currentSet[file] = true
	}

	for _, file := range p.FilesToUnstage {
		if !currentSet[file] {
			return fmt.Errorf("%w: file %s is not currently staged", utils.ErrRestorationPlanInvalid, file)
		}
	}

	return nil
}

// GetFilesToUnstage returns the list of files to unstage
func (p *RestorationPlan) GetFilesToUnstage() []string {
	return p.FilesToUnstage
}
