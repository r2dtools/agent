package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/r2dtools/agent/config"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates R2DTools agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		aConfig := config.GetConfig()
		scriptsPath := aConfig.GetScriptsDirAbsPath()
		command := exec.Command("bash", filepath.Join(scriptsPath, "update.sh"))
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()

		if err != nil {
			return err
		}

		aConfig.Merge()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
