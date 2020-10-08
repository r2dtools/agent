package server

import (
	"errors"
	"fmt"

	"github.com/r2dtools/agent/certificate"
	"github.com/r2dtools/agent/logger"
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

func refresh(data interface{}) (agentintegration.ServerData, error) {
	return agentintegration.ServerData{AgentVersion: "1.0.0", OsCode: "ubuntu", OsVersion: "18.04"}, nil
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

func getVhostCertificate(data interface{}) (interface{}, error) {
	mData, ok := data.(map[string]string)

	if !ok {
		return nil, errors.New("invalid request data format")
	}

	serverName, ok := mData["serverName"]

	if !ok {
		return nil, errors.New("invalid request data: server name is not specified")
	}

	cert, err := certificate.GetDomainCertificateFromHTTPRequest(serverName)

	if err != nil {
		return nil, err
	}

	return cert, nil
}
