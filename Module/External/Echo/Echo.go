package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var event string
	var extra string
	flag.Usage = func() {
		fmt.Println(`Yarbes-Echo by Philip Griesbacher @ 2015
Commandline Parameter:
-event The event which should be printed`)
	}
	flag.StringVar(&event, "event", "", "the event")
	flag.StringVar(&extra, "extra", "", "the extra")
	flag.Parse()
	if len(os.Args) > 1 {
		fmt.Println(`{"Event": ` + event + `, "Messages" :[
			{
			"Timestamp" :"now",
			"Severity"  :"Debug",
			"Message"   :"Event dump: ` + strings.Replace(strings.Replace(event, `"`, "'", -1), `\`, "", -1) + `.",
			"Source"    :"echo module"
			},
			{
			"Timestamp" :"now",
			"Severity"  :"Debug",
			"Message"   :"Extra Flag: ` + extra + `.",
			"Source"    :"echo module"
			}
			]
		}`)
	}
}
