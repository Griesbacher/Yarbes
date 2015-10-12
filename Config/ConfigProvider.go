package Config

import (
	"gopkg.in/gcfg.v1"
	"sync"
)

//ServerConfigProvider provides the server configfile as an object
type ServerConfigProvider struct {
	Cfg ServerConfig
}

//ClientConfigProvider provides the client configfile as an object
type ClientConfigProvider struct {
	Cfg ClientConfig
}

var singleServerConfigProvider *ServerConfigProvider
var singleClientConfigProvider *ClientConfigProvider
var mutex = &sync.Mutex{}

//InitServerConfigProvider parses the server config file
func InitServerConfigProvider(configPath string) {
	initConfigProvider(configPath, singleServerConfigProvider)
}

//InitClientConfigProvider parses the client config file
func InitClientConfigProvider(configPath string) {
	initConfigProvider(configPath, singleClientConfigProvider)
}

func initConfigProvider(configPath string, provider interface{}) {
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

//GetServerConfig simple getter for the server config, panics if not initialized
func GetServerConfig() ServerConfig {
	if singleServerConfigProvider == nil {
		panic("Call InitServerConfigProvider first!")
	}
	return singleServerConfigProvider.Cfg
}

//GetClientConfig simple getter for the client config, panics if not initialized
func GetClientConfig() ClientConfig {
	if singleClientConfigProvider == nil {
		panic("Call InitClientonfigProvider first!")
	}
	return singleClientConfigProvider.Cfg
}
