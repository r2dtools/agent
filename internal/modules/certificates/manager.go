package certificates

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates/deploy"
	"github.com/r2dtools/agent/internal/pkg/certificate"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/agentintegration"
	"github.com/unknwon/com"
)

const (
	caServer = "https://acme-v02.api.letsencrypt.org/directory"
	httpPort = 80
	tlsPort  = 443
)

// CertificateManager manages certificates: issue, renew, ....
type CertificateManager struct {
	legoBinPath, dataPath string
	CertStorage           *Storage
	logger                logger.Logger
	Config                *config.Config
}

// Issue issues a certificate
func (c *CertificateManager) Issue(certData agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	serverName := certData.ServerName
	var challengeType ChallengeType

	wServer, err := webserver.GetWebServer(certData.WebServer, c.Config.ToMap())

	if err != nil {
		return nil, err
	}

	if certData.ChallengeType == "" {
		return nil, errors.New("challenge type is not specified")
	}

	if certData.ChallengeType == HttpChallengeTypeCode {
		challengeType = &HTTPChallengeType{
			HTTPPort: httpPort,
			TLSPort:  tlsPort,
			WebRoot:  certData.DocRoot,
		}
	} else if certData.ChallengeType == DnsChallengeTypeCode {
		provider := certData.GetAdditionalParam("provider")
		if provider == "" {
			return nil, errors.New("dns provider is not specified")
		}

		challengeType = &DNSChallengeType{provider}
	} else {
		return nil, fmt.Errorf("unsupported challenge type: %s", certData.ChallengeType)
	}

	params := []string{"--email=" + certData.Email, "--domains=" + serverName}

	for _, subject := range certData.Subjects {
		if subject != serverName {
			params = append(params, "--domains="+subject)
		}
	}

	params = append(params, challengeType.GetParams()...)
	_, err = c.execCmd("run", params)

	if err != nil {
		c.logger.Debug("%v", err)
		return nil, err
	}

	if certData.Assign {
		certPath := c.CertStorage.GetVhostCertificatePath(serverName, "pem")

		return c.deployCertificate(wServer, serverName, certPath, certPath)
	}

	return c.CertStorage.GetCertificate(serverName)
}

// Assign assign certificate from the storage to a domain
func (c *CertificateManager) Assign(certData agentintegration.CertificateAssignRequestData) (*agentintegration.Certificate, error) {
	certPath, err := c.CertStorage.GetCertificatePath(certData.CertName)
	if err != nil {
		return nil, fmt.Errorf("could not assign certificate to the domain '%s': %v", certData.ServerName, err)
	}

	wServer, err := webserver.GetWebServer(certData.WebServer, c.Config.ToMap())

	if err != nil {
		return nil, err
	}

	return c.deployCertificate(wServer, certData.ServerName, certPath, certPath)
}

// Upload deploys an existed certificate
func (c *CertificateManager) Upload(certName, webServer, pemData string) (*agentintegration.Certificate, error) {
	var certPath string
	var err error
	if certPath, err = c.CertStorage.AddPemCertificate(certName, pemData); err != nil {
		return nil, err
	}

	wServer, err := webserver.GetWebServer(webServer, c.Config.ToMap())

	if err != nil {
		return nil, err
	}

	return c.deployCertificate(wServer, certName, certPath, certPath)
}

// GetStorageCertList returns names of all certificates in the storage
func (c *CertificateManager) GetStorageCertList() ([]string, error) {
	return c.CertStorage.GetCertificateNameList()
}

// GetStorageCertData returns certificate by name
func (c *CertificateManager) GetStorageCertData(certName string) (*agentintegration.Certificate, error) {
	return c.CertStorage.GetCertificate(certName)
}

// GetStorageCertData returns certificate by name
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

func (c *CertificateManager) execCmd(command string, params []string) ([]byte, error) {
	c.ensureDataPathExists()
	aParams := []string{"--server=" + c.getCAServer(), "--accept-tos", "--path=" + c.dataPath, "--pem"}
	params = append(params, aParams...)
	params = append(params, command)
	cmd := exec.Command(c.legoBinPath, params...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if len(output) == 0 {
			return nil, err
		}

		return nil, errors.New(getOutputError(string(output)))
	}

	return output, nil
}

func (c *CertificateManager) getCAServer() string {
	if !c.Config.IsSet("CAServer") {
		return caServer
	}

	return c.Config.GetString("CAServer")
}

func GetCertificateManager(config *config.Config, logger logger.Logger) (*CertificateManager, error) {
	legoBinPath := filepath.Join(config.RootPath, "lego")
	dataPath := config.GetModuleVarAbsDir("certificates")

	if config.IsSet("LegoBinPath") {
		legoBinPath = config.GetString("LegoBinPath")
	}

	certManager := &CertificateManager{
		logger:      logger,
		CertStorage: GetDefaultCertStorage(config),
		Config:      config,
		legoBinPath: legoBinPath,
		dataPath:    dataPath,
	}

	return certManager, nil
}

// GetOutputError returns error message from the stdout
func getOutputError(output string) string {
	errIndex := strings.Index(output, "error: ")

	if errIndex != -1 {
		output = output[errIndex:]
	}

	output = strings.ReplaceAll(output, "error: ", "")
	parts := strings.Split(output, "\n")
	var errorParts []string

	for _, part := range parts {
		if strings.Contains(part, "[INFO]") || strings.Contains(part, "[WARN]") {
			continue
		}

		// Skip log time: xxxx/xx/xx xx:xx:xx
		part = removeRegexString(part, `^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} (.*)`)

		if part == "" {
			continue
		}

		errorParts = append(errorParts, part)
	}

	output = strings.Join(errorParts, "\n")

	// Skip ", url:" string. Seems it is a bug in lego library
	// https://github.com/go-acme/lego/blob/master/acme/errors.go#L47
	return removeRegexString(output, `(?s)(.*), url:$`)
}

func removeRegexString(str string, regex string) string {
	re, err := regexp.Compile(regex)

	if err == nil {
		rParts := re.FindStringSubmatch(str)

		if len(rParts) > 1 {
			str = rParts[1]
		}
	}

	return strings.TrimSpace(str)
}

func (c *CertificateManager) ensureDataPathExists() {
	if !com.IsExist(c.dataPath) {
		os.MkdirAll(c.dataPath, 0755)
	}
}
