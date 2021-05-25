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
	fmt.Print("Listening for https://localhost:8443/hello\n")
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
		log.Fatalf("Failed to start server with error: %v", err)
	}
}
