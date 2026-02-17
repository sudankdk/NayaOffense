package tools

import "context"

type Tool interface {
	Name() string
	Validate(input map[string]string) error
	Run(ctx context.Context, input map[string]string) (any, error)
}
