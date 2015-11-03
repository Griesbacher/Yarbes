package Incoming

import (
	"github.com/abbot/go-http-auth"
	"github.com/griesbacher/Yarbes/Config"
	"net/http"
)

//HTTPInterface represents a incoming HTTP interface
type HTTPInterface struct {
	quit          chan bool
	isRunning     bool
	HTTPListenTo  string
	authenticator *auth.BasicAuth
}

//NewHTTPInterface creates a new HTTPInterface
func NewHTTPInterface(listenTo string) *HTTPInterface {
	authenticator := auth.NewBasicAuthenticator("Yarbes", auth.HtpasswdFileProvider(Config.GetServerConfig().LogServer.HtpasswdPath))
	http := &HTTPInterface{quit: make(chan bool), isRunning: false, HTTPListenTo: listenTo, authenticator: authenticator}
	return http
}

//Start starts listening for requests
func (httpI *HTTPInterface) Start() {
	if !httpI.isRunning {
		go httpI.serve()
		httpI.isRunning = true
	}
}

//Stop closes the port
func (httpI HTTPInterface) Stop() {
	//do nothing because http closes at program end automatically
}

//IsRunning returns true if the daemon is running
func (httpI HTTPInterface) IsRunning() bool {
	return httpI.isRunning
}

func (httpI *HTTPInterface) serve() {
	err := http.ListenAndServeTLS(httpI.HTTPListenTo, Config.GetServerConfig().TLS.Cert, Config.GetServerConfig().TLS.Key, nil)
	if err != nil {
		panic(err)
	}
}

//PublishHandler add a handler to the given path, basic auth will be used
func (httpI HTTPInterface) PublishHandler(path string, handler auth.AuthenticatedHandlerFunc) {
	http.HandleFunc(path, httpI.authenticator.Wrap(handler))
}
