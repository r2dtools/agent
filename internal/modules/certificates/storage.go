package certificates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/pkg/certificate"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agentintegration"
	"github.com/unknwon/com"
)

const certExtension = "pem"

type Storage struct {
	path   string
	logger logger.Logger
}

func (s *Storage) AddPemCertificate(certName, pemData string) (string, error) {
	certPath := s.getFilePathByNameWithExt(certName, certExtension)

	if err := os.WriteFile(certPath, []byte(pemData), 0644); err != nil {
		return "", fmt.Errorf("could not save certificate data to the storage: %v", err)
	}

	return certPath, nil
}

func (s *Storage) RemoveCertificate(certName string) error {
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

func (s *Storage) GetCertificate(certName string) (*agentintegration.Certificate, error) {
	certPath, err := s.GetCertificatePath(certName)

	if err != nil {
		return nil, err
	}

	return certificate.GetCertificateFromFile(certPath)
}

func (s *Storage) GetCertificateAsString(certName string) (string, string, error) {
	certPath, err := s.GetCertificatePath(certName)

	if err != nil {
		return "", "", err
	}

	certContent, err := os.ReadFile(certPath)

	if err != nil {
		return "", "", fmt.Errorf("could not read certificate content: %v", err)
	}

	return certPath, string(certContent), nil
}

func (s *Storage) GetCertificates() (map[string]*agentintegration.Certificate, error) {
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

func (s *Storage) GetCertificatePath(certName string) (string, error) {
	certNameMap, err := s.getStorageCertNameMap()

	if err != nil {
		return "", err
	}

	_, ok := certNameMap[certName]

	if !ok {
		return "", fmt.Errorf("could not find certificate '%s'", certName)
	}

	certPath := s.getFilePathByNameWithExt(certName, certExtension)

	return certPath, nil
}

func (s *Storage) getFilePathByNameWithExt(fileName, extension string) string {
	return filepath.Join(s.path, fileName+"."+extension)
}

func (s *Storage) getCertificateKeyPath(certName string) string {
	return s.getFilePathByNameWithExt(certName, "key")
}

func (s *Storage) getStorageCertNameMap() (map[string]struct{}, error) {
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

func GetDefaultCertStorage(config *config.Config, logger logger.Logger) (*Storage, error) {
	dataPath := config.GetPathInsideVarDir("certificates")

	if !com.IsExist(dataPath) {
		err := os.MkdirAll(dataPath, 0755)

		if err != nil {
			return nil, err
		}
	}

	return &Storage{path: dataPath, logger: logger}, nil
}
