package repository

import (
	"errors"
	"fmt"
	"strings"
)

// Git CLI error types for categorized error handling (FR-006)

var (
	// ErrGitNotFound indicates git executable is not in PATH
	ErrGitNotFound = errors.New("git executable not found in PATH")

	// ErrGitVersionTooOld indicates git version < 2.34.0
	ErrGitVersionTooOld = errors.New("git version 2.34.0 or higher required")

	// ErrGitPermissionDenied indicates a filesystem permission error
	ErrGitPermissionDenied = errors.New("permission denied for git operation")

	// ErrGitSigningFailed indicates commit signing failed
	ErrGitSigningFailed = errors.New("git commit signing failed")

	// ErrGitFileNotFound indicates a file was not found in the repository
	ErrGitFileNotFound = errors.New("file not found in git repository")
)

// ErrGitCommandFailed is a generic error for git command failures
type ErrGitCommandFailed struct {
	Command  string   // Git subcommand (e.g., "status", "commit", "add")
	Args     []string // Command arguments
	ExitCode int      // Process exit code
	Stderr   string   // Captured stderr output
}

// Error implements the error interface
func (e *ErrGitCommandFailed) Error() string {
	detail := strings.TrimSpace(e.Stderr)
	if detail == "" {
		detail = "No additional details from git. Check repository state or run the command manually."
	}
	return fmt.Sprintf("git %s failed (exit %d): %s", e.Command, e.ExitCode, detail)
}
