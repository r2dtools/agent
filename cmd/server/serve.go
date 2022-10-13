package server

import (
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules"
	"github.com/r2dtools/agent/internal/pkg/service"
	"github.com/r2dtools/agent/internal/server"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agent/pkg/router"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts TCP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.GetConfig()

		logger, err := logger.NewLogger(config)
		if err != nil {
			return err
		}

		serviceManager := service.ServiceManager{
			Logger: logger,
		}
		if err := modules.RegisterSercices(&serviceManager, logger); err != nil {
			return err
		}

		router := router.Router{}
		router.RegisterHandler("main", &server.MainHandler{Logger: logger})
		modules.RegisterHandlers(&router, logger)

		server := &server.Server{
			Port:           config.Port,
			ServiceManager: serviceManager,
			Router:         router,
			Logger:         logger,
		}

		if err := server.Serve(); err != nil {
			return err
		}

		return nil
	},
}
