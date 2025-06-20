package crwauth

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Auth_mTls struct {
	addr string
}

func MakeAuth(addr string) *Auth_mTls {
	amtls := Auth_mTls{addr: addr}
	return &amtls
}

func (m *Auth_mTls) BeginAuth(servcert, servkey, ca string) error {
	// Load server certificate and key
	var ret_err string
	cert, err := tls.LoadX509KeyPair(servcert, servkey)
	if err != nil {
		log.Fatalf("failed to load server key pair: %v", err)
	}

	// Load CA cert to verify client certs
	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		ret_err = fmt.Sprintf("failed to read CA cert: %v", err)
		return errors.New(ret_err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS with client cert verification
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate() // deprecated in Go 1.14+, safe to omit if Go >=1.14

	server := &http.Server{
		Addr:      m.addr, // ":8443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Optionally access client certificate
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			clientCert := r.TLS.PeerCertificates[0]
			fmt.Fprintf(w, "Hello, %s!\n", clientCert.Subject.CommonName)
		} else {
			fmt.Fprint(w, "Hello, client!")
		}
	})
	err = server.ListenAndServeTLS("", "") // certs are provided via tls.Config
	if err != nil {
		ret_err = fmt.Sprintf("server failed: %v", err)
		return errors.New(ret_err)
	}
	return nil
}
