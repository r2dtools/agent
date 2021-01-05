package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/r2dtools/agentintegration"
)

var client http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// GetX509CertificateFromHTTPRequest retrieves certificate from http request to domain
func GetX509CertificateFromHTTPRequest(domain string) ([]*x509.Certificate, error) {
	request, err := http.NewRequest("GET", "https://"+domain, nil)

	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)

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

// GetCertificateForDomainFromHTTPRequest returns a certificate for a domain
func GetCertificateForDomainFromHTTPRequest(domain string) (*agentintegration.Certificate, error) {
	certs, err := GetX509CertificateFromHTTPRequest(domain)

	if err != nil {
		return nil, err
	}

	if len(certs) == 0 {
		return nil, nil
	}

	var roots []*x509.Certificate

	if len(certs) > 1 {
		roots = certs[1:]
	}

	return ConvertX509CertificateToIntCert(certs[0], roots), nil
}
