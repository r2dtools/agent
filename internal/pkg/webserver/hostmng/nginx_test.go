package hostmng

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNginxHost(t *testing.T) {
	hostManager := NginxHostManager{
		EnabledConfigRootPath: "/etc/nginx/sites-enabled",
	}

	enabledConfigFilePath := "/etc/nginx/sites-enabled/example3.com.conf"
	availableConfigFilePath := "/etc/nginx/sites-available/example3.com.conf"

	_, err := os.Lstat(enabledConfigFilePath)
	assert.Nil(t, err)

	err = hostManager.Disable(enabledConfigFilePath)
	assert.Nil(t, err)
	_, err = os.Lstat(enabledConfigFilePath)
	assert.NotNil(t, err)

	err = hostManager.Enable(availableConfigFilePath, filepath.Dir(enabledConfigFilePath))
	assert.Nil(t, err)
	_, err = os.Lstat(enabledConfigFilePath)
	assert.Nil(t, err)
}
