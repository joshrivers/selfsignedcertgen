package selfsignedcertgen_test

import (
	"net/http"

	"github.com/joshrivers/selfsignedcertgen"
)

func ExampleNewSelfSigner() {
	keyLocations := selfsignedcertgen.NewSelfSigner().GenerateKeyAndCertificate()
	http.ListenAndServeTLS(":443", keyLocations.CertPath, keyLocations.KeyPath, nil)
}
