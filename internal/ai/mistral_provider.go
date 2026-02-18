package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gage-technologies/mistral-go"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/golgoth31/gitcomm/pkg/ai/prompt"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

// MistralProvider implements AIProvider for Mistral AI
type MistralProvider struct {
	config    *model.AIProviderConfig
	client    *mistral.MistralClient
	generator prompt.PromptGenerator
	validator conventional.MessageValidator
}

// NewMistralProvider creates a new Mistral provider
func NewMistralProvider(config *model.AIProviderConfig) AIProvider {
	if config.APIKey == "" {
		utils.Logger.Debug().Msg("Mistral API key not provided")
	}

	// Initialize Mistral SDK client
	// Use custom endpoint constructor when endpoint is configured (e.g., for testing or self-hosted)
	var client *mistral.MistralClient
	if config.Endpoint != "" {
		timeout := config.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		// Use 1 retry for custom endpoints (self-hosted, testing) to avoid
		// excessive retries against non-standard servers
		client = mistral.NewMistralClient(config.APIKey, config.Endpoint, 1, timeout)
	} else {
		client = mistral.NewMistralClientDefault(config.APIKey)
	}

	return &MistralProvider{
		config:    config,
		client:    client,
		generator: prompt.NewUnifiedPromptGenerator(),
		validator: conventional.NewValidator(),
	}
}

// GenerateCommitMessage generates a commit message using Mistral AI
func (p *MistralProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) {
	if p.config.APIKey == "" {
		return "", fmt.Errorf("%w: Mistral API key not configured", utils.ErrAIProviderUnavailable)
	}

	// Generate unified system and user messages
	systemMsg, err := p.generator.GenerateSystemMessage(p.validator)
	if err != nil {
		return "", fmt.Errorf("failed to generate system message: %w", err)
	}

	userMsg, err := p.generator.GenerateUserMessage(repoState)
	if err != nil {
		return "", fmt.Errorf("failed to generate user message: %w", err)
	}

	// Prepare model
	modelName := p.config.Model
	if modelName == "" {
		modelName = "mistral-large-latest"
	}

	maxTokens := p.config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 500
	}

	// Create chat request using SDK
	messages := []mistral.ChatMessage{
		{
			Role:    mistral.RoleSystem,
			Content: systemMsg,
		},
		{
			Role:    mistral.RoleUser,
			Content: userMsg,
		},
	}

	params := mistral.DefaultChatRequestParams
	params.MaxTokens = maxTokens

	// Execute SDK API call with context support
	// The Mistral SDK doesn't accept context.Context, so we wrap the call
	// in a goroutine and use select to respect context cancellation/deadline.
	type chatResult struct {
		resp *mistral.ChatCompletionResponse
		err  error
	}
	resultCh := make(chan chatResult, 1)
	go func() {
		resp, err := p.client.Chat(modelName, messages, &params)
		resultCh <- chatResult{resp: resp, err: err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-resultCh:
		if result.err != nil {
			return "", p.mapSDKError(result.err)
		}

		// Extract message content from SDK response
		if len(result.resp.Choices) == 0 {
			return "", fmt.Errorf("%w: no response from API", utils.ErrAIProviderUnavailable)
		}

		content := result.resp.Choices[0].Message.Content
		if content == "" {
			return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
		}

		return content, nil
	}
}

// mapSDKError maps SDK-specific errors to existing error types
func (p *MistralProvider) mapSDKError(err error) error {
	errStr := err.Error()

	// Check for context cancellation/deadline
	if strings.Contains(strings.ToLower(errStr), "timeout") ||
		strings.Contains(strings.ToLower(errStr), "deadline") ||
		strings.Contains(strings.ToLower(errStr), "context canceled") {
		return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
	}

	// Map HTTP error codes from SDK format "(HTTP Error NNN) ..."
	if strings.Contains(errStr, "HTTP Error") {
		return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
	}

	// Check for authentication errors
	if strings.Contains(strings.ToLower(errStr), "authentication") ||
		strings.Contains(strings.ToLower(errStr), "invalid") ||
		strings.Contains(errStr, "401") {
		return fmt.Errorf("%w: check API key and network connection: API key invalid", utils.ErrAIProviderUnavailable)
	}

	// Generic error mapping - preserve user-facing message, wrap with ErrAIProviderUnavailable
	return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
}
