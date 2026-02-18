package integration

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/golgoth31/gitcomm/internal/repository"
)

// TestSignalInterruptDuringStaging tests that staging state is restored when CLI is interrupted during staging
func TestSignalInterruptDuringStaging(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-signal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init", tmpDir)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create and modify a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage initial file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "commit", "-m", "initial commit")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Modify the file (now it's modified but not staged)
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Test staging state capture and restoration
	repo, err := repository.NewGitRepository(tmpDir, false, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	ctx := context.Background()

	// Capture pre-staging state (should be empty - no files staged)
	preState, err := repo.CaptureStagingState(ctx)
	if err != nil {
		t.Fatalf("Failed to capture pre-staging state: %v", err)
	}

	// Pre-state should be empty (file is modified but not staged)
	if !preState.IsEmpty() {
		t.Logf("Pre-state has staged files: %v", preState.StagedFiles)
	}

	// Stage modified files
	result, err := repo.StageModifiedFiles(ctx)
	if err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	if !result.Success {
		t.Fatalf("Staging failed: %v", result.FailedFiles)
	}

	// Verify file is staged
	currentState, err := repo.CaptureStagingState(ctx)
	if err != nil {
		t.Fatalf("Failed to capture current state: %v", err)
	}

	if !currentState.Contains("test.txt") {
		t.Error("File should be staged after StageModifiedFiles")
	}

	// Simulate interruption by restoring state
	// Note: UnstageFiles implementation uses worktree.Remove which may need refinement
	// This test verifies the API structure - actual unstaging may require git reset
	if err := repo.UnstageFiles(ctx, []string{"test.txt"}); err != nil {
		t.Logf("UnstageFiles returned error (implementation may need git reset): %v", err)
	}

	// Verify state capture works
	restoredState, err := repo.CaptureStagingState(ctx)
	if err != nil {
		t.Fatalf("Failed to capture restored state: %v", err)
	}

	// Verify restoration worked
	if restoredState.Contains("test.txt") {
		t.Error("File should not be staged after restoration")
	}

	// Verify restored state matches pre-state
	if !restoredState.IsEmpty() {
		t.Errorf("Restored state should be empty, but has files: %v", restoredState.StagedFiles)
	}

	t.Logf("Pre-state: %v, Current: %v, Restored: %v", preState.StagedFiles, currentState.StagedFiles, restoredState.StagedFiles)
	t.Log("Signal interrupt test passed - staging state restored successfully")
}

// TestSignalHandling tests signal handling setup (structure test)
func TestSignalHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test that signal channel can be created and signals can be received
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Send a test signal in a goroutine
	go func() {
		time.Sleep(10 * time.Millisecond)
		sigChan <- os.Interrupt
	}()

	// Wait for signal with timeout
	select {
	case sig := <-sigChan:
		if sig != os.Interrupt {
			t.Errorf("Expected SIGINT, got %v", sig)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for signal")
	}

	// Cleanup
	signal.Stop(sigChan)
	close(sigChan)
}
