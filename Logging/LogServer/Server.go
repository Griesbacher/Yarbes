package LogServer

import (
	"fmt"
)

//Server receives LogMessages and does something with them
type Server struct {
	LogQueue  chan LogMessage
	quit      chan bool
	isRunning bool
}

//NewLogServer constructs a new LogServer
func NewLogServer() *Server {
	return &Server{LogQueue: make(chan LogMessage, 100), quit: make(chan bool), isRunning: false}
}

//Start starts the LogServer
func (log Server) Start() {
	if !log.isRunning {
		go log.handleLog()
	}
}

//Stop stops the LogServer
func (log Server) Stop() {
	if log.isRunning {
		log.quit <- true
		<-log.quit
	}
}

//IsRunning returns true if the daemon is running
func (log Server) IsRunning() bool {
	return log.isRunning
}

func (log *Server) handleLog() {
	log.isRunning = true
	var message LogMessage
	for {
		select {
		case <-log.quit:
			log.quit <- true
			return
		case message = <-log.LogQueue:
			fmt.Println(message.String())
		}
	}
}
