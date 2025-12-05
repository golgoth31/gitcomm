/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golgoth31/gitcomm/internal/config"
	"github.com/golgoth31/gitcomm/internal/model"
	"github.com/golgoth31/gitcomm/internal/repository"
	"github.com/golgoth31/gitcomm/internal/service"
	"github.com/golgoth31/gitcomm/internal/utils"
	"github.com/spf13/cobra"
)

var (
	debug      bool
	addAll     bool
	noSignoff  bool
	noSign     bool
	provider   string
	skipAI     bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "gitcomm",
	Short: "Automate git commit message creation with Conventional Commits",
	Long: `gitcomm is a CLI tool that helps you create properly formatted
commit messages following the Conventional Commits specification.
It supports both manual input and AI-assisted generation.

Examples:
  # Create commit with manual input
  gitcomm

  # Auto-stage files and create commit
  gitcomm -a

  # Create commit without signoff
  gitcomm -s

  # Use AI to generate commit message
  gitcomm --provider openai

  # Skip AI and use manual input
  gitcomm --skip-ai

For more information, visit: https://github.com/golgoth31/gitcomm`,
	Run: runCommand,
}

func runCommand(cmd *cobra.Command, args []string) {
	// Initialize logger
	utils.InitLogger(debug)

	// Create context with cancellation for signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful interruption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		utils.Logger.Debug().Err(err).Msg("Failed to load configuration, continuing with defaults")
		cfg = &config.Config{}
	}

	// Initialize git repository early (needed for restoration)
	gitRepo, err := repository.NewGitRepository("", noSign)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize git repository: %v\n", err)
		os.Exit(1)
	}

	// Create commit options
	options := &model.CommitOptions{
		AutoStage:  addAll,
		NoSignoff:  noSignoff,
		AIProvider: provider,
		SkipAI:     skipAI,
	}

	// Log CLI options
	utils.Logger.Debug().
		Bool("auto_stage", options.AutoStage).
		Bool("no_signoff", options.NoSignoff).
		Bool("no_sign", noSign).
		Str("ai_provider", options.AIProvider).
		Bool("skip_ai", options.SkipAI).
		Msg("CLI options")

	// Channel to signal restoration completion
	restoreDone := make(chan struct{})

	// Create commit service
	commitService := service.NewCommitService(gitRepo, options, cfg)

	// Set restoration completion channel
	commitService.SetRestoreDoneChannel(restoreDone)

	// Handle signals in a goroutine
	signalReceived := false
	go func() {
		sig := <-sigChan
		if signalReceived {
			// Ignore subsequent signals (multiple Ctrl+C handling)
			utils.Logger.Debug().Msg("Ignoring subsequent interrupt signal")
			return
		}
		signalReceived = true
		utils.Logger.Debug().Str("signal", sig.String()).Msg("Received interrupt signal")
		cancel() // Cancel context to stop ongoing operations
		os.Exit(0)
	}()

	// Execute commit workflow
	var commitErr error
	if err := commitService.CreateCommit(ctx); err != nil {
		commitErr = err
	}

	// Check if error is due to context cancellation (Ctrl+C)
	if ctx.Err() == context.Canceled {
		// Signal was received - wait for restoration to complete or timeout
		utils.Logger.Debug().Msg("Workflow cancelled by signal - waiting for restoration")

		// Wait for restoration with overall 5-second timeout
		select {

		case <-restoreDone:
			// Restoration completed
			utils.Logger.Debug().Msg("Restoration completed")
		case <-time.After(5 * time.Second):
			// Overall timeout exceeded
			utils.Logger.Debug().Msg("Overall timeout exceeded - exiting")
			fmt.Printf("Warning: Restoration did not complete in time.\n")
		}

		close(restoreDone)
		os.Exit(130) // Exit code for SIGINT
	}

	if commitErr != nil {
		if commitErr == utils.ErrNoChanges {
			fmt.Println("No changes to commit.")
			return
		}
		fmt.Fprintf(os.Stderr, "Error: commit failed: %v\n", commitErr)
		os.Exit(1)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging (raw text format, no timestamps)")
	rootCmd.Flags().BoolVarP(&addAll, "add-all", "a", false, "Automatically stage all unstaged files")
	rootCmd.Flags().BoolVarP(&noSignoff, "no-signoff", "s", false, "Disable commit signoff")
	rootCmd.Flags().BoolVar(&noSign, "no-sign", false, "Disable commit signing")
	rootCmd.Flags().StringVar(&provider, "provider", "", "Override default AI provider")
	rootCmd.Flags().BoolVar(&skipAI, "skip-ai", false, "Skip AI generation and proceed directly to manual input")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to configuration file (default: ~/.gitcomm/config.yaml)")
}
