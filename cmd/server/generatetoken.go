package server

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/google/uuid"
)

var GenerateTokenCmd = &cobra.Command{
	Use:   "generate-token",
	Short: "Generate new token",
	RunE: func(cmd *cobra.Command, args []string) error {
		randomUuid, err := uuid.NewRandom()

		if err != nil {
			return err
		}

		token := randomUuid.String()

		err = os.Setenv("TOKEN", token)

		if err != nil {
			return err
		}

		fmt.Printf("Token: %s\n", token)
		fmt.Println("Please restart sslbot service")

		return nil
	},
}
