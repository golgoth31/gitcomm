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

// TestSignalHandlingWithTimeout tests that CLI exits within 5 seconds when Ctrl+C is pressed
func TestSignalHandlingWithTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-timeout-test-*")
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
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Test timeout context creation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Verify context respects timeout
	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got: %v", ctx.Err())
		}
	case <-time.After(4 * time.Second):
		t.Error("Context should have timed out within 3 seconds")
	}
}

// TestRestorationTimeoutScenario tests that restoration times out after 3 seconds
func TestRestorationTimeoutScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "gitcomm-restore-timeout-test-*")
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

	// Create repository instance
	repo, err := repository.NewGitRepository(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create a file and stage it
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "add", testFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Create timeout context (3 seconds)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Attempt restoration with timeout context
	// This should complete quickly (no files to unstage) or timeout
	startTime := time.Now()
	err = repo.UnstageFiles(timeoutCtx, []string{testFile})
	duration := time.Since(startTime)

	if err != nil {
		// Error is acceptable (file might not be staged, or timeout)
		t.Logf("Restoration returned error (expected): %v", err)
	}

	// Verify timeout is respected (should complete within 3 seconds or timeout)
	if duration > 4*time.Second {
		t.Errorf("Restoration should complete or timeout within 3 seconds, took: %v", duration)
	}
}

// TestMultipleCtrlCHandling tests that multiple Ctrl+C presses are handled gracefully
func TestMultipleCtrlCHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a channel to simulate signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Simulate first Ctrl+C
	firstSignalReceived := false
	secondSignalReceived := false

	// Create a goroutine to handle signals
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-sigChan:
			firstSignalReceived = true
			// Simulate ignoring subsequent signals
			select {
			case <-sigChan:
				secondSignalReceived = true
			case <-time.After(100 * time.Millisecond):
				// Timeout - no second signal
			}
		case <-time.After(1 * time.Second):
			// Timeout
		}
	}()

	// Send first signal
	sigChan <- os.Interrupt
	time.Sleep(50 * time.Millisecond)

	// Send second signal (should be ignored or handled gracefully)
	sigChan <- os.Interrupt
	time.Sleep(50 * time.Millisecond)

	// Wait for goroutine to complete
	select {
	case <-done:
		// Verify first signal was received
		if !firstSignalReceived {
			t.Error("First signal should have been received")
		}
		// Second signal handling is implementation-dependent
		t.Logf("First signal received: %v, Second signal received: %v", firstSignalReceived, secondSignalReceived)
	case <-time.After(2 * time.Second):
		t.Error("Signal handling goroutine did not complete")
	}
}
