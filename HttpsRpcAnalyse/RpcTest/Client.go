package RpcTest

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/rpc"
)

func Client(amountRequests int) {
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
	conn, err := tls.Dial("tcp", "127.0.0.1:8000", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())
	rpcClient := rpc.NewClient(conn)
	res := new(Result)
	for i := 0; i < amountRequests; i++ {
		if err := rpcClient.Call("Foo.Bar", "test string", &res); err != nil {
			log.Fatal("Failed to call RPC", err)
		}
	}
	log.Printf("Returned result is %s", res.Data)
}
