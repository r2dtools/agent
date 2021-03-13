package certificates

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agent/system"
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

	if err := system.GetPrivilege().IncreasePrivilege(); err != nil {
		logger.Error(fmt.Sprintf("certificate issue: increase privilege failed: %v", err))
	}

	defer (func() {
		if err := system.GetPrivilege().DropPrivilege(); err != nil {
			logger.Error(fmt.Sprintf("certificate issue: drop privilege failed: %v", err))
		}
	})()
	certManager, err := GetCertificateManager(certData)

	if err != nil {
		return nil, err
	}

	return certManager.Issue(certData)
}
