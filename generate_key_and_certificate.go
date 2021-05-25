package selfsignedcertgen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
)

// Generates a key and certificate from SelfSigner spec and
// stores them in temporary files. Returns a struct containing
// the path to the generated temporary files.
func (ss SelfSigner) GenerateKeyAndCertificate() KeyLocations {
	notAfter := ss.ValidFrom.Add(ss.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{ss.Country},
			CommonName:         ss.Hosts[0],
			Locality:           []string{ss.Locality},
			Organization:       []string{ss.Organization},
			OrganizationalUnit: []string{ss.OrganizationUnit},
		},
		NotBefore:             ss.ValidFrom,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range ss.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	priv, err := rsa.GenerateKey(rand.Reader, ss.RsaKeyBits)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}
	certOut, err := ioutil.TempFile("", "cert.*.pem")
	if err != nil {
		log.Fatalf("Failed to open certificate file for writing: %v", err)
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		log.Fatalf("Failed to write data to certificate file: %v", err)
	}
	err = certOut.Close()
	if err != nil {
		log.Fatalf("Error closing certificate file: %v", err)
	}

	keyOut, err := ioutil.TempFile("", "key.*.pem")
	if err != nil {
		log.Fatalf("Failed to open key file: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to format private key: %v", err)
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		log.Fatalf("Unable to write data to key file: %v", err)
	}
	err = keyOut.Close()
	if err != nil {
		log.Fatalf("Error closing key file: %v", err)
	}
	err = os.Chmod(keyOut.Name(), 0600)
	if err != nil {
		log.Fatalf("Unable to restrict key file premissions: %v", err)
	}

	return KeyLocations{KeyPath: keyOut.Name(), CertPath: certOut.Name()}
}
