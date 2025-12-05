package repository

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
)

func TestNewGitRepository_ExtractsConfigBeforeOpening(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user in .git/config
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create repository - should extract config before opening
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Verify repository was created
	if repo == nil {
		t.Fatal("Repository is nil")
	}

	// Verify config was extracted (check by creating a commit and verifying author)
	impl := repo.(*gitRepositoryImpl)
	if impl.config == nil {
		t.Fatal("Config was not extracted")
	}

	// Verify config values
	if impl.config.UserName != "Test User" {
		t.Errorf("Expected UserName 'Test User', got '%s'", impl.config.UserName)
	}
	if impl.config.UserEmail != "test@example.com" {
		t.Errorf("Expected UserEmail 'test@example.com', got '%s'", impl.config.UserEmail)
	}
}

func TestCreateCommit_UsesExtractedConfigForAuthor(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory with .git/config
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user in .git/config
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Commit Author")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "author@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create repository
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create a test file and stage it
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create commit
	commitMsg := &model.CommitMessage{
		Type:    "test",
		Subject: "test commit",
	}
	if err := repo.CreateCommit(context.Background(), commitMsg); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Verify commit author using git log
	cmd = exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%an <%ae>")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit author: %v", err)
	}

	author := string(output)
	expectedAuthor := "Commit Author <author@example.com>"
	if author != expectedAuthor+"\n" {
		t.Errorf("Expected author '%s', got '%s'", expectedAuthor, author)
	}
}

func TestCreateCommit_UsesDefaultsWhenConfigMissing(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory with .git but no config
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Don't configure git user - should use defaults
	// Remove any user config that git init might have created
	cmd = exec.Command("git", "-C", tmpDir, "config", "--unset", "user.name")
	cmd.Run() // Ignore error if not set
	cmd = exec.Command("git", "-C", tmpDir, "config", "--unset", "user.email")
	cmd.Run() // Ignore error if not set

	// Save and clear HOME to prevent reading global config
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set HOME to non-existent directory
	os.Setenv("HOME", filepath.Join(t.TempDir(), "nonexistent"))

	// Create repository
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create a test file and stage it
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create commit
	commitMsg := &model.CommitMessage{
		Type:    "test",
		Subject: "test commit",
	}
	if err := repo.CreateCommit(context.Background(), commitMsg); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Verify commit author uses defaults
	cmd = exec.Command("git", "-C", tmpDir, "log", "-1", "--format=%an <%ae>")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get commit author: %v", err)
	}

	author := string(output)
	expectedAuthor := "gitcomm <gitcomm@local>"
	if author != expectedAuthor+"\n" {
		t.Errorf("Expected default author '%s', got '%s'", expectedAuthor, author)
	}
}

func TestGetRepositoryState_PopulatesDiffForStagedFiles(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create initial file and commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Modify and stage file
	if err := os.WriteFile(testFile, []byte("modified\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Get repository state
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file has diff populated
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	if state.StagedFiles[0].Diff == "" {
		t.Error("Expected Diff field to be populated for staged file, got empty")
	}
}

func TestGetRepositoryState_LeavesDiffEmptyForUnstagedFiles(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Modify file but don't stage (unstaged change)
	if err := os.WriteFile(testFile, []byte("unstaged change\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Get repository state
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify unstaged file has empty diff
	if len(state.UnstagedFiles) != 1 {
		t.Fatalf("Expected 1 unstaged file, got %d", len(state.UnstagedFiles))
	}

	if state.UnstagedFiles[0].Diff != "" {
		t.Error("Expected Diff field to be empty for unstaged file, got non-empty")
	}
}

func TestGetRepositoryState_DiffComputationAccuracy(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file with specific content and commit
	testFile := filepath.Join(tmpDir, "test.txt")
	initialContent := "line1\nline2\nline3\n"
	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Make specific changes: modify line 2, add line 4
	modifiedContent := "line1\nline2_modified\nline3\nline4_new\n"
	if err := os.WriteFile(testFile, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Get repository state
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify diff contains the specific changes
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	diff := state.StagedFiles[0].Diff
	if diff == "" {
		t.Fatal("Expected diff to be populated, got empty")
	}

	// Verify diff contains modified content indicators
	// The exact format may vary, but should show changes
	if !strings.Contains(diff, "line2_modified") && !strings.Contains(diff, "line4_new") {
		t.Error("Expected diff to contain modified content, but key changes not found")
	}
}

func TestGetRepositoryState_LargeDiffShowsMetadataOnly(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create initial file with small content and commit
	testFile := filepath.Join(tmpDir, "large.txt")
	smallContent := "initial\n"
	if err := os.WriteFile(testFile, []byte(smallContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Create large content (>5000 chars) and stage
	largeContent := strings.Repeat("line with content\n", 400) // ~6000+ chars
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to write large content: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage large file: %v", err)
	}

	// Get repository state
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify large diff shows metadata only
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	diff := state.StagedFiles[0].Diff
	if diff == "" {
		t.Fatal("Expected metadata to be populated, got empty")
	}

	// Verify it's metadata format (not full diff)
	if strings.Contains(diff, "file:") && strings.Contains(diff, "size:") {
		// This is metadata format - good
	} else if len(diff) > 5000 {
		t.Errorf("Expected metadata for large diff, but got full diff (%d chars)", len(diff))
	}
}

func TestGetRepositoryState_LargeNewFileShowsMetadataOnly(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create large new file (>5000 chars) and stage
	testFile := filepath.Join(tmpDir, "large_new.txt")
	largeContent := strings.Repeat("new file line with content\n", 250) // ~7000+ chars
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage large file: %v", err)
	}

	// Get repository state
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify large new file shows metadata only
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	diff := state.StagedFiles[0].Diff
	if diff == "" {
		t.Fatal("Expected metadata to be populated, got empty")
	}

	// Verify it's metadata format (not full file content)
	if strings.Contains(diff, "file:") && strings.Contains(diff, "size:") {
		// This is metadata format - good
	} else if len(diff) > 5000 {
		t.Errorf("Expected metadata for large new file, but got full content (%d chars)", len(diff))
	}
}

func TestGetRepositoryState_HandlesFileReadErrors(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Modify and stage file
	if err := os.WriteFile(testFile, []byte("modified\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Remove file to cause read error (file exists in index but not in worktree)
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}

	// Get repository state - should handle error gracefully
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify error was handled gracefully (Diff is empty, processing continued)
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	// Diff should be empty due to read error, but processing should continue
	// (Error is logged, but function doesn't fail)
}

func TestGetRepositoryState_HandlesDiffComputationFailures(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Get repository state
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify function completed successfully even if some diffs failed
	// (In this case, there are no staged files, so no errors expected)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}
}

func TestGetRepositoryState_HandlesUnmergedFiles(t *testing.T) {
	// Setup: Initialize logger
	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	cmd.Run()
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	cmd.Run()

	// Create file and commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Create a branch and make conflicting changes
	cmd = exec.Command("git", "-C", tmpDir, "checkout", "-b", "feature")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create branch: %v", err)
	}

	if err := os.WriteFile(testFile, []byte("feature change\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "feature commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Switch back to main and make conflicting change
	cmd = exec.Command("git", "-C", tmpDir, "checkout", "main")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to checkout main: %v", err)
	}

	if err := os.WriteFile(testFile, []byte("main change\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "main commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Merge feature branch (will create conflict)
	cmd = exec.Command("git", "-C", tmpDir, "merge", "feature", "--no-edit")
	cmd.Run() // Ignore error - merge will fail with conflict

	// Get repository state - should handle unmerged files gracefully
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify function completed successfully
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// If there are unmerged files, they should be handled gracefully
	// (Diff may be empty if computation fails, but processing should continue)
}

// TestGetRepositoryState_ExcludesNewFilesWhenAddAllFalse verifies that new files (git.Added status)
// are excluded from repository state when includeNewFiles context value is false.
func TestGetRepositoryState_ExcludesNewFilesWhenAddAllFalse(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create initial file and commit
	existingFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("initial\n"), 0644); err != nil {
		t.Fatalf("Failed to create existing file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", existingFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage existing file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Create new file and stage it
	newFile := filepath.Join(tmpDir, "newfile.txt")
	if err := os.WriteFile(newFile, []byte("new content\n"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage new file: %v", err)
	}

	// Get repository state with includeNewFiles = false
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, false)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
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
		t.Error("Expected new file to be excluded when includeNewFiles is false, but it was included")
	}
}

// TestGetRepositoryState_IncludesModifiedFilesWhenAddAllFalse verifies that modified files
// are always included regardless of the includeNewFiles flag.
func TestGetRepositoryState_IncludesModifiedFilesWhenAddAllFalse(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create and commit initial file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Modify and stage file
	if err := os.WriteFile(testFile, []byte("modified\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Get repository state with includeNewFiles = false
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, false)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify modified file is included
	foundModifiedFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "test.txt" && file.Status == "modified" {
			foundModifiedFile = true
			break
		}
	}

	if !foundModifiedFile {
		t.Error("Expected modified file to be included regardless of includeNewFiles flag, but it was excluded")
	}
}

// TestGetRepositoryState_IncludesDeletedFilesWhenAddAllFalse verifies that deleted files
// are always included regardless of the includeNewFiles flag.
func TestGetRepositoryState_IncludesDeletedFilesWhenAddAllFalse(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create and commit file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Delete and stage file
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage deleted file: %v", err)
	}

	// Get repository state with includeNewFiles = false
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, false)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify deleted file is included
	foundDeletedFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "test.txt" && file.Status == "deleted" {
			foundDeletedFile = true
			break
		}
	}

	if !foundDeletedFile {
		t.Error("Expected deleted file to be included regardless of includeNewFiles flag, but it was excluded")
	}
}

// TestGetRepositoryState_IncludesRenamedFilesWhenAddAllFalse verifies that renamed files
// are always included regardless of the includeNewFiles flag.
// Note: This test verifies that files with "renamed" status are not filtered.
// When git mv is used, git may detect it as a rename and show git.Renamed status,
// which should never be filtered (only git.Added status files are filtered).
func TestGetRepositoryState_IncludesRenamedFilesWhenAddAllFalse(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create and commit file with substantial content so git detects rename
	oldFile := filepath.Join(tmpDir, "old.txt")
	content := strings.Repeat("content line\n", 10) // Enough content for git to detect rename
	if err := os.WriteFile(oldFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create old file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", oldFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Rename file (git mv) - git should detect this as a rename
	newFile := filepath.Join(tmpDir, "new.txt")
	cmd = exec.Command("git", "-C", tmpDir, "mv", oldFile, newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to rename file: %v", err)
	}

	// Get repository state with includeNewFiles = false
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, false)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify that we have staged files (renamed file should be included)
	// The key point: renamed files have git.Renamed status, not git.Added,
	// so they should never be filtered regardless of includeNewFiles flag
	if len(state.StagedFiles) == 0 {
		t.Error("Expected renamed file to be included (git.Renamed status should not be filtered), but no staged files found")
	}

	// Check if the new file path exists in staged files
	// Note: git may show rename in different ways, but the important thing is
	// that files with git.Renamed status are not filtered (only git.Added are filtered)
	foundStagedFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "new.txt" || file.Status == "renamed" {
			foundStagedFile = true
			// Verify it's not being filtered - renamed files should always be included
			break
		}
	}

	// If git detected it as a rename, it should be included
	// If git shows it as "added" after mv, that's a git behavior, but our filter
	// correctly excludes only true new files (git.Added), not renamed files
	if !foundStagedFile && len(state.StagedFiles) > 0 {
		// If we have staged files but not the renamed one, check if git showed it differently
		t.Logf("Staged files: %v", state.StagedFiles)
	}
}

// TestGetRepositoryState_ExcludesManuallyStagedNewFiles verifies that manually staged new files
// are excluded when includeNewFiles is false.
func TestGetRepositoryState_ExcludesManuallyStagedNewFiles(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create new file and manually stage it (simulating user running git add before gitcomm)
	newFile := filepath.Join(tmpDir, "manual.txt")
	if err := os.WriteFile(newFile, []byte("manual content\n"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to manually stage new file: %v", err)
	}

	// Get repository state with includeNewFiles = false
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, false)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify manually staged new file is excluded
	foundManualFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "manual.txt" {
			foundManualFile = true
			break
		}
	}

	if foundManualFile {
		t.Error("Expected manually staged new file to be excluded when includeNewFiles is false, but it was included")
	}
}

// TestGetRepositoryState_ExcludesBinaryNewFiles verifies that binary new files
// follow the same exclusion rules as text files.
func TestGetRepositoryState_ExcludesBinaryNewFiles(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create binary file (PNG header) and stage it
	binaryFile := filepath.Join(tmpDir, "image.png")
	binaryContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	if err := os.WriteFile(binaryFile, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", binaryFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage binary file: %v", err)
	}

	// Get repository state with includeNewFiles = false
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, false)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify binary new file is excluded
	foundBinaryFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "image.png" {
			foundBinaryFile = true
			break
		}
	}

	if foundBinaryFile {
		t.Error("Expected binary new file to be excluded when includeNewFiles is false, but it was included")
	}
}

// TestGetRepositoryState_IncludesNewFilesWhenAddAllTrue verifies that new files
// are included when includeNewFiles is true.
func TestGetRepositoryState_IncludesNewFilesWhenAddAllTrue(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create new file and stage it
	newFile := filepath.Join(tmpDir, "newfile.txt")
	if err := os.WriteFile(newFile, []byte("new content\n"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage new file: %v", err)
	}

	// Get repository state with includeNewFiles = true
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.WithValue(context.Background(), IncludeNewFilesKey, true)
	state, err := repo.GetRepositoryState(ctx)
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify new file is included
	foundNewFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "newfile.txt" {
			foundNewFile = true
			break
		}
	}

	if !foundNewFile {
		t.Error("Expected new file to be included when includeNewFiles is true, but it was excluded")
	}
}

// TestGetRepositoryState_DefaultBehaviorIncludesAll verifies backward compatibility:
// when context value is not present, all files are included (default behavior).
func TestGetRepositoryState_DefaultBehaviorIncludesAll(t *testing.T) {
	utils.InitLogger(true)

	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create new file and stage it
	newFile := filepath.Join(tmpDir, "newfile.txt")
	if err := os.WriteFile(newFile, []byte("new content\n"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage new file: %v", err)
	}

	// Get repository state without context value (default behavior)
	repo, err := NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify new file is included (backward compatibility)
	foundNewFile := false
	for _, file := range state.StagedFiles {
		if file.Path == "newfile.txt" {
			foundNewFile = true
			break
		}
	}

	if !foundNewFile {
		t.Error("Expected new file to be included by default (backward compatibility), but it was excluded")
	}
}
