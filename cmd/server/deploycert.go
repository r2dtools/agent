package server

import (
	"fmt"
	"path/filepath"

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
			return fmt.Errorf("invalid webserver %s", webServerCode)
		}

		webServer, err := webserver.GetWebServer(webServerCode, map[string]string{})

		if err != nil {
			return err
		}

		processManager, err := webServer.GetProcessManager()

		if err != nil {
			return err
		}

		vhost, err := webServer.GetVhostByName(serverName)

		if err != nil {
			return err
		}

		webServerReverter := &reverter.Reverter{
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

		sslConfigFilePath, originEnabledConfigFilePath, err := deployer.DeployCertificate(vhost, certPath, certKeyPath)

		if err != nil {
			if rErr := webServerReverter.Rollback(); rErr != nil {
				log.Error(fmt.Sprintf("failed to rallback webserver configuration on cert deploy: %v", rErr))
			}

			return err
		}

		if err = webServer.GetVhostManager().Enable(sslConfigFilePath, filepath.Dir(originEnabledConfigFilePath)); err != nil {
			if rErr := webServerReverter.Rollback(); rErr != nil {
				log.Error(fmt.Sprintf("failed to rallback webserver configuration on host enabling: %v", rErr))
			}

			return err
		}

		if err = processManager.Reload(); err != nil {
			if rErr := webServerReverter.Rollback(); rErr != nil {
				log.Error(fmt.Sprintf("failed to rallback webserver configuration on webserver reload: %v", rErr))
			}

			return err
		}

		if err = webServerReverter.Commit(); err != nil {
			if rErr := webServerReverter.Rollback(); rErr != nil {
				log.Error(fmt.Sprintf("failed to commit webserver configuration: %v", rErr))
			}
		}

		return nil
	},
}

var certPath string
var certKeyPath string

func init() {
	DeployCertificateCmd.PersistentFlags().StringVarP(&serverName, "domain", "d", "", "domain to deploy a certificate")
	DeployCertificateCmd.PersistentFlags().StringVarP(&certPath, "cert", "c", "", "path to a certificate file")
	DeployCertificateCmd.PersistentFlags().StringVarP(&certKeyPath, "key", "k", "", "path to a certificate key path")
}
