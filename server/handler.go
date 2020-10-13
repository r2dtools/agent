package server

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/r2dtools/agent/certificate"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/utils"
	"github.com/r2dtools/agent/webserver"
	"github.com/r2dtools/agentintegration"
)

// MainHandler handles common agent requests
type MainHandler struct{}

// Handle handles request
func (h *MainHandler) Handle(request Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "refresh":
		response, err = refresh(request.Data)
	case "getVhosts":
		response, err = getVhosts(request.Data)
	case "getVhostCertificate":
		response, err = getVhostCertificate(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func refresh(data interface{}) (*agentintegration.ServerData, error) {
	cmd := exec.Command("bash", "scripts/os.sh")
	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(output), "/")

	// Linux/Ubuntu/20.04/focal
	if len(parts) < 4 {
		logger.Debug(fmt.Sprintf("os.sh script output %s", string(output)))

		return nil, fmt.Errorf("could not get OS data")
	}

	version, err := utils.GetAgentVersion()

	if err != nil {
		return nil, err
	}

	return &agentintegration.ServerData{AgentVersion: version, OsCode: strings.ToLower(parts[1]), OsVersion: parts[2]}, nil
}

func getVhosts(data interface{}) ([]agentintegration.VirtualHost, error) {
	webServerCodes := webserver.GetSupportedWebServers()
	var vhosts []agentintegration.VirtualHost

	for _, webServerCode := range webServerCodes {
		webserver, err := webserver.GetWebServer(webServerCode, nil)

		if err != nil {
			logger.Error(err.Error())
			continue
		}

		wVhosts, err := webserver.GetVhosts()

		if err != nil {
			logger.Error(err.Error())
			continue
		}

		vhosts = append(vhosts, wVhosts...)
	}

	return vhosts, nil
}

func getVhostCertificate(data interface{}) (*agentintegration.Certificate, error) {
	mData, ok := data.(map[string]interface{})

	if !ok {
		return nil, errors.New("invalid request data format")
	}

	vhostNameRaw, ok := mData["vhostName"]

	if !ok {
		return nil, errors.New("invalid request data: vhost name is not specified")
	}

	vhostName, ok := vhostNameRaw.(string)

	if !ok {
		return nil, errors.New("invalid request data: vhost name is invalid")
	}

	certs, err := certificate.GetX509CertificateFromHTTPRequest(vhostName)

	if err != nil {
		logger.Info(fmt.Sprintf("could not get vhost '%s' certificate: %v", vhostName, err))
		return nil, nil
	}

	if len(certs) == 0 {
		return nil, nil
	}

	var roots []*x509.Certificate

	if len(certs) > 1 {
		roots = certs[1:]
	}

	return certificate.ConvertX509CertificateToIntCert(certs[0], roots), nil
}
