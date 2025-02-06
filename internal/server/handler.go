package server

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/pkg/agent"
	"github.com/r2dtools/agent/internal/pkg/certificate"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/router"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agentintegration"
	"github.com/shirou/gopsutil/host"
)

// MainHandler handles common agent requests
type MainHandler struct {
	Config *config.Config
	Logger logger.Logger
}

// Handle handles request
func (h *MainHandler) Handle(request router.Request) (interface{}, error) {
	var response interface{}
	var err error

	switch action := request.GetAction(); action {
	case "refresh":
		response, err = h.refresh()
	case "getVhosts":
		response, err = h.getVhosts()
	case "getVhostCertificate":
		response, err = h.getVhostCertificate(request.Data)
	case "getvhostconfig":
		response, err = h.getVhostConfig(request.Data)
	default:
		response, err = nil, fmt.Errorf("invalid action '%s' for module '%s'", action, request.GetModule())
	}

	return response, err
}

func (h *MainHandler) refresh() (*agentintegration.ServerData, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("could not get system info: %v", err)
	}
	var serverData agentintegration.ServerData
	serverData.BootTime = info.BootTime
	serverData.Uptime = info.Uptime
	serverData.KernelArch = info.KernelArch
	serverData.KernelVersion = info.KernelVersion
	serverData.HostName = info.Hostname
	serverData.Platform = info.Platform
	serverData.PlatformFamily = info.PlatformFamily
	serverData.PlatformVersion = info.PlatformVersion
	serverData.Os = info.OS

	version, err := agent.GetAgentVersion(h.Config)
	if err != nil {
		return nil, err
	}
	serverData.AgentVersion = version

	return &serverData, nil
}

func (h *MainHandler) getVhosts() ([]agentintegration.VirtualHost, error) {
	webServerCodes := webserver.GetSupportedWebServers()
	var vhosts []agentintegration.VirtualHost
	options := h.Config.ToMap()

	for _, webServerCode := range webServerCodes {
		webserver, err := webserver.GetWebServer(webServerCode, options)

		if err != nil {
			h.Logger.Error(err.Error())
			continue
		}

		wVhosts, err := webserver.GetVhosts()

		if err != nil {
			h.Logger.Error(err.Error())
			continue
		}

		vhosts = append(vhosts, wVhosts...)
	}

	return vhosts, nil
}

func (h *MainHandler) getVhostCertificate(data interface{}) (*agentintegration.Certificate, error) {
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

	cert, err := certificate.GetCertificateForDomainFromRequest(vhostName)

	if err != nil {
		message := "could not get vhost '%s' certificate: %v"
		h.Logger.Info(message, vhostName, err)

		return nil, fmt.Errorf(message, vhostName, err)
	}

	return cert, nil
}

func (h *MainHandler) getVhostConfig(data interface{}) (agentintegration.VirtualHostConfigResponseData, error) {
	var response agentintegration.VirtualHostConfigResponseData
	var request agentintegration.VirtualHostConfigRequestData

	err := mapstructure.Decode(data, &request)

	if err != nil {
		return response, fmt.Errorf("invalid vhodt config request data: %v", err)
	}

	options := h.Config.ToMap()
	wServer, err := webserver.GetWebServer(request.WebServer, options)

	if err != nil {
		return response, err
	}

	vhost, err := wServer.GetVhostByName(request.ServerName)

	if err != nil {
		return response, err
	}

	if vhost == nil {
		return response, fmt.Errorf("vhost %s not found", request.ServerName)
	}

	configFile, err := os.Open(vhost.FilePath)

	if err != nil {
		return response, err
	}

	content, err := io.ReadAll(configFile)

	if err != nil {
		return response, err
	}

	response.Content = string(content)

	return response, nil
}
