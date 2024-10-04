package agent

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agent/config"
)

func GetAgentVersion(config *config.Config) (string, error) {
	path := filepath.Join(config.ExecutablePath, ".version")

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "dev", nil
	}

	output, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("could not detect agent version: %v", err)
	}

	return strings.Trim(string(output), " \n"), nil
}
