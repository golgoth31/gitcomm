package gitcomm

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/golgoth31/gitcomm/pkg/gptscript"
)

// GitExec function performs git workflow
func GitExec(addAll, show bool, noSignoff bool, msg string) {
	ctx := context.Background()
	if addAll {
		gitColorCmd("add", "-A")
	}
	gitColorCmd("add", "-u")

	gpt := gptscript.GPTScript{}
	if err := gpt.Run(ctx, nil); err != nil {
		log.Fatalf("error: %v\n", err)
	}

	os.Exit(0)

	gitColorCmd("commit", "-s", "-m", msg)
	if noSignoff {
		gitColorCmd("commit", "-m", msg)
	}
	if show {
		gitColorCmd("show", "-s")
	}
}

// CheckForUncommited function checks if there are changes that need commit
func CheckForUncommited() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	CheckIfError(err)
	return len(out) != 0
}

// CheckIsGitDir function checks is dir inside git worktree
func CheckIsGitDir() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	isGitDir := strings.TrimSpace(string(out))
	if err == nil && isGitDir == "true" {
		return true
	}
	return false
}
