package certificates

import (
	"path/filepath"
	"testing"

	"github.com/r2dtools/sslbot/config"
	"github.com/r2dtools/sslbot/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetCertificates(t *testing.T) {
	storage := getStorage()

	certs, err := storage.GetCertificates()
	assert.Nil(t, err)
	assert.Len(t, certs, 2)

	cert, ok := certs["example.com"]
	assert.True(t, ok)
	assert.Equal(t, "example.com", cert.CN)

	cert, ok = certs["example2.com"]
	assert.True(t, ok)
	assert.Equal(t, "example.com", cert.CN)
}

func TestGetCertificate(t *testing.T) {
	storage := getStorage()

	cert, err := storage.GetCertificate("example.com")
	assert.Nil(t, err)

	assert.Equal(t, "example.com", cert.CN)
}

func TestAddRemoveCertificate(t *testing.T) {
	storage := getStorage()

	certPath, data, err := storage.GetCertificateAsString("example2.com")
	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(storage.path, "example2.com.pem"), certPath)

	certPath, err = storage.AddPemCertificate("example3.com", data)
	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(storage.path, "example3.com.pem"), certPath)

	cert, err := storage.GetCertificate("example3.com")
	assert.Nil(t, err)
	assert.Equal(t, "example.com", cert.CN)

	err = storage.RemoveCertificate("example3.com")
	assert.Nil(t, err)

	certs, err := storage.GetCertificates()
	assert.Nil(t, err)

	_, ok := certs["example3.com"]
	assert.False(t, ok)
}

func getStorage() Storage {
	config := config.Config{}

	return Storage{
		path:   config.GetPathInsideVarDir("certificates"),
		logger: &logger.NilLogger{},
	}
}
