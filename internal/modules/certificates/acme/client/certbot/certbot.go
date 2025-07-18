package certbot

import (
	"fmt"
	"os/exec"

	"github.com/r2dtools/agentintegration"
	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/modules/certificates/acme"
)

type CertBot struct {
	bin string
}

func (b CertBot) Issue(docRoot string, certData agentintegration.CertificateIssueRequestData) error {
	var challengeType acme.ChallengeType
	serverName := certData.ServerName
	params := []string{"certonly", "-m " + certData.Email, "-n"}

	switch certData.ChallengeType {
	case acme.HttpChallengeTypeCode:
		challengeType = HTTPChallengeType{WebRoot: docRoot}
	default:
		return fmt.Errorf("unsupported challenge type: %s", certData.ChallengeType)
	}

	params = append(params, challengeType.GetParams()...)
	params = append(params, "-d "+serverName)

	for _, subject := range certData.Subjects {
		if subject != serverName {
			params = append(params, "-d "+subject)
		}
	}

	params = append(params, "--agree-tos")
	cmdName := b.bin

	if cmdName == "" {
		cmdName = "certbot"
	}

	cmd := exec.Command(cmdName, params...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if len(output) == 0 {
			return err
		}

		return fmt.Errorf("%s\n%s", output, err.Error())
	}

	return nil
}

func CreateCertBot(config *config.Config) CertBot {
	return CertBot{bin: config.CertBotBin}
}
