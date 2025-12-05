package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/golgoth31/gitcomm/pkg/ai/prompt"
	"github.com/golgoth31/gitcomm/pkg/conventional"
)

// LocalProvider implements AIProvider for local models
type LocalProvider struct {
	config    *model.AIProviderConfig
	client    *http.Client
	generator prompt.PromptGenerator
	validator conventional.MessageValidator
}

// NewLocalProvider creates a new local model provider
func NewLocalProvider(config *model.AIProviderConfig) AIProvider {
	if config.Endpoint == "" {
		utils.Logger.Debug().Msg("Local provider endpoint not configured")
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &LocalProvider{
		config:    config,
		client:    &http.Client{Timeout: timeout},
		generator: prompt.NewUnifiedPromptGenerator(),
		validator: conventional.NewValidator(),
	}
}

// GenerateCommitMessage generates a commit message using a local model
func (p *LocalProvider) GenerateCommitMessage(ctx context.Context, repoState *model.RepositoryState) (string, error) {
	if p.config.Endpoint == "" {
		return "", fmt.Errorf("%w: local provider endpoint not configured", utils.ErrAIProviderUnavailable)
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

	// Prepare request (OpenAI-compatible format for local models)
	requestBody := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemMsg,
			},
			{
				"role":    "user",
				"content": userMsg,
			},
		},
		"max_tokens": p.config.MaxTokens,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", p.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	}

	// Execute request
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", utils.ErrAIProviderUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%w: API returned status %d: %s", utils.ErrAIProviderUnavailable, resp.StatusCode, string(body))
	}

	// Parse response (OpenAI-compatible format)
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("%w: no response from API", utils.ErrAIProviderUnavailable)
	}

	return response.Choices[0].Message.Content, nil
}
