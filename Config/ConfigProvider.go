package Config
import (
	"sync"
	"gopkg.in/gcfg.v1"
)

type ConfigProvider struct {
	configPath string
	Cfg        Config
}

var singleConfigProvider *ConfigProvider = nil
var mutex = &sync.Mutex{}

func InitConfigProvider(configPath string) *ConfigProvider{
	mutex.Lock()

	if singleConfigProvider == nil {
		var cfg Config
		err := gcfg.ReadFileInto(&cfg, configPath)
		if err != nil {
			panic(err)
		}
		singleConfigProvider = &ConfigProvider{configPath:configPath, Cfg:cfg}
	}
	mutex.Unlock()
	return singleConfigProvider
}

func GetConfig() *Config{
	if singleConfigProvider == nil {
		panic("Call GetConfigProvider first!")
	}
	return singleConfigProvider.Cfg
}