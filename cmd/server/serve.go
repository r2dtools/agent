package server

import (
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates"
	"github.com/r2dtools/agent/internal/modules/servermonitor"
	serverMonitorService "github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/router"
	"github.com/r2dtools/agent/internal/pkg/service"
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

		serviceManager := service.ServiceManager{
			Logger: logger,
		}

		if err := registerSercices(&serviceManager, config, logger); err != nil {
			return err
		}

		certificatesHandler, err := certificates.GetHandler(config, logger)

		if err != nil {
			return err
		}

		servermonitorHandler := servermonitor.GetHandler(config, logger)

		router := router.Router{}
		router.RegisterHandler("main", &server.MainHandler{
			Config: config,
			Logger: logger,
		})
		router.RegisterHandler("certificates", certificatesHandler)
		router.RegisterHandler("servermonitor", servermonitorHandler)

		server := &server.Server{
			Port:           config.Port,
			ServiceManager: serviceManager,
			Router:         router,
			Logger:         logger,
			Config:         config,
		}

		if err := server.Serve(); err != nil {
			return err
		}

		return nil
	},
}

func registerSercices(serviceManager *service.ServiceManager, config *config.Config, logger logger.Logger) error {
	smService, err := serverMonitorService.GetStatCollectorService(config, logger)

	if err != nil {
		return err
	}

	scService, err := serverMonitorService.GetStatCleanerService(config, logger)

	if err != nil {
		return err
	}

	serviceManager.AddService("servermonitor.statcollector", smService)
	serviceManager.AddService("servermonitor.statcleaner", scService)

	return nil
}
