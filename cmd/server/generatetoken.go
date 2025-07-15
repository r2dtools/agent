package server

import (
	"fmt"
	"os"

	"github.com/r2dtools/sslbot/config"
	"github.com/spf13/cobra"

	"github.com/google/uuid"
)

var GenerateTokenCmd = &cobra.Command{
	Use:   "generate-token",
	Short: "Generate new token",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.GetConfig()

		if err != nil {
			return err
		}

		randomUuid, err := uuid.NewRandom()

		if err != nil {
			return err
		}

		token := randomUuid.String()
		tokenPath := conf.GetTokenPath()

		tokenFile, err := os.Create(tokenPath)

		if err != nil {
			return err
		}

		defer tokenFile.Close()
		_, err = tokenFile.WriteString(token)

		if err != nil {
			return err
		}

		fmt.Printf("Token: %s\n", token)
		fmt.Println("Please restart the SSLBot service: systemctl restart sslbot.service")

		return nil
	},
}
