package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/gcfg/v2"
	"github.com/golgoth31/gitcomm/internal/utils"
)

// GitConfig represents extracted git configuration values from .git/config and ~/.gitconfig files
type GitConfig struct {
	UserName      string
	UserEmail     string
	SigningKey    string
	GPGFormat     string
	CommitGPGSign bool
}

// CommitSigner represents the configured commit signer extracted from git config and prepared for use with go-git
type CommitSigner struct {
	PrivateKeyPath string
	PublicKeyPath  string
	Format         string
	Signer         interface{} // Will be git.Signer (e.g., *ssh.PublicKeys from go-git)
}

// ConfigExtractor defines the interface for extracting git configuration
type ConfigExtractor interface {
	// Extract reads git configuration from .git/config and ~/.gitconfig
	// Returns extracted config values, with local config taking precedence
	// Errors are logged but not returned (silent ignore per FR-009)
	Extract(repoPath string) *GitConfig
}

// gitConfigFile represents the structure of a git config file for gcfg parsing
type gitConfigFile struct {
	User struct {
		Name       string
		Email      string
		SigningKey string
	}
	GPG struct {
		Format string
	}
	Commit struct {
		GPGSign string
	}
}

// FileConfigExtractor implements ConfigExtractor by reading git config files directly
type FileConfigExtractor struct{}

// NewFileConfigExtractor creates a new FileConfigExtractor instance
func NewFileConfigExtractor() ConfigExtractor {
	return &FileConfigExtractor{}
}

// Extract reads git configuration from .git/config and ~/.gitconfig
// Returns extracted config values, with local config taking precedence
func (e *FileConfigExtractor) Extract(repoPath string) *GitConfig {
	config := &GitConfig{
		UserName:      "gitcomm",
		UserEmail:     "gitcomm@local",
		SigningKey:    "",
		GPGFormat:     "",
		CommitGPGSign: false,
	}

	// Try to read local config first
	localConfigPath := filepath.Join(repoPath, ".git", "config")
	if err := e.readConfigFile(localConfigPath, config, true); err != nil {
		utils.Logger.Debug().Err(err).Str("path", localConfigPath).Msg("Failed to read local git config, will try global config")
	}

	// Try to read global config (fallback or for missing values)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		utils.Logger.Debug().Err(err).Msg("Failed to get home directory for global git config")
		return config
	}

	globalConfigPath := filepath.Join(homeDir, ".gitconfig")
	if err := e.readConfigFile(globalConfigPath, config, false); err != nil {
		utils.Logger.Debug().Err(err).Str("path", globalConfigPath).Msg("Failed to read global git config")
	}

	return config
}

// readConfigFile reads a git config file and merges values into config
// If isLocal is true, values override existing config (precedence)
// If isLocal is false, values only fill in missing fields
func (e *FileConfigExtractor) readConfigFile(path string, config *GitConfig, isLocal bool) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ConfigError{Message: "config file not found", Err: err}
	}

	// Read and parse config file
	// gcfg may return warnings for unknown sections, but we only care about [user], [gpg], [commit]
	// So we'll try to parse and ignore warnings about other sections
	var cfg gitConfigFile
	err := gcfg.ReadFileInto(&cfg, path)
	// gcfg returns errors for warnings, but we can still extract values we care about
	// Check if the error is just about unknown sections (warnings)
	if err != nil {
		// If we got user config values, the parse was partially successful
		// Try to read the file manually to extract just the sections we need
		// For now, log the error but continue - we'll fall back to defaults if needed
		utils.Logger.Debug().Err(err).Str("path", path).Msg("gcfg parsing returned warnings/errors, attempting manual extraction")

		// Try manual extraction as fallback
		return e.readConfigFileManual(path, config, isLocal)
	}

	// Merge values into config
	if isLocal || config.UserName == "gitcomm" {
		if cfg.User.Name != "" {
			config.UserName = cfg.User.Name
		}
	}
	if isLocal || config.UserEmail == "gitcomm@local" {
		if cfg.User.Email != "" {
			config.UserEmail = cfg.User.Email
		}
	}
	if isLocal || config.SigningKey == "" {
		if cfg.User.SigningKey != "" {
			config.SigningKey = cfg.User.SigningKey
		}
	}
	if isLocal || config.GPGFormat == "" {
		if cfg.GPG.Format != "" {
			config.GPGFormat = cfg.GPG.Format
		}
	}
	// For commit.gpgsign: local config takes precedence
	// Only update if isLocal (local config) or if global and value hasn't been set by local
	// Since we can't track if local set it, we use a simpler approach:
	// - Local: always update
	// - Global: only update if we're reading global AND local config file didn't exist
	//   But we don't have that info here. For now, we'll only update from global in the gcfg path
	//   if local didn't set it. In manual path, we'll track it differently.
	if isLocal {
		// Local config: always update (takes precedence)
		if cfg.Commit.GPGSign != "" {
			config.CommitGPGSign = strings.ToLower(cfg.Commit.GPGSign) == "true"
		}
	}
	// Note: We don't update commit.gpgsign from global config here to avoid overwriting local values
	// Global commit.gpgsign will be handled in the manual parser with proper precedence tracking

	return nil
}

// readConfigFileManual reads config file manually to extract only the sections we need
// This is a fallback when gcfg fails due to unknown sections
func (e *FileConfigExtractor) readConfigFileManual(path string, config *GitConfig, isLocal bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return &ConfigError{Message: "failed to read config file", Err: err}
	}

	lines := strings.Split(string(data), "\n")
	var currentSection string
	var inUserSection, inGPGSection, inCommitSection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(strings.Trim(line, "[]"))
			inUserSection = currentSection == "user"
			inGPGSection = currentSection == "gpg"
			inCommitSection = currentSection == "commit"
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(strings.ToLower(parts[0]))
			value := strings.TrimSpace(parts[1])

			if inUserSection {
				if key == "name" && (isLocal || config.UserName == "gitcomm") {
					config.UserName = value
				} else if key == "email" && (isLocal || config.UserEmail == "gitcomm@local") {
					config.UserEmail = value
				} else if key == "signingkey" && (isLocal || config.SigningKey == "") {
					config.SigningKey = value
				}
			} else if inGPGSection {
				if key == "format" && (isLocal || config.GPGFormat == "") {
					config.GPGFormat = value
				}
			} else if inCommitSection {
				if key == "gpgsign" {
					// Parse commit.gpgsign (can be "true", "false", or empty)
					// Local config always takes precedence
					if isLocal {
						lowerValue := strings.ToLower(value)
						if lowerValue == "true" {
							config.CommitGPGSign = true
						}
					}
					// Don't update from global config here - it will be handled by
					// readCommitGPGSignFromLocal() to ensure local precedence
				}
			}
		}
	}

	return nil
}
