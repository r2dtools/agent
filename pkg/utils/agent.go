package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agent/config"
)

// GetAgentVersion returns agent version
func GetAgentVersion() (string, error) {
	config := config.GetConfig()
	output, err := ioutil.ReadFile(filepath.Join(config.ExecutablePath, ".version"))

	if err != nil {
		return "", fmt.Errorf("could not detect agent version: %v", err)
	}

	return strings.Trim(string(output), " \n"), nil
}
