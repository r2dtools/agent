package deploy

import (
	"github.com/r2dtools/agent/webserver"
	"github.com/r2dtools/agentintegration"
)

// ApacheCertificateDeployer certificate deployer to apache virtual host
type ApacheCertificateDeployer struct {
	webServer *webserver.ApacheWebServer
}

// DeployCertificate deploys certificate to apache domain
func (d *ApacheCertificateDeployer) DeployCertificate(vhost agentintegration.VirtualHost) {}
