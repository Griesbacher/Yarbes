package Incoming

import (
	"crypto/tls"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/TLS"
	"log"
	"net/rpc"
)

type RpcInterface struct {
	eventQueue chan Event.Event
	quit       chan bool
	isRunning  bool
}

func NewRpcInterface(eventQueue chan Event.Event) *RpcInterface {
	rpc := &RpcInterface{eventQueue: eventQueue, quit: make(chan bool), isRunning: false}
	return rpc
}

func (rpcI RpcInterface) Start() {
	if !rpcI.isRunning {
		rpcI.serve()
	}
}

func (rpcI RpcInterface) Stop() {
	rpcI.quit <- true
	<-rpcI.quit
}

func (rpcI RpcInterface) serve() {
	if err := rpc.Register(&RpcHandler{rpcI}); err != nil {
		panic(err)
	}
	config := TLS.GenerateServerTLSConfig(Config.GetServerConfig().RuleSystem.TLSCert, Config.GetServerConfig().RuleSystem.TLSKey, Config.GetServerConfig().RuleSystem.TLSCaCert)
	listenTo := Config.GetServerConfig().RuleSystem.RpcInterface
	listener, err := tls.Listen("tcp", listenTo, &config)
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

type RpcHandler struct {
	inter RpcInterface
}

type Result struct {
	Err error
}

func (handler *RpcHandler) CreateEvent(args *string, result *Result) error {
	event, err := Event.NewEvent([]byte(*args))
	if err == nil {
		handler.inter.eventQueue <- *event
	}
	result.Err = err
	return err
}