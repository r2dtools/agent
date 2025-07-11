package server

import (
	"fmt"
	"slices"

	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/modules/certificates/commondir"
	"github.com/r2dtools/sslbot/internal/pkg/logger"
	"github.com/r2dtools/sslbot/internal/pkg/webserver"
	"github.com/r2dtools/sslbot/internal/pkg/webserver/reverter"
	"github.com/spf13/cobra"
)

var CommonDirCmd = &cobra.Command{
	Use:   "common-dir",
	Short: "Manage ACME common directory for a host",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.GetConfig()

		if err != nil {
			return err
		}

		log, err := logger.NewLogger(config)

		if err != nil {
			return err
		}

		if serverName == "" {
			return fmt.Errorf("domain is not specified")
		}

		supportedWebServerCodes := webserver.GetSupportedWebServers()

		if webServerCode == "" {
			return fmt.Errorf("webserver is not specified")
		}

		if !slices.Contains(supportedWebServerCodes, webServerCode) {
			return fmt.Errorf("invalid webserver %s", webServerCode)
		}

		webServer, err := webserver.GetWebServer(webServerCode, map[string]string{})

		if err != nil {
			return err
		}

		webServerReverter := &reverter.Reverter{
			HostMng: webServer.GetVhostManager(),
			Logger:  log,
		}

		commonDirManager, err := commondir.GetCommonDirManager(webServer, webServerReverter, log, config.ToMap())

		if err != nil {
			return err
		}

		if enableCommonDir {
			err = commonDirManager.EnableCommonDir(serverName)
		} else if disableCommonDir {
			err = commonDirManager.DisableCommonDir(serverName)
		} else {
			fmt.Printf("Common directory status for host %s: %t\n", serverName, commonDirManager.GetCommonDirStatus(serverName).Enabled)

			return nil
		}

		return err
	},
}

var enableCommonDir bool
var disableCommonDir bool

func init() {
	CommonDirCmd.PersistentFlags().StringVarP(&serverName, "domain", "d", "", "domain to enable common directory")
	CommonDirCmd.PersistentFlags().BoolVar(&enableCommonDir, "enable", false, "enable common directory")
	CommonDirCmd.PersistentFlags().BoolVar(&disableCommonDir, "disable", false, "disable common directory")
}
