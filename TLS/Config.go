package TLS

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

//GenerateServerTLSConfig generates a TLS Config which should be used by the server, min TLS Version 1.2
func GenerateServerTLSConfig(serverCrt, serverKey, caCert string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		panic(err)
	}

	pem, err := ioutil.ReadFile(caCert)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
		Rand:         rand.Reader,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			//tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,	//disable for circleci.com
			//tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	return &config
}

//GenerateClientTLSConfig generates a TLS Client Config, min TLS Version 1.2
func GenerateClientTLSConfig(clientCrt, clientKey, caCert string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		panic(err)
	}

	pem, err := ioutil.ReadFile(caCert)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
		Rand:         rand.Reader,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			//tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			//tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}
	config.BuildNameToCertificate()
	return &config
}

func findTlsCipherSuites() (result []uint16) {
	result = []uint16{}
	defer func() {
		if rec := recover(); rec != nil {
		}
	}()

	return result
}
