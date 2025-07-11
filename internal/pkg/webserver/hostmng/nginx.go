package hostmng

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/r2dtools/sslbot/internal/pkg/utils"
)

type NginxHostManager struct {
	EnabledConfigRootPath string
}

func (m *NginxHostManager) Enable(configFilePath, enabledConfigRootPath string) error {
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

	if enabledConfigRootPath == "" {
		enabledConfigFilePath = filepath.Join(m.EnabledConfigRootPath, fileName)
	} else {
		enabledConfigFilePath = filepath.Join(enabledConfigRootPath, fileName)
	}

	if _, err := os.Lstat(enabledConfigFilePath); errors.Is(err, os.ErrNotExist) {
		return os.Symlink(configFilePath, enabledConfigFilePath)
	}

	return nil
}

func (m *NginxHostManager) Disable(enabledConfigFilePath string) error {
	var err error

	if _, err = os.Lstat(enabledConfigFilePath); err == nil {
		if err = os.Remove(enabledConfigFilePath); err != nil {
			return fmt.Errorf("failed to remove config file symlink: %s", enabledConfigFilePath)
		}
	}

	return err
}
