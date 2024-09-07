package deploy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/agentintegration"
	nginxConfig "github.com/r2dtools/gonginx/config"
)

// NginxCertificateDeployer certificate deployer to apache virtual host
type NginxCertificateDeployer struct {
	logger    logger.Logger
	webServer *webserver.NginxWebServer
	reverter  reverter.Reverter
}

// DeployCertificate deploys certificate to nginx domain
func (d *NginxCertificateDeployer) DeployCertificate(vhost *agentintegration.VirtualHost, certPath, certKeyPath, chainPath, fullChainPath string) error {
	wConfig := d.webServer.Config
	serverBlocks := wConfig.FindServerBlocksByServerName(vhost.ServerName)

	if len(serverBlocks) == 0 {
		return fmt.Errorf("nginx host %s does not exixst", vhost.ServerName)
	}

	var sslServerBlock *nginxConfig.ServerBlock
	var err error
	serverBlock := serverBlocks[0]

	for _, serverBlock := range serverBlocks {
		if serverBlock.HasSSL() {
			sslServerBlock = &serverBlock
		}
	}

	if sslServerBlock == nil {
		sslServerBlock, err = d.createSslhost(vhost, serverBlock)

		if err != nil {
			return err
		}
	}

	certKeyDirective := nginxConfig.NewDirective("ssl_certificate_key", []string{certKeyPath})
	sslServerBlock.AddDirective(certKeyDirective, false)

	certDirective := nginxConfig.NewDirective("ssl_certificate", []string{certPath})
	sslServerBlock.AddDirective(certDirective, false)

	sslServerBlockFileName := filepath.Base(sslServerBlock.FilePath)
	configFile := wConfig.GetConfigFile(sslServerBlockFileName)

	err = configFile.Dump()

	if err != nil {
		return err
	}

	return nil
}

func (d *NginxCertificateDeployer) createSslhost(
	vhost *agentintegration.VirtualHost,
	serverBlock nginxConfig.ServerBlock,
) (*nginxConfig.ServerBlock, error) {
	content := serverBlock.Dump()

	filePath := filepath.Clean(serverBlock.FilePath)
	extension := filepath.Ext(filePath)
	fileName := strings.TrimSuffix(filepath.Base(filePath), extension)
	directory := filepath.Dir(filePath)

	sslFileName := fmt.Sprintf("%s-ssl.%s", fileName, extension)
	sslFilePath := filepath.Join(directory, sslFileName)

	// todo: implment reverter
	if _, err := os.Stat(sslFilePath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(sslFilePath)

		if err != nil {
			return nil, err
		}

		_, err = file.Write([]byte(content))

		if err != nil {
			return nil, err
		}

		err = d.webServer.Config.ParseFile(sslFilePath)

		if err != nil {
			return nil, err
		}

		configFile := d.webServer.Config.GetConfigFile(sslFileName)
		serverBlocks := configFile.FindServerBlocksByServerName(vhost.ServerName)

		if len(serverBlocks) == 0 {
			return nil, fmt.Errorf("nginx ssl host %s not found", vhost.ServerName)
		}

		serverBlock := serverBlocks[0]
		serverBlock.DeleteDirectiveByName("listen")

		if serverBlock.IsIpv4Enabled() {
			listenDirective := nginxConfig.NewDirective("listen", []string{"443", "ssl"})
			serverBlock.AddDirective(listenDirective, true)
		}

		if serverBlock.IsIpv6Enabled() {
			listenDirective := nginxConfig.NewDirective("listen", []string{"[::]:443", "ssl"})
			serverBlock.AddDirective(listenDirective, true)
		}

		return &serverBlock, nil
	}

	return nil, fmt.Errorf("config file already exists %s", filePath)
}
