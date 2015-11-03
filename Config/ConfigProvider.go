package Config

import (
	"github.com/griesbacher/Yarbes/Config/ConfigLayouts"
	"gopkg.in/gcfg.v1"
	"sync"
)

var singleServerConfig ConfigLayouts.Server
var singleClientConfig ConfigLayouts.Client
var singleMailConfig ConfigLayouts.Mail
var singleEventsPerTimeConfig ConfigLayouts.EventsPerTime
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

//InitEventsPerTimeConfig parses the EventsPerTime config file
func InitEventsPerTimeConfig(configPath string) {
	initConfig(configPath, &singleEventsPerTimeConfig)
}

func initConfig(configPath string, config interface{}) {
	var err error
	mutex.Lock()
	switch real := config.(type) {
	default:
		err = gcfg.ReadFileInto(real, configPath)
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

//GetEventsPerTimeConfig simple getter for the EventsPerTime config, panics if not initialized
func GetEventsPerTimeConfig() *ConfigLayouts.EventsPerTime {
	if &singleEventsPerTimeConfig == nil {
		panic("Call InitMailConfig first!")
	}
	return &singleEventsPerTimeConfig
}
