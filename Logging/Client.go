package Logging

import (
	"fmt"
	"github.com/griesbacher/Yarbes/Event"
	"github.com/griesbacher/Yarbes/Logging/Local"
	"github.com/griesbacher/Yarbes/Logging/LogServer"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/RPC/Outgoing"
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

//NewClientOwnName creates a localClient or a RPC if a address is given and uses the given name as source
func NewClientOwnName(target, name string) (*Client, error) {
	if target == "" {
		//use local logger
		return NewLocalClient(), nil
	}
	logLocal := false
	logRPC := Outgoing.NewRPCInterface(target)
	if err := logRPC.Connect(); err != nil {
		fmt.Println(err)
		logLocal = true
	}

	var clientName string
	if name == "" {
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
	} else {
		clientName = name
	}
	err := logRPC.SendMessage(LogServer.NewDebugLogMessage(clientName, clientName+"-connected"))
	if err != nil {
		fmt.Println(err)
		logLocal = true
	}

	return &Client{logRPC: logRPC, name: clientName, localLogger: Local.GetLogger(), logLocal: logLocal}, nil
}

//NewClient creates a localClient or a RPC if a address is given
func NewClient(target string) (*Client, error) {
	return NewClientOwnName(target, "")
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
func (client Client) Log(event *Event.Event, message *LogServer.LogMessage) {
	if client.logLocal {
		client.localLogger.Println(message)
	} else {
		if event != nil {
			message.Event = *event
		}
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

//DebugEvent logs the event and the message local/remote to on debug level
func (client Client) DebugEvent(event *Event.Event, v ...interface{}) {
	if client.logLocal {
		client.localLogger.Debug(v)
	} else {
		client.Log(event, LogServer.NewDebugLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Debug logs the message local/remote to on debug level
func (client Client) Debug(v ...interface{}) {
	client.DebugEvent(nil, v)
}

//Info logs the message local/remote to on info level
func (client Client) Info(v ...interface{}) {
	if client.logLocal {
		client.localLogger.Info(v)
	} else {
		client.Log(nil, LogServer.NewInfoLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Warn logs the message local/remote to on warn level
func (client Client) Warn(v ...interface{}) {
	if client.logLocal {
		client.localLogger.Warn(v)
	} else {
		client.Log(nil, LogServer.NewWarnLogMessage(client.name, fmt.Sprint(v)))
	}
}

//Error logs the message local/remote to on error level
func (client Client) Error(v ...interface{}) {
	message := appendStackToMessage(v)
	if client.logLocal {
		client.localLogger.Error(message)
	} else {
		client.Log(nil, LogServer.NewErrorLogMessage(client.name, message))
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
