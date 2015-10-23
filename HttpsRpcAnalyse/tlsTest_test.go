package main_test

import (
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/HttpsTest"
	"github.com/griesbacher/SystemX/HttpsRpcAnalyse/RpcTest"
	"testing"
)

func BenchmarkRPC(b *testing.B) {
	RPCTest.Client(b.N)
}

func BenchmarkHTTP(b *testing.B) {
	HttpsTest.Client(b.N)
}
