package server

import (
	"github.com/r2dtools/agent/config"
	"github.com/spf13/cobra"
)

var GenerateTokenCmd = &cobra.Command{
	Use:   "generate-token",
	Short: "Generate new token",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetConfig()

		if err != nil {
			return err
		}

		return nil
	},
}
