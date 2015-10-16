package NetworkInterfaces

import (
	"time"
)

//RPCEvent is used to create events over RPC, if Delay is nil it creates a basic delay
type RPCEvent struct {
	EventAsString string
	Delay         *time.Duration
}
