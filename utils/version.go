package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// GetAgentVersion returns agent version
func GetAgentVersion() (string, error) {
	output, err := ioutil.ReadFile(".version")

	if err != nil {
		return "", fmt.Errorf("could not detect agent version: %v", err)
	}

	return strings.Trim(string(output), " \n"), nil
}
