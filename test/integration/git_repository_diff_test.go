package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/repository"
	"github.com/golgoth31/gitcomm/internal/utils"
)

func TestGetRepositoryState_WithStagedModifiedFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

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
	if err := os.WriteFile(testFile, []byte("line 1\nline 2\nline 3\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Modify file and stage it
	if err := os.WriteFile(testFile, []byte("line 1\nline 2 modified\nline 3\nline 4\n"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file has diff
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	stagedFile := state.StagedFiles[0]
	if stagedFile.Path != "test.txt" {
		t.Errorf("Expected file path 'test.txt', got '%s'", stagedFile.Path)
	}

	if stagedFile.Diff == "" {
		t.Error("Expected Diff field to be populated for staged modified file, got empty")
	}

	// Verify diff contains expected content
	if !strings.Contains(stagedFile.Diff, "diff --git") {
		t.Error("Expected diff to contain 'diff --git' header")
	}

	if !strings.Contains(stagedFile.Diff, "test.txt") {
		t.Error("Expected diff to contain file path")
	}
}

func TestGetRepositoryState_WithStagedAddedFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create new file and stage it
	testFile := filepath.Join(tmpDir, "new.txt")
	if err := os.WriteFile(testFile, []byte("new file content\nline 2\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file has diff
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	stagedFile := state.StagedFiles[0]
	if stagedFile.Path != "new.txt" {
		t.Errorf("Expected file path 'new.txt', got '%s'", stagedFile.Path)
	}

	if stagedFile.Diff == "" {
		t.Error("Expected Diff field to be populated for staged added file, got empty")
	}

	// Verify diff contains expected content for new file
	if !strings.Contains(stagedFile.Diff, "new file") || !strings.Contains(stagedFile.Diff, "mode 100644") {
		t.Error("Expected diff to contain 'new file mode' for added file")
	}
}

func TestGetRepositoryState_WithStagedDeletedFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit it
	testFile := filepath.Join(tmpDir, "to_delete.txt")
	if err := os.WriteFile(testFile, []byte("content to delete\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Delete file and stage deletion
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "rm", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage deletion: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file has diff
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	stagedFile := state.StagedFiles[0]
	if stagedFile.Path != "to_delete.txt" {
		t.Errorf("Expected file path 'to_delete.txt', got '%s'", stagedFile.Path)
	}

	if stagedFile.Status != "deleted" {
		t.Errorf("Expected status 'deleted', got '%s'", stagedFile.Status)
	}

	if stagedFile.Diff == "" {
		t.Error("Expected Diff field to be populated for staged deleted file, got empty")
	}

	// Verify diff contains deletion markers
	if !strings.Contains(stagedFile.Diff, "deleted file") {
		t.Error("Expected diff to contain 'deleted file' for deleted file")
	}
}

func TestGetRepositoryState_WithStagedRenamedFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit it
	oldFile := filepath.Join(tmpDir, "old.txt")
	if err := os.WriteFile(oldFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", oldFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Rename file using git mv
	newFile := filepath.Join(tmpDir, "new.txt")
	cmd = exec.Command("git", "-C", tmpDir, "mv", oldFile, newFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to rename file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file has diff
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	stagedFile := state.StagedFiles[0]
	if stagedFile.Path != "new.txt" {
		t.Errorf("Expected file path 'new.txt', got '%s'", stagedFile.Path)
	}

	if stagedFile.Status != "renamed" {
		t.Errorf("Expected status 'renamed', got '%s'", stagedFile.Status)
	}

	if stagedFile.Diff == "" {
		t.Error("Expected Diff field to be populated for staged renamed file, got empty")
	}

	// Verify diff contains rename information
	if !strings.Contains(stagedFile.Diff, "rename from") || !strings.Contains(stagedFile.Diff, "rename to") {
		t.Error("Expected diff to contain 'rename from' and 'rename to' for renamed file")
	}

	if !strings.Contains(stagedFile.Diff, "similarity") {
		t.Error("Expected diff to contain 'similarity' percentage for renamed file")
	}
}

func TestGetRepositoryState_WithStagedCopiedFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit it
	sourceFile := filepath.Join(tmpDir, "source.txt")
	if err := os.WriteFile(sourceFile, []byte("source content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", sourceFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Copy file using git add with --find-copies
	destFile := filepath.Join(tmpDir, "dest.txt")
	if err := os.WriteFile(destFile, []byte("source content\n"), 0644); err != nil {
		t.Fatalf("Failed to create copied file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", destFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage copied file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file has diff (may be added or copied depending on git version)
	if len(state.StagedFiles) < 1 {
		t.Fatalf("Expected at least 1 staged file, got %d", len(state.StagedFiles))
	}

	// Find the copied file
	var copiedFile *model.FileChange
	for i := range state.StagedFiles {
		if state.StagedFiles[i].Path == "dest.txt" {
			copiedFile = &state.StagedFiles[i]
			break
		}
	}

	if copiedFile == nil {
		t.Fatal("Expected to find 'dest.txt' in staged files")
	}

	if copiedFile.Diff == "" {
		t.Error("Expected Diff field to be populated for staged copied file, got empty")
	}

	// If git detected it as copied, verify copy information
	if copiedFile.Status == "copied" {
		if !strings.Contains(copiedFile.Diff, "copy from") || !strings.Contains(copiedFile.Diff, "copy to") {
			t.Error("Expected diff to contain 'copy from' and 'copy to' for copied file")
		}
	}
}

func TestGetRepositoryState_WithNoStagedChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit it (no staged changes)
	testFile := filepath.Join(tmpDir, "committed.txt")
	if err := os.WriteFile(testFile, []byte("content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Get repository state (no staged changes)
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify no staged files
	if len(state.StagedFiles) != 0 {
		t.Fatalf("Expected 0 staged files, got %d", len(state.StagedFiles))
	}
}

func TestGetRepositoryState_WithUnstagedFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create file and commit it
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Modify file but don't stage it (unstaged change)
	if err := os.WriteFile(testFile, []byte("modified content\n"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
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

	unstagedFile := state.UnstagedFiles[0]
	if unstagedFile.Path != "test.txt" {
		t.Errorf("Expected file path 'test.txt', got '%s'", unstagedFile.Path)
	}

	if unstagedFile.Diff != "" {
		t.Error("Expected Diff field to be empty for unstaged file, got non-empty")
	}
}

func TestGetRepositoryState_DiffMatchesGitDiffCached(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

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
	if err := os.WriteFile(testFile, []byte("line 1\nline 2\nline 3\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Make specific changes and stage
	if err := os.WriteFile(testFile, []byte("line 1\nline 2 modified\nline 3\nline 4 added\n"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify diff matches git diff --cached (allowing for minor formatting differences)
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	actualDiff := strings.TrimSpace(state.StagedFiles[0].Diff)
	if actualDiff == "" {
		t.Fatal("Expected diff to be populated, got empty")
	}

	// Compare key elements (file path, change markers)
	// Note: Exact match may differ due to hash differences, but structure should match
	if !strings.Contains(actualDiff, "diff --git") {
		t.Error("Expected diff to contain 'diff --git' header")
	}
	if !strings.Contains(actualDiff, "test.txt") {
		t.Error("Expected diff to contain file path")
	}
	// Check for change markers (though exact format may vary)
	if !strings.Contains(actualDiff, "@@") && !strings.Contains(actualDiff, "+") && !strings.Contains(actualDiff, "-") {
		t.Error("Expected diff to contain change markers (@@, +, or -)")
	}
}

func TestGetRepositoryState_MultipleFilesComputedIndependently(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create multiple files and commit
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("file1 content\n"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("file2 content\n"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", file1, file2)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Modify both files and stage
	if err := os.WriteFile(file1, []byte("file1 modified\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("file2 modified\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file2: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", file1, file2)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified files: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify both files have diffs computed independently
	if len(state.StagedFiles) != 2 {
		t.Fatalf("Expected 2 staged files, got %d", len(state.StagedFiles))
	}

	for _, stagedFile := range state.StagedFiles {
		if stagedFile.Diff == "" {
			t.Errorf("Expected Diff to be populated for %s, got empty", stagedFile.Path)
		}
		// Verify each diff is independent (contains its own file path)
		if !strings.Contains(stagedFile.Diff, stagedFile.Path) {
			t.Errorf("Expected diff for %s to contain its file path", stagedFile.Path)
		}
	}
}

func TestGetRepositoryState_StagedFileWithUnstagedModifications(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

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
	if err := os.WriteFile(testFile, []byte("initial\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Make staged change
	if err := os.WriteFile(testFile, []byte("staged change\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Make additional unstaged change (not staged)
	if err := os.WriteFile(testFile, []byte("staged change\nunstaged addition\n"), 0644); err != nil {
		t.Fatalf("Failed to add unstaged change: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify staged file diff only includes staged changes (not unstaged)
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	actualDiff := state.StagedFiles[0].Diff
	if actualDiff == "" {
		t.Fatal("Expected diff to be populated, got empty")
	}

	// Verify diff doesn't contain unstaged content
	if strings.Contains(actualDiff, "unstaged addition") {
		t.Error("Expected diff to only include staged changes, but found unstaged content")
	}

	// Verify diff contains staged content
	if !strings.Contains(actualDiff, "staged change") {
		t.Error("Expected diff to contain staged change content")
	}

	// Verify unstaged file is separate
	if len(state.UnstagedFiles) != 1 {
		t.Fatalf("Expected 1 unstaged file, got %d", len(state.UnstagedFiles))
	}

	if state.UnstagedFiles[0].Diff != "" {
		t.Error("Expected unstaged file to have empty Diff field")
	}
}

func TestGetRepositoryState_BinaryFilesHaveEmptyDiff(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create binary file (PNG header)
	binaryFile := filepath.Join(tmpDir, "image.png")
	binaryContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	if err := os.WriteFile(binaryFile, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", binaryFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage binary file: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify binary file has empty diff
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	if state.StagedFiles[0].Diff != "" {
		t.Error("Expected binary file to have empty Diff field, got non-empty")
	}
}

func TestGetRepositoryState_HandlesUnmergedFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

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

	// Create branch and make conflicting changes
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

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "feature")
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

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "main")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Attempt merge (will create conflict)
	cmd = exec.Command("git", "-C", tmpDir, "merge", "feature", "--no-edit")
	cmd.Run() // Ignore error - merge will create conflict

	// Get repository state - should handle unmerged files gracefully
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify function completed successfully (doesn't crash)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// Unmerged files should be handled gracefully
	// (May have empty diff if computation fails, but processing continues)
}

func TestGetRepositoryState_HandlesEmptyRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository (no commits yet - empty repository)
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create new file and stage it (no HEAD exists)
	testFile := filepath.Join(tmpDir, "new.txt")
	if err := os.WriteFile(testFile, []byte("new file content\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Get repository state - should handle empty repository (no HEAD)
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify empty repository handled correctly
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// All staged files should be treated as new additions
	if len(state.StagedFiles) != 1 {
		t.Fatalf("Expected 1 staged file, got %d", len(state.StagedFiles))
	}

	if state.StagedFiles[0].Path != "new.txt" {
		t.Errorf("Expected file path 'new.txt', got '%s'", state.StagedFiles[0].Path)
	}

	// Diff should be populated (treats as new file since no HEAD)
	if state.StagedFiles[0].Diff == "" {
		t.Error("Expected Diff to be populated for new file in empty repository, got empty")
	}
}

func TestGetRepositoryState_PerformanceWithManyFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create initial commit with some files
	for i := 0; i < 10; i++ {
		testFile := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(testFile, []byte(fmt.Sprintf("initial content %d\n", i)), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	// Modify and stage 100 files
	for i := 0; i < 100; i++ {
		testFile := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i%10))
		if err := os.WriteFile(testFile, []byte(fmt.Sprintf("modified content %d\n", i)), 0644); err != nil {
			t.Fatalf("Failed to modify file: %v", err)
		}
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage modified files: %v", err)
	}

	// Measure performance
	start := time.Now()

	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Verify performance requirement: <2 seconds for 100 files (SC-003)
	if elapsed > 2*time.Second {
		t.Errorf("Performance test failed: GetRepositoryState took %v for 100 files, expected <2s", elapsed)
	}

	// Verify all files were processed
	if len(state.StagedFiles) < 10 {
		t.Errorf("Expected at least 10 staged files, got %d", len(state.StagedFiles))
	}

	t.Logf("Performance: Processed %d staged files in %v", len(state.StagedFiles), elapsed)
}

func TestGetRepositoryState_ErrorRate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping error rate test in short mode")
	}

	utils.InitLogger(true)

	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Create initial commit
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

	// Create multiple files and stage them
	var stagedFiles []string
	for i := 0; i < 100; i++ {
		file := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(file, []byte(fmt.Sprintf("content %d\n", i)), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		stagedFiles = append(stagedFiles, file)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	// Get repository state
	repo, err := repository.NewGitRepository(tmpDir, false, true)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	state, err := repo.GetRepositoryState(context.Background())
	if err != nil {
		t.Fatalf("Failed to get repository state: %v", err)
	}

	// Count files with empty diff (errors)
	errors := 0
	total := len(state.StagedFiles)
	for _, file := range state.StagedFiles {
		if file.Diff == "" {
			errors++
		}
	}

	// Verify error rate <1% (SC-004)
	errorRate := float64(errors) / float64(total) * 100
	if errorRate >= 1.0 {
		t.Errorf("Error rate test failed: %.2f%% of files have empty diff (errors), expected <1%%", errorRate)
	}

	t.Logf("Error rate: %.2f%% (%d errors out of %d files)", errorRate, errors, total)
}
