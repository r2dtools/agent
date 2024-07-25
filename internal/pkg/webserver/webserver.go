package webserver

import (
	"fmt"

	"github.com/r2dtools/agentintegration"
)

const (
	WebServerNginxCode  = "nginx"
	WebServerApacheCode = "apache"
)

// GetSupportedWebServers returns the codes of supported web servers
func GetSupportedWebServers() []string {
	return []string{WebServerApacheCode}
}

// WebServer is an interface for a webserver like nginx, apache
type WebServer interface {
	GetVhostByName(serverName string) (*agentintegration.VirtualHost, error)
	GetVhosts() ([]agentintegration.VirtualHost, error)
	GetCode() string
}

// GetWebServer returns WebServer object by code
func GetWebServer(webServerCode string, options map[string]string) (WebServer, error) {
	var webServer WebServer
	var err error

	switch webServerCode {
	case WebServerApacheCode:
		webServer, err = GetApacheWebServer(options)
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
