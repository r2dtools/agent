package cmd

import (
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/server"
	"github.com/r2dtools/agent/system"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts TCP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		system.GetPrivilege().Init()
		system.GetPrivilege().DropPrivilege()
		config := config.GetConfig()
		server := &server.Server{
			Port: config.Port,
		}
		err := server.Serve()

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
