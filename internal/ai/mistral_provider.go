package ai

import (
	"context"
	"fmt"
	"strings"

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
	client := mistral.NewMistralClientDefault(config.APIKey)

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

	// Execute SDK API call
	// Note: Mistral SDK doesn't support context directly, but respects timeout from HTTP client
	resp, err := p.client.Chat(modelName, messages, &params)
	if err != nil {
		// Map SDK errors to existing error types
		return "", p.mapSDKError(err)
	}

	// Extract message content from SDK response
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%w: no response from API", utils.ErrAIProviderUnavailable)
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
	}

	return content, nil
}

// mapSDKError maps SDK-specific errors to existing error types
func (p *MistralProvider) mapSDKError(err error) error {
	// Check for authentication errors
	errStr := err.Error()
	// Map common SDK error patterns to existing error types
	if strings.Contains(strings.ToLower(errStr), "authentication") ||
		strings.Contains(strings.ToLower(errStr), "invalid") ||
		strings.Contains(errStr, "401") {
		return fmt.Errorf("%w: API key invalid", utils.ErrAIProviderUnavailable)
	}
	if strings.Contains(strings.ToLower(errStr), "rate limit") ||
		strings.Contains(errStr, "429") {
		return fmt.Errorf("%w: rate limit exceeded", utils.ErrAIProviderUnavailable)
	}
	if strings.Contains(strings.ToLower(errStr), "timeout") ||
		strings.Contains(strings.ToLower(errStr), "deadline") {
		return fmt.Errorf("%w: timeout", utils.ErrAIProviderUnavailable)
	}

	// Generic error mapping - preserve user-facing message, wrap with ErrAIProviderUnavailable
	return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
}
