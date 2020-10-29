package deploy

import (
	"github.com/r2dtools/agent/webserver"
	"github.com/r2dtools/agentintegration"
)

// CertificateDeployer is an interface for a deployer of certificate to a virtual host
type CertificateDeployer interface {
	DeployCertificate(vhost agentintegration.VirtualHost)
}

// GetCertificateDeployer returns certificate deployer for a web server
func GetCertificateDeployer(webServer webserver.WebServer) CertificateDeployer {
	return nil
}
