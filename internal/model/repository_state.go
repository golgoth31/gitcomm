package model

// RepositoryState represents the current state of the git repository for commit message generation
type RepositoryState struct {
	// StagedFiles is the list of staged file changes
	StagedFiles []FileChange

	// UnstagedFiles is the list of unstaged file changes
	UnstagedFiles []FileChange

	// RawDiff is the condensed diff output from rtk (when rtk is active).
	// When non-empty, this replaces per-file FileChange.Diff for AI prompt generation.
	RawDiff string
}

// FileChange represents a single file change in the repository
type FileChange struct {
	// Path is the file path relative to repository root
	Path string

	// Status is the change status (added, modified, deleted, renamed)
	Status string

	// Diff is the optional unified diff content for the change
	Diff string
}

// IsEmpty returns true if there are no staged or unstaged changes
func (r *RepositoryState) IsEmpty() bool {
	return len(r.StagedFiles) == 0 && len(r.UnstagedFiles) == 0
}

// HasChanges returns true if there are staged or unstaged changes
func (r *RepositoryState) HasChanges() bool {
	return !r.IsEmpty()
}
