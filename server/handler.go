package server

import (
	"errors"
	"fmt"

	"github.com/r2dtools/agent/certificate"
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/logger"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agent/system"
	"github.com/r2dtools/agent/utils"
	"github.com/r2dtools/agent/webserver"
	"github.com/r2dtools/agentintegration"
	"github.com/shirou/gopsutil/host"
)

// MainHandler handles common agent requests
type MainHandler struct{}

// Handle handles request
func (h *MainHandler) Handle(request router.Request) (interface{}, error) {
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

	version, err := utils.GetAgentVersion()
	if err != nil {
		return nil, err
	}
	serverData.AgentVersion = version

	return &serverData, nil
}

func getVhosts(data interface{}) ([]agentintegration.VirtualHost, error) {
	webServerCodes := webserver.GetSupportedWebServers()
	var vhosts []agentintegration.VirtualHost
	options := config.GetConfig().ToMap()

	if err := system.GetPrivilege().IncreasePrivilege(); err != nil {
		logger.Error(fmt.Sprintf("getVhosts: increase privilege failed: %v", err))
	}

	defer (func() {
		if err := system.GetPrivilege().DropPrivilege(); err != nil {
			logger.Error(fmt.Sprintf("getVhosts: drop privilege failed: %v", err))
		}
	})()

	for _, webServerCode := range webServerCodes {
		webserver, err := webserver.GetWebServer(webServerCode, options)

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
		vhosts = utils.FilterVhosts(vhosts)
		vhosts = utils.MergeVhosts(vhosts)
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

	cert, err := certificate.GetCertificateForDomainFromRequest(vhostName)

	if err != nil {
		message := fmt.Sprintf("could not get vhost '%s' certificate: %v", vhostName, err)
		logger.Info(fmt.Sprintf(message))
		return nil, fmt.Errorf(message)
	}

	return cert, nil
}
