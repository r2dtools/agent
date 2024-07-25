package webserver

import (
	"fmt"

	"github.com/r2dtools/agentintegration"
	nginxConfig "github.com/r2dtools/gonginx/config"
)

const (
	defaultNginxRoot = "/etc/nginx"
)

type NginxWebServer struct {
	config  *nginxConfig.Config
	options map[string]string
}

func (nws *NginxWebServer) GetCode() string {
	return WebServerNginxCode
}

func (nws *NginxWebServer) GetVhostByName(serverName string) (*agentintegration.VirtualHost, error) {
	vhosts, err := nws.GetVhosts()

	if err != nil {
		return nil, err
	}

	return getVhostByName(vhosts, serverName), nil
}

func (nws *NginxWebServer) GetVhosts() ([]agentintegration.VirtualHost, error) {
	var vhosts []agentintegration.VirtualHost

	nVhosts := nws.config.FindServerBlocks()

	for _, nVhost := range nVhosts {
		var addresses []agentintegration.VirtualHostAddress

		for _, address := range nVhost.GetAddresses() {
			addresses = append(addresses, agentintegration.VirtualHostAddress{
				IsIpv6: address.IsIpv6,
				Host:   address.Host,
				Port:   address.Port,
			})
		}

		serverNames := nVhost.GetServerNames()

		if len(serverNames) == 0 {
			continue
		}

		aliases := []string{}

		if len(serverNames) > 1 {
			aliases = serverNames[1:]
		}

		vhost := agentintegration.VirtualHost{
			FilePath:   "",
			ServerName: serverNames[0],
			DocRoot:    nVhost.GetDocumentRoot(),
			Aliases:    aliases,
			Ssl:        nVhost.HasSSL(),
			WebServer:  WebServerApacheCode,
			Addresses:  addresses,
		}
		vhosts = append(vhosts, vhost)
	}

	return vhosts, nil
}

func GetNginxWebServer(options map[string]string) (*NginxWebServer, error) {
	root, ok := options["NginxRoot"]

	if !ok {
		root = defaultNginxRoot
	}

	config, err := nginxConfig.GetConfig(root, "", false)

	if err != nil {
		return nil, fmt.Errorf("could not parse nginx config: %v", err)
	}

	return &NginxWebServer{config, options}, nil
}
