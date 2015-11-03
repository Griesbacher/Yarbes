package bin

import (
	"fmt"
	"github.com/griesbacher/Yarbes/Client/Livestatus"
	"github.com/griesbacher/Yarbes/Config"
	"github.com/griesbacher/Yarbes/Logging"
	"github.com/griesbacher/Yarbes/NetworkInterfaces/Outgoing"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

//Client starts a example client
func Client(configPath, cpuProfile string) {

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	Config.InitClientConfig(configPath)

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

	delayed(eventRPC, 1)

	logger.Debug("Start")
	useLivestatus(logger, eventRPC)
	logger.Debug("Fertig")
	logger.Disconnect()
	eventRPC.Disconnect()
}

func useLivestatus(logger *Logging.Client, eventRPC *Outgoing.RPCInterface) {
	livestatus := Livestatus.NewCollector(*logger, eventRPC)
	livestatus.Start()
	time.Sleep(time.Duration(30) * time.Second)
	livestatus.Stop()
}

func delayed(eventRPC *Outgoing.RPCInterface, wait int) {
	delay := time.Duration(wait) * time.Second
	var event = []byte(`{"Hallo": "Delayed", "Start":"` + time.Now().Format(time.RFC3339) + `", "hostname":"localhost", "type":"ALERT", "time":` + fmt.Sprintf("%d", time.Now().Unix()) + `}`)
	eventRPC.CreateDelayedEvent(event, &delay)
	time.Sleep(delay * 2)
}
