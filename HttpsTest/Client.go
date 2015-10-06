package HttpsTest
import (
	"crypto/tls"
	"log"
	"io/ioutil"
	"crypto/x509"
	"net/http"
	"fmt"
	"bytes"
	"time"
)

func Client() *http.Client {
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
	return &http.Client{Transport: tr}

}

func Request(client *http.Client, data string) {
	req, err := http.NewRequest("POST", "https://127.0.0.1:8090/", bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	robots, err := ioutil.ReadAll(resp.Body)
	time.Sleep(time.Duration(400)*time.Millisecond)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("<- %s", robots)
}