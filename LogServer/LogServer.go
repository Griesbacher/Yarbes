package LogServer

import "fmt"

type LogServer struct {
	LogQueue  chan LogMessage
	quit      chan bool
	isRunning bool
}

func NewLogServer() *LogServer {
	return &LogServer{LogQueue: make(chan LogMessage, 100), quit: make(chan bool), isRunning: false}
}

func (log LogServer) Start() {
	if !log.isRunning {
		go log.handleLog()
	}
}

func (log LogServer) Stop() {
	if log.isRunning {
		log.quit <- true
		<-log.quit
	}
}

func (log *LogServer) handleLog() {
	log.isRunning = true
	var message LogMessage
	for {
		select {
		case <-log.quit:
			log.quit <- true
			return
		case message = <-log.LogQueue:
			fmt.Printf("[%s]@[%s] %s\n", message.Source, message.Timestamp, message.Message)
		}
	}
}
