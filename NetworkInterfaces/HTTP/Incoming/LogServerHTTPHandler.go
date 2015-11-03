package Incoming

import (
	"io"
	"net/http"
)

//LogServerHTTPHandler is a RPC handler which accepts LogMessages
type LogServerHTTPHandler struct {
}

//LogView displays the basic logs
func (handler LogServerHTTPHandler) LogView(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
