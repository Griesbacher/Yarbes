package bin

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Client/Livestatus"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Logging"
	"github.com/griesbacher/SystemX/NetworkInterfaces/Outgoing"
	"os"
	"time"
)

//Client starts a example client
func Client() {
	var configPath string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-configPath Path to the config file. If no file path is given the default is ./serverConfig.gcfg.
		`)
	}
	flag.StringVar(&configPath, "configPath", "clientConfig.gcfg", "path to the config file")
	flag.Parse()
	Config.InitClientConfigProvider(configPath)

	var logger *Logging.Client
	logger, err := Logging.NewClient(Config.GetClientConfig().LogServer.RPCInterface)
	if err != nil {
		fmt.Println("using local logger")
		logger = Logging.NewLocalClient()
	}

	eventRPC := Outgoing.NewRPCInterface(Config.GetClientConfig().Backend.RPCInterface)
	err = eventRPC.Connect()
	if err != nil {
		logger.Error(err)
		os.Exit(2)
	}
	/*delay := time.Duration(5) * time.Second
	var event = []byte(`{"Hallo": "Delayed", "Start":"` + time.Now().Format(time.RFC3339) + `"}`)
	eventRPC.CreateDelayedEvent(event, &delay)
*/
	logger.Debug("Start")
	livestatus := Livestatus.NewCollector(*logger, eventRPC)
	livestatus.Start()
	time.Sleep(time.Duration(30) * time.Second)
	livestatus.Stop()
	logger.Debug("Fertig")
	logger.Disconnect()

	eventRPC.Disconnect()
}
