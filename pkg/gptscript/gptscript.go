package gptscript

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gptscript-ai/gptscript/pkg/cache"
	"github.com/gptscript-ai/gptscript/pkg/engine"
	"github.com/gptscript-ai/gptscript/pkg/llm"
	"github.com/gptscript-ai/gptscript/pkg/loader"
	"github.com/gptscript-ai/gptscript/pkg/monitor"
	"github.com/gptscript-ai/gptscript/pkg/mvl"
	"github.com/gptscript-ai/gptscript/pkg/openai"
	"github.com/gptscript-ai/gptscript/pkg/repos/runtimes"
	"github.com/gptscript-ai/gptscript/pkg/runner"
	"github.com/gptscript-ai/gptscript/pkg/types"

	"github.com/rs/zerolog/log"
)

type (
	DisplayOptions monitor.Options
	CacheOptions   cache.Options
	OpenAIOptions  openai.Options
)

type GPTScript struct {
	CacheOptions
	OpenAIOptions
	DisplayOptions

	_client llm.Client `usage:"-"`
}

const (
	gptProgram = `tools: gitstatus, sys.abort

Create well formed git commit message based of off the currently staged file
contents. The message should convey why something was changed and not what
changed. Use the well known format that has the prefix fix, etc but not chore.

Do not use markdown format for the output.

If there are no changes abort.

---
name: gitstatus

#!/bin/sh

git diff --staged`
)

func (r *GPTScript) Run(ctx context.Context, args []string) error {
	log.Debug().Msg("Starting GPTScript")
	defer engine.CloseDaemons()

	var (
		prg types.Program
		err error
	)

	mvl.SetSimpleFormat()
	mvl.SetError()

	log.Debug().Msg("Loading program")
	prg, err = loader.ProgramFromSource(ctx, gptProgram, "")
	if err != nil {
		return err
	}

	log.Debug().Msg("Getting client")
	client, err := r.getClient(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Creating runner")
	r.DisplayOptions.DisplayProgress = false
	runner, err := runner.New(client, runner.Options{
		MonitorFactory: monitor.NewConsole(monitor.Options(r.DisplayOptions), monitor.Options{
			DisplayProgress: false,
		}),
		RuntimeManager: runtimes.Default(cache.Complete(cache.Options(r.CacheOptions)).CacheDir),
	})
	if err != nil {
		return err
	}

	fmt.Println("Running program")
	s, err := runner.Run(ctx, prg, os.Environ(), "")
	if err != nil {
		return err
	}

	fmt.Print(s)
	if !strings.HasSuffix(s, "\n") {
		fmt.Println()

	}

	return nil
}

func (r *GPTScript) getClient(ctx context.Context) (llm.Client, error) {
	if r._client != nil {
		return r._client, nil
	}

	cacheClient, err := cache.New(cache.Options(r.CacheOptions))
	if err != nil {
		return nil, err
	}

	oaClient, err := openai.NewClient(openai.Options(r.OpenAIOptions), openai.Options{
		Cache: cacheClient,
	})
	if err != nil {
		return nil, err
	}

	registry := llm.NewRegistry()

	if err := registry.AddClient(ctx, oaClient); err != nil {
		return nil, err
	}

	r._client = registry
	return r._client, nil
}
