package webserver

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agentintegration"
	nginxConfig "github.com/r2dtools/gonginxconf/config"
	"github.com/r2dtools/sslbot/internal/pkg/certificate"
	"github.com/r2dtools/sslbot/internal/pkg/webserver/hostmng"
	"github.com/r2dtools/sslbot/internal/pkg/webserver/processmng"
)

const (
	defaultNginxRoot      = "/etc/nginx"
	NginxCertKeyDirective = "ssl_certificate_key"
	NginxCertDirective    = "ssl_certificate"
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
			FilePath:    strings.Trim(nVhost.FilePath, "\""),
			ServerName:  strings.Trim(serverNames[0], "\""),
			DocRoot:     strings.Trim(nVhost.GetDocumentRoot(), "\""),
			Aliases:     aliases,
			Ssl:         nVhost.HasSSL(),
			WebServer:   WebServerNginxCode,
			Addresses:   addresses,
			Certificate: getCertificate(nVhost),
		}
		vhosts = append(vhosts, vhost)
	}

	vhosts = filterVhosts(vhosts)
	vhosts = mergeVhosts(vhosts)

	return vhosts, nil
}

func (nws *NginxWebServer) GetVhostManager() HostManager {
	return &hostmng.NginxHostManager{
		EnabledConfigRootPath: filepath.Join(nws.root, "sites-enabled"),
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

func getCertificate(serverBlock nginxConfig.ServerBlock) *agentintegration.Certificate {
	certDirectives := serverBlock.FindDirectives(NginxCertDirective)

	if len(certDirectives) == 0 {
		return nil
	}

	certDirective := certDirectives[len(certDirectives)-1]
	cert, _ := certificate.GetCertificateFromFile(certDirective.GetFirstValue())

	return cert
}
