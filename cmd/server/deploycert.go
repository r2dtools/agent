package server

import (
	"fmt"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates/deploy"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var DeployCertificateCmd = &cobra.Command{
	Use:   "deploy-cert",
	Short: "Deploy certificate to a domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.GetConfig()

		if err != nil {
			return err
		}

		log, err := logger.NewLogger(config)

		if err != nil {
			return err
		}

		supportedWebServerCodes := webserver.GetSupportedWebServers()

		if webServerCode == "" {
			return fmt.Errorf("webserver is not specified")
		}

		if serverName == "" {
			return fmt.Errorf("domain is not specified")
		}

		if !slices.Contains(supportedWebServerCodes, webServerCode) {
			return fmt.Errorf("invalid webserver code %s", webServerCode)
		}

		webServer, err := webserver.GetWebServer(webServerCode, map[string]string{})

		if err != nil {
			return err
		}

		vhost, err := webServer.GetVhostByName(serverName)

		if err != nil {
			return err
		}

		webServerReverter := reverter.Reverter{
			HostMng: webServer.GetVhostManager(),
			Logger:  log,
		}

		if vhost == nil {
			return fmt.Errorf("could not find virtual host '%s'", serverName)
		}

		deployer, err := deploy.GetCertificateDeployer(webServer, webServerReverter, log)

		if err != nil {
			return err
		}

		return deployer.DeployCertificate(vhost, "/path/to/cert", "/path/to/cert-key", "", "")
	},
}

var serverName string

func init() {
	DeployCertificateCmd.PersistentFlags().StringVarP(&serverName, "domain", "s", "", "domain to deploy a certificate")
}
