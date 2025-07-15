package commondir

import (
	"fmt"

	"github.com/r2dtools/sslbot/internal/pkg/logger"
	"github.com/r2dtools/sslbot/internal/pkg/webserver"
	"github.com/r2dtools/sslbot/internal/pkg/webserver/reverter"
)

type CommonDir struct {
	Enabled bool
	Root    string
}

type CommonDirManager interface {
	EnableCommonDir(serverName string) error
	DisableCommonDir(serverName string) error
	GetCommonDirStatus(serverName string) CommonDir
}

func GetCommonDirManager(webServer webserver.WebServer, reverter *reverter.Reverter, logger logger.Logger, options map[string]string) (CommonDirManager, error) {
	nginxCommonDir := options["nginx_acme_common_dir"]

	switch w := webServer.(type) {
	case *webserver.NginxWebServer:
		return &NginxCommonDirManager{logger: logger, webServer: w, reverter: reverter, commonDir: nginxCommonDir}, nil
	default:
		return nil, fmt.Errorf("could not create common directory manager: webserver '%s' is not supported", webServer.GetCode())
	}
}
