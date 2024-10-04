package deploy

import (
	"fmt"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/agentintegration"
)

type CertificateDeployer interface {
	DeployCertificate(vhost *agentintegration.VirtualHost, certPath, certKeyPath string) (string, string, error)
}

func GetCertificateDeployer(webServer webserver.WebServer, reverter *reverter.Reverter, logger logger.Logger) (CertificateDeployer, error) {
	switch w := webServer.(type) {
	case *webserver.NginxWebServer:
		return &NginxCertificateDeployer{logger: logger, webServer: w, reverter: reverter}, nil
	default:
		return nil, fmt.Errorf("could not create deployer: webserver '%s' is not supported", webServer.GetCode())
	}
}
