package tools

import (
	"fmt"
	"os"
	"runtime"
)

const basePath = "opt"

func Resolve(tool string) (string, error) {
	path := fmt.Sprintf("%s/%s/%s", basePath, tool, tool)
	if runtime.GOOS == "windows" {
		path += ".exe"
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("binary not found: %s", tool)
	}

	return path, nil
}
