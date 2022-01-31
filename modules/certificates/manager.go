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
	legoBinPath, dataPath string
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

	certPath := c.getVhostCertificatePath(serverName, "crt")
	keyPath := c.getVhostCertificateKeyPath(serverName)

	return c.deployCertificate(serverName, certData.WebServer, certPath, keyPath)
}

// Upload deploys an existed certificate
func (c *CertificateManager) Upload(certData *agentintegration.CertificateUploadRequestData) (*agentintegration.Certificate, error) {
	cert, err := utils.LoadCertficateAndKeyFromPem(certData.PemCertificate)
	if err != nil {
		return nil, fmt.Errorf("uploaded certificate is invalid: %v", err)
	}

	certPath := c.getVhostCertificatePath(certData.ServerName, "pem")
	keyPath := c.getVhostCertificateKeyPath(certData.ServerName)
	c.ensureCertificatesDirPathExists()

	if err = os.WriteFile(certPath, []byte(certData.PemCertificate), 0644); err != nil {
		return nil, fmt.Errorf("could not save certificate data: %v", err)
	}

	if err = os.WriteFile(keyPath, cert.PrivateKey, 0644); err != nil {
		return nil, fmt.Errorf("could not save certificate private key: %v", err)
	}

	return c.deployCertificate(certData.ServerName, certData.WebServer, certPath, keyPath)
}

// GetStorageCertList returns names of all certificates in the storage
func (c *CertificateManager) GetStorageCertList() ([]string, error) {
	certNameList := []string{}
	certNameMap, err := c.getStorageCertNameMap()
	if err != nil {
		return certNameList, err
	}
	for name, _ := range certNameMap {
		certNameList = append(certNameList, name)
	}

	return certNameList, err
}

// GetStorageCertData returns certificate by name
func (c *CertificateManager) GetStorageCertData(certName string) (*agentintegration.Certificate, error) {
	certNameMap, err := c.getStorageCertNameMap()
	if err != nil {
		return nil, err
	}
	certExt, ok := certNameMap[certName]
	if !ok {
		return nil, fmt.Errorf("could not find certificate '%s'", certName)
	}
	certPath := filepath.Join(c.getCertificatesDirPath(), certName+certExt)

	return certificate.GetCertificateFromFile(certPath)
}

func (c *CertificateManager) getStorageCertNameMap() (map[string]string, error) {
	certExtensions := []string{".crt", ".pem"}
	certNameMap := make(map[string]string)
	certPath := c.getCertificatesDirPath()
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

		certNameMap[name[:len(name)-len(certExt)]] = certExt
	}
	return certNameMap, nil
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
	c.ensureDataPathExists()
	aParams := []string{"--server=" + c.getCAServer(), "--accept-tos", "--path=" + c.dataPath}
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
	dataPath := aConfig.GetModuleVarAbsDir("certificates")

	if aConfig.IsSet("LegoBinPath") {
		legoBinPath = aConfig.GetString("LegoBinPath")
	}

	certManager := &CertificateManager{
		legoBinPath: legoBinPath,
		dataPath:    dataPath,
	}

	return certManager, nil
}

// getCertificatesDirPath returns path to directory where certificates are stored
func (c *CertificateManager) getCertificatesDirPath() string {
	return filepath.Join(c.dataPath, "certificates")
}

func (c *CertificateManager) getVhostCertificatePath(serverName, extension string) string {
	return filepath.Join(c.getCertificatesDirPath(), serverName+"."+extension)
}

func (c *CertificateManager) getVhostCertificateKeyPath(serverName string) string {
	return filepath.Join(c.getCertificatesDirPath(), serverName+".key")
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

func (c *CertificateManager) ensureDataPathExists() {
	if !com.IsExist(c.dataPath) {
		os.MkdirAll(c.dataPath, 0755)
	}
}

func (c *CertificateManager) ensureCertificatesDirPathExists() {
	certsDirPath := c.getCertificatesDirPath()

	if !com.IsExist(certsDirPath) {
		os.MkdirAll(certsDirPath, 0755)
	}
}
