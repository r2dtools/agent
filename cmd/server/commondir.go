package server

import (
	"fmt"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates/commondir"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
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

		webServer, _, err := findWebServerHost(serverName, log)

		if err != nil {
			return err
		}

		processManager, err := webServer.GetProcessManager()

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
			fmt.Printf("Common directory status for host %s: %t\n", serverName, commonDirManager.IsCommonDirEnabled(serverName))

			return nil
		}

		if err != nil {
			if rErr := webServerReverter.Rollback(); rErr != nil {
				log.Error(fmt.Sprintf("failed to rallback webserver configuration on common directory switching: %v", rErr))
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

var enableCommonDir bool
var disableCommonDir bool

func init() {
	CommonDirCmd.PersistentFlags().StringVarP(&serverName, "domain", "d", "", "domain to enable common directory")
	CommonDirCmd.PersistentFlags().BoolVar(&enableCommonDir, "enable", false, "enable common directory")
	CommonDirCmd.PersistentFlags().BoolVar(&disableCommonDir, "disable", false, "disable common directory")
}
