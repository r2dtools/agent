package certificates

import (
	"fmt"
	"path/filepath"

	"github.com/r2dtools/agentintegration"
	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/modules/certificates/acme/client"
	"github.com/r2dtools/sslbot/internal/modules/certificates/commondir"
	"github.com/r2dtools/sslbot/internal/modules/certificates/deploy"
	"github.com/r2dtools/sslbot/internal/pkg/certificate"
	"github.com/r2dtools/sslbot/internal/pkg/logger"
	"github.com/r2dtools/sslbot/internal/pkg/webserver"
	"github.com/r2dtools/sslbot/internal/pkg/webserver/reverter"
)

type CertificateManager struct {
	CertStorage client.CertStorage
	acmeClient  client.AcmeClient
	logger      logger.Logger
	config      *config.Config
}

func (c *CertificateManager) Issue(certData agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	serverName := certData.ServerName

	options := c.config.ToMap()
	wServer, err := webserver.GetWebServer(certData.WebServer, options)

	if err != nil {
		return nil, err
	}

	webServerReverter := &reverter.Reverter{
		HostMng: wServer.GetVhostManager(),
		Logger:  c.logger,
	}
	commonDirManager, err := commondir.GetCommonDirManager(wServer, webServerReverter, c.logger, options)

	if err != nil {
		return nil, err
	}

	vhost, err := wServer.GetVhostByName(serverName)

	if err != nil {
		return nil, err
	}

	if vhost == nil {
		return nil, fmt.Errorf("host %s not found", serverName)
	}

	docRoot := vhost.DocRoot
	commonDir := commonDirManager.GetCommonDirStatus(serverName)

	if commonDir.Enabled {
		docRoot = commonDir.Root
	}

	err = c.acmeClient.Issue(docRoot, certData)

	if err != nil {
		c.logger.Debug("%v", err)

		return nil, err
	}

	if certData.Assign {
		certPath, err := c.CertStorage.GetCertificatePath(serverName)

		if err != nil {
			return nil, err
		}

		return c.deployCertificate(wServer, serverName, certPath, certPath)
	}

	return c.CertStorage.GetCertificate(serverName)
}

func (c *CertificateManager) Assign(certData agentintegration.CertificateAssignRequestData) (*agentintegration.Certificate, error) {
	certPath, err := c.CertStorage.GetCertificatePath(certData.CertName)
	if err != nil {
		return nil, fmt.Errorf("could not assign certificate to the domain '%s': %v", certData.ServerName, err)
	}

	wServer, err := webserver.GetWebServer(certData.WebServer, c.config.ToMap())

	if err != nil {
		return nil, err
	}

	return c.deployCertificate(wServer, certData.ServerName, certPath, certPath)
}

func (c *CertificateManager) Upload(certName, webServer, pemData string) (*agentintegration.Certificate, error) {
	var certPath string
	var err error
	if certPath, err = c.CertStorage.AddPemCertificate(certName, pemData); err != nil {
		return nil, err
	}

	wServer, err := webserver.GetWebServer(webServer, c.config.ToMap())

	if err != nil {
		return nil, err
	}

	return c.deployCertificate(wServer, certName, certPath, certPath)
}

func (c *CertificateManager) GetStorageCertificates() (map[string]*agentintegration.Certificate, error) {
	return c.CertStorage.GetCertificates()
}

func (c *CertificateManager) GetStorageCertData(certName string) (*agentintegration.Certificate, error) {
	return c.CertStorage.GetCertificate(certName)
}

func (c *CertificateManager) RemoveCertificate(certName string) error {
	return c.CertStorage.RemoveCertificate(certName)
}

func (c *CertificateManager) deployCertificate(wServer webserver.WebServer, serverName, certPath, keyPath string) (*agentintegration.Certificate, error) {
	processManager, err := wServer.GetProcessManager()

	if err != nil {
		return nil, err
	}

	vhost, err := wServer.GetVhostByName(serverName)
	if err != nil {
		return nil, err
	}

	webServerReverter := &reverter.Reverter{
		HostMng: wServer.GetVhostManager(),
		Logger:  c.logger,
	}

	if vhost == nil {
		return nil, fmt.Errorf("could not find virtual host '%s'", serverName)
	}

	deployer, err := deploy.GetCertificateDeployer(wServer, webServerReverter, c.logger)
	if err != nil {
		return nil, err
	}

	sslConfigFilePath, originEnabledConfigFilePath, err := deployer.DeployCertificate(vhost, certPath, keyPath)

	if err != nil {
		if rErr := webServerReverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rallback webserver configuration on cert deploy: %v", rErr))
		}

		return nil, err
	}

	if err = wServer.GetVhostManager().Enable(sslConfigFilePath, filepath.Dir(originEnabledConfigFilePath)); err != nil {
		if rErr := webServerReverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rallback webserver configuration on host enabling: %v", rErr))
		}

		return nil, err
	}

	if err = processManager.Reload(); err != nil {
		if rErr := webServerReverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to rallback webserver configuration on webserver reload: %v", rErr))
		}

		return nil, err
	}

	if err = webServerReverter.Commit(); err != nil {
		if rErr := webServerReverter.Rollback(); rErr != nil {
			c.logger.Error(fmt.Sprintf("failed to commit webserver configuration: %v", rErr))
		}
	}

	return certificate.GetCertificateFromFile(certPath)
}

func GetCertificateManager(config *config.Config, logger logger.Logger) (*CertificateManager, error) {
	storage, err := client.CreateCertStorage(config, logger)

	if err != nil {
		return nil, err
	}

	acmeClient, err := client.CreateAcmeClient(config)

	if err != nil {
		return nil, err
	}

	certManager := &CertificateManager{
		logger:      logger,
		CertStorage: storage,
		config:      config,
		acmeClient:  acmeClient,
	}

	return certManager, nil
}
