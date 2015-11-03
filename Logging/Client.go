package Logging

import (
	"fmt"
	"github.com/griesbacher/Yarbes/Logging/Local"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/Outgoing"
	"github.com/kdar/factorlog"
	"os"
	"runtime"
)

//Client combines locallogging with factorlog and remote logging via RPC
type Client struct {
	logRPC      *Outgoing.RPCInterface
	name        string
	localLogger *factorlog.FactorLog
	logLocal    bool
}

//NewLocalClient constructs a new client, which logs to stdout
func NewLocalClient() *Client {
	return &Client{localLogger: Local.GetLogger()}
}

//NewClient creates a localClient or a RPC if a address is given
func NewClient(target string) (*Client, error) {
	if target == "" {
		//use local logger
		return NewLocalClient(), nil
	}
	logLocal := false
	logRPC := Outgoing.NewRPCInterface(target)
	if err := logRPC.Connect(); err != nil {
		logLocal = true
	}

	var clientName string
	for name := range logRPC.Config.NameToCertificate {
		clientName = name
		break
	}
	if clientName == "" {
		var err error
		clientName, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	err := logRPC.SendMessage(LogServer.NewDebugLogMessage(clientName, "connected"))
	if err != nil {
		logLocal = true
	}

	return &Client{logRPC: logRPC, name: clientName, localLogger: Local.GetLogger(), logLocal: logLocal}, nil

}

//LogMultiple sends the logMessages to the remote logServer, log an error to stdout
func (client Client) LogMultiple(messages *[]*LogServer.LogMessage) {
	if client.logLocal {
		for message := range *messages {
			client.localLogger.Println(message)
		}
	} else {
		err := client.logRPC.SendMessages(messages)
		if err != nil {
			client.localLogger.Error(appendStackToMessage(err))
		}
	}
}

//Log sends the logMessage to the remote logServer, log an error to stdout
func (client Client) Log(message *LogServer.LogMessage) {
	if client.logLocal {
		client.localLogger.Println(message)
	} else {
		err := client.logRPC.SendMessage(message)
		if err != nil {
			client.localLogger.Error(appendStackToMessage(err))
		}
	}
}

//Disconnect closes the connection to the remote logServer
func (client Client) Disconnect() {
	if client.logRPC != nil {
		client.logRPC.Disconnect()
	}
}

//Debug logs the message local/remote to on debug level
func (client Client) Debug(v ...interface{}) {
	if client.logLocal {
		client.localLogger.Debug(v)
	} else {
		client.Log(LogServer.NewDebugLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Info logs the message local/remote to on info level
func (client Client) Info(v ...interface{}) {
	if client.logLocal {
		client.localLogger.Info(v)
	} else {
		client.Log(LogServer.NewInfoLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Warn logs the message local/remote to on warn level
func (client Client) Warn(v ...interface{}) {
	if client.logLocal {
		client.localLogger.Warn(v)
	} else {
		client.Log(LogServer.NewWarnLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Error logs the message local/remote to on error level
func (client Client) Error(v ...interface{}) {
	message := appendStackToMessage(v)

	if client.logLocal {
		client.localLogger.Error(message)
	} else {
		client.Log(LogServer.NewErrorLogMessage(client.name, message))
	}
}

func appendStackToMessage(v ...interface{}) string {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	stack := fmt.Sprintf("%s\n", string(buf[0:stackSize]))
	if stackSize > 1000 {
		f, _ := os.Create("error_dump")
		f.WriteString(fmt.Sprintf("%s\n%s", v, stack))
		return fmt.Sprintf("%s\ndumped to file", v)
	}
	return fmt.Sprintf("%s\n%s", v, stack)
}
