package server

import (
	"fmt"
	"os"

	"github.com/r2dtools/sslbot/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ShowTokenCmd = &cobra.Command{
	Use:   "show-token",
	Short: "Show current sslbot token",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.GetConfig()

		if err != nil {
			return err
		}

		data, err := os.ReadFile(conf.ConfigFilePath)

		if err != nil {
			return err
		}

		confMap := make(map[string]any)
		err = yaml.Unmarshal(data, confMap)

		if err != nil {
			return err
		}

		token, ok := confMap["Token"]

		if ok {
			fmt.Printf("Token: %s\n", token)

			return nil
		}

		fmt.Println("Token not generated")

		return nil
	},
}
