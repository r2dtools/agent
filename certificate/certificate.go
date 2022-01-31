package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/r2dtools/agentintegration"
	"github.com/unknwon/com"
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

// GetCertificateFromFile read and parse certificate from file
func GetCertificateFromFile(path string) (*agentintegration.Certificate, error) {
	if !com.IsFile(path) {
		return nil, fmt.Errorf("certificate file '%s' does not exists", path)
	}
	certContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate content: %v", err)
	}

	var bCert []byte
	for {
		block, rest := pem.Decode([]byte(certContent))
		if block == nil {
			break
		}
		// get first certificate in the chain
		if block.Type == "CERTIFICATE" {
			bCert = block.Bytes
			break
		}
		certContent = rest
	}

	if len(bCert) == 0 {
		return nil, errors.New("could not parse certificate")
	}

	x509Cert, err := x509.ParseCertificate(bCert)
	if err != nil {
		return nil, err
	}
	cert := ConvertX509CertificateToIntCert(x509Cert, nil)

	return cert, nil
}
