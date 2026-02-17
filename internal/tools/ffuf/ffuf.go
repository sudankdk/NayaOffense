package ffuf

import (
	"context"
	"errors"

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
	return runner.Execute(ctx, bin,
		"-u", input["url"],
		"-w", input["wordlist"],
		"-of", "json")

}

func init() {
	tools.RegisterTool(&Tool{})
}
