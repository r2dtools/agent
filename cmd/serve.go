package cmd

import (
	"github.com/spf13/cobra"
	"github.com/r2dtools/agent/server"
	"github.com/r2dtools/agent/config"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts TSP server",
	RunE: func(cmd *cobra.Command, args []string) error {
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
