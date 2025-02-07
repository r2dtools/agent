package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agentintegration"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var HostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Show virtual hosts of web servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.GetConfig()

		if err != nil {
			return err
		}

		log, err := logger.NewLogger(conf)

		if err != nil {
			return err
		}

		supportedWebServerCodes := webserver.GetSupportedWebServers()
		webServerCodes := supportedWebServerCodes

		if webServerCode != "" {
			if !slices.Contains(supportedWebServerCodes, webServerCode) {
				return fmt.Errorf("invalid webserver %s", webServerCode)
			}

			webServerCodes = []string{webServerCode}
		}

		var vhosts []agentintegration.VirtualHost

		for _, webServerCode := range webServerCodes {
			webServer, err := webserver.GetWebServer(webServerCode, map[string]string{})

			if err != nil {
				log.Info(fmt.Sprintf("failed to get %s webserver: %v", webServerCode, err))

				continue
			}

			hosts, err := webServer.GetVhosts()

			if err != nil {
				return err
			}

			vhosts = append(vhosts, hosts...)
		}

		if isJson {
			output, err := json.Marshal(vhosts)

			if err != nil {
				return err
			}

			return writeOutput(cmd, string(output))
		}

		var outputParts []string

		for _, host := range vhosts {
			output, err := yaml.Marshal(host)

			if err != nil {
				return err
			}

			outputParts = append(outputParts, string(output))
		}

		return writeOutput(cmd, strings.Join(outputParts, "\n"))
	},
}
