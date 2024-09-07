package hostmng

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type NginxHostManager struct {
	AvailableConfigRootPath string
}

func (m *NginxHostManager) Enable(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file not found: %s", configFilePath)
	}

	fileName := filepath.Base(configFilePath)
	availableConfigFilePath := filepath.Join(m.AvailableConfigRootPath, fileName)

	return os.Symlink(configFilePath, availableConfigFilePath)
}

func (m *NginxHostManager) Disable(configFilePath string) error {
	fileName := filepath.Base(configFilePath)
	availableConfigFilePath := filepath.Join(m.AvailableConfigRootPath, fileName)

	if _, err := os.Lstat(availableConfigFilePath); err == nil {
		if err = os.Remove(availableConfigFilePath); err != nil {
			return fmt.Errorf("failed to remove config file symlink: %s", configFilePath)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file symlink not found: %s", configFilePath)
	}

	return nil
}
