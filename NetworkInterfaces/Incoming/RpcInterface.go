package Incoming

import (
	"crypto/tls"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/TLS"
	"log"
	"net/rpc"
)

type RPCInterface struct {
	quit        chan bool
	isRunning   bool
	RPCListenTo string
}

func NewRPCInterface(listenTo string) *RPCInterface {
	rpc := &RPCInterface{quit: make(chan bool), isRunning: false, RPCListenTo: listenTo}
	return rpc
}

func (rpcI RPCInterface) Start() {
	if !rpcI.isRunning {
		go rpcI.serve()
	}
}

func (rpcI RPCInterface) Stop() {
	if rpcI.isRunning {
		rpcI.quit <- true
		<-rpcI.quit
	}
}

func (rpcI *RPCInterface) serve() {
	rpcI.isRunning = true
	config := TLS.GenerateServerTLSConfig(Config.GetServerConfig().TLS.Cert, Config.GetServerConfig().TLS.Key, Config.GetServerConfig().TLS.CaCert)
	listener, err := tls.Listen("tcp", rpcI.RPCListenTo, config)
	if err != nil {
		panic(err)
	}
	firstByte := make([]byte, 1)
	for {
		conn, err := listener.Accept()
		log.Printf("server: connection from %s", conn.RemoteAddr())
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		bytesRead, err := conn.Read(firstByte)
		if err != nil {
			panic(err)
		}
		if tlscon, ok := conn.(*tls.Conn); bytesRead == 1 && ok {
			state := tlscon.ConnectionState()
			sub := state.PeerCertificates[0].Subject
			log.Println(state)
			log.Println(sub)
		}
		go func() {
			log.Printf("server: accepted from %s", conn.RemoteAddr())
			defer conn.Close()
			rpc.ServeConn(conn)
		}()
	}
}

func (rpcI RPCInterface) publishHandler(rcvr interface{}) {
	if err := rpc.Register(rcvr); err != nil {
		panic(err)
	}
}
