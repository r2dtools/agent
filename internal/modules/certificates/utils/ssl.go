package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type certificateData struct {
	Certificate [][]byte
	PrivateKey  []byte
	Request     []byte
}

// LoadCertficateAndKeyFromFile reads file, divides into key and certificates
func LoadCertficateAndKeyFromPem(certPem string) (*certificateData, error) {
	raw := []byte(certPem)
	var err error
	var cert certificateData

	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, pem.EncodeToMemory(block))
		} else if block.Type == "CERTIFICATE REQUEST" {
			cert.Request = pem.EncodeToMemory(block)
		} else {
			cert.PrivateKey = pem.EncodeToMemory(block)
			_, err = parsePrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("could not read private key: %s", err)
			}
		}
		raw = rest
	}

	if len(cert.Certificate) == 0 {
		return nil, fmt.Errorf("no certificate found")
	} else if cert.PrivateKey == nil {
		return nil, fmt.Errorf("no private key found")
	}

	return &cert, nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("failed to parse private key")
}
