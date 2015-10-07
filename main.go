package main

import (
	"flag"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem"
	"github.com/griesbacher/SystemX/RuleSystem/Incoming"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StartAndStoppable interface {
	Start()
	Stop()
}

func main() {
	var configPath string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-configPath Path to the config file. If no file path is given the default is ./serverConfig.gcfg.
		`)
	}
	flag.StringVar(&configPath, "configPath", "serverConfig.gcfg", "path to the config file")
	flag.Parse()
	Config.InitServerConfigProvider(configPath)

	b := []byte(`{
   "k1" : "v1",
   "k2" : 10,
   "k3" : ["v4",12.3,{"k11" : "v11", "k22" : "v22"}]
}`)

	event, err := Event.NewEvent(b)
	if err != nil {
		panic(err)
	}
	ruleSystem := RuleSystem.NewRuleSystem()
	ruleSystem.Start()
	ruleSystem.EventQueue <- *event

	rpcI := Incoming.NewRpcInterface(ruleSystem.EventQueue)
	rpcI.Start()

	time.Sleep(time.Duration(5) * time.Second)
	//Listen for Interrupts
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT)
	signal.Notify(interruptChannel, syscall.SIGTERM)
	go func() {
		<-interruptChannel
		cleanUp([]StartAndStoppable{ruleSystem, rpcI})
		os.Exit(1)
	}()

	//wait for the end to come
	for {
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func cleanUp(itemsToStop []StartAndStoppable) {
	for _, item := range itemsToStop {
		item.Stop()
		time.Sleep(500 * time.Millisecond)
	}
}
