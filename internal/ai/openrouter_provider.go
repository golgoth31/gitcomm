package ai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/golgoth31/gitcomm/pkg/ai/prompt"
	"github.com/golgoth31/gitcomm/pkg/conventional"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// openRouterDefaultBaseURL is the default OpenRouter API base URL
const openRouterDefaultBaseURL = "https://openrouter.ai/api/v1"

// openRouterDefaultModel is the default OpenRouter model used when none is configured
const openRouterDefaultModel = "openrouter/auto"

// OpenRouterProvider implements AIProvider for OpenRouter
type OpenRouterProvider struct {
	config    *model.AIProviderConfig
	client    openai.Client
	generator prompt.PromptGenerator
	validator conventional.MessageValidator
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(config *model.AIProviderConfig) AIProvider {
	if config.APIKey == "" {
		utils.Logger.Debug().Msg("OpenRouter API key not provided")
	}

	baseURL := config.Endpoint
	if baseURL == "" {
		baseURL = openRouterDefaultBaseURL
	}

	clientOptions := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(baseURL),
	}

	// Honor configured timeout (default: 30s) via a custom HTTP client
	if config.Timeout > 0 {
		clientOptions = append(clientOptions, option.WithHTTPClient(&http.Client{Timeout: config.Timeout}))
	}

	// Initialize OpenAI SDK client pointed at OpenRouter
	client := openai.NewClient(clientOptions...)

	return &OpenRouterProvider{
		config:    config,
		client:    client,
		generator: prompt.NewUnifiedPromptGenerator(),
		validator: conventional.NewValidator(),
	}
}

// GenerateCommitMessage generates a commit message using OpenRouter Chat Completions API
func (p *OpenRouterProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) {
	if p.config.APIKey == "" {
		return "", fmt.Errorf("%w: OpenRouter API key not configured", utils.ErrAIProviderUnavailable)
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
		modelName = openRouterDefaultModel
	}

	// Convert messages to Chat Completions input format
	inputItems := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String(systemMsg),
				},
			},
		},
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(userMsg),
				},
			},
		},
	}

	// Create Chat Completions API request using SDK
	req := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(modelName),
		Messages: inputItems,
	}

	if p.config.MaxTokens > 0 {
		req.MaxTokens = openai.Int(int64(p.config.MaxTokens))
	}

	// Execute Chat Completions API call with context (respects cancellation/timeout)
	resp, err := p.client.Chat.Completions.New(ctx, req)
	if err != nil {
		utils.Logger.Debug().Err(err).Msg("Error generating commit message")
		return "", mapOpenAIError(err)
	}

	utils.Logger.Debug().Msgf("Chat Completions API response: %+v", resp)

	// Extract message content from Chat Completions response
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
	}

	return content, nil
}
