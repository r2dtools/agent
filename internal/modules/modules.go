package modules

import (
	"github.com/r2dtools/agent/internal/modules/certificates"
	"github.com/r2dtools/agent/internal/modules/servermonitor"
	serverMonitorService "github.com/r2dtools/agent/internal/modules/servermonitor/service"
	"github.com/r2dtools/agent/internal/pkg/service"
	"github.com/r2dtools/agent/pkg/logger"
	"github.com/r2dtools/agent/pkg/router"
)

func RegisterHandlers(router *router.Router, logger logger.LoggerInterface) {
	router.RegisterHandler("certificates", &certificates.Handler{Logger: logger})
	router.RegisterHandler("servermonitor", &servermonitor.Handler{Logger: logger})
}

func RegisterSercices(serviceManager *service.ServiceManager, logger logger.LoggerInterface) error {
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
