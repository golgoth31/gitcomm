package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/golgoth31/gitcomm/pkg/ai/prompt"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

// AnthropicProvider implements AIProvider for Anthropic
type AnthropicProvider struct {
	config    *model.AIProviderConfig
	client    anthropic.Client
	generator prompt.PromptGenerator
	validator conventional.MessageValidator
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(config *model.AIProviderConfig) AIProvider {
	if config.APIKey == "" {
		utils.Logger.Debug().Msg("Anthropic API key not provided")
	}

	// Initialize Anthropic SDK client
	// NewClient doesn't return an error - it reads from environment or uses provided options
	client := anthropic.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	return &AnthropicProvider{
		config:    config,
		client:    client,
		generator: prompt.NewUnifiedPromptGenerator(),
		validator: conventional.NewValidator(),
	}
}

// GenerateCommitMessage generates a commit message using Anthropic
func (p *AnthropicProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) {
	if p.config.APIKey == "" {
		return "", fmt.Errorf("%w: Anthropic API key not configured", utils.ErrAIProviderUnavailable)
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

	// Anthropic doesn't support system messages, so prepend system to user message
	combinedMsg := systemMsg + "\n\n" + userMsg

	// Prepare model
	modelName := p.config.Model
	if modelName == "" {
		modelName = "claude-3-opus-20240229"
	}

	maxTokens := p.config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 500
	}

	// Create message request using SDK
	req := anthropic.MessageNewParams{
		Model: anthropic.Model(modelName),
		Messages: []anthropic.MessageParam{
			{
				Role: anthropic.MessageParamRoleUser,
				Content: []anthropic.ContentBlockParamUnion{
					{
						OfText: &anthropic.TextBlockParam{
							Text: combinedMsg,
						},
					},
				},
			},
		},
		MaxTokens: int64(maxTokens),
	}

	// Execute SDK API call with context (respects cancellation/timeout)
	resp, err := p.client.Messages.New(ctx, req)
	if err != nil {
		// Map SDK errors to existing error types
		return "", p.mapSDKError(err)
	}

	// Extract message content from SDK response
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("%w: no response from API", utils.ErrAIProviderUnavailable)
	}

	// Extract text content from the first content block
	// ContentBlockUnion is a union type with Text field directly accessible
	contentBlock := resp.Content[0]
	content := contentBlock.Text
	if content == "" {
		return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
	}

	return content, nil
}

// mapSDKError maps SDK-specific errors to existing error types
func (p *AnthropicProvider) mapSDKError(err error) error {
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
