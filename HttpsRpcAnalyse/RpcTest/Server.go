package RPCTest

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
)

//Server to test rpc
func Server() {
	if err := rpc.Register(new(Foo)); err != nil {
		log.Fatal("Failed to register RPC method")
	}
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
	service := "127.0.0.1:8000"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	rpc.ServeConn(conn)
}

//Result struct
type Result struct {
	Data string
}

//Foo handler
type Foo bool

//Bar handlerfunction
func (f *Foo) Bar(args *string, res *Result) error {
	res.Data = *args
	return nil
}
