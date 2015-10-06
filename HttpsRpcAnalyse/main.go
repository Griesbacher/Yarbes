package main

import (
	"time"
	"os"
	"strconv"
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/HttpsTest"
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/RpcTest"
)

func main() {
	if len(os.Args != 3) {
		panic("arg1: http|rpc ,arg2:rounds")
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
	}else {
		go RpcTest.Server()
		time.Sleep(time.Duration(5) * time.Second)
		RpcTest.Client(loops)

	}
}
