package certificates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/pkg/certificate"
	"github.com/r2dtools/agentintegration"
	"github.com/unknwon/com"
)

type Storage struct {
	Path string
}

// AddCertificate add .pem certificate to the storage
func (s *Storage) AddPemCertificate(certName, pemData string) (string, error) {
	certPath := s.GetVhostCertificatePath(certName, "pem")
	s.ensureCertificatesDirPathExists()

	if err := os.WriteFile(certPath, []byte(pemData), 0644); err != nil {
		return "", fmt.Errorf("could not save certificate data to the storage: %v", err)
	}

	return certPath, nil
}

// RemoveCertificate remove certificate from the storage
func (s *Storage) RemoveCertificate(certName string) error {
	certPemPath := s.GetVhostCertificatePath(certName, "pem")
	certCrtPath := s.GetVhostCertificatePath(certName, "crt")
	certIssuerCrtPath := s.GetVhostCertificatePath(certName, "issuer.crt")
	certJsonData := s.GetVhostCertificatePath(certName, "json")
	keyPath := s.GetVhostCertificateKeyPath(certName)
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

// GetStorageCertList returns names of all certificates in the storage
func (s *Storage) GetCertificateNameList() ([]string, error) {
	certNameList := []string{}
	certNameMap, err := s.getStorageCertNameMap()
	if err != nil {
		return certNameList, err
	}
	for name := range certNameMap {
		certNameList = append(certNameList, name)
	}

	return certNameList, err
}

// GetStorageCertData returns certificate by name
func (s *Storage) GetCertificate(certName string) (*agentintegration.Certificate, error) {
	certPath, err := s.GetCertificatePath(certName)
	if err != nil {
		return nil, err
	}

	return certificate.GetCertificateFromFile(certPath)
}

// GetCertificateAsString returns certificate file path and content
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

func (s *Storage) getStorageCertNameMap() (map[string]string, error) {
	certExtensions := []string{".pem"}
	certNameMap := make(map[string]string)
	certPath := s.GetCertificatesDirPath()
	if !com.IsExist(certPath) {
		return nil, nil
	}
	entries, err := os.ReadDir(certPath)
	if err != nil {
		return nil, fmt.Errorf("could not get the list of certificates in the storage: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		certExt := filepath.Ext(name)

		if !com.IsSliceContainsStr(certExtensions, certExt) {
			continue
		}

		baseName := strings.TrimSuffix(name, certExt)
		if filepath.Ext(baseName) == ".issuer" {
			continue
		}
		certNameMap[name[:len(name)-len(certExt)]] = certExt
	}
	return certNameMap, nil
}

// getCertificatesDirPath returns path to directory where certificates are stored
func (s *Storage) GetCertificatesDirPath() string {
	return filepath.Join(s.Path, "certificates")
}

func (s *Storage) GetVhostCertificatePath(certName, extension string) string {
	return filepath.Join(s.GetCertificatesDirPath(), certName+"."+extension)
}

func (s *Storage) GetVhostCertificateKeyPath(certName string) string {
	return s.GetVhostCertificatePath(certName, "key")
}

func (s *Storage) GetCertificatePath(certName string) (string, error) {
	certNameMap, err := s.getStorageCertNameMap()
	if err != nil {
		return "", err
	}
	certExt, ok := certNameMap[certName]
	if !ok {
		return "", fmt.Errorf("could not find certificate '%s'", certName)
	}
	certPath := s.GetVhostCertificatePath(certName, strings.TrimPrefix(certExt, "."))

	return certPath, nil
}

func (s *Storage) ensureCertificatesDirPathExists() {
	certsDirPath := s.GetCertificatesDirPath()
	if !com.IsExist(certsDirPath) {
		os.MkdirAll(certsDirPath, 0755)
	}
}

func GetDefaultCertStorage(config *config.Config) *Storage {
	dataPath := config.GetModuleVarAbsDir("certificates")

	return &Storage{dataPath}
}
