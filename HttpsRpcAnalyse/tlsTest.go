package main

import (
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/HttpsTest"
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/RpcTest"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		panic("arg1: http|rpc|server ,arg2:rounds, cpuprofile")
	}

	if len(os.Args) > 3 && os.Args[3] != "" {
		cpu, err := os.Create(os.Args[3] + ".cpu")
		if err != nil {
			log.Fatal(err)
		}
		heap, err := os.Create(os.Args[3] + ".heap")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpu)
		pprof.WriteHeapProfile(heap)
		defer pprof.StopCPUProfile()
	}

	loops, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	if os.Args[1] == "http" {
		go HttpsTest.Server()
		time.Sleep(time.Duration(5) * time.Second)
		HttpsTest.Client(loops)
	} else if os.Args[1] == "rpc" {
		go RPCTest.Server()
		time.Sleep(time.Duration(5) * time.Second)
		RPCTest.Client(loops)
	} else if os.Args[1] == "server" {
		go HttpsTest.Server()
		RPCTest.Server()
	}
}
