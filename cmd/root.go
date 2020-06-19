package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "r2dtools",
	Short: "R2 server tools agent",
	Run: func (cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute entry point for cli commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
