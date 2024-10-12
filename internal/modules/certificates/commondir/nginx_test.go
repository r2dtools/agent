package commondir

import (
	"strings"
	"testing"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/gonginx/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestNginxCommonDir(t *testing.T) {
	host := "example2.com"

	manager, nginxWebServer, rv := getNginxCommonDirManager(t)
	defer rv.Rollback()

	assert.False(t, manager.IsCommonDirEnabled(host))

	err := manager.EnableCommonDir(host)
	assert.Nil(t, err)
	assert.True(t, manager.IsCommonDirEnabled(host))

	blocks := nginxWebServer.Config.FindServerBlocksByServerName(host)
	assert.Len(t, blocks, 1)

	block := blocks[0]
	locations := block.FindLocationBlocks()
	assert.Len(t, locations, 2)

	acmeBlockExists := slices.ContainsFunc(locations, func(block config.LocationBlock) bool {
		return strings.Contains(block.GetLocationMatch(), acmeLocation)
	})
	assert.True(t, acmeBlockExists)

	err = manager.DisableCommonDir(host)
	assert.Nil(t, err)
	assert.False(t, manager.IsCommonDirEnabled(host))

	locations = block.FindLocationBlocks()
	assert.Len(t, locations, 1)

	acmeBlockExists = slices.ContainsFunc(locations, func(block config.LocationBlock) bool {
		return strings.Contains(block.GetLocationMatch(), acmeLocation)
	})
	assert.False(t, acmeBlockExists)
}

func getNginxCommonDirManager(t *testing.T) (CommonDirManager, webserver.NginxWebServer, *reverter.Reverter) {
	nginxWebServer, err := webserver.GetNginxWebServer(nil)
	assert.Nil(t, err)

	rv := &reverter.Reverter{
		Logger: &logger.NilLogger{},
	}
	manager, err := GetCommonDirManager(nginxWebServer, rv, &logger.NilLogger{})
	assert.Nil(t, err)

	return manager, *nginxWebServer, rv
}
