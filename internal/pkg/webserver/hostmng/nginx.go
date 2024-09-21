package hostmng

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type NginxHostManager struct {
	AvailableConfigRootPath string
	EnabledConfigRootPath   string
}

func (m *NginxHostManager) Enable(configFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file not found: %s", configFilePath)
	}

	fileName := filepath.Base(configFilePath)
	enabledConfigFilePath := filepath.Join(m.EnabledConfigRootPath, fileName)

	if _, err := os.Lstat(enabledConfigFilePath); errors.Is(err, os.ErrNotExist) {
		return os.Symlink(configFilePath, enabledConfigFilePath)
	}

	return nil
}

func (m *NginxHostManager) Disable(configFilePath string) error {
	fileName := filepath.Base(configFilePath)
	enabledConfigFilePath := filepath.Join(m.EnabledConfigRootPath, fileName)

	if _, err := os.Lstat(enabledConfigFilePath); err == nil {
		if err = os.Remove(enabledConfigFilePath); err != nil {
			return fmt.Errorf("failed to remove config file symlink: %s", configFilePath)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file symlink not found: %s", configFilePath)
	}

	return nil
}
