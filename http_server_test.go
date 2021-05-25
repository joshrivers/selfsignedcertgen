package selfsignedcertgen

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHTTPSServerResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	keyLocations := NewSelfSigner().GenerateKeyAndCertificate()
	srv := http.Server{Addr: "localhost:8443", Handler: mux}
	go srv.ListenAndServeTLS(keyLocations.CertPath, keyLocations.KeyPath)
	defer srv.Close()
	time.Sleep(100 * time.Millisecond)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	res, err := client.Get("https://localhost:8443")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.StatusCode)
	}
}

func TestTLSCertificateHandshake(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	keyLocations := NewSelfSigner().GenerateKeyAndCertificate()
	srv := http.Server{Addr: "localhost:8444", Handler: mux}
	go srv.ListenAndServeTLS(keyLocations.CertPath, keyLocations.KeyPath)
	defer srv.Close()
	cfg := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "localhost:8444", &cfg)
	if err != nil {
		t.Fatalf("Server does not support TLS. Error: %v", err.Error())
	}
	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	if expiry.Before(time.Now().Add(364 * 24 * time.Hour)) {
		t.Errorf("Expiry too short: %v", expiry)
	}
	cn := conn.ConnectionState().PeerCertificates[0].Issuer.CommonName
	country := conn.ConnectionState().PeerCertificates[0].Issuer.Country[0]
	org := conn.ConnectionState().PeerCertificates[0].Issuer.Organization[0]
	if country != "JP" {
		t.Errorf("Issuer Country not set correctly. Value: %v", country)
	}
	if org != "Weyland Yutani" {
		t.Errorf("Issuer Organization not set correctly. Value: %v", org)
	}
	if cn != "www.example.com" {
		t.Errorf("Issuer Common Name not set correctly. Value: %v", cn)
	}
	t.Log(conn.ConnectionState().PeerCertificates[0].Issuer)
	time.Sleep(100 * time.Millisecond)
}
