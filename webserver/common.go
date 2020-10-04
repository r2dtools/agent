package webserver

import (
	"fmt"

	"github.com/r2dtools/agentintegration"
)

const (
	// WebServerNginxCode nginx web server code
	WebServerNginxCode = "nginx"
	// WebServerApacheCode apache web server code
	WebServerApacheCode = "apache"
)

// GetSupportedWebServers returns the codes of supported web servers
func GetSupportedWebServers() []string {
	return []string{WebServerApacheCode}
}

// WebServer is an interface for a webserver like nginx, apache
type WebServer interface {
	GetVhosts() ([]agentintegration.VirtualHost, error)
}

// GetWebServer returns WebServer object by code
func GetWebServer(webServerCode string, options map[string]string) (WebServer, error) {
	var webServer WebServer
	var err error

	switch webServerCode {
	case WebServerApacheCode:
		webServer, err = GetApacheWebServer(options)
	default:
		err = fmt.Errorf("web server '%s' is not supported", webServerCode)
	}

	return webServer, err
}
