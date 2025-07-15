package server

import (
	"fmt"

	"github.com/r2dtools/sslbot/config"
	"github.com/spf13/cobra"
)

var ShowTokenCmd = &cobra.Command{
	Use:   "show-token",
	Short: "Show SSLBot current token",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.GetConfig()

		if err != nil {
			return err
		}

		token := conf.GetToken()

		if token != "" {
			fmt.Printf("Token: %s\n", token)

			return nil
		}

		fmt.Println("Token not generated")

		return nil
	},
}
