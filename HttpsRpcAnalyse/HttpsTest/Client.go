package HttpsTest

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

//Client to test https
func Client(loops int) {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}
	pem, err := ioutil.ReadFile("certs/ca.crt")
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
	tr := &http.Transport{
		TLSClientConfig:    &config,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	data := "test string"
	for i := 0; i < loops; i++ {
		req, err := http.NewRequest("POST", "https://127.0.0.1:8090/", bytes.NewBuffer([]byte(data)))
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
	}
}
