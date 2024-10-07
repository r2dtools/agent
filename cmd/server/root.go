package server

import (
	"fmt"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agentintegration"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "r2dtools",
	Short: "R2DTools agent",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var isJson bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&isJson, "json", "j", false, "show result in json format")
	RootCmd.AddCommand(ServeCmd)
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(HostsCmd)
	RootCmd.AddCommand(DeployCertificateCmd)
	RootCmd.AddCommand(IssueCertificateCmd)
	RootCmd.AddCommand(GenerateTokenCmd)
	RootCmd.AddCommand(CommonDirCmd)
}

func writeOutput(cmd *cobra.Command, output string) error {
	_, err := cmd.OutOrStdout().Write([]byte(output))

	if err != nil {
		return err
	}

	return nil
}

var webServerCode string
var serverName string

func init() {
	RootCmd.PersistentFlags().StringVarP(&webServerCode, "webserver", "w", "", "webserver code (nginx|apache)")
}

func findWebServerHost(serverName string, log logger.Logger) (webserver.WebServer, *agentintegration.VirtualHost, error) {
	supportedWebServerCodes := webserver.GetSupportedWebServers()

	for _, webServerCode := range supportedWebServerCodes {
		webServer, err := webserver.GetWebServer(webServerCode, map[string]string{})

		if err != nil {
			log.Error("failed to get webserver %s", webServerCode)

			continue
		}

		vhost, err := webServer.GetVhostByName(serverName)

		if err != nil {
			log.Error("failed to get webserver %s host %s", webServerCode, serverName)

			continue
		}

		if vhost != nil {
			return webServer, vhost, nil
		}
	}

	return nil, nil, fmt.Errorf("could not find virtual host '%s'", serverName)
}
