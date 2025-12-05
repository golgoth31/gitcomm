package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	AI AIConfig
}

// AIConfig represents AI provider configuration
type AIConfig struct {
	DefaultProvider string
	Providers       map[string]model.AIProviderConfig
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// Set default config path
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".gitcomm", "config.yaml")
	}

	// Configure viper
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("GITCOMM")
	v.AutomaticEnv()

	// Read config file (optional - may not exist)
	if err := v.ReadInConfig(); err != nil {
		// Config file is optional, continue with defaults
	}

	config := &Config{
		AI: AIConfig{
			DefaultProvider: v.GetString("ai.default_provider"),
			Providers:       make(map[string]model.AIProviderConfig),
		},
	}

	// Load provider configurations
	providers := v.GetStringMap("ai.providers")
	for name := range providers {
		providerConfig := model.AIProviderConfig{
			Name:      name,
			APIKey:    v.GetString(fmt.Sprintf("ai.providers.%s.api_key", name)),
			Model:     v.GetString(fmt.Sprintf("ai.providers.%s.model", name)),
			Endpoint:  v.GetString(fmt.Sprintf("ai.providers.%s.endpoint", name)),
			Timeout:   30 * time.Second,
			MaxTokens: 500,
		}

		// Override timeout if specified
		if timeoutStr := v.GetString(fmt.Sprintf("ai.providers.%s.timeout", name)); timeoutStr != "" {
			if timeout, err := time.ParseDuration(timeoutStr); err == nil {
				providerConfig.Timeout = timeout
			}
		}

		// Override max tokens if specified
		if maxTokens := v.GetInt(fmt.Sprintf("ai.providers.%s.max_tokens", name)); maxTokens > 0 {
			providerConfig.MaxTokens = maxTokens
		}

		config.AI.Providers[name] = providerConfig
	}

	return config, nil
}

// GetProviderConfig returns the configuration for a specific provider
func (c *Config) GetProviderConfig(name string) (*model.AIProviderConfig, error) {
	if name == "" {
		name = c.AI.DefaultProvider
	}

	provider, ok := c.AI.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not configured", name)
	}

	return &provider, nil
}
