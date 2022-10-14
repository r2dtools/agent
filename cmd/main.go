package main

import "github.com/spf13/cobra"
import "github.com/r2dtools/agent/cmd/server"

var rootCmd = &cobra.Command{
	Use:   "r2dtools",
	Short: "R2D tools agent",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func main() {
	rootCmd.AddCommand(server.ServeCmd)
	rootCmd.AddCommand(server.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
