package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
		modelName = shared.ChatModelGPT4_1Nano
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
		Store: openai.Bool(false), // Stateless mode
	}

	// Execute Responses API call with context (respects cancellation/timeout)
	resp, err := p.client.Responses.New(ctx, req)
	if err != nil {
		// Map Responses API errors to existing error types
		utils.Logger.Debug().Err(err).Msg("Error generating commit message")
		return "", mapOpenAIError(err)
	}

	// Extract message content from Responses API response
	// Use OutputText() method to extract text from Output array
	content := resp.OutputText()
	if content == "" {
		return "", fmt.Errorf("%w: empty response from API", utils.ErrAIProviderUnavailable)
	}

	utils.Logger.Debug().Str("output", content).Msg("Responses API response received")
	return content, nil
}

// mapOpenAIError maps OpenAI SDK API errors to existing error types
func mapOpenAIError(err error) error {
	// Preserve context errors so callers can detect cancellation/deadline via errors.Is
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// Classify typed API errors by HTTP status code first (reliable and
	// format-independent compared to string matching on the SDK message)
	var apiErr *openai.Error
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return fmt.Errorf("%w: API key invalid: %w", utils.ErrAIProviderUnavailable, err)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%w: rate limit exceeded: %w", utils.ErrAIProviderUnavailable, err)
		case http.StatusNotFound:
			return fmt.Errorf("%w: model not found or no access: %w", utils.ErrAIProviderUnavailable, err)
		}
	}

	// String matching as a fallback for non-typed transport errors (network, DNS, etc.)
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "authentication") ||
		strings.Contains(errStr, "401") ||
		strings.Contains(errStr, "invalid api key") {
		return fmt.Errorf("%w: API key invalid: %w", utils.ErrAIProviderUnavailable, err)
	}
	if strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "429") {
		return fmt.Errorf("%w: rate limit exceeded: %w", utils.ErrAIProviderUnavailable, err)
	}
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline") {
		return fmt.Errorf("%w: timeout: %w", utils.ErrAIProviderUnavailable, err)
	}

	// Generic error mapping - preserve original error chain for debugging
	return fmt.Errorf("%w: %w", utils.ErrAIProviderUnavailable, err)
}
