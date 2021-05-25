package selfsignedcertgen

// Return value of the GenerateKeyAndCertificate function
// containing the path to the generated key and certificate
// files.
type KeyLocations struct {
	KeyPath  string // path to the generated pem-encoded RSA key
	CertPath string // path to the self-signed x509 certificate
}
