# selfsignedcertgen

[![codecov]([https://codecov.io/gh/joshrivers/selfsignedcertgen/branch/main/graph/badge.svg?token=81X65EV65H](https://codecov.io/gh/joshrivers/selfsignedcertgen/branch/main/graph/badge.svg?token=81X65EV65H))]([https://codecov.io/gh/joshrivers/selfsignedcertgen](https://codecov.io/gh/joshrivers/selfsignedcertgen))
[![GoDoc]([https://img.shields.io/badge/pkg.go.dev-doc-blue](https://img.shields.io/badge/pkg.go.dev-doc-blue))]([http://pkg.go.dev/github.com/joshrivers/selfsignedcertgen](http://pkg.go.dev/github.com/joshrivers/selfsignedcertgen))

Generate self signed certificates for use in a Golang web server.

You can generate a new TLS private key and sign it with a self-signed certificate authority with a simple one-liner:

```go
openssl req -x509 -newkey rsa:4096 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/C=US/ST=OR/L=Portland/O=test/OU=example/CN=www.example.com"
```

Unfortunately this requires openssl to be installed and some pre-launch execution for a containerized application.

This library will generate the equivalent key and certificate files in Go at runtime to allow a TLS server to start. It is implemented using RSA keys for simplicity and current compatibility requirements.

## Note on security

If possible, you should never use self-signed certificates. These days it is pretty easy to use Let's Encrypt and there are a number of Go libraries that will autoprovision an TLS certificate using the Let's Encrypt api. If possible you should use one of those and not self-signed certificates. Self-signed certificates often create a number of security problems that leave you open to a lower level of security than using plain old HTTP. Avoid them if possible. Remember that you need `chrome://flags/#allow-insecure-localhost` enabled to hit insecure HTTPS on localhost now.

Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joshrivers/selfsignedcertgen"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func main() {
	fmt.Print("Listening for [https://localhost:8443/hello](https://localhost:8443/hello)\n")
	http.HandleFunc("/hello", hello)

	signer := selfsignedcertgen.NewSelfSigner()
	signer.Hosts = []string{"www.example.net", "replica.example.net"}
	signer.Country = "US"
	signer.Locality = "Los Angeles"
	signer.Organization = "Tyrell Corporation"
	signer.OrganizationUnit = "Nexus Design"
	signer.RsaKeyBits = 4096
	signer.ValidFor = 10 * 365 * 24 * time.Hour
	keyLocations := signer.GenerateKeyAndCertificate()
	err := http.ListenAndServeTLS("localhost:8443", keyLocations.CertPath, keyLocations.KeyPath, nil)
	if err != nil {
		log.Fatalf("Failed to start server with error: %!(NOVERB)v", err)
	}
}
```
