package repository

import (
	"errors"
	"fmt"
	"strings"
)

const maxStderrDisplayLen = 1500

// FormatErrorForDisplay formats an error for user display, surfacing git stderr
// and applying truncation per FR-001–FR-005. Callers must not pass nil; returns "" for nil.
func FormatErrorForDisplay(err error) string {
	if err == nil {
		return ""
	}
	var gitErr *ErrGitCommandFailed
	if errors.As(err, &gitErr) {
		detail := strings.TrimSpace(gitErr.Stderr)
		if detail == "" {
			detail = "No additional details from git. Check repository state or run the command manually."
		} else if len(detail) > maxStderrDisplayLen {
			detail = detail[:maxStderrDisplayLen] + fmt.Sprintf("… (%d additional characters)", len(gitErr.Stderr)-maxStderrDisplayLen)
		}
		return fmt.Sprintf("git %s failed (exit %d). Details: %s", gitErr.Command, gitErr.ExitCode, detail)
	}
	return err.Error()
}
