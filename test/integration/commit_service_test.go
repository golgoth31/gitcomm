package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/golgoth31/gitcomm/internal/repository"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// TestCommitService_ExcludesNewFilesWithoutAddAllFlag verifies that when AutoStage is false,
// new files are excluded from the repository state used for commit message generation.
func TestCommitService_ExcludesNewFilesWithoutAddAllFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup: Initialize logger
	utils.InitLogger(false)

	// Create temporary directory for test repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-exclude-new-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user (required for commits)
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create and commit initial file
	existingFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", existingFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage existing file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Modify existing file
	if err := os.WriteFile(existingFile, []byte("modified content\n"), 0644); err != nil {
		t.Fatalf("Failed to modify existing file: %v", err)
	}

	// Create new file
	newFile := filepath.Join(tmpDir, "newfile.txt")
	if err := os.WriteFile(newFile, []byte("new content\n"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Create repository
	gitRepo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Get repository state directly to verify filtering
	ctx := context.Background()
	state, err := gitRepo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify new file is not in repository state (it's not staged, so it shouldn't be there anyway)
	// But if it were staged, it should be excluded when AutoStage is false
	// Since we're not auto-staging, let's manually stage both files to test the filtering
	cmd = exec.Command("git", "-C", tmpDir, "add", existingFile, newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	// Now get repository state with includeNewFiles = false (simulating AutoStage = false)
	// We need to set the context value manually since we're not going through CreateCommit
	ctx = context.WithValue(ctx, repository.IncludeNewFilesKey, false)
	state, err = gitRepo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify modified file is included
	foundModifiedFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "existing.txt" {
			foundModifiedFile = true
			break
		}
	}

	if !foundModifiedFile {
		t.Error("Expected modified file (existing.txt) to be included in repository state, but it was not found")
	}

	// Verify new file is excluded
	foundNewFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "newfile.txt" {
			foundNewFile = true
			break
		}
	}

	if foundNewFile {
		t.Error("Expected new file (newfile.txt) to be excluded from repository state when AutoStage is false, but it was included")
	}
}

// TestCommitService_IncludesNewFilesWithAddAllFlag verifies that when AutoStage is true,
// new files are included in the repository state used for commit message generation.
func TestCommitService_IncludesNewFilesWithAddAllFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup: Initialize logger
	utils.InitLogger(false)

	// Create temporary directory for test repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-include-new-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user (required for commits)
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create and commit initial file
	existingFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", existingFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage existing file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Modify existing file
	if err := os.WriteFile(existingFile, []byte("modified content\n"), 0644); err != nil {
		t.Fatalf("Failed to modify existing file: %v", err)
	}

	// Create new file
	newFile := filepath.Join(tmpDir, "newfile.txt")
	if err := os.WriteFile(newFile, []byte("new content\n"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Manually stage both files to test filtering behavior
	cmd = exec.Command("git", "-C", tmpDir, "add", existingFile, newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	// Create repository
	gitRepo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Get repository state with includeNewFiles = true (simulating AutoStage = true)
	ctx := context.WithValue(context.Background(), repository.IncludeNewFilesKey, true)
	state, err := gitRepo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify modified file is included
	foundModifiedFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "existing.txt" {
			foundModifiedFile = true
			break
		}
	}

	if !foundModifiedFile {
		t.Error("Expected modified file (existing.txt) to be included in repository state, but it was not found")
	}

	// Verify new file is included when AutoStage is true
	foundNewFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "newfile.txt" {
			foundNewFile = true
			break
		}
	}

	if !foundNewFile {
		t.Error("Expected new file (newfile.txt) to be included in repository state when AutoStage is true, but it was excluded")
	}
}
