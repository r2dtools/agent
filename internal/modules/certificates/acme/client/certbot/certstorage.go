package certbot

import (
	"os"

	"github.com/r2dtools/agentintegration"
	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/pkg/logger"
	"github.com/unknwon/com"
)

type CertStorage struct {
	path   string
	logger logger.Logger
}

func (s CertStorage) AddPemCertificate(certName, pemData string) (certPath string, err error) {
	return "", nil
}

func (s CertStorage) RemoveCertificate(certName string) error {
	return nil
}

func (s CertStorage) GetCertificate(certName string) (*agentintegration.Certificate, error) {
	return nil, nil
}

func (s CertStorage) GetCertificateAsString(certName string) (certPath string, certContent string, err error) {
	return "", "", nil
}

func (s CertStorage) GetCertificates() (map[string]*agentintegration.Certificate, error) {
	return nil, nil
}

func (s CertStorage) GetCertificatePath(certName string) (certPath string, err error) {
	return "", nil
}

func CreateCertStorage(config *config.Config, logger logger.Logger) (CertStorage, error) {
	workDir := config.CertBotWokrDir

	if !com.IsExist(workDir) {
		err := os.MkdirAll(workDir, 0755)

		if err != nil {
			return CertStorage{}, err
		}
	}

	return CertStorage{path: workDir, logger: logger}, nil
}
