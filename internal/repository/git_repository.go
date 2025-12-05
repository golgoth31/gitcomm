package repository

import (
	"context"

	"github.com/golgoth31/gitcomm/internal/model"
)

// GitRepository defines the interface for git operations
type GitRepository interface {
	// GetRepositoryState retrieves the current repository state (staged and unstaged changes)
	GetRepositoryState(ctx context.Context) (*model.RepositoryState, error)

	// CreateCommit creates a git commit with the given message
	CreateCommit(ctx context.Context, message *model.CommitMessage) error

	// StageAllFiles stages all unstaged files (equivalent to git add -A)
	StageAllFiles(ctx context.Context) error

	// CaptureStagingState captures the current staging state of the repository for restoration purposes
	CaptureStagingState(ctx context.Context) (*model.StagingState, error)

	// StageModifiedFiles stages all modified (but not untracked) files in the repository
	StageModifiedFiles(ctx context.Context) (*model.AutoStagingResult, error)

	// StageAllFilesIncludingUntracked stages all modified and untracked files in the repository (equivalent to git add -A)
	StageAllFilesIncludingUntracked(ctx context.Context) (*model.AutoStagingResult, error)

	// UnstageFiles unstages the specified files, restoring them to their pre-staged state
	UnstageFiles(ctx context.Context, files []string) error
}
