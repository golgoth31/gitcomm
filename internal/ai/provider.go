package ai

import (
	"context"

	"github.com/golgoth31/gitcomm/internal/model"
)

// AIProvider defines the interface for AI providers that generate commit messages
type AIProvider interface {
	// GenerateCommitMessage generates a commit message based on repository state
	GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error)
}
