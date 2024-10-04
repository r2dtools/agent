package hostmng

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/r2dtools/agent/internal/pkg/utils"
)

type NginxHostManager struct {
	AvailableConfigRootPath string
	EnabledConfigRootPath   string
}

func (m *NginxHostManager) Enable(configFilePath, originConfigFilePath string) error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file not found: %s", configFilePath)
	}

	isSymlink, err := utils.IsSymlink(configFilePath)

	if err != nil {
		return err
	}

	if isSymlink {
		// if symlink - host already enabled
		return nil
	}

	var enabledConfigFilePath string
	fileName := filepath.Base(configFilePath)

	if originConfigFilePath == "" {
		enabledConfigFilePath = filepath.Join(m.EnabledConfigRootPath, fileName)
	} else {
		enabledConfigDir := filepath.Dir(originConfigFilePath)
		enabledConfigFilePath = filepath.Join(enabledConfigDir, fileName)
	}

	if _, err := os.Lstat(enabledConfigFilePath); errors.Is(err, os.ErrNotExist) {
		return os.Symlink(configFilePath, enabledConfigFilePath)
	}

	return nil
}

func (m *NginxHostManager) Disable(enabledConfigFilePath string) error {
	if _, err := os.Lstat(enabledConfigFilePath); err == nil {
		if err = os.Remove(enabledConfigFilePath); err != nil {
			return fmt.Errorf("failed to remove config file symlink: %s", enabledConfigFilePath)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file symlink not found: %s", enabledConfigFilePath)
	}

	return nil
}
