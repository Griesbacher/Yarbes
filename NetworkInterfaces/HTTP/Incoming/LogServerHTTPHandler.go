package Incoming

import (
	"github.com/abbot/go-http-auth"
	"io"
	"net/http"
)

//LogServerHTTPHandler is a HTTP handler which displays LogMessages
type LogServerHTTPHandler struct {
}

//LogView displays the basic logs
func (handler LogServerHTTPHandler) LogView(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	io.WriteString(w, "Hello world2!")
}
