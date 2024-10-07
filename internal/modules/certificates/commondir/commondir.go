package commondir

import (
	"fmt"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
)

type CommonDirManager interface {
	EnableCommonDir(serverName string) error
	DisableCommonDir(serverName string) error
	IsCommonDirEnabled(serverName string) bool
}

func GetCommonDirManager(webServer webserver.WebServer, reverter *reverter.Reverter, logger logger.Logger) (CommonDirManager, error) {
	switch w := webServer.(type) {
	case *webserver.NginxWebServer:
		return &NginxCommonDirManager{logger: logger, webServer: w, reverter: reverter}, nil
	default:
		return nil, fmt.Errorf("could not create common directory manager: webserver '%s' is not supported", webServer.GetCode())
	}
}
