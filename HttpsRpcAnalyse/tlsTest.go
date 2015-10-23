package main

import (
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/HttpsTest"
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/RpcTest"
	"os"
	"strconv"
	"time"
	"log"
"runtime/pprof"
)

func main() {

	if len(os.Args) != 4 {
		panic("arg1: http|rpc ,arg2:rounds, cpuprofile")
	}

	if os.Args[3] != "" {
		f, err := os.Create(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	loops, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	if os.Args[1] == "http" {
		go HttpsTest.Server()
		time.Sleep(time.Duration(5) * time.Second)
		client := HttpsTest.Client()
		for i := 0; i < loops; i++ {
			HttpsTest.Request(client, "test string")
		}
	} else {
		go RPCTest.Server()
		time.Sleep(time.Duration(5) * time.Second)
		RPCTest.Client(loops)
	}
}
