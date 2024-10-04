package webserver

import (
	"fmt"

	"github.com/r2dtools/agentintegration"
)

const (
	WebServerNginxCode  = "nginx"
	WebServerApacheCode = "apache"
)

type HostManager interface {
	Enable(configFilePath, originConfigFilePath string) error
	Disable(configFilePath string) error
}

type ProcessManager interface {
	Reload() error
}

func GetSupportedWebServers() []string {
	return []string{WebServerNginxCode}
}

type WebServer interface {
	GetVhostByName(serverName string) (*agentintegration.VirtualHost, error)
	GetVhosts() ([]agentintegration.VirtualHost, error)
	GetCode() string
	GetVhostManager() HostManager
	GetProcessManager() (ProcessManager, error)
}

func GetWebServer(webServerCode string, options map[string]string) (WebServer, error) {
	var webServer WebServer
	var err error

	switch webServerCode {
	case WebServerNginxCode:
		webServer, err = GetNginxWebServer(options)
	default:
		err = fmt.Errorf("web server '%s' is not supported", webServerCode)
	}

	return webServer, err
}

func getVhostByName(vhosts []agentintegration.VirtualHost, serverName string) *agentintegration.VirtualHost {
	for _, vhost := range vhosts {
		if vhost.ServerName == serverName {
			return &vhost
		}
	}

	return nil
}
