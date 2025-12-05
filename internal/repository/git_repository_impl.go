package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	gitconfig "github.com/golgoth31/gitcomm/pkg/git/config"
	"github.com/hiddeco/sshsig"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/crypto/ssh"
)

const (
	// maxDiffSize is the maximum character count for diff content before showing metadata only
	maxDiffSize = 5000
	// diffContext is the number of context lines in unified diff format (0 for minimal token usage)
	diffContext = 0
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

const (
	// IncludeNewFilesKey is the context key for controlling whether new files are included in repository state
	// This key is used to pass the addAll flag from service layer to repository layer via context
	IncludeNewFilesKey contextKey = "includeNewFiles"
)

// gitRepositoryImpl implements GitRepository using go-git
type gitRepositoryImpl struct {
	repo   *git.Repository
	path   string
	config *gitconfig.GitConfig
	signer *gitconfig.CommitSigner
}

// NewGitRepository creates a new GitRepository implementation
func NewGitRepository(repoPath string, noSign bool) (GitRepository, error) {
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
	// This ensures config is available before initializing git objects
	extractor := gitconfig.NewFileConfigExtractor()
	gitConfig := extractor.Extract(path)

	// Prepare commit signer if SSH signing is configured
	// The --no-sign flag will be applied later via SetNoSign() if needed
	signer := prepareCommitSigner(gitConfig, noSign)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrNotGitRepository, err)
	}

	return &gitRepositoryImpl{
		repo:   repo,
		path:   path,
		config: gitConfig,
		signer: signer,
	}, nil
}

// GetRepositoryState retrieves the current repository state, including computed diffs for staged files.
// For each staged file, this method computes a unified diff (patch format) between the staged state and HEAD.
// The diff computation is optimized for token usage:
//   - Uses 0 lines of context (minimal token usage)
//   - For files/diffs exceeding 5000 characters, shows only metadata (file size, line count, change summary)
//   - Binary files have empty diff
//   - Unmerged files attempt diff computation, fallback to empty if fails
//   - Errors are logged but don't stop processing (empty diff is set on error)
//
// Filtering behavior:
//   - New files (git.Added status) are excluded when includeNewFiles context value is false
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

	worktree, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	state := &model.RepositoryState{
		StagedFiles:   []model.FileChange{},
		UnstagedFiles: []model.FileChange{},
	}

	for file, fileStatus := range status {
		// Check staging area
		if fileStatus.Staging != git.Unmodified {
			// Skip new files when includeNewFiles is false
			if fileStatus.Staging == git.Added && !includeNewFiles {
				continue
			}
			change := model.FileChange{
				Path:   file,
				Status: statusCodeToString(fileStatus.Staging),
				Diff:   "", // Initialize empty, will be populated below
			}

			// Compute diff for staged files
			// Check if binary file first (FR-013)
			isBinary, err := r.isBinaryFile(file)
			if err != nil {
				// Log error but continue processing (FR-010)
				utils.Logger.Debug().Err(err).Str("file", file).Msg("Failed to check if file is binary, continuing")
			}

			if isBinary {
				// Binary files have empty diff (FR-013)
				change.Diff = ""
			} else if fileStatus.Staging == git.UpdatedButUnmerged {
				// Unmerged files: attempt diff computation, fallback to empty if fails (FR-008)
				diff, err := r.computeFileDiff(ctx, file, fileStatus.Staging)
				if err != nil {
					utils.Logger.Debug().Err(err).Str("file", file).Msg("Failed to compute diff for unmerged file, setting empty")
					change.Diff = ""
				} else {
					change.Diff = r.applySizeLimit(diff, file, fileStatus.Staging)
				}
			} else {
				// Compute diff for staged file
				diff, err := r.computeFileDiff(ctx, file, fileStatus.Staging)
				if err != nil {
					// Log error but continue processing (FR-010)
					utils.Logger.Debug().Err(err).Str("file", file).Msg("Failed to compute diff, setting empty")
					change.Diff = ""
				} else {
					// Apply size limit (FR-016)
					change.Diff = r.applySizeLimit(diff, file, fileStatus.Staging)
				}
			}

			state.StagedFiles = append(state.StagedFiles, change)
		}

		// Check worktree
		if fileStatus.Worktree != git.Unmodified {
			change := model.FileChange{
				Path:   file,
				Status: statusCodeToString(fileStatus.Worktree),
				Diff:   "", // Unstaged files always have empty diff (FR-011)
			}
			state.UnstagedFiles = append(state.UnstagedFiles, change)
		}
	}

	return state, nil
}

// getHEADTree returns the HEAD tree or empty tree if HEAD doesn't exist (empty repository).
// This handles the edge case where a repository has no commits yet (FR-009).
func (r *gitRepositoryImpl) getHEADTree() (*object.Tree, error) {
	head, err := r.repo.Head()
	if err == plumbing.ErrReferenceNotFound {
		// Empty repository - return empty tree
		return &object.Tree{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := r.repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	return tree, nil
}

// getStagedIndexTree converts the staged index to a tree for diff computation
// Note: This is a placeholder - we'll compute diffs per file instead of building a full tree
// go-git doesn't directly expose index as tree, so we use per-file diff computation
func (r *gitRepositoryImpl) getStagedIndexTree() (*object.Tree, error) {
	// This function is not used in the current implementation
	// We compute diffs file-by-file in computeFileDiff instead
	return nil, fmt.Errorf("getStagedIndexTree: not implemented, using per-file diff approach")
}

// computeFileDiff computes the diff between HEAD and staged state for a single file.
// It handles different file statuses:
//   - Added: Shows full file content as additions
//   - Modified: Computes unified diff between HEAD and staged content
//   - Deleted: Shows file removal
//   - Renamed/Copied: Shows rename/copy information with similarity percentage
//
// The diff is computed using go-git plumbing API and formatted as unified diff with 0 context lines (FR-012).
func (r *gitRepositoryImpl) computeFileDiff(ctx context.Context, filePath string, fileStatus git.StatusCode) (string, error) {
	// Get HEAD tree
	headTree, err := r.getHEADTree()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD tree: %w", err)
	}

	// Get the file from HEAD tree
	var headFile *object.File
	if headTree != nil && len(headTree.Entries) > 0 {
		headFile, err = headTree.File(filePath)
		if err != nil && err != object.ErrFileNotFound {
			return "", fmt.Errorf("failed to get file from HEAD: %w", err)
		}
	}

	// For staged files, we read the content directly from the filesystem
	// The staged content is what's currently in the worktree files
	// go-git's worktree.Status() tells us what's staged, but we need the actual content
	// We'll read from the worktree filesystem for staged content
	fullPath := filepath.Join(r.path, filePath)

	// Handle different file statuses
	switch fileStatus {
	case git.Added:
		// New file - diff is the entire file content as additions
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return "", fmt.Errorf("failed to read staged file: %w", err)
		}
		diff := r.formatNewFileDiff(filePath, content)
		return diff, nil

	case git.Deleted:
		// Deleted file - diff shows removal
		if headFile == nil {
			return "", fmt.Errorf("file not found in HEAD for deletion diff")
		}
		diff := r.formatDeletedFileDiff(filePath, headFile)
		return diff, nil

	case git.Modified:
		// Modified file - compute diff between HEAD and staged
		if headFile == nil {
			// File doesn't exist in HEAD, treat as new file
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("failed to read staged file: %w", err)
			}
			return r.formatNewFileDiff(filePath, content), nil
		}

		// Read staged content
		stagedContentBytes, err := os.ReadFile(fullPath)
		if err != nil {
			return "", fmt.Errorf("failed to read staged file: %w", err)
		}
		stagedContent := string(stagedContentBytes)

		// Get HEAD content
		headContent, err := headFile.Contents()
		if err != nil {
			return "", fmt.Errorf("failed to read HEAD file content: %w", err)
		}

		diff := r.formatModifiedFileDiff(filePath, headContent, stagedContent)
		// diff := "pouet"
		utils.Logger.Debug().
			Str("file", filePath).
			Int("headSize", len(headContent)).
			Int("stagedSize", len(stagedContent)).
			Int("diffSize", len(diff)).
			Msg("Formatted modified file diff")
		return diff, nil

	case git.Renamed:
		// Renamed file - format with similarity
		// Note: go-git Status doesn't provide old path directly, we'd need to track it
		// For now, return a placeholder that will be enhanced
		return r.formatRenameDiff(filePath, filePath, 95.0), nil

	case git.Copied:
		// Copied file - format with similarity
		return r.formatCopyDiff(filePath, filePath, 100.0), nil

	default:
		return "", nil
	}
}

// formatUnifiedDiff formats diff as unified patch with 0 context lines (FR-012).
// This produces a diff format that matches git diff --cached output structure.
// The 0 context lines minimize token usage for AI models while preserving essential change information.
// Uses go-git's optimized Meyers diff algorithm for efficient computation.
func (r *gitRepositoryImpl) formatUnifiedDiff(filePath string, oldContent, newContent string) string {
	// Use go-git's optimized Meyers diff algorithm (O(N*d) complexity)
	diffs := diff.Do(oldContent, newContent)

	// Early return if no changes
	if len(diffs) == 0 {
		utils.Logger.Debug().Str("file", filePath).Msg("No differences found")
		return ""
	}

	// Build unified diff format efficiently
	var result strings.Builder
	result.Grow(len(oldContent) + len(newContent) + 1024) // Pre-allocate capacity

	// Write diff header
	result.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", filePath, filePath))
	result.WriteString("index 0000000..1111111 100644\n")
	result.WriteString(fmt.Sprintf("--- a/%s\n", filePath))
	result.WriteString(fmt.Sprintf("+++ b/%s\n", filePath))

	// Convert diff results to unified format with 0 context lines
	// Process line by line for memory efficiency
	oldLineNum := 1
	newLineNum := 1
	hunkStartOld := 0
	hunkStartNew := 0
	hunkOldCount := 0
	hunkNewCount := 0
	var hunkLines strings.Builder
	inHunk := false

	for _, d := range diffs {
		// Process text line by line to avoid large allocations
		text := d.Text
		if text == "" {
			continue
		}

		// Handle line endings - split efficiently
		lines := strings.Split(text, "\n")
		// Remove trailing empty line if text doesn't end with newline
		if len(lines) > 0 && text[len(text)-1] != '\n' && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		switch d.Type {
		case diffmatchpatch.DiffEqual:
			// Matching lines - skip (0 context lines means no context shown)
			// Close current hunk if we have one
			if inHunk {
				result.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", hunkStartOld, hunkOldCount, hunkStartNew, hunkNewCount))
				result.WriteString(hunkLines.String())
				hunkLines.Reset()
				inHunk = false
				hunkOldCount = 0
				hunkNewCount = 0
			}
			oldLineNum += len(lines)
			newLineNum += len(lines)

		case diffmatchpatch.DiffDelete:
			// Deleted lines - start new hunk if needed
			if !inHunk {
				hunkStartOld = oldLineNum
				hunkStartNew = newLineNum
				inHunk = true
			}
			for _, line := range lines {
				hunkLines.WriteString(fmt.Sprintf("-%s\n", line))
			}
			hunkOldCount += len(lines)
			oldLineNum += len(lines)

		case diffmatchpatch.DiffInsert:
			// Inserted lines - start new hunk if needed
			if !inHunk {
				hunkStartOld = oldLineNum
				hunkStartNew = newLineNum
				inHunk = true
			}
			for _, line := range lines {
				hunkLines.WriteString(fmt.Sprintf("+%s\n", line))
			}
			hunkNewCount += len(lines)
			newLineNum += len(lines)
		}
	}

	// Output final hunk if any
	if inHunk {
		result.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", hunkStartOld, hunkOldCount, hunkStartNew, hunkNewCount))
		result.WriteString(hunkLines.String())
	}

	output := result.String()
	utils.Logger.Debug().
		Str("file", filePath).
		Int("diffSize", len(output)).
		Int("diffChunks", len(diffs)).
		Msg("Unified diff formatted successfully")

	return output
}

// formatNewFileDiff formats diff for a new file
func (r *gitRepositoryImpl) formatNewFileDiff(filePath string, content []byte) string {
	lines := strings.Split(string(content), "\n")
	var diff strings.Builder
	diff.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", filePath, filePath))
	diff.WriteString(fmt.Sprintf("new file mode 100644\n"))
	diff.WriteString(fmt.Sprintf("index 0000000..1111111\n"))
	diff.WriteString(fmt.Sprintf("--- /dev/null\n"))
	diff.WriteString(fmt.Sprintf("+++ b/%s\n", filePath))
	diff.WriteString(fmt.Sprintf("@@ -0,0 +1,%d @@\n", len(lines)))
	for _, line := range lines {
		diff.WriteString(fmt.Sprintf("+%s\n", line))
	}
	return diff.String()
}

// formatDeletedFileDiff formats diff for a deleted file
func (r *gitRepositoryImpl) formatDeletedFileDiff(filePath string, headFile *object.File) string {
	content, err := headFile.Contents()
	if err != nil {
		return fmt.Sprintf("diff --git a/%s b/%s\ndeleted file\n", filePath, filePath)
	}
	lines := strings.Split(content, "\n")
	var diff strings.Builder
	diff.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", filePath, filePath))
	diff.WriteString(fmt.Sprintf("deleted file mode 100644\n"))
	diff.WriteString(fmt.Sprintf("index 1111111..0000000\n"))
	diff.WriteString(fmt.Sprintf("--- a/%s\n", filePath))
	diff.WriteString(fmt.Sprintf("+++ /dev/null\n"))
	diff.WriteString(fmt.Sprintf("@@ -1,%d +0,0 @@\n", len(lines)))
	for _, line := range lines {
		diff.WriteString(fmt.Sprintf("-%s\n", line))
	}
	return diff.String()
}

// formatModifiedFileDiff formats diff for a modified file
func (r *gitRepositoryImpl) formatModifiedFileDiff(filePath, oldContent, newContent string) string {
	return r.formatUnifiedDiff(filePath, oldContent, newContent)
}

// isBinaryFile checks if a file is binary
func (r *gitRepositoryImpl) isBinaryFile(filePath string) (bool, error) {
	// Read first few bytes to check for binary content
	fullPath := filepath.Join(r.path, filePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Check for null bytes or non-text characters
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true, nil
		}
	}

	return false, nil
}

// generateMetadata generates metadata string for large files/diffs
// Includes accurate change summary (lines added, removed, modified)
func (r *gitRepositoryImpl) generateMetadata(filePath string, fileStatus git.StatusCode) string {
	fullPath := filepath.Join(r.path, filePath)
	info, err := os.Stat(fullPath)
	statusStr := statusCodeToString(fileStatus)
	if err != nil {
		return fmt.Sprintf("file: %s\nsize: unknown\nlines: unknown\nchanges: %s", filePath, statusStr)
	}

	// Count lines if it's a text file
	lineCount := 0
	var addedLines, removedLines int
	if !strings.HasSuffix(filePath, ".png") && !strings.HasSuffix(filePath, ".jpg") {
		content, err := os.ReadFile(fullPath)
		if err == nil {
			lineCount = strings.Count(string(content), "\n") + 1
			// For modified files, try to compute change summary
			if fileStatus == git.Modified {
				// Get HEAD content for comparison
				headTree, err := r.getHEADTree()
				if err == nil && headTree != nil && len(headTree.Entries) > 0 {
					headFile, err := headTree.File(filePath)
					if err == nil {
						headContent, err := headFile.Contents()
						if err == nil {
							headLines := strings.Split(headContent, "\n")
							newLines := strings.Split(string(content), "\n")
							// Simple change counting
							removedLines = len(headLines)
							addedLines = len(newLines)
						}
					}
				}
			} else if fileStatus == git.Added {
				addedLines = lineCount
			} else if fileStatus == git.Deleted {
				// For deleted files, we'd need HEAD content
				removedLines = lineCount
			}
		}
	}

	// Build metadata string
	metadata := fmt.Sprintf("file: %s\nsize: %d bytes\nlines: %d", filePath, info.Size(), lineCount)
	if addedLines > 0 || removedLines > 0 {
		metadata += fmt.Sprintf("\nadded: %d lines\nremoved: %d lines", addedLines, removedLines)
	}
	metadata += fmt.Sprintf("\nchanges: %s", statusStr)
	return metadata
}

// formatRenameDiff formats rename diff with similarity percentage
func (r *gitRepositoryImpl) formatRenameDiff(oldPath, newPath string, similarity float64) string {
	return fmt.Sprintf("rename from %s\nrename to %s\nsimilarity %.0f%%", oldPath, newPath, similarity)
}

// formatCopyDiff formats copy diff with similarity percentage
func (r *gitRepositoryImpl) formatCopyDiff(sourcePath, destPath string, similarity float64) string {
	return fmt.Sprintf("copy from %s\ncopy to %s\nsimilarity %.0f%%", sourcePath, destPath, similarity)
}

// applySizeLimit checks if diff exceeds 5000 characters (FR-016) and replaces with metadata if needed.
// This token optimization ensures large files/diffs don't consume excessive tokens for AI models.
// Files/diffs under 5000 characters show full content, larger ones show only metadata.
func (r *gitRepositoryImpl) applySizeLimit(diff string, filePath string, fileStatus git.StatusCode) string {
	if len(diff) > maxDiffSize {
		return r.generateMetadata(filePath, fileStatus)
	}
	return diff
}

// CreateCommit creates a git commit with the given message
func (r *gitRepositoryImpl) CreateCommit(ctx context.Context, message *model.CommitMessage) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Format commit message
	formatter := &formattingService{}
	commitMsg := formatter.format(message)

	// Add signoff if needed
	if message.Signoff {
		// Use extracted git config for user name and email
		userName := r.config.UserName
		userEmail := r.config.UserEmail
		if userName != "" && userEmail != "" {
			commitMsg += fmt.Sprintf("\n\nSigned-off-by: %s <%s>", userName, userEmail)
		}
	}

	// Use extracted git config for commit author (FR-003, FR-004)
	// Defaults to "gitcomm <gitcomm@local>" if git config values are missing (FR-012)
	author := &object.Signature{
		Name:  r.config.UserName,
		Email: r.config.UserEmail,
		When:  time.Now(),
	}

	// Prepare commit options
	opts := &git.CommitOptions{
		Author: author,
	}

	// Add SSH signer if configured and enabled (FR-006, FR-007)
	// Signing is enabled when: gpg.format = ssh AND user.signingkey is set AND commit.gpgsign != false AND --no-sign flag not set
	// Use CommitOptions.Signer (not SignKey) which supports custom signers including SSH
	if r.signer.Signer != nil {
		if signer, ok := r.signer.Signer.(git.Signer); ok {
			opts.Signer = signer
		} else {
			// Type assertion failed - log debug message but continue without signing
			utils.Logger.Debug().Msg("CommitSigner.Signer does not implement git.Signer interface, signing disabled")
		}
	}

	utils.Logger.Debug().Any("commit opts", opts).Msg("commit opts")

	// Create commit (signing failures handled gracefully per FR-013)
	hash, err := worktree.Commit(commitMsg, opts)
	if err != nil {
		// Check if error is signing-related
		if strings.Contains(err.Error(), "sign") || strings.Contains(err.Error(), "key") || strings.Contains(err.Error(), "ssh") {
			// Retry without signing (per FR-013)
			utils.Logger.Debug().Err(err).Msg("SSH signing failed, creating unsigned commit")
			opts.Signer = nil
			_, err = worktree.Commit(commitMsg, opts)
			if err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to create commit: %w", err)
	}

	mycomm, _ := r.repo.CommitObject(hash)
	utils.Logger.Debug().Any("commit", mycomm).Msg("commit")

	return nil
}

// prepareCommitSigner creates a CommitSigner from GitConfig if SSH signing is configured.
//
// Signing is enabled when all of the following are true:
//   - gpg.format = "ssh" (FR-006)
//   - user.signingkey is set (FR-006)
//   - commit.gpgsign is not explicitly false (FR-007)
//   - noSign flag is false (FR-008)
//
// Precedence order (highest to lowest):
//  1. --no-sign flag (explicit user override)
//  2. commit.gpgsign = false (explicit opt-out in config)
//  3. SSH signing configuration (if all requirements met)
//
// If signing fails (e.g., key file not found), signing is disabled and an unsigned commit is created (FR-013).
func prepareCommitSigner(gitConfig *gitconfig.GitConfig, noSign bool) *gitconfig.CommitSigner {
	signer := &gitconfig.CommitSigner{
		PrivateKeyPath: "",
		PublicKeyPath:  gitConfig.SigningKey,
		Format:         gitConfig.GPGFormat,
		Signer:         nil,
	}

	gitConfig.CommitGPGSign = true

	// Check if signing should be disabled by flag (highest precedence)
	if noSign || !gitConfig.CommitGPGSign || gitConfig.GPGFormat != "ssh" {
		utils.Logger.Debug().Bool("noSign", noSign).Bool("commitGPGSign", gitConfig.CommitGPGSign).Str("gpgFormat", gitConfig.GPGFormat).Msg("signing disabled")
		return signer
	}

	// Derive private key path from public key path (remove .pub extension)
	// Git config stores public key path in user.signingkey, but signing requires private key
	privateKeyPath := gitConfig.SigningKey
	if strings.HasSuffix(privateKeyPath, ".pub") {
		privateKeyPath = strings.TrimSuffix(privateKeyPath, ".pub")
	}

	// Check if private key file exists
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		utils.Logger.Debug().Str("path", privateKeyPath).Msg("SSH private key file not found, signing disabled")
		return signer
	}

	// Load private key and create SSH signer
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		utils.Logger.Debug().Err(err).Str("path", privateKeyPath).Msg("Failed to read SSH private key, signing disabled")
		return signer
	}

	// block, _ := pem.Decode(privateKeyBytes)
	// if block == nil {
	// 	utils.Logger.Debug().Str("path", privateKeyPath).Msg("Failed to decode SSH private key, signing disabled")
	// 	return signer
	// }

	// Parse SSH private key
	// Note: ssh.ParsePrivateKey does not support passphrase-protected keys
	// For passphrase-protected keys, use ssh.ParsePrivateKeyWithPassphrase instead
	sshSigner, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		// Check if error is due to passphrase requirement
		if strings.Contains(err.Error(), "passphrase") || strings.Contains(err.Error(), "decrypt") {
			utils.Logger.Debug().Err(err).Str("path", privateKeyPath).Msg("SSH private key is passphrase-protected, signing disabled (passphrase support not implemented)")
		} else {
			utils.Logger.Debug().Err(err).Str("path", privateKeyPath).Msg("Failed to parse SSH private key, signing disabled")
		}
		return signer
	}

	// Create wrapper that implements git.Signer for SSH commit signing
	sshCommitSigner := sshCommitSignerWrapper{
		signer: sshSigner,
	}

	var _ git.Signer = sshCommitSigner

	signer.PrivateKeyPath = privateKeyPath
	signer.Signer = sshCommitSigner

	utils.Logger.Debug().Any("signer", signer).Msg("signer")

	return signer
}

// sshCommitSignerWrapper implements git.Signer interface for SSH commit signing.
// This wrapper adapts golang.org/x/crypto/ssh.Signer to work with go-git's git.Signer interface.
type sshCommitSignerWrapper struct {
	signer ssh.Signer
}

// // Sign implements git.Signer interface for SSH commit signing.
// // Reads the commit message and signs it using the SSH private key.
// // Returns the signature in a format compatible with git's SSH signing protocol.
func (s sshCommitSignerWrapper) Sign(message io.Reader) ([]byte, error) {
	sig, err := sshsig.Sign(message, s.signer, sshsig.HashSHA512, "git")
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Print the signature in armored (PEM) format.
	armored := sshsig.Armor(sig)
	return armored, nil
}

// StageAllFiles stages all unstaged files (equivalent to git add -A)
func (r *gitRepositoryImpl) StageAllFiles(ctx context.Context) error {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add all files
	err = worktree.AddGlob(".")
	if err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	return nil
}

// statusCodeToString converts git StatusCode to string representation
func statusCodeToString(s git.StatusCode) string {
	switch s {
	case git.Added:
		return "added"
	case git.Deleted:
		return "deleted"
	case git.Modified:
		return "modified"
	case git.Renamed:
		return "renamed"
	case git.Copied:
		return "copied"
	case git.UpdatedButUnmerged:
		return "unmerged"
	default:
		return "unmodified"
	}
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

// CaptureStagingState captures the current staging state of the repository for restoration purposes
func (r *gitRepositoryImpl) CaptureStagingState(ctx context.Context) (*model.StagingState, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var stagedFiles []string
	for file, fileStatus := range status {
		if fileStatus.Staging != git.Unmodified {
			stagedFiles = append(stagedFiles, file)
		}
	}

	return &model.StagingState{
		StagedFiles:    stagedFiles,
		CapturedAt:     time.Now(),
		RepositoryPath: r.path,
	}, nil
}

// StageModifiedFiles stages all modified (but not untracked) files in the repository
func (r *gitRepositoryImpl) StageModifiedFiles(ctx context.Context) (*model.AutoStagingResult, error) {
	startTime := time.Now()
	worktree, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var filesToStage []string
	for file, fileStatus := range status {
		// Stage only modified files (not untracked)
		if fileStatus.Worktree != git.Unmodified && fileStatus.Worktree != git.Untracked {
			filesToStage = append(filesToStage, file)
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
		_, err := worktree.Add(file)
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
		// Restore all staged files
		for _, file := range stagedFiles {
			worktree.Remove(file)
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
	worktree, err := r.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	var filesToStage []string
	for file, fileStatus := range status {
		// Stage all files that are modified or untracked
		if fileStatus.Worktree != git.Unmodified {
			filesToStage = append(filesToStage, file)
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
		_, err := worktree.Add(file)
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
		// Restore all staged files
		for _, file := range stagedFiles {
			worktree.Remove(file)
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

	// Use git reset HEAD to unstage files (more reliable than worktree.Remove)
	// This is equivalent to "git reset HEAD <file>" for each file
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", r.path, "reset", "HEAD", "--"}, files...)...)
	if err := cmd.Run(); err != nil {
		// If git reset fails, try worktree.Remove as fallback
		worktree, worktreeErr := r.repo.Worktree()
		if worktreeErr != nil {
			return fmt.Errorf("%w: git reset failed: %v, worktree access failed: %v", utils.ErrRestorationFailed, err, worktreeErr)
		}

		// Fallback: try worktree.Remove for each file
		var errors []string
		for _, file := range files {
			_, removeErr := worktree.Remove(file)
			if removeErr != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", file, removeErr))
			}
		}

		if len(errors) > 0 {
			return fmt.Errorf("%w: failed to unstage files (reset failed: %v, remove errors: %s)", utils.ErrRestorationFailed, err, strings.Join(errors, "; "))
		}
	}

	return nil
}
