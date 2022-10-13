package certificates

import "fmt"

const (
	httpChallengeType = "http"
	dnsChallengeType  = "dns"
)

// ChallengeType should be implemented by object to specify challenge type during certificate issuence
type ChallengeType interface {
	GetParams() []string
}

// HTTPChallengeType implements http challenge type
type HTTPChallengeType struct {
	HTTPPort,
	TLSPort int
	WebRoot string
}

// DNSChallengeType implements http challenge type
type DNSChallengeType struct {
	Provider string
}

// GetParams retuns the list of parameters for http challenge type
func (ct *HTTPChallengeType) GetParams() []string {
	return []string{"--http", fmt.Sprintf("--http.port=%d", ct.HTTPPort), fmt.Sprintf("--tls.port=%d", ct.TLSPort), "--http.webroot=" + ct.WebRoot}
}

// GetParams retuns the list of parameters for dns challenge type
func (ct *DNSChallengeType) GetParams() []string {
	return []string{"--dns=" + ct.Provider}
}
