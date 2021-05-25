package selfsignedcertgen

import "time"

// Struct containing parameters for key and certificate generation
type SelfSigner struct {
	Country          string
	Hosts            []string // first host will be used as CN, all hosts will be added as SAN
	Locality         string
	Organization     string
	OrganizationUnit string
	RsaKeyBits       int           // default: 2048
	ValidFor         time.Duration // duraation after ValidFrom the certificate will be valid
	ValidFrom        time.Time     // earliest date for certificate validity
}

// Create a new SelfSigner struct with dummy values
func NewSelfSigner() *SelfSigner {
	ss := new(SelfSigner)
	ss.Country = "JP"
	ss.Hosts = []string{"www.example.com"}
	ss.Locality = "Calpamos"
	ss.Organization = "Weyland Yutani"
	ss.OrganizationUnit = "Special Services"
	ss.RsaKeyBits = 2048
	ss.ValidFrom = time.Now()
	ss.ValidFor = 365 * 24 * time.Hour
	return ss
}
