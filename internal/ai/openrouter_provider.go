package ai

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

// openRouterDefaultTimeout is the default request timeout used when none is configured
const openRouterDefaultTimeout = 30 * time.Second

// OpenRouterProvider implements AIProvider for OpenRouter
type OpenRouterProvider struct {
	config    *model.AIProviderConfig
	client    openai.Client
	generator prompt.PromptGenerator
	validator conventional.MessageValidator
}

// Compile-time assertion that OpenRouterProvider implements the AIProvider interface
var _ AIProvider = (*OpenRouterProvider)(nil)

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(config *model.AIProviderConfig) AIProvider {
	if config == nil {
		config = &model.AIProviderConfig{}
	}

	if config.APIKey == "" {
		utils.Logger.Debug().Msg("OpenRouter API key not provided")
	}

	baseURL := config.Endpoint
	if baseURL == "" {
		baseURL = openRouterDefaultBaseURL
	}

	// Honor configured timeout via a custom HTTP client; default to 30s when unset
	// so the provider never relies on the SDK's default client (which has no timeout)
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = openRouterDefaultTimeout
	}

	clientOptions := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(baseURL),
		option.WithHTTPClient(&http.Client{Timeout: timeout}),
	}

	// Initialize OpenAI SDK client pointed at OpenRouter
	// NewClient doesn't return an error - it reads from environment or uses provided options
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
		utils.Logger.Warn().Err(err).Msg("Error generating commit message from OpenRouter")
		return "", mapOpenAIError(err)
	}

	// Extract message content from Chat Completions response
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("%w: empty response from API (no choices)", utils.ErrAIProviderUnavailable)
	}

	choice := resp.Choices[0]

	// Surface content-filter refusals instead of a generic empty-response error
	if choice.Message.Refusal != "" {
		return "", fmt.Errorf("%w: model refused to respond: %s", utils.ErrAIProviderUnavailable, choice.Message.Refusal)
	}

	content := choice.Message.Content
	if content == "" {
		// A length stop means MaxTokens cut the response off, not that it was empty
		if choice.FinishReason == "length" {
			return "", fmt.Errorf("%w: response truncated (increase max_tokens)", utils.ErrAIProviderUnavailable)
		}
		return "", fmt.Errorf("%w: empty response from API (finish_reason: %s)", utils.ErrAIProviderUnavailable, choice.FinishReason)
	}

	utils.Logger.Debug().Str("choice_content", content).Msg("Chat Completions response received")
	return content, nil
}
