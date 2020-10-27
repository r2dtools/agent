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
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agentintegration"
)

const (
	caServer = "https://acme-v02.api.letsencrypt.org/directory"
	httpPort = 80
	tlsPort  = 443
)

// CertificateManager manages certificates: issue, renew, ....
type CertificateManager struct {
	legoBinPath   string
	challengeType ChallengeType
}

// Issue issues a certificate
func (c *CertificateManager) Issue(certData agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	subjects := []string{certData.ServerName}

	for _, subject := range certData.Subjects {
		if subject != certData.ServerName {
			subjects = append(subjects, subject)
		}
	}

	params := []string{"--email=" + certData.Email, "--domains=" + strings.Join(subjects, " ")}
	_, err := c.execCmd("run", params)

	if err != nil {
		logger.Debug(fmt.Sprintf("%v", err))
		return nil, err
	}

	return nil, nil
}

func (c *CertificateManager) execCmd(command string, params []string) ([]byte, error) {
	dataPath := getDataPath()

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.MkdirAll(dataPath, 0755)
	}

	aParams := []string{"--server=" + c.getCAServer(), "--accept-tos", "--path=" + dataPath}
	aParams = append(aParams, c.challengeType.GetParams()...)
	params = append(params, aParams...)
	params = append(params, command)
	cmd := exec.Command(c.legoBinPath, params...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.Debug(string(output))
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
func GetCertificateManager(certData agentintegration.CertificateIssueRequestData) (*CertificateManager, error) {
	aConfig := config.GetConfig()
	legoBinPath := filepath.Join(aConfig.ExecutablePath, "lego")

	if aConfig.IsSet("LegoBinPath") {
		legoBinPath = aConfig.GetString("LegoBinPath")
	}

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

	certManager := &CertificateManager{
		legoBinPath:   legoBinPath,
		challengeType: challengeType,
	}

	return certManager, nil
}

// getDataPath returns directory path to store data
func getDataPath() string {
	return filepath.Join(config.GetConfig().GetVarDirAbsPath(), "modules", "certificates-module")
}

// getCertificatesDirPath returns path to directory where certificates are stored
func getCertificatesDirPath() string {
	dataPath := getDataPath()

	return filepath.Join(dataPath, "certificates")
}

func getVhostCertificatePath(serverName string) string {
	return filepath.Join(getCertificatesDirPath(), serverName+".crt")
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

	parts := strings.Split(output, "\n")
	var errorParts []string

	for _, part := range parts {
		if strings.Index(part, "[INFO]") != -1 || strings.Index(part, "[WARN]") != -1 {
			continue
		}

		// Skip log time: xxxx/xx/xx xx:xx:xx
		re, err := regexp.Compile(`^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} (.*)`)

		if err == nil {
			rParts := re.FindStringSubmatch(part)

			if len(rParts) > 1 {
				part = rParts[1]
			}
		}

		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		errorParts = append(errorParts, part)
	}

	return strings.Join(errorParts, "\n")
}
