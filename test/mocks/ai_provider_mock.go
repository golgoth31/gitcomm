package mocks

import (
	"context"
	"errors"

	"github.com/golgoth31/gitcomm/internal/model"
)

// MockAIProvider is a mock implementation of AIProvider for testing
type MockAIProvider struct {
	GenerateFunc func(ctx context.Context, repoState *model.RepositoryState) (string, error)
	ShouldFail   bool
	FailError    error
	Response     string
}

// GenerateCommitMessage implements the AIProvider interface
func (m *MockAIProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) {
	if m.ShouldFail {
		if m.FailError != nil {
			return "", m.FailError
		}
		return "", errors.New("mock AI provider error")
	}

	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, repoState)
	}

	if m.Response != "" {
		return m.Response, nil
	}

	// Default response
	return "feat: default commit message from mock", nil
}
