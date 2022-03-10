package modules

import (
	"github.com/r2dtools/agent/modules/certificates"
	"github.com/r2dtools/agent/modules/servermonitor"
	serverMonitorService "github.com/r2dtools/agent/modules/servermonitor/service"
	"github.com/r2dtools/agent/router"
	"github.com/r2dtools/agent/service"
)

// RegisterHandlers register modules handlers
func RegisterHandlers(router *router.Router) {
	router.RegisterHandler("certificates", &certificates.Handler{})
	router.RegisterHandler("servermonitor", &servermonitor.Handler{})
}

func RegisterSercices(serviceManager *service.ServiceManager) error {
	smService, err := serverMonitorService.GetStatCollectorService()
	if err != nil {
		return err
	}
	scService, err := serverMonitorService.GetStatCleanerService()
	if err != nil {
		return err
	}

	serviceManager.AddService("servermonitor.statcollector", smService)
	serviceManager.AddService("servermonitor.statcleaner", scService)

	return nil
}
