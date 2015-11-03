package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var event string
	flag.Usage = func() {
		fmt.Println(`Yarbes-Echo by Philip Griesbacher @ 2015
Commandline Parameter:
-event The event which should be printed`)
	}
	flag.StringVar(&event, "event", "", "the event")
	flag.Parse()
	if len(os.Args) > 1 {
		/*	fmt.Println(`{"Event": ` + os.Args[1] + `, "LogMessages" :[{
			"Timestamp" :"now",
			"Severity"  :"Debug",
			"Message"   :"hallo from module",
			"Source"    :"echo module"
			}]
		}`)*/
		fmt.Println(`{"Event": ` + event + `}`)
	}
}
