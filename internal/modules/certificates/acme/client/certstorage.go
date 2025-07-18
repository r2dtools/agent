package client

import (
	"github.com/r2dtools/agentintegration"
	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/modules/certificates/acme/client/certbot"
	"github.com/r2dtools/sslbot/internal/modules/certificates/acme/client/lego"
	"github.com/r2dtools/sslbot/internal/pkg/logger"
)

type CertStorage interface {
	AddPemCertificate(certName, pemData string) (certPath string, err error)
	RemoveCertificate(certName string) error
	GetCertificate(certName string) (*agentintegration.Certificate, error)
	GetCertificateAsString(certName string) (certPath string, certContent string, err error)
	GetCertificates() (map[string]*agentintegration.Certificate, error)
	GetCertificatePath(certName string) (certPath string, err error)
}

func CreateCertStorage(config *config.Config, logger logger.Logger) (CertStorage, error) {
	if config.CertBotEnabled {
		return certbot.CreateCertStorage(config, logger)
	}

	return lego.CreateCertStorage(config, logger)
}
