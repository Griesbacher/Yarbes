package LogServer

import (
	"fmt"
)

type Server struct {
	LogQueue  chan LogMessage
	quit      chan bool
	isRunning bool
}

func NewLogServer() *Server {
	return &Server{LogQueue: make(chan LogMessage, 100), quit: make(chan bool), isRunning: false}
}

func (log Server) Start() {
	if !log.isRunning {
		go log.handleLog()
	}
}

func (log Server) Stop() {
	if log.isRunning {
		log.quit <- true
		<-log.quit
	}
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
			fmt.Printf("[%s]@[%s]-[%d] %s\n", message.Source, message.Timestamp, message.LogLevel, message.Message)
		}
	}
}
