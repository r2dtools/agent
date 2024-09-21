package webserver

import (
	"fmt"
	"path/filepath"

	"github.com/r2dtools/agent/internal/pkg/webserver/hostmng"
	"github.com/r2dtools/agent/internal/pkg/webserver/processmng"
	"github.com/r2dtools/agentintegration"
	nginxConfig "github.com/r2dtools/gonginx/config"
)

const (
	defaultNginxRoot = "/etc/nginx"
)

type NginxWebServer struct {
	Config  *nginxConfig.Config
	root    string
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

	nVhosts := nws.Config.FindServerBlocks()

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
			FilePath:   nVhost.FilePath,
			ServerName: serverNames[0],
			DocRoot:    nVhost.GetDocumentRoot(),
			Aliases:    aliases,
			Ssl:        nVhost.HasSSL(),
			WebServer:  WebServerNginxCode,
			Addresses:  addresses,
		}
		vhosts = append(vhosts, vhost)
	}

	vhosts = filterVhosts(vhosts)
	vhosts = mergeVhosts(vhosts)

	return vhosts, nil
}

func (nws *NginxWebServer) GetVhostManager() HostManager {
	return &hostmng.NginxHostManager{
		AvailableConfigRootPath: filepath.Join(nws.root, "sites-available"),
		EnabledConfigRootPath:   filepath.Join(nws.root, "sites-enabled"),
	}
}

func (nws *NginxWebServer) GetProcessManager() (ProcessManager, error) {
	return processmng.GetNginxProcessManager()
}

func GetNginxWebServer(options map[string]string) (*NginxWebServer, error) {
	root := getNginxRoot(options)
	config, err := nginxConfig.GetConfig(root, "", false)

	if err != nil {
		return nil, fmt.Errorf("could not parse nginx config: %v", err)
	}

	return &NginxWebServer{
		Config:  config,
		root:    root,
		options: options,
	}, nil
}

func getNginxRoot(options map[string]string) string {
	root, ok := options["NginxRoot"]

	if !ok {
		root = defaultNginxRoot
	}

	return root
}
