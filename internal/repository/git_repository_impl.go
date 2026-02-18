package repository

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	gitconfig "github.com/golgoth31/gitcomm/pkg/git/config"
)

const (
	// maxDiffSize is the maximum character count for diff content before showing metadata only
	maxDiffSize = 5000
	// minGitMajor is the minimum required git major version
	minGitMajor = 2
	// minGitMinor is the minimum required git minor version (for SSH signing support)
	minGitMinor = 34
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

const (
	// IncludeNewFilesKey is the context key for controlling whether new files are included in repository state
	// This key is used to pass the addAll flag from service layer to repository layer via context
	IncludeNewFilesKey contextKey = "includeNewFiles"
)

// gitRepositoryImpl implements GitRepository using external git CLI commands
type gitRepositoryImpl struct {
	path   string                 // Repository root path
	gitBin string                 // Resolved path to git executable
	rtkBin string                 // Resolved path to rtk executable (empty if not found)
	useRTK bool                   // Whether to proxy git commands through rtk
	config *gitconfig.GitConfig   // Git configuration
	signer *gitconfig.CommitSigner // Commit signer configuration
}

// NewGitRepository creates a new GitRepository implementation using external git CLI
func NewGitRepository(repoPath string, noSign bool) (GitRepository, error) {
	// Lookup git executable (FR-016)
	gitBin, err := exec.LookPath("git")
	if err != nil {
		return nil, ErrGitNotFound
	}

	// Validate git version >= 2.34.0 (FR-016)
	if err := validateGitVersion(gitBin); err != nil {
		return nil, err
	}

	// Check if rtk is available for proxying git commands
	var rtkBin string
	var useRTK bool
	if rtkPath, rtkErr := exec.LookPath("rtk"); rtkErr == nil {
		rtkBin = rtkPath
		useRTK = true
		utils.Logger.Debug().Str("rtk", rtkPath).Msg("rtk found, proxying git commands through rtk")
	} else {
		utils.Logger.Debug().Msg("rtk not found, using git directly")
	}

	// Find git repository root
	path := repoPath
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Walk up to find .git directory
	gitPath := path
	for {
		gitDir := filepath.Join(gitPath, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			path = gitPath
			break
		}
		parent := filepath.Dir(gitPath)
		if parent == gitPath {
			return nil, utils.ErrNotGitRepository
		}
		gitPath = parent
	}

	// Extract git config BEFORE opening repository (FR-001, FR-002)
	extractor := gitconfig.NewFileConfigExtractor()
	gitConfig := extractor.Extract(path)

	// Prepare commit signer if SSH signing is configured
	signer := prepareCommitSigner(gitConfig, noSign)

	return &gitRepositoryImpl{
		path:   path,
		gitBin: gitBin,
		rtkBin: rtkBin,
		useRTK: useRTK,
		config: gitConfig,
		signer: signer,
	}, nil
}

// validateGitVersion checks that git version is >= 2.34.0 (required for SSH signing support)
func validateGitVersion(gitBin string) error {
	cmd := exec.Command(gitBin, "--version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: failed to run git --version: %v", ErrGitVersionTooOld, err)
	}

	// Parse "git version X.Y.Z" or "git version X.Y.Z (Apple Git-NNN)"
	output := strings.TrimSpace(stdout.String())
	parts := strings.Fields(output)
	if len(parts) < 3 {
		return fmt.Errorf("%w: unexpected git version output: %s", ErrGitVersionTooOld, output)
	}

	versionStr := parts[2]
	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 2 {
		return fmt.Errorf("%w: unable to parse git version: %s", ErrGitVersionTooOld, versionStr)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("%w: unable to parse major version: %s", ErrGitVersionTooOld, versionStr)
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return fmt.Errorf("%w: unable to parse minor version: %s", ErrGitVersionTooOld, versionStr)
	}

	if major < minGitMajor || (major == minGitMajor && minor < minGitMinor) {
		return fmt.Errorf("%w: found %s, need >= %d.%d.0", ErrGitVersionTooOld, versionStr, minGitMajor, minGitMinor)
	}

	utils.Logger.Debug().Str("version", versionStr).Msg("Git version validated")
	return nil
}

// resolveCommand returns the binary and arguments to use for a git command.
// When rtk is available, it proxies: rtk git -- <args...>
// The "--" separates rtk options from git arguments (e.g. -C).
// Otherwise, it uses git directly: git <args...>
func (r *gitRepositoryImpl) resolveCommand(args []string) (string, []string) {
	if r.useRTK {
		return r.rtkBin, append([]string{"git", "--"}, args...)
	}
	return r.gitBin, args
}

// execGit executes a git command and returns stdout, stderr, and error.
// It logs command execution details for debugging (FR-018) and categorizes errors (FR-006).
func (r *gitRepositoryImpl) execGit(ctx context.Context, args ...string) (string, string, error) {
	// Handle nil context gracefully
	if ctx == nil {
		ctx = context.Background()
	}

	// Prepend -C <path> to run in repo directory
	allArgs := append([]string{"-C", r.path}, args...)

	bin, resolvedArgs := r.resolveCommand(allArgs)
	cmd := exec.CommandContext(ctx, bin, resolvedArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	// Determine git subcommand for logging
	subcommand := ""
	if len(args) > 0 {
		subcommand = args[0]
	}

	// Log execution (FR-018)
	logEvent := utils.Logger.Debug().
		Str("command", subcommand).
		Strs("args", args[1:]).
		Dur("duration", duration)

	if err != nil {
		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		logEvent.Int("exit_code", exitCode).
			Str("stderr", strings.TrimSpace(stderr.String())).
			Msg("git command failed")

		// Categorize the error
		return stdout.String(), stderr.String(), categorizeError(subcommand, args[1:], exitCode, stderr.String())
	}

	logEvent.Int("exit_code", 0).Msg("git command succeeded")
	return stdout.String(), stderr.String(), nil
}

// categorizeError parses stderr and exit code to produce a categorized error type (FR-006)
func categorizeError(command string, args []string, exitCode int, stderr string) error {
	stderrLower := strings.ToLower(stderr)

	// Check for specific error patterns
	if strings.Contains(stderrLower, "not a git repository") {
		return utils.ErrNotGitRepository
	}

	if strings.Contains(stderrLower, "permission denied") {
		return fmt.Errorf("%w: %s", ErrGitPermissionDenied, strings.TrimSpace(stderr))
	}

	if strings.Contains(stderrLower, "signing failed") ||
		strings.Contains(stderrLower, "gpg failed") ||
		strings.Contains(stderrLower, "error signing") {
		return fmt.Errorf("%w: %s", ErrGitSigningFailed, strings.TrimSpace(stderr))
	}

	if strings.Contains(stderrLower, "pathspec") ||
		strings.Contains(stderrLower, "does not exist") {
		return fmt.Errorf("%w: %s", ErrGitFileNotFound, strings.TrimSpace(stderr))
	}

	// Generic command failure
	return &ErrGitCommandFailed{
		Command:  command,
		Args:     args,
		ExitCode: exitCode,
		Stderr:   strings.TrimSpace(stderr),
	}
}

// parseStatus parses `git status --porcelain=v1` output into staged and unstaged file lists.
// Porcelain v1 format: "XY PATH" or "XY ORIG_PATH -> PATH" for renames.
// X = staging area status, Y = worktree status.
func parseStatus(output string) (staged []model.FileChange, unstaged []model.FileChange) {
	staged = []model.FileChange{}
	unstaged = []model.FileChange{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) < 4 {
			// Minimum: "XY P" (2 status chars + space + at least 1 char path)
			continue
		}

		x := line[0] // Staging area status
		y := line[1] // Worktree status
		rawPath := line[3:]

		// Handle renames/copies: "ORIG_PATH -> PATH"
		filePath := rawPath
		if strings.Contains(rawPath, " -> ") {
			parts := strings.SplitN(rawPath, " -> ", 2)
			filePath = parts[1]
		}

		// Staged files: X is not ' ', not '?', not '!'
		if x != ' ' && x != '?' && x != '!' {
			staged = append(staged, model.FileChange{
				Path:   filePath,
				Status: porcelainStatusToString(x),
				Diff:   "",
			})
		}

		// Unstaged/worktree files: Y is not ' '
		if y != ' ' {
			status := porcelainStatusToString(y)
			// Untracked files ('?') are mapped as "added" for worktree display
			if y == '?' {
				status = "added"
			}
			unstaged = append(unstaged, model.FileChange{
				Path:   filePath,
				Status: status,
				Diff:   "", // Unstaged files always have empty diff (FR-011)
			})
		}
	}

	return staged, unstaged
}

// parseDiff parses `git diff --cached --unified=0` output into a per-file diff map.
// Splits on "diff --git" boundaries, detects binary files, returns map[filepath]diffContent.
func parseDiff(output string) map[string]string {
	result := make(map[string]string)

	if strings.TrimSpace(output) == "" {
		return result
	}

	// Split on "diff --git" boundaries
	// Each chunk starts with "diff --git a/... b/..."
	chunks := strings.Split(output, "diff --git ")

	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		// Re-add the header for proper diff format
		fullChunk := "diff --git " + chunk

		// Extract file path from "a/<path> b/<path>"
		firstLine := strings.SplitN(chunk, "\n", 2)[0]
		filePath := extractPathFromDiffHeader(firstLine)
		if filePath == "" {
			continue
		}

		// Check for binary file
		if strings.Contains(fullChunk, "Binary files") || strings.Contains(fullChunk, "GIT binary patch") {
			result[filePath] = "" // Binary files have empty diff
			continue
		}

		result[filePath] = fullChunk
	}

	return result
}

// extractPathFromDiffHeader extracts the file path from "a/<path> b/<path>" header line
func extractPathFromDiffHeader(header string) string {
	// Header format: "a/<path> b/<path>"
	parts := strings.SplitN(header, " b/", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

// isBinaryFile checks if a file is binary by reading the first 512 bytes
// and checking for NUL bytes or known binary file extensions
func (r *gitRepositoryImpl) isBinaryFile(filePath string) bool {
	// Check known binary extensions first
	ext := strings.ToLower(filepath.Ext(filePath))
	binaryExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".bmp": true,
		".ico": true, ".webp": true, ".svg": true, ".tiff": true, ".tif": true,
		".exe": true, ".dll": true, ".so": true, ".dylib": true,
		".zip": true, ".tar": true, ".gz": true, ".bz2": true, ".xz": true,
		".7z": true, ".rar": true, ".jar": true, ".war": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true, ".odt": true, ".ods": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mkv": true, ".mov": true,
		".wav": true, ".flac": true, ".ogg": true, ".webm": true,
		".ttf": true, ".otf": true, ".woff": true, ".woff2": true, ".eot": true,
		".class": true, ".pyc": true, ".pyo": true, ".o": true, ".a": true,
		".wasm": true, ".bin": true, ".dat": true, ".db": true, ".sqlite": true,
	}
	if binaryExts[ext] {
		return true
	}

	// Read first 512 bytes and check for NUL bytes
	fullPath := filepath.Join(r.path, filePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return false
	}

	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true
		}
	}

	return false
}

// porcelainStatusToString converts a porcelain status character to string representation
func porcelainStatusToString(c byte) string {
	switch c {
	case 'A':
		return "added"
	case 'D':
		return "deleted"
	case 'M':
		return "modified"
	case 'R':
		return "renamed"
	case 'C':
		return "copied"
	case 'U':
		return "unmerged"
	case '?':
		return "untracked"
	default:
		return "unmodified"
	}
}

// GetRepositoryState retrieves the current repository state, including computed diffs for staged files.
// For each staged file, this method uses git diff to compute a unified diff (patch format) between staged and HEAD.
// The diff computation is optimized for token usage:
//   - Uses 0 lines of context (minimal token usage)
//   - For files/diffs exceeding 5000 characters, shows only metadata (file size, line count, change summary)
//   - Binary files have empty diff
//   - Errors are logged but don't stop processing (empty diff is set on error)
//
// Filtering behavior:
//   - New files (added status) are excluded when includeNewFiles context value is false
//   - Modified, deleted, renamed files are always included regardless of flag
//   - When context value is not present, defaults to including all files (backward compatible)
//
// Unstaged files always have empty diff field (FR-011).
func (r *gitRepositoryImpl) GetRepositoryState(ctx context.Context) (*model.RepositoryState, error) {
	// Extract includeNewFiles from context (default: true for backward compatibility)
	includeNewFiles := true
	if val := ctx.Value(IncludeNewFilesKey); val != nil {
		if includeNewFilesVal, ok := val.(bool); ok {
			includeNewFiles = includeNewFilesVal
		}
	}

	// Get status
	statusOut, _, err := r.execGit(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	staged, unstaged := parseStatus(statusOut)

	// Get diffs for staged files
	diffOut, _, err := r.execGit(ctx, "diff", "--cached", "--unified=0")
	if err != nil {
		// Log error but continue with empty diffs
		utils.Logger.Debug().Err(err).Msg("Failed to get staged diffs, continuing with empty diffs")
		diffOut = ""
	}

	diffs := parseDiff(diffOut)

	// Apply diffs and filtering to staged files
	state := &model.RepositoryState{
		StagedFiles:   []model.FileChange{},
		UnstagedFiles: unstaged,
	}

	for _, file := range staged {
		// Skip new files when includeNewFiles is false
		if file.Status == "added" && !includeNewFiles {
			continue
		}

		// Check if binary file first (FR-013)
		if r.isBinaryFile(file.Path) {
			file.Diff = "" // Binary files have empty diff
		} else if diff, ok := diffs[file.Path]; ok {
			// Assign diff from parsed output
			file.Diff = r.applySizeLimit(diff, file.Path, file.Status)
		}

		state.StagedFiles = append(state.StagedFiles, file)
	}

	return state, nil
}

// CaptureStagingState captures the current staging state of the repository for restoration purposes
func (r *gitRepositoryImpl) CaptureStagingState(ctx context.Context) (*model.StagingState, error) {
	statusOut, _, err := r.execGit(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	staged, _ := parseStatus(statusOut)

	var stagedFiles []string
	for _, file := range staged {
		stagedFiles = append(stagedFiles, file.Path)
	}

	return &model.StagingState{
		StagedFiles:    stagedFiles,
		CapturedAt:     time.Now(),
		RepositoryPath: r.path,
	}, nil
}

// CreateCommit creates a git commit with the given message
func (r *gitRepositoryImpl) CreateCommit(ctx context.Context, message *model.CommitMessage) error {
	// Format commit message
	formatter := &formattingService{}
	commitMsg := formatter.format(message)

	// Add signoff if needed
	if message.Signoff {
		userName := r.config.UserName
		userEmail := r.config.UserEmail
		if userName != "" && userEmail != "" {
			commitMsg += fmt.Sprintf("\n\nSigned-off-by: %s <%s>", userName, userEmail)
		}
	}

	// Build commit command with author env vars
	commitEnv := append(os.Environ(),
		"GIT_AUTHOR_NAME="+r.config.UserName,
		"GIT_AUTHOR_EMAIL="+r.config.UserEmail,
		"GIT_COMMITTER_NAME="+r.config.UserName,
		"GIT_COMMITTER_EMAIL="+r.config.UserEmail,
	)

	// If signing is enabled, try signed commit first
	if r.signer.Enabled {
		signArgs := []string{
			"-c", "gpg.format=ssh",
			"-c", "user.signingkey=" + r.signer.PublicKeyPath,
			"-c", "commit.gpgsign=true",
			"commit", "-S", "-m", commitMsg,
		}

		err := r.execGitWithEnv(ctx, commitEnv, signArgs...)
		if err != nil {
			// Check if error is signing-related; if so, retry without signing
			errStr := err.Error()
			if strings.Contains(errStr, "signing") ||
				strings.Contains(errStr, "gpg") ||
				strings.Contains(errStr, "sign") {
				utils.Logger.Debug().Err(err).Msg("SSH signing failed, creating unsigned commit")
			} else {
				return fmt.Errorf("failed to create signed commit: %w", err)
			}
		} else {
			return nil // Signed commit succeeded
		}
	}

	// Unsigned commit (or signing fallback)
	unsignedArgs := []string{"commit", "-m", commitMsg}
	if err := r.execGitWithEnv(ctx, commitEnv, unsignedArgs...); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	return nil
}

// execGitWithEnv executes a git command with custom environment variables.
// Used for commit commands that need GIT_AUTHOR_NAME/EMAIL and signing config.
func (r *gitRepositoryImpl) execGitWithEnv(ctx context.Context, env []string, args ...string) error {
	// Handle nil context gracefully
	if ctx == nil {
		ctx = context.Background()
	}

	allArgs := append([]string{"-C", r.path}, args...)

	bin, resolvedArgs := r.resolveCommand(allArgs)
	cmd := exec.CommandContext(ctx, bin, resolvedArgs...)
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	startTime := time.Now()
	err := cmd.Run()
	duration := time.Since(startTime)

	// Determine subcommand for logging
	subcommand := ""
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			subcommand = arg
			break
		}
	}

	logEvent := utils.Logger.Debug().
		Str("command", subcommand).
		Strs("args", args).
		Dur("duration", duration)

	if err != nil {
		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		logEvent.Int("exit_code", exitCode).
			Str("stderr", strings.TrimSpace(stderr.String())).
			Msg("git command failed")
		return categorizeError(subcommand, args, exitCode, stderr.String())
	}

	logEvent.Int("exit_code", 0).Msg("git command succeeded")
	return nil
}

// StageAllFiles stages all unstaged files (equivalent to git add -A)
func (r *gitRepositoryImpl) StageAllFiles(ctx context.Context) error {
	_, _, err := r.execGit(ctx, "add", "-A")
	if err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}
	return nil
}

// StageModifiedFiles stages all modified (but not untracked) files in the repository
func (r *gitRepositoryImpl) StageModifiedFiles(ctx context.Context) (*model.AutoStagingResult, error) {
	startTime := time.Now()

	// Get current status
	statusOut, _, err := r.execGit(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	// Filter modified files (not untracked) from worktree
	var filesToStage []string
	lines := strings.Split(statusOut, "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		y := line[1]
		// Stage only modified worktree files (not untracked '?' or unmodified ' ')
		if y != ' ' && y != '?' {
			rawPath := line[3:]
			if strings.Contains(rawPath, " -> ") {
				parts := strings.SplitN(rawPath, " -> ", 2)
				rawPath = parts[1]
			}
			filesToStage = append(filesToStage, rawPath)
		}
	}

	if len(filesToStage) == 0 {
		return &model.AutoStagingResult{
			StagedFiles: []string{},
			FailedFiles: []model.StagingFailure{},
			Success:     true,
			Duration:    time.Since(startTime),
		}, nil
	}

	var stagedFiles []string
	var failedFiles []model.StagingFailure

	for _, file := range filesToStage {
		_, _, err := r.execGit(ctx, "add", "--", file)
		if err != nil {
			failedFiles = append(failedFiles, model.StagingFailure{
				FilePath:  file,
				Error:     err,
				ErrorType: "other",
			})
		} else {
			stagedFiles = append(stagedFiles, file)
		}
	}

	// If any file failed, abort and restore
	if len(failedFiles) > 0 {
		if len(stagedFiles) > 0 {
			rollbackArgs := append([]string{"reset", "HEAD", "--"}, stagedFiles...)
			_, _, _ = r.execGit(ctx, rollbackArgs...)
		}
		return &model.AutoStagingResult{
			StagedFiles: []string{},
			FailedFiles: failedFiles,
			Success:     false,
			Duration:    time.Since(startTime),
		}, fmt.Errorf("%w: failed to stage %d file(s)", utils.ErrStagingFailed, len(failedFiles))
	}

	return &model.AutoStagingResult{
		StagedFiles: stagedFiles,
		FailedFiles: []model.StagingFailure{},
		Success:     true,
		Duration:    time.Since(startTime),
	}, nil
}

// StageAllFilesIncludingUntracked stages all modified and untracked files in the repository (equivalent to git add -A)
func (r *gitRepositoryImpl) StageAllFilesIncludingUntracked(ctx context.Context) (*model.AutoStagingResult, error) {
	startTime := time.Now()

	// Get current status
	statusOut, _, err := r.execGit(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	// Filter all changed files from worktree (including untracked)
	var filesToStage []string
	lines := strings.Split(statusOut, "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		y := line[1]
		// Stage all worktree files that are not unmodified
		if y != ' ' {
			rawPath := line[3:]
			if strings.Contains(rawPath, " -> ") {
				parts := strings.SplitN(rawPath, " -> ", 2)
				rawPath = parts[1]
			}
			filesToStage = append(filesToStage, rawPath)
		}
	}

	if len(filesToStage) == 0 {
		return &model.AutoStagingResult{
			StagedFiles: []string{},
			FailedFiles: []model.StagingFailure{},
			Success:     true,
			Duration:    time.Since(startTime),
		}, nil
	}

	var stagedFiles []string
	var failedFiles []model.StagingFailure

	for _, file := range filesToStage {
		_, _, err := r.execGit(ctx, "add", "--", file)
		if err != nil {
			failedFiles = append(failedFiles, model.StagingFailure{
				FilePath:  file,
				Error:     err,
				ErrorType: "other",
			})
		} else {
			stagedFiles = append(stagedFiles, file)
		}
	}

	// If any file failed, abort and restore
	if len(failedFiles) > 0 {
		if len(stagedFiles) > 0 {
			rollbackArgs := append([]string{"reset", "HEAD", "--"}, stagedFiles...)
			_, _, _ = r.execGit(ctx, rollbackArgs...)
		}
		return &model.AutoStagingResult{
			StagedFiles: []string{},
			FailedFiles: failedFiles,
			Success:     false,
			Duration:    time.Since(startTime),
		}, fmt.Errorf("%w: failed to stage %d file(s)", utils.ErrStagingFailed, len(failedFiles))
	}

	return &model.AutoStagingResult{
		StagedFiles: stagedFiles,
		FailedFiles: []model.StagingFailure{},
		Success:     true,
		Duration:    time.Since(startTime),
	}, nil
}

// UnstageFiles unstages the specified files, restoring them to their pre-staged state
func (r *gitRepositoryImpl) UnstageFiles(ctx context.Context, files []string) error {
	if len(files) == 0 {
		return nil
	}

	// Use git reset HEAD to unstage files
	resetArgs := append([]string{"reset", "HEAD", "--"}, files...)
	_, _, err := r.execGit(ctx, resetArgs...)
	if err != nil {
		return fmt.Errorf("%w: git reset failed: %v", utils.ErrRestorationFailed, err)
	}

	return nil
}

// prepareCommitSigner creates a CommitSigner from GitConfig if SSH signing is configured.
//
// Signing is enabled when all of the following are true:
//   - gpg.format = "ssh"
//   - user.signingkey is set
//   - commit.gpgsign is not explicitly false
//   - noSign flag is false
//
// Signing is delegated to git CLI, so no private key loading is needed here.
func prepareCommitSigner(gitConfig *gitconfig.GitConfig, noSign bool) *gitconfig.CommitSigner {
	signer := &gitconfig.CommitSigner{
		PrivateKeyPath: "",
		PublicKeyPath:  gitConfig.SigningKey,
		Format:         gitConfig.GPGFormat,
		Enabled:        false,
	}

	gitConfig.CommitGPGSign = true

	// Check if signing should be disabled by flag (highest precedence)
	if noSign || !gitConfig.CommitGPGSign || gitConfig.GPGFormat != "ssh" {
		utils.Logger.Debug().Bool("noSign", noSign).Bool("commitGPGSign", gitConfig.CommitGPGSign).Str("gpgFormat", gitConfig.GPGFormat).Msg("signing disabled")
		return signer
	}

	// Check if signing key is configured
	if gitConfig.SigningKey == "" {
		utils.Logger.Debug().Msg("No signing key configured, signing disabled")
		return signer
	}

	// Derive private key path from public key path (remove .pub extension)
	privateKeyPath := strings.TrimSuffix(gitConfig.SigningKey, ".pub")

	signer.PrivateKeyPath = privateKeyPath
	signer.Enabled = true

	utils.Logger.Debug().
		Str("publicKey", signer.PublicKeyPath).
		Str("privateKey", signer.PrivateKeyPath).
		Str("format", signer.Format).
		Bool("enabled", signer.Enabled).
		Msg("SSH commit signing configured (delegated to git CLI)")

	return signer
}

// generateMetadata generates metadata string for large files/diffs
func (r *gitRepositoryImpl) generateMetadata(filePath string, status string) string {
	fullPath := filepath.Join(r.path, filePath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Sprintf("file: %s\nsize: unknown\nlines: unknown\nchanges: %s", filePath, status)
	}

	// Count lines if it's a text file
	lineCount := 0
	content, err := os.ReadFile(fullPath)
	if err == nil {
		lineCount = strings.Count(string(content), "\n") + 1
	}

	return fmt.Sprintf("file: %s\nsize: %d bytes\nlines: %d\nchanges: %s", filePath, info.Size(), lineCount, status)
}

// applySizeLimit checks if diff exceeds 5000 characters and replaces with metadata if needed.
// This token optimization ensures large files/diffs don't consume excessive tokens for AI models.
func (r *gitRepositoryImpl) applySizeLimit(diff string, filePath string, status string) string {
	if len(diff) > maxDiffSize {
		return r.generateMetadata(filePath, status)
	}
	return diff
}

// formattingService is a temporary helper for formatting
// TODO: Use the actual service once dependency injection is set up
type formattingService struct{}

func (f *formattingService) format(message *model.CommitMessage) string {
	var parts []string

	header := message.Type
	if message.Scope != "" {
		header = fmt.Sprintf("%s(%s)", header, message.Scope)
	}
	header = fmt.Sprintf("%s: %s", header, message.Subject)
	parts = append(parts, header)

	if message.Body != "" {
		parts = append(parts, "")
		parts = append(parts, message.Body)
	}

	if message.Footer != "" {
		parts = append(parts, "")
		parts = append(parts, message.Footer)
	}

	return strings.Join(parts, "\n")
}
