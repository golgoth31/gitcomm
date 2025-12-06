package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/utils"
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

	// T013: Validate path is not a directory
	if info, err := os.Stat(configPath); err == nil && info.Mode().IsDir() {
		return nil, fmt.Errorf("config path is a directory, not a file: %s", configPath)
	}

	// T012: Check if config file exists, create if missing
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// T014: Create parent directories
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory %s: %w", configDir, err)
		}

		// T015: Create empty file with O_CREATE|O_WRONLY|O_EXCL flags
		// T016: Set file permissions to 0600
		file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
		if err != nil {
			// T017: Handle race condition (file created by another process)
			if os.IsExist(err) {
				// File was created by another process, treat as success
			} else {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			file.Close()
			// T018: Log file creation
			utils.Logger.Debug().Str("path", configPath).Msg("Created config file")
		}
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
