package Module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/griesbacher/SystemX/Config"
	"io/ioutil"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

//ExternalModule caches all the files within the search path
type ExternalModule struct {
	searchPaths []string
	modules     map[string]string
}

var singleExternalModule *ExternalModule
var mutex = &sync.Mutex{}

//NewExternalModule constructs a new ExternalModule, this is a singleton
func NewExternalModule() *ExternalModule {
	mutex.Lock()
	if singleExternalModule == nil {
		modulePath := Config.GetServerConfig().RuleSystem.ModulePath
		singleExternalModule = &ExternalModule{[]string{modulePath}, map[string]string{}}
	}
	mutex.Unlock()
	return singleExternalModule
}

//Call tries to execute the given Module with the given Event and returns the whole output as Result
func (external ExternalModule) Call(moduleName, args, event string) (*Result, error) {
	if !external.doesModuleExist(moduleName) {
		external.searchModules()
		if !external.doesModuleExist(moduleName) {
			return nil, fmt.Errorf("Module: %s not found", moduleName)
		}
	}

	cmd := exec.Command(external.modules[moduleName], event, args)
	var out bytes.Buffer
	cmd.Stdout = &out
	runtimeErr := cmd.Run()
	var moduleResult Result

	if len(out.Bytes()) != 0 {
		if err := json.Unmarshal(out.Bytes(), &moduleResult); err != nil {
			return nil, err
		}
	}
	if runtimeErr != nil {
		returnCode := runtimeErr.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus()
		if &moduleResult == nil {
			moduleResult = Result{ReturnCode: returnCode}
		} else {
			moduleResult.ReturnCode = returnCode
		}
	}

	return &moduleResult, runtimeErr
}

func (external ExternalModule) doesModuleExist(moduleName string) bool {
	for command := range external.modules {
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
			if runtime.GOOS == "windows" && filepath.Ext(file.Name()) == "bin" {
				continue
			} else if filepath.Ext(file.Name()) == "exe" {
				continue
			}
			moduleName := getFilename(file.Name())
			external.modules[moduleName] = path.Join(searchPath, file.Name())
		}
	}
}

func getFilename(filename string) string {
	extension := filepath.Ext(filename)
	return strings.ToLower(filename[0 : len(filename)-len(extension)])
}
