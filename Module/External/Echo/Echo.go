package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(`{"Event": ` + os.Args[1] + `, "LogMessages" :[{
			"Timestamp" :"now",
			"Severity"  :"Debug",
			"Message"   :"hallo from module",
			"Source"    :"echo module"
			}]
		}`)
	}
}
