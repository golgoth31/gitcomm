package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/golgoth31/gitcomm/pkg/ai/prompt"
	"github.com/golgoth31/gitcomm/pkg/conventional"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

// OpenAIProvider implements AIProvider for OpenAI
type OpenAIProvider struct {
	config    *model.AIProviderConfig
	client    openai.Client
	generator prompt.PromptGenerator
	validator conventional.MessageValidator
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config *model.AIProviderConfig) AIProvider {
	if config.APIKey == "" {
		utils.Logger.Debug().Msg("OpenAI API key not provided")
	}

	// Initialize OpenAI SDK v3 client
	// NewClient doesn't return an error - it reads from environment or uses provided options
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	return &OpenAIProvider{
		config:    config,
		client:    client,
		generator: prompt.NewUnifiedPromptGenerator(),
		validator: conventional.NewValidator(),
	}
}

// GenerateCommitMessage generates a commit message using OpenAI Responses API
func (p *OpenAIProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) {
	if p.config.APIKey == "" {
		return "", fmt.Errorf("%w: OpenAI API key not configured", utils.ErrAIProviderUnavailable)
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
		modelName = shared.ChatModelGPT4_1
	}

	maxTokens := p.config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 500
	}

	// Convert messages to Responses API input format
	// Use EasyInputMessage for system and user messages
	inputItems := []responses.ResponseInputItemUnionParam{
		{
			OfMessage: &responses.EasyInputMessageParam{
				Role: responses.EasyInputMessageRoleSystem,
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(systemMsg),
				},
			},
		},
		{
			OfMessage: &responses.EasyInputMessageParam{
				Role: responses.EasyInputMessageRoleUser,
				Content: responses.EasyInputMessageContentUnionParam{
					OfString: openai.String(userMsg),
				},
			},
		},
	}

	// Create Responses API request using SDK v3
	req := responses.ResponseNewParams{
		Model: shared.ResponsesModel(modelName),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: inputItems,
		},
		MaxOutputTokens: openai.Int(int64(maxTokens)),
		Store:           openai.Bool(false), // Stateless mode
	}

	// Execute Responses API call with context (respects cancellation/timeout)
	resp, err := p.client.Responses.New(ctx, req)
	if err != nil {
		// Map Responses API errors to existing error types
		utils.Logger.Debug().Err(err).Msg("Error generating commit message")
		return "", p.mapSDKError(err)
	}

	// Extract message content from Responses API response
	// Use OutputText() method to extract text from Output array
	content := resp.OutputText()
	if content == "" {
		return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
	}

	return content, nil
}

// mapSDKError maps Responses API-specific errors to existing error types
func (p *OpenAIProvider) mapSDKError(err error) error {
	// Check for authentication errors
	errStr := err.Error()
	// Map common Responses API error patterns to existing error types
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

	// Generic error mapping for unmappable errors - preserve user-facing message, wrap with ErrAIProviderUnavailable
	// Original SDK error message preserved in wrapped error for debugging
	return fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
}
