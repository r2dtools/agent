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
	Enable(configFilePath string) error
	Disable(configFilePath string) error
}

// GetSupportedWebServers returns the codes of supported web servers
func GetSupportedWebServers() []string {
	return []string{WebServerNginxCode}
}

// WebServer is an interface for a webserver like nginx, apache
type WebServer interface {
	GetVhostByName(serverName string) (*agentintegration.VirtualHost, error)
	GetVhosts() ([]agentintegration.VirtualHost, error)
	GetCode() string
	GetVhostManager() HostManager
}

// GetWebServer returns WebServer object by code
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
