package main_test

import (
	"github.com/griesbacher/Yarbes/HttpsRpcAnalyse/HttpsTest"
	"github.com/griesbacher/Yarbes/HttpsRpcAnalyse/RpcTest"
	"testing"
)

func BenchmarkRPC(b *testing.B) {
	RPCTest.Client(b.N)
}

func BenchmarkHTTP(b *testing.B) {
	HttpsTest.Client(b.N)
}

func TestMain(m *testing.M){

}