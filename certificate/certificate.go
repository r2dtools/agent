package certificate

import (
	"crypto/x509"
	"errors"
	"net/http"
	"time"

	"github.com/r2dtools/agentintegration"
)

// GetX509CertificateFromHTTPRequest retrieves certificate from http request to domain
func GetX509CertificateFromHTTPRequest(domain string) ([]*x509.Certificate, error) {
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

// ConvertX509CertificateToIntCert converts x509 certificate to agentintegration.Certificate
func ConvertX509CertificateToIntCert(certificate *x509.Certificate, roots []*x509.Certificate) *agentintegration.Certificate {
	certPool := x509.NewCertPool()

	for _, root := range roots {
		certPool.AddCert(root)
	}

	opts := x509.VerifyOptions{
		Roots: certPool,
	}
	_, err := certificate.Verify(opts)
	isValid := err == nil
	cert := agentintegration.Certificate{
		DNSNames:       certificate.DNSNames,
		CN:             certificate.Subject.CommonName,
		EmailAddresses: certificate.EmailAddresses,
		Organization:   certificate.Subject.Organization,
		IsCA:           certificate.IsCA,
		ValidFrom:      certificate.NotBefore.Format(time.RFC822Z),
		ValidTo:        certificate.NotAfter.Format(time.RFC822Z),
		Issuer: agentintegration.Issuer{
			CN:           certificate.Issuer.CommonName,
			Organization: certificate.Issuer.Organization,
		},
		IsValid: isValid,
	}

	return &cert
}
