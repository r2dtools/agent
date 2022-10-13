package deploy

import (
	"fmt"

	"github.com/r2dtools/a2conf"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agentintegration"
)

// ApacheCertificateDeployer certificate deployer to apache virtual host
type ApacheCertificateDeployer struct {
	logger    logger.LoggerInterface
	webServer *webserver.ApacheWebServer
}

// DeployCertificate deploys certificate to apache domain
func (d *ApacheCertificateDeployer) DeployCertificate(vhost *agentintegration.VirtualHost, certPath, certKeyPath, chainPath, fullChainPath string) error {
	configurator := d.webServer.GetApacheConfigurator()

	if err := configurator.DeployCertificate(vhost.ServerName, certPath, certKeyPath, chainPath, fullChainPath); err != nil {
		d.rollback(configurator)
		return fmt.Errorf("could not deploy certificate to virtual host '%s': %v", vhost.ServerName, err)
	}

	if err := configurator.Save(); err != nil {
		message := fmt.Sprintf("could not deploy certificate for virtual host '%s': could not save changes for apache configuration: %v", vhost.ServerName, err)
		d.logger.Error(message)
		d.rollback(configurator)

		return fmt.Errorf(message)
	}

	if !configurator.CheckConfiguration() {
		message := fmt.Sprintf("could not deploy certificate for virtual host '%s': apache configuration is invalid.", vhost.ServerName)
		d.logger.Error(message)
		d.rollback(configurator)

		return fmt.Errorf(message)
	}

	if err := configurator.Commit(); err != nil {
		d.logger.Error(fmt.Sprintf("error while deploying certificate to virtual host '%s': %v", vhost.ServerName, err))
	}

	if err := configurator.RestartWebServer(); err != nil {
		return err
	}

	return nil
}

func (d *ApacheCertificateDeployer) rollback(configurator a2conf.ApacheConfigurator) {
	if err := configurator.Rollback(); err != nil {
		d.logger.Error(fmt.Sprintf("could not rollback apache configuration: %v", err))
	}
}
