package deploy

import (
	"fmt"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agentintegration"
)

type CertificateDeployer interface {
	DeployCertificate(vhost *agentintegration.VirtualHost, certPath, certKeyPath, chainPath, fullChainPath string) error
}

func GetCertificateDeployer(webServer webserver.WebServer, logger logger.Logger) (CertificateDeployer, error) {
	if aWebServer, ok := webServer.(*webserver.ApacheWebServer); ok {
		return &ApacheCertificateDeployer{logger: logger, webServer: aWebServer}, nil
	}

	return nil, fmt.Errorf("could not create deployer: webserver '%s' is not supported", webServer.GetCode())
}
