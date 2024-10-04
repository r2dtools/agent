package deploy

import (
	"testing"

	"github.com/r2dtools/agent/internal/pkg/logger"
	"github.com/r2dtools/agent/internal/pkg/webserver"
	"github.com/r2dtools/agent/internal/pkg/webserver/reverter"
	"github.com/r2dtools/agentintegration"
	"github.com/stretchr/testify/assert"
)

func TestDeployCertificateToNonSslHost(t *testing.T) {
	deployer, nginxWebServer, rv := getNginxDeployer(t)
	defer rv.Rollback()

	hosts, err := nginxWebServer.GetVhosts()
	assert.Nilf(t, err, "get nginx hosts error: %v", err)

	servername := "example3.com"
	host := findHost(servername, hosts)
	assert.NotNilf(t, host, "host %s not found", servername)

	configPath, _, err := deployer.DeployCertificate(host, "test/certificate/example.com.key", "test/certificate/example.com.crt")
	assert.Nilf(t, err, "deploy certificate error: %v", err)
	assert.Equal(t, "/etc/nginx/sites-available/example3.com-ssl.conf", configPath)

	hosts, err = nginxWebServer.GetVhosts()
	assert.Nilf(t, err, "get nginx hosts after deploy error: %v", err)

	host = findHost("example3.com", hosts)
	assert.True(t, host.Ssl)
}

func TestDeployCertificateToSslHost(t *testing.T) {
	deployer, nginxWebServer, rv := getNginxDeployer(t)
	defer rv.Rollback()

	hosts, err := nginxWebServer.GetVhosts()
	assert.Nilf(t, err, "get nginx hosts error: %v", err)

	servername := "example2.com"
	host := findHost(servername, hosts)
	assert.NotNilf(t, host, "host %s not found", servername)
	assert.True(t, host.Ssl)

	configPath, _, err := deployer.DeployCertificate(host, "test/certificate/example2.com.key", "test/certificate/example2.com.crt")
	assert.Nilf(t, err, "deploy certificate error: %v", err)
	assert.Equal(t, "/etc/nginx/sites-enabled/example2.com.conf", configPath)

	hosts, err = nginxWebServer.GetVhosts()
	assert.Nilf(t, err, "get nginx hosts after deploy error: %v", err)

	host = findHost("example2.com", hosts)
	assert.True(t, host.Ssl)
}

func getNginxDeployer(t *testing.T) (CertificateDeployer, webserver.NginxWebServer, *reverter.Reverter) {
	nginxWebServer, err := webserver.GetNginxWebServer(nil)
	assert.Nil(t, err)

	rv := &reverter.Reverter{}
	deployer, err := GetCertificateDeployer(nginxWebServer, rv, &logger.NilLogger{})
	assert.Nil(t, err)

	return deployer, *nginxWebServer, rv
}

func findHost(servername string, hosts []agentintegration.VirtualHost) *agentintegration.VirtualHost {
	for _, host := range hosts {
		if host.ServerName == servername {
			return &host
		}
	}

	return nil
}
