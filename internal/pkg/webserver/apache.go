package webserver

import (
	"fmt"

	"github.com/r2dtools/a2conf"
	"github.com/r2dtools/agentintegration"
)

// ApacheWebServer provides functionality to work with apache web server
type ApacheWebServer struct {
	configurator a2conf.ApacheConfigurator
	options      map[string]string
}

// GetApacheConfigurator returns webserver code
func (aws *ApacheWebServer) GetApacheConfigurator() a2conf.ApacheConfigurator {
	return aws.configurator
}

// GetCode returns webserver code
func (aws *ApacheWebServer) GetCode() string {
	return WebServerApacheCode
}

// GetVhostByName returns virtual host by name
func (aws *ApacheWebServer) GetVhostByName(serverName string) (*agentintegration.VirtualHost, error) {
	vhosts, err := aws.GetVhosts()

	if err != nil {
		return nil, err
	}

	return getVhostByName(vhosts, serverName), nil
}

// GetVhosts returns apache web server vitual hosts
func (aws *ApacheWebServer) GetVhosts() ([]agentintegration.VirtualHost, error) {
	var vhosts []agentintegration.VirtualHost
	aVhosts, err := aws.configurator.GetVhosts()

	if err != nil {
		return nil, fmt.Errorf("could not get apache virtual hosts %v", err)
	}

	for _, aVhost := range aVhosts {
		if !aVhost.Enabled || aVhost.ModMacro {
			continue
		}

		var addresses []agentintegration.VirtualHostAddress

		for _, address := range aVhost.Addresses {
			addresses = append(addresses, agentintegration.VirtualHostAddress{
				IsIpv6: address.IsIpv6,
				Host:   address.Host,
				Port:   address.Port,
			})
		}

		vhost := agentintegration.VirtualHost{
			FilePath:   aVhost.FilePath,
			ServerName: aVhost.ServerName,
			DocRoot:    aVhost.DocRoot,
			Aliases:    aVhost.Aliases,
			Ssl:        aVhost.Ssl,
			WebServer:  WebServerApacheCode,
			Addresses:  addresses,
		}
		vhosts = append(vhosts, vhost)
	}

	vhosts = filterVhosts(vhosts)
	vhosts = mergeVhosts(vhosts)

	return vhosts, nil
}

// GetApacheWebServer creates an instance of ApacheWebServer
func GetApacheWebServer(options map[string]string) (*ApacheWebServer, error) {
	configurator, err := a2conf.GetApacheConfigurator(options)

	if err != nil {
		return nil, fmt.Errorf("could not create apache configurator: %v", err)
	}

	return &ApacheWebServer{configurator, options}, nil
}
