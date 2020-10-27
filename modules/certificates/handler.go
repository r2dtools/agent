package certificates

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agentintegration"
)

// Handler handles requests to the module
type Handler struct{}

// Handle handles request to the module
func (h *Handler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "issue":
		response, err = issue(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func issue(data interface{}) (*agentintegration.Certificate, error) {
	var certData agentintegration.CertificateIssueRequestData

	err := mapstructure.Decode(data, &certData)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate request data: %v", err)
	}

	certManager, err := GetCertificateManager(certData)

	if err != nil {
		return nil, err
	}

	certificate, err := certManager.Issue(certData)

	if err != nil {
		return nil, err
	}

	return certificate, nil
}
