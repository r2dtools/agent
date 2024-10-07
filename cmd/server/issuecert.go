package server

import (
	"encoding/json"
	"fmt"

	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agentintegration"
	"github.com/spf13/cobra"
)

var IssueCertificateCmd = &cobra.Command{
	Use:   "issue-cert",
	Short: "Secure domain with a certificate",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.GetConfig()

		if err != nil {
			return err
		}

		log, err := logger.NewLogger(config)

		if err != nil {
			return err
		}

		if email == "" {
			return fmt.Errorf("email is not specified")
		}

		if serverName == "" {
			return fmt.Errorf("domain is not specified")
		}

		webServer, vhost, err := findWebServerHost(serverName, log)

		if err != nil {
			return err
		}

		certManager, err := certificates.GetCertificateManager(config, log)

		if err != nil {
			return err
		}

		certData := agentintegration.CertificateIssueRequestData{
			Email:         email,
			ServerName:    serverName,
			DocRoot:       vhost.DocRoot,
			WebServer:     webServer.GetCode(),
			ChallengeType: certificates.HttpChallengeTypeCode,
			Subjects:      aliases,
			Assign:        assign,
		}
		cert, err := certManager.Issue(certData)

		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(cert, "", " ")

		if err != nil {
			return err
		}

		fmt.Println(string(data))

		return nil
	},
}

var email string
var assign bool
var aliases []string

func init() {
	aliases = make([]string, 0)
	IssueCertificateCmd.PersistentFlags().StringVarP(&serverName, "domain", "d", "", "domain to secure")
	IssueCertificateCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "certificate email address")
	IssueCertificateCmd.PersistentFlags().BoolVarP(&assign, "assign", "s", true, "assignt certificate to the domain")
	IssueCertificateCmd.PersistentFlags().StringSliceVarP(&aliases, "alias", "a", nil, "domain aliases that need to be included in the certificate")
}
