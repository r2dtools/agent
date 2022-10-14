package server

import (
	"fmt"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/pkg/utils"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show agent version",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.GetConfig()

		version, err := utils.GetAgentVersion(config)
		if err != nil {
			return err
		}

		fmt.Println(version)

		return nil
	},
}
