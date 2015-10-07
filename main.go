package main

import (
	"github.com/griesbacher/SystemX/Event"
	"github.com/griesbacher/SystemX/RuleSystem"
	"time"
	"os"
	"os/signal"
	"syscall"
	"github.com/griesbacher/SystemX/Config"
	"flag"
	"fmt"
)

type Stoppable interface {
	Stop()
}

func main() {
	var configPath string
	flag.Usage = func() {
		fmt.Println(`SystemX by Philip Griesbacher @ 2015
Commandline Parameter:
-configPath Path to the config file. If no file path is given the default is ./config.gcfg.
		`)
	}
	flag.StringVar(&configPath, "configPath", "config.gcfg", "path to the config file")
	flag.Parse()
	Config.InitConfigProvider(configPath)

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

	time.Sleep(time.Duration(5) * time.Second)
	//Listen for Interrupts
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT)
	signal.Notify(interruptChannel, syscall.SIGTERM)
	go func() {
		<-interruptChannel
		cleanUp([]Stoppable{ruleSystem})
		os.Exit(1)
	}()

	//wait for the end to come
	for {
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func cleanUp(itemsToStop []Stoppable) {
	for _, item := range itemsToStop {
		item.Stop()
		time.Sleep(500 * time.Millisecond)
	}
}