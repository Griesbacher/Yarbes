package TLS

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func GenerateServerTLSConfig(serverCrt, serverKey, caCert string) tls.Config {
	cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		panic("server: loadkeys")
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
	}
	config.Rand = rand.Reader

	return config
}

func GenerateClientTLSConfig(clientCrt, clientKey, caCert string) tls.Config {
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
	}
	return config
}
