package ffuf

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/sudankdk/offense/internal/runner"
	"github.com/sudankdk/offense/internal/tools"
)

type Tool struct{}

func (t *Tool) Name() string {
	return "ffuf"
}

func (t *Tool) Validate(input map[string]string) error {
	url, okUrl := input["url"]
	wordlist, okWordlist := input["wordlist"]
	if !okUrl || !okWordlist || url == "" || wordlist == "" {
		return errors.New("url and wordlist required")
	}
	return nil
}

func (t *Tool) Run(ctx context.Context, input map[string]string) (any, error) {
	if err := t.Validate(input); err != nil {
		return nil, err
	}

	bin, err := tools.Resolve(t.Name())
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "ffuf-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp output file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	_, err = runner.Execute(ctx, bin,
		"-u", input["url"],
		"-w", input["wordlist"],
		"-of", "json",
		"-o", tmpPath,
		"-s",
	)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ffuf output: %w", err)
	}
	if len(raw) == 0 {
		return nil, errors.New("ffuf produced no output")
	}

	return parseOutput(string(raw))
}

func init() {
	tools.RegisterTool(&Tool{})
}
