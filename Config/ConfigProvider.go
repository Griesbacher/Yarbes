package Config

import (
	"errors"
	"github.com/griesbacher/SystemX/Config/ConfigLayouts"
	"gopkg.in/gcfg.v1"
	"sync"
)

var singleServerConfig ConfigLayouts.Server
var singleClientConfig ConfigLayouts.Client
var singleMailConfig ConfigLayouts.Mail
var mutex = &sync.Mutex{}

//InitServerConfig parses the server config file
func InitServerConfig(configPath string) {
	initConfig(configPath, &singleServerConfig)
}

//InitClientConfig parses the client config file
func InitClientConfig(configPath string) {
	initConfig(configPath, &singleClientConfig)
}

//InitMailConfig parses the mail config file
func InitMailConfig(configPath string) {
	initConfig(configPath, &singleMailConfig)
}

func initConfig(configPath string, config interface{}) {
	var err error
	mutex.Lock()
	switch real := config.(type) {
	case *ConfigLayouts.Server, *ConfigLayouts.Client, *ConfigLayouts.Mail:
		err = gcfg.ReadFileInto(real, configPath)
	default:
		err = errors.New("Unkown config layout")
	}
	mutex.Unlock()
	if err != nil {
		panic(err)
	}
}

//GetServerConfig simple getter for the server config, panics if not initialized
func GetServerConfig() *ConfigLayouts.Server {
	if &singleServerConfig == nil {
		panic("Call InitServerConfig first!")
	}
	return &singleServerConfig
}

//GetClientConfig simple getter for the client config, panics if not initialized
func GetClientConfig() *ConfigLayouts.Client {
	if &singleClientConfig == nil {
		panic("Call InitClientonfig first!")
	}
	return &singleClientConfig
}

//GetMailConfig simple getter for the mail config, panics if not initialized
func GetMailConfig() *ConfigLayouts.Mail {
	if &singleMailConfig == nil {
		panic("Call InitMailConfig first!")
	}
	return &singleMailConfig
}
