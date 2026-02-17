package runner

import (
	"context"
	"os/exec"
	"time"
)

func Execute(ctx context.Context, bin string, args ...string) (string, error) {
	defaultTimeout := 60 * time.Second
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
