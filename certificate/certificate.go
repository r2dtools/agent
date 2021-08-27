package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"time"

	"github.com/r2dtools/agentintegration"
)

// GetX509CertificateFromRequest retrieves certificate from http request to domain
func GetX509CertificateFromRequest(domain string) ([]*x509.Certificate, error) {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Minute}, "tcp", domain+":443", &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
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

// GetCertificateForDomainFromRequest returns a certificate for a domain
func GetCertificateForDomainFromRequest(domain string) (*agentintegration.Certificate, error) {
	certs, err := GetX509CertificateFromRequest(domain)

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
