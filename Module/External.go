package Module
import (
	"bytes"
	"os/exec"
	"io/ioutil"
	"path"
	"fmt"
	"sync"
	"strings"
	"path/filepath"
	"github.com/griesbacher/SystemX/Event"
)

type ExternalModule struct {
	searchPaths []string
	modules     map[string]string
}

var singleExternalModule *ExternalModule = nil
var mutex = &sync.Mutex{}

func GetExternalModule() *ExternalModule {
	mutex.Lock()
	if singleExternalModule == nil {
		//TODO: durch Config ersetzen
		singleExternalModule = &ExternalModule{[]string{"Module"}, map[string]string{}}
	}
	mutex.Unlock()
	return singleExternalModule
}

func (external ExternalModule) Call(moduleName string, event Event.Event) (*Event.Event, error) {
	if !external.doesModuleExist(moduleName) {
		external.searchModules()
		if !external.doesModuleExist(moduleName) {
			panic(fmt.Sprintf("Module: %s not found", moduleName))
		}
	}

	cmd := exec.Command(external.modules[moduleName], string(event.DataRaw))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var newEvent *Event.Event
	newEvent, err = Event.NewEvent(out.Bytes())
	return newEvent, err
}

func (external ExternalModule) doesModuleExist(moduleName string) bool {
	for command, _ := range external.modules {
		if command == moduleName {
			return true
		}
	}
	return false
}

func (external *ExternalModule) searchModules() {
	for _, searchPath := range external.searchPaths {
		files, _ := ioutil.ReadDir(searchPath)
		for _, file := range files {
			moduleName := getFilename(file.Name())
			external.modules[moduleName] = path.Join(searchPath, file.Name())
		}
	}
}

func getFilename(filename string) string {
	extension := filepath.Ext(filename)
	return strings.ToLower(filename[0:len(filename) - len(extension)])
}