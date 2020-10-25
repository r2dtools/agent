package certificates

import (
	"fmt"
	"os/exec"

	"github.com/r2dtools/agentintegration"
)

// CertificateManager manages certificates: issue, renew, ....
type CertificateManager struct {
	LegoBinPath string
}

// Issue issues a certificate
func (c *CertificateManager) Issue(agentintegration.CertificateIssueRequestData) (*agentintegration.Certificate, error) {
	return nil, nil
}

func (c *CertificateManager) execCmd(params []string) ([]byte, error) {
	cmd := exec.Command(c.LegoBinPath, params...)
	output, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("could not execute lego command: %v", err)
	}

	return output, nil
}
