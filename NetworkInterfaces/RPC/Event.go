package RPC

import (
	"time"
)

//Event is used to create events over RPC, if Delay is nil it creates a basic delay
type Event struct {
	EventAsString string
	Delay         *time.Duration
}
