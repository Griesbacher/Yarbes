package RpcTest

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"reflect"
)

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
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		/*		tlsconn,ok := conn.(*tls.Conn)
				if ok{
					buf := make([]byte, 512)
					conn.Read(buf)
					fmt.Println(tlsconn.ConnectionState().PeerCertificates[0].Subject)
				}
		*/go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	rpc.ServeConn(conn)
	fmt.Println("---")
	fmt.Println(reflect.TypeOf(conn))
	tlsconn, ok := conn.(*tls.Conn)
	if ok {
		//			buf := make([]byte, 512)
		//			conn.Read(buf)
		fmt.Println(tlsconn.ConnectionState())
	}
	log.Println("server: conn: closed")
}

type Result struct {
	Data string
}

type Foo bool

func (f *Foo) Bar(args *string, res *Result) error {
	res.Data = *args
	log.Printf("Received %q, send %s", *args, res.Data)
	//return fmt.Errorf("Whoops, error happened")
	return nil
}
