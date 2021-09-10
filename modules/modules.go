package modules

import (
	"github.com/r2dtools/agent/modules/certificates"
	"github.com/r2dtools/agent/modules/servermonitor"
	"github.com/r2dtools/agent/router"
)

// RegisterHandlers register modules handlers
func RegisterHandlers(router *router.Router) {
	router.RegisterHandler("certificates", &certificates.Handler{})
	router.RegisterHandler("servermonitor", &servermonitor.Handler{})
}
