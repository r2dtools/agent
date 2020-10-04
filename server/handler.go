package server

import (
	"fmt"

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
		response, err = getVhosts(request.Data)
	case "getVhosts":
		response, err = getVhosts(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func refresh(data interface{}) (interface{}, error) {
	var respData struct {
		OsCode,
		OsVersion,
		AgentVersion string
	}

	respData.AgentVersion = "1.0.0"
	respData.OsCode = "ubuntu"
	respData.OsVersion = "18.04"

	return respData, nil
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
