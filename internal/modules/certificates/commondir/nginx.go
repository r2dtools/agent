package commondir

import (
	"fmt"
	"path/filepath"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/gonginxconf/config"
)

const (
	acmeLocation    = "/.well-known/acme-challenge/"
	defaulCommonDir = "/var/www/html/"
)

type NginxCommonDirManager struct {
	webServer *webserver.NginxWebServer
	reverter  *reverter.Reverter
	logger    logger.Logger
	commonDir string
}

func (c *NginxCommonDirManager) EnableCommonDir(serverName string) error {
	wConfig := c.webServer.Config
	serverBlock := c.findServerBlock(serverName)

	if serverBlock == nil {
		return fmt.Errorf("nginx host %s on 80 or 443 port does not exist", serverName)
	}

	serverBlockFileName := filepath.Base(serverBlock.FilePath)
	configFile := wConfig.GetConfigFile(serverBlockFileName)

	if configFile == nil {
		return fmt.Errorf("failed to find config file for host %s", serverName)
	}

	processManager, err := c.webServer.GetProcessManager()

	if err != nil {
		return err
	}

	if c.findCommonDirBlock(serverBlock) != nil {
		c.logger.Info("common directory is already enabled for %s host", serverName)

		return nil
	}

	commonDir := c.commonDir

	if commonDir == "" {
		commonDir = defaulCommonDir
	}

	commonDirLocationBlock := serverBlock.AddLocationBlock("^~", acmeLocation, true)
	commonDirLocationBlock.AddDirective(config.NewDirective("root", []string{commonDir}), true, false)
	commonDirLocationBlock.AddDirective(config.NewDirective("default_type", []string{`"text/plain"`}), true, false)

	if err := c.reverter.BackupConfig(serverBlock.FilePath); err != nil {
		return err
	}

	err = configFile.Dump()

	if err != nil {
		if rErr := c.reverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rollback webserver configuration on common directory switching: %v", rErr))

		}

		return err
	}

	if err := processManager.Reload(); err != nil {
		if rErr := c.reverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rollback webserver configuration on webserver reload: %v", rErr))
		}

		return err
	}

	c.reverter.Commit()

	return nil
}

func (c *NginxCommonDirManager) DisableCommonDir(serverName string) error {
	wConfig := c.webServer.Config
	serverBlock := c.findServerBlock(serverName)

	if serverBlock == nil {
		return fmt.Errorf("nginx host %s on 80 or 443 port does not exist", serverName)
	}

	serverBlockFileName := filepath.Base(serverBlock.FilePath)
	configFile := wConfig.GetConfigFile(serverBlockFileName)

	if configFile == nil {
		return fmt.Errorf("failed to find config file for host %s", serverName)
	}

	processManager, err := c.webServer.GetProcessManager()

	if err != nil {
		return err
	}

	commonDirBlock := c.findCommonDirBlock(serverBlock)

	if commonDirBlock == nil {
		return nil
	}

	serverBlock.DeleteLocationBlock(*commonDirBlock)

	if err := c.reverter.BackupConfig(serverBlock.FilePath); err != nil {
		return err
	}

	err = configFile.Dump()

	if err != nil {
		if rErr := c.reverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rollback webserver configuration on common directory switching: %v", rErr))

		}

		return err
	}

	if err := processManager.Reload(); err != nil {
		if rErr := c.reverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rollback webserver configuration on webserver reload: %v", rErr))
		}

		return err
	}

	c.reverter.Commit()

	return nil
}

func (c *NginxCommonDirManager) IsCommonDirEnabled(serverName string) bool {
	serverBlock := c.findServerBlock(serverName)

	if serverBlock == nil {
		return false
	}

	return c.findCommonDirBlock(serverBlock) != nil
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

func (c *NginxCommonDirManager) findServerBlock(serverName string) *config.ServerBlock {
	wConfig := c.webServer.Config
	serverBlocks := wConfig.FindServerBlocksByServerName(serverName)

	if len(serverBlocks) == 0 {
		return nil
	}

	var nonSslServerBlock *config.ServerBlock

	for _, serverBlock := range serverBlocks {
		for _, address := range serverBlock.GetAddresses() {
			if address.Port == "80" && nonSslServerBlock == nil {
				nonSslServerBlock = &serverBlock

				continue
			}

			if address.Port == "443" {
				return &serverBlock
			}
		}
	}

	return nonSslServerBlock
}
