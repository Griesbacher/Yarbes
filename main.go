package main

import (
	"time"
"github.com/griesbacher/SystemX/HttpsTest"
)

func main() {
	go HttpsTest.Server()
	time.Sleep(time.Duration(5)*time.Second)
	client := HttpsTest.Client()

	HttpsTest.Request(client, "test string")

}
