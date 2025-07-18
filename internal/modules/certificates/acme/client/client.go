package client

import (
	"github.com/r2dtools/agentintegration"
	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/modules/certificates/acme/client/certbot"
	"github.com/r2dtools/sslbot/internal/modules/certificates/acme/client/lego"
)

type AcmeClient interface {
	Issue(docRoot string, certData agentintegration.CertificateIssueRequestData) error
}

func CreateAcmeClient(config *config.Config) (AcmeClient, error) {
	if config.CertBotEnabled {
		return certbot.CreateCertBot(config), nil
	}

	return lego.CreateClient(config)
}
