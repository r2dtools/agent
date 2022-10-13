package server

import (
	"github.com/r2dtools/agent/config"
	"github.com/r2dtools/agent/internal/modules/certificates"
	"github.com/r2dtools/agent/internal/modules/servermonitor"
	serverMonitorService "github.com/r2dtools/agent/internal/modules/servermonitor/service"
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

		if err := registerSercices(&serviceManager, logger); err != nil {
			return err
		}

		certificatesHandler, err := certificates.GetHandler(logger)
		if err != nil {
			return err
		}

		servermonitorHandler := servermonitor.GetHandler(logger)

		router := router.Router{}
		router.RegisterHandler("main", &server.MainHandler{Logger: logger})
		router.RegisterHandler("certificates", certificatesHandler)
		router.RegisterHandler("servermonitor", servermonitorHandler)

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

func registerSercices(serviceManager *service.ServiceManager, logger logger.LoggerInterface) error {
	smService, err := serverMonitorService.GetStatCollectorService(logger)
	if err != nil {
		return err
	}

	scService, err := serverMonitorService.GetStatCleanerService(logger)
	if err != nil {
		return err
	}

	serviceManager.AddService("servermonitor.statcollector", smService)
	serviceManager.AddService("servermonitor.statcleaner", scService)

	return nil
}
