package Incoming

import (
	"crypto/tls"
	"fmt"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/TLS"
	"io"
	"log"
	"net/rpc"
)

//RPCInterface represents a incoming RPC interface
type RPCInterface struct {
	quit        chan bool
	isRunning   bool
	RPCListenTo string
}

//NewRPCInterface creates a new RPCInterface
func NewRPCInterface(listenTo string) *RPCInterface {
	rpc := &RPCInterface{quit: make(chan bool), isRunning: false, RPCListenTo: listenTo}
	return rpc
}

//Start starts listening for requests
func (rpcI *RPCInterface) Start() {
	if !rpcI.isRunning {
		go rpcI.serve()
		rpcI.isRunning = true
	}
}

//Stop closes the port
func (rpcI RPCInterface) Stop() {
	//do nothing because rpc closes at program end automatically
}

//IsRunning returns true if the daemon is running
func (rpcI RPCInterface) IsRunning() bool {
	return rpcI.isRunning
}

func (rpcI *RPCInterface) serve() {
	config := TLS.GenerateServerTLSConfig(Config.GetServerConfig().TLS.Cert, Config.GetServerConfig().TLS.Key, Config.GetServerConfig().TLS.CaCert)
	listener, err := tls.Listen("tcp", rpcI.RPCListenTo, config)
	if err != nil {
		fmt.Println("listener")
		panic(err)
	}
	firstByte := make([]byte, 1)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		bytesRead, err := conn.Read(firstByte)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				continue
			}
			//TODO: durch log austauschen
			fmt.Println("first byte")
			panic(err)
		}
		if tlscon, ok := conn.(*tls.Conn); bytesRead == 1 && ok {
			state := tlscon.ConnectionState()
			sub := state.PeerCertificates[0].Subject
			if isCommonNameOnBlackList(sub.CommonName) || isDNSNameOnBlackList(state.PeerCertificates[0].DNSNames) {
				fmt.Println(sub.CommonName, " is blacklisted")
				conn.Close()
				continue
			}
		}
		go func() {
			defer conn.Close()
			rpc.ServeConn(conn)
		}()
	}
}

func (rpcI RPCInterface) publishHandler(rcvr interface{}) {
	if err := rpc.Register(rcvr); err != nil {
		fmt.Println("publish")
		panic(err)
	}
}

func isCommonNameOnBlackList(clientName string) bool {
	for _, name := range Config.GetServerConfig().TLS.BlackList {
		if clientName == name {
			return true
		}
	}
	return false
}

func isDNSNameOnBlackList(dnsNames []string) bool {
	for _, dnsName := range dnsNames {
		if isCommonNameOnBlackList(dnsName) {
			return true
		}
	}
	return false
}
