package certificate

import (
	"crypto/x509"
	"errors"
	"net/http"
)

// GetDomainCertificateFromHTTPRequest retrieves certificate from http request to domain
func GetDomainCertificateFromHTTPRequest(domain string) ([]*x509.Certificate, error) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest("GET", "https://"+domain, nil)

	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)

	var hnErr x509.HostnameError

	if errors.As(err, &hnErr) {
		return []*x509.Certificate{hnErr.Certificate}, nil
	}

	if err != nil {
		return nil, err
	}

	if response.TLS == nil {
		return nil, nil
	}

	certificates := response.TLS.PeerCertificates

	return certificates, nil
}
