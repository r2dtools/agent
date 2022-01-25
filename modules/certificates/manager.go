package certificates

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/r2dtools/agent/certificate"
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/modules/certificates/deploy"
	"github.com/r2dtools/agent/modules/certificates/utils"
	"github.com/r2dtools/agent/webserver"
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
	legoBinPath string
}

// Issue issues a certificate
func (c *CertificateManager) Issue(certData agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	serverName := certData.ServerName
	var challengeType ChallengeType

	if certData.ChallengeType == "" {
		return nil, errors.New("challenge type is not specified")
	}

	if certData.ChallengeType == httpChallengeType {
		challengeType = &HTTPChallengeType{
			HTTPPort: httpPort,
			TLSPort:  tlsPort,
			WebRoot:  certData.DocRoot,
		}
	} else if certData.ChallengeType == dnsChallengeType {
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
	_, err := c.execCmd("run", params)

	if err != nil {
		logger.Debug(fmt.Sprintf("%v", err))
		return nil, err
	}

	certPath := getVhostCertificatePath(serverName, "crt")
	keyPath := getVhostCertificateKeyPath(serverName)

	return c.deployCertificate(serverName, certData.WebServer, certPath, keyPath)
}

// Upload deploys an existed certificate
func (c *CertificateManager) Upload(certData *agentintegration.CertificateUploadRequestData) (*agentintegration.Certificate, error) {
	cert, err := utils.LoadCertficateAndKeyFromPem(certData.PemCertificate)
	if err != nil {
		return nil, fmt.Errorf("uploaded certificate is invalid: %v", err)
	}

	certPath := getVhostCertificatePath(certData.ServerName, "pem")
	keyPath := getVhostCertificateKeyPath(certData.ServerName)
	ensureCertificatesDirPathExists()

	if err = os.WriteFile(certPath, []byte(certData.PemCertificate), 0644); err != nil {
		return nil, fmt.Errorf("could not save certificate data: %v", err)
	}

	if err = os.WriteFile(keyPath, cert.PrivateKey, 0644); err != nil {
		return nil, fmt.Errorf("could not save certificate private key: %v", err)
	}

	return c.deployCertificate(certData.ServerName, certData.WebServer, certPath, keyPath)
}

func (c *CertificateManager) deployCertificate(serverName, webServer, certPath, keyPath string) (*agentintegration.Certificate, error) {
	webserver, err := webserver.GetWebServer(webServer, config.GetConfig().ToMap())
	if err != nil {
		return nil, err
	}

	vhost, err := webserver.GetVhostByName(serverName)
	if err != nil {
		return nil, err
	}

	if vhost == nil {
		return nil, fmt.Errorf("could not find virtual host '%s'", serverName)
	}

	deployer, err := deploy.GetCertificateDeployer(webserver)
	if err != nil {
		return nil, err
	}

	if err = deployer.DeployCertificate(vhost, certPath, keyPath, "", certPath); err != nil {
		return nil, err
	}

	return certificate.GetCertificateForDomainFromRequest(serverName)
}

func (c *CertificateManager) execCmd(command string, params []string) ([]byte, error) {
	ensureDataPathExists()
	dataPath := getDataPath()
	aParams := []string{"--server=" + c.getCAServer(), "--accept-tos", "--path=" + dataPath}
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
	aConfig := config.GetConfig()

	if !aConfig.IsSet("CAServer") {
		return caServer
	}

	return aConfig.GetString("CAServer")
}

// GetCertificateManager creates CertificateManager instance
func GetCertificateManager() (*CertificateManager, error) {
	aConfig := config.GetConfig()
	legoBinPath := filepath.Join(aConfig.ExecutablePath, "lego")

	if aConfig.IsSet("LegoBinPath") {
		legoBinPath = aConfig.GetString("LegoBinPath")
	}

	certManager := &CertificateManager{
		legoBinPath: legoBinPath,
	}

	return certManager, nil
}

// getDataPath returns directory path to store data
func getDataPath() string {
	return config.GetConfig().GetModuleVarAbsDir("certificates")
}

// getCertificatesDirPath returns path to directory where certificates are stored
func getCertificatesDirPath() string {
	dataPath := getDataPath()
	return filepath.Join(dataPath, "certificates")
}

func getVhostCertificatePath(serverName, extension string) string {
	return filepath.Join(getCertificatesDirPath(), serverName+"."+extension)
}

func getVhostCertificateKeyPath(serverName string) string {
	return filepath.Join(getCertificatesDirPath(), serverName+".key")
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
		if strings.Index(part, "[INFO]") != -1 || strings.Index(part, "[WARN]") != -1 {
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

func ensureDataPathExists() {
	dataPath := getDataPath()

	if !com.IsExist(dataPath) {
		os.MkdirAll(dataPath, 0755)
	}
}

func ensureCertificatesDirPathExists() {
	certsDirPath := getCertificatesDirPath()

	if !com.IsExist(certsDirPath) {
		os.MkdirAll(certsDirPath, 0755)
	}
}
