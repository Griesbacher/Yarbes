package main

import (
	"github.com/griesbacher/SystemX/bin"
	"time"
)

func main() {
	go bin.Server()
	time.Sleep(time.Duration(1) * time.Second)
	bin.Client()
	time.Sleep(time.Duration(1) * time.Second)
}
