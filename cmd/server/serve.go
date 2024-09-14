package server

import (
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/router"
	"github.com/r2dtools/agent/internal/server"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts TCP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.GetConfig()

		if err != nil {
			return err
		}

		logger, err := logger.NewLogger(config)

		if err != nil {
			return err
		}

		certificatesHandler, err := certificates.GetHandler(config, logger)

		if err != nil {
			return err
		}

		router := router.Router{}
		router.RegisterHandler("main", &server.MainHandler{
			Config: config,
			Logger: logger,
		})
		router.RegisterHandler("certificates", certificatesHandler)

		server := &server.Server{
			Port:   config.Port,
			Router: router,
			Logger: logger,
			Config: config,
		}

		if err := server.Serve(); err != nil {
			return err
		}

		return nil
	},
}
