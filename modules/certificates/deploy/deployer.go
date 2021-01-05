package deploy

import (
	"fmt"

	"github.com/r2dtools/agent/webserver"
	"github.com/r2dtools/agentintegration"
)

// CertificateDeployer is an interface for a deployer of certificate to a virtual host
type CertificateDeployer interface {
	DeployCertificate(vhost *agentintegration.VirtualHost, certPath, certKeyPath, chainPath, fullChainPath string) error
}

// GetCertificateDeployer returns certificate deployer for a webserver
func GetCertificateDeployer(webServer webserver.WebServer) (CertificateDeployer, error) {
	if aWebServer, ok := webServer.(*webserver.ApacheWebServer); ok {
		return &ApacheCertificateDeployer{webServer: aWebServer}, nil
	}

	return nil, fmt.Errorf("could not create deployer: webserver '%s' is not supported", webServer.GetCode())
}
