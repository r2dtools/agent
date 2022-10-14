package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agent/config"
)

func GetAgentVersion(config *config.Config) (string, error) {
	output, err := os.ReadFile(filepath.Join(config.ExecutablePath, ".version"))

	if err != nil {
		return "", fmt.Errorf("could not detect agent version: %v", err)
	}

	return strings.Trim(string(output), " \n"), nil
}
