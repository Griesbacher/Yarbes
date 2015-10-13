package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(`{"Event": ` + os.Args[1] + `, "LogMessages" :[{
			"Timestamp" :"now",
			"Level"     :"debug",
			"Message"   :"hallo from module",
			"Source"    :"echo module"
			}]
		}`)
	}
}
