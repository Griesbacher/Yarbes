package HttpsTest

import (
	"net/http"
	"crypto/tls"

	"fmt"
	"log"
	"io/ioutil"
	"crypto/x509"
	"crypto/rand"
	"time"
)



func Server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(200)*time.Millisecond)
		robots, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(robots))
	})
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	config.Rand = rand.Reader
	server := &http.Server{
		Addr: ":8090",
		TLSConfig: &config,
	}


	server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
}