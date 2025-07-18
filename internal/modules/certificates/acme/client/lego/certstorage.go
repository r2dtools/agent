package lego

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agentintegration"
	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/pkg/certificate"
	"github.com/r2dtools/sslbot/internal/pkg/logger"
	"github.com/unknwon/com"
)

const certExtension = "pem"

type CertStorage struct {
	path   string
	logger logger.Logger
}

func (s CertStorage) AddPemCertificate(certName, pemData string) (certPath string, err error) {
	certPath = s.getFilePathByNameWithExt(certName, certExtension)

	if err := os.WriteFile(certPath, []byte(pemData), 0644); err != nil {
		return "", fmt.Errorf("could not save certificate data to the storage: %v", err)
	}

	return certPath, nil
}

func (s CertStorage) RemoveCertificate(certName string) error {
	certPemPath := s.getFilePathByNameWithExt(certName, certExtension)
	certCrtPath := s.getFilePathByNameWithExt(certName, "crt")
	certIssuerCrtPath := s.getFilePathByNameWithExt(certName, "issuer.crt")
	certJsonData := s.getFilePathByNameWithExt(certName, "json")
	keyPath := s.getCertificateKeyPath(certName)
	rPaths := []string{certPemPath, certCrtPath}
	nrPaths := []string{certIssuerCrtPath, keyPath, certJsonData}

	for _, path := range rPaths {
		if com.IsFile(path) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("could not remove certificate %s: %v", certName, err)
			}
		}
	}

	for _, path := range nrPaths {
		if com.IsFile(path) {
			os.Remove(path)
		}
	}

	return nil
}

func (s CertStorage) GetCertificate(certName string) (*agentintegration.Certificate, error) {
	certPath, err := s.GetCertificatePath(certName)

	if err != nil {
		return nil, err
	}

	return certificate.GetCertificateFromFile(certPath)
}

func (s CertStorage) GetCertificateAsString(certName string) (certPath string, certContent string, err error) {
	certPath, err = s.GetCertificatePath(certName)

	if err != nil {
		return "", "", err
	}

	certContentBytes, err := os.ReadFile(certPath)

	if err != nil {
		return "", "", fmt.Errorf("could not read certificate content: %v", err)
	}

	return certPath, string(certContentBytes), nil
}

func (s CertStorage) GetCertificates() (map[string]*agentintegration.Certificate, error) {
	certNameMap, err := s.getStorageCertNameMap()

	if err != nil {
		return nil, err
	}

	certsMap := map[string]*agentintegration.Certificate{}

	for certName := range certNameMap {
		certPath := s.getFilePathByNameWithExt(certName, certExtension)
		cert, err := certificate.GetCertificateFromFile(certPath)

		if err != nil {
			s.logger.Error("failed to parse certificate %s: %v", certName, err)

			continue
		}

		certsMap[certName] = cert
	}

	return certsMap, nil
}

func (s CertStorage) GetCertificatePath(certName string) (certPath string, err error) {
	certNameMap, err := s.getStorageCertNameMap()

	if err != nil {
		return "", err
	}

	_, ok := certNameMap[certName]

	if !ok {
		return "", fmt.Errorf("could not find certificate '%s'", certName)
	}

	certPath = s.getFilePathByNameWithExt(certName, certExtension)

	return certPath, nil
}

func (s CertStorage) getFilePathByNameWithExt(fileName, extension string) string {
	return filepath.Join(s.path, fileName+"."+extension)
}

func (s CertStorage) getCertificateKeyPath(certName string) string {
	return s.getFilePathByNameWithExt(certName, "key")
}

func (s CertStorage) getStorageCertNameMap() (map[string]struct{}, error) {
	certNameMap := make(map[string]struct{})
	entries, err := os.ReadDir(s.path)

	if err != nil {
		return nil, fmt.Errorf("could not get certificate list in the storage: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		certExt := filepath.Ext(name)

		if strings.Trim(certExt, ".") != certExtension {
			continue
		}

		certNameMap[name[:(len(name)-len(certExt))]] = struct{}{}
	}

	return certNameMap, nil
}

func CreateCertStorage(config *config.Config, logger logger.Logger) (CertStorage, error) {
	dataPath := config.GetPathInsideVarDir("ssl", "certificates")

	if !com.IsExist(dataPath) {
		err := os.MkdirAll(dataPath, 0755)

		if err != nil {
			return CertStorage{}, err
		}
	}

	return CertStorage{path: dataPath, logger: logger}, nil
}
