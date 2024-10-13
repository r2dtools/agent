package commondir

import (
	"fmt"
	"path/filepath"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/gonginx/config"
)

const (
	acmeLocation = "/.well-known/acme-challenge/"
	acmeRoot     = "/var/www/html/"
)

type NginxCommonDirManager struct {
	webServer *webserver.NginxWebServer
	reverter  *reverter.Reverter
	logger    logger.Logger
}

func (c *NginxCommonDirManager) EnableCommonDir(serverName string) error {
	wConfig := c.webServer.Config
	nonSslServerBlock := c.findNonSslServerBlock(serverName)

	if nonSslServerBlock == nil {
		return fmt.Errorf("nginx host %s on 80 port does not exist", serverName)
	}

	nonSslServerBlockFileName := filepath.Base(nonSslServerBlock.FilePath)
	configFile := wConfig.GetConfigFile(nonSslServerBlockFileName)

	if configFile == nil {
		return fmt.Errorf("failed to find config file for host %s", serverName)
	}

	if c.findCommonDirBlock(nonSslServerBlock) != nil {
		c.logger.Info("common directory is already enabled for %s host", serverName)

		return nil
	}

	commonDirLocationBlock := nonSslServerBlock.AddLocationBlock("^~", acmeLocation, true)
	commonDirLocationBlock.AddDirective(config.NewDirective("root", []string{acmeRoot}), true, false)
	commonDirLocationBlock.AddDirective(config.NewDirective("default_type", []string{`"text/plain"`}), true, false)

	if err := c.reverter.BackupConfig(nonSslServerBlock.FilePath); err != nil {
		return err
	}

	return configFile.Dump()
}

func (c *NginxCommonDirManager) DisableCommonDir(serverName string) error {
	wConfig := c.webServer.Config
	nonSslServerBlock := c.findNonSslServerBlock(serverName)

	if nonSslServerBlock == nil {
		return fmt.Errorf("nginx host %s on 80 port does not exist", serverName)
	}

	nonSslServerBlockFileName := filepath.Base(nonSslServerBlock.FilePath)
	configFile := wConfig.GetConfigFile(nonSslServerBlockFileName)

	if configFile == nil {
		return fmt.Errorf("failed to find config file for host %s", serverName)
	}

	commonDirBlock := c.findCommonDirBlock(nonSslServerBlock)

	if commonDirBlock == nil {
		return nil
	}

	nonSslServerBlock.DeleteLocationBlock(*commonDirBlock)

	if err := c.reverter.BackupConfig(nonSslServerBlock.FilePath); err != nil {
		return err
	}

	return configFile.Dump()
}

func (c *NginxCommonDirManager) IsCommonDirEnabled(serverName string) bool {
	nonSslServerBlock := c.findNonSslServerBlock(serverName)

	if nonSslServerBlock == nil {
		return false
	}

	return c.findCommonDirBlock(nonSslServerBlock) != nil
}

func (c *NginxCommonDirManager) findCommonDirBlock(serverBlock *config.ServerBlock) *config.LocationBlock {
	locationBlocks := serverBlock.FindLocationBlocks()

	for _, locationBlock := range locationBlocks {
		if locationBlock.GetLocationMatch() == acmeLocation {
			return &locationBlock
		}
	}

	return nil
}

func (c *NginxCommonDirManager) findNonSslServerBlock(serverName string) *config.ServerBlock {
	wConfig := c.webServer.Config
	serverBlocks := wConfig.FindServerBlocksByServerName(serverName)

	if len(serverBlocks) == 0 {
		return nil
	}

	var nonSslServerBlock *config.ServerBlock

out:
	for _, serverBlock := range serverBlocks {
		for _, address := range serverBlock.GetAddresses() {
			if address.Port == "80" {
				nonSslServerBlock = &serverBlock

				break out
			}
		}
	}

	return nonSslServerBlock
}
