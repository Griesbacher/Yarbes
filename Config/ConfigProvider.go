package Config

import (
	"gopkg.in/gcfg.v1"
	"sync"
)

type ServerConfigProvider struct {
	Cfg ServerConfig
}

type ClientConfigProvider struct {
	Cfg ClientConfig
}

var singleServerConfigProvider *ServerConfigProvider = nil
var singleClientConfigProvider *ClientConfigProvider = nil
var mutex = &sync.Mutex{}

func InitServerConfigProvider(configPath string) {
	InitConfigProvider(configPath, singleServerConfigProvider)
}

func InitClientonfigProvider(configPath string) {
	InitConfigProvider(configPath, singleClientConfigProvider)
}

func InitConfigProvider(configPath string, provider interface{}) {
	mutex.Lock()
	var err error
	switch provider.(type) {
	case *ServerConfigProvider:
		var cfg ServerConfig
		err = gcfg.ReadFileInto(&cfg, configPath)
		if err != nil {
			panic(err)
		}
		singleServerConfigProvider = &ServerConfigProvider{Cfg: cfg}
	case *ClientConfigProvider:
		var cfg ClientConfig
		err = gcfg.ReadFileInto(&cfg, configPath)
		if err != nil {
			panic(err)
		}
		singleClientConfigProvider = &ClientConfigProvider{Cfg: cfg}
	}
	mutex.Unlock()
}

func GetServerConfig() ServerConfig {
	if singleServerConfigProvider == nil {
		panic("Call InitServerConfigProvider first!")
	}
	return singleServerConfigProvider.Cfg
}

func GetClientConfig() ClientConfig {
	if singleClientConfigProvider == nil {
		panic("Call InitClientonfigProvider first!")
	}
	return singleClientConfigProvider.Cfg
}
