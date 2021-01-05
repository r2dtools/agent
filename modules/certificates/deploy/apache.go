package deploy

import (
	"fmt"

	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/webserver"
	"github.com/r2dtools/agentintegration"
)

// ApacheCertificateDeployer certificate deployer to apache virtual host
type ApacheCertificateDeployer struct {
	webServer *webserver.ApacheWebServer
}

// DeployCertificate deploys certificate to apache domain
func (d *ApacheCertificateDeployer) DeployCertificate(vhost *agentintegration.VirtualHost, certPath, certKeyPath, chainPath, fullChainPath string) error {
	configurator := d.webServer.GetApacheConfigurator()

	if err := configurator.DeployCertificate(vhost.ServerName, certPath, certKeyPath, chainPath, fullChainPath); err != nil {
		return fmt.Errorf("could not deploy certificate to virtual host '%s': %v", vhost.ServerName, err)
	}

	if err := configurator.Save(); err != nil {
		message := fmt.Sprintf("could not deploy certificate for virtual host '%s': could not save changes for apache configuration: %v", vhost.ServerName, err)
		logger.Error(message)

		return fmt.Errorf(message)
	}

	if !configurator.CheckConfiguration() {
		message := fmt.Sprintf("could not deploy certificate for virtual host '%s': apache configuration is invalid.", vhost.ServerName)
		logger.Error(message)

		if err := configurator.Rollback(); err != nil {
			logger.Error(fmt.Sprintf("could not rollback apache configuration: %v", err))
		}

		return fmt.Errorf(message)
	}

	if err := configurator.Commit(); err != nil {
		logger.Error(fmt.Sprintf("error while deploying certificate to virtual host '%s': %v", vhost.ServerName, err))
	}

	if err := configurator.RestartWebServer(); err != nil {
		return err
	}

	return nil
}
