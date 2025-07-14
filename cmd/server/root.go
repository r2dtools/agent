package server

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "sslbot",
	Short: "R2DTools SSLBot",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var isJson bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&isJson, "json", "j", false, "show result in json format")
	RootCmd.AddCommand(ServeCmd)
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(HostsCmd)
	RootCmd.AddCommand(DeployCertificateCmd)
	RootCmd.AddCommand(IssueCertificateCmd)
	RootCmd.AddCommand(GenerateTokenCmd)
	RootCmd.AddCommand(CommonDirCmd)
	RootCmd.AddCommand(ShowTokenCmd)
}

func writeOutput(cmd *cobra.Command, output string) error {
	_, err := cmd.OutOrStdout().Write([]byte(output))

	if err != nil {
		return err
	}

	return nil
}

var webServerCode string
var serverName string

func init() {
	RootCmd.PersistentFlags().StringVarP(&webServerCode, "webserver", "w", "", "webserver (nginx|apache)")
}
