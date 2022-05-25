package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
)

func setUpTLS(pkey, cert, ca string, mtls bool) (*tls.Config, error) {
	var c *tls.Config

	if cert != "" && pkey != "" {
		certs, err := tls.LoadX509KeyPair(cert, pkey)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to read client certificate or private key: %v", err))
		}
		c = &tls.Config{
			Certificates: []tls.Certificate{certs},
			ClientAuth:   tls.VerifyClientCertIfGiven,
		}
	}

	if mtls && (cert == "" || pkey == "") {
		return nil, errors.New("in order to enable MTLS a certificate and private key are required")
	}

	// if CA is not provided there is an automatic load of local CA in os
	if mtls && ca != "" {
		certPool := x509.NewCertPool()
		// Trusted store, contain the list of trusted certificates. client has to use one of this certificate to be trusted by this server
		bs, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("failed to read client ca cert: %v", err)
		}

		ok := certPool.AppendCertsFromPEM(bs)
		if !ok {
			return nil, fmt.Errorf("failed to append client certs: %v", err)
		}
		c.ClientCAs = certPool
	}

	return c, nil
}
