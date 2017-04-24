package RPC

import (
	"encoding/json"
	"time"
)

//EventCreatedField is the field in the event which contains the unix time of the event creation
const EventCreatedField = "__created"

//Event is used to create events over RPC, if Delay is nil it creates a basic delay
type Event struct {
	EventAsString string
	Delay         *time.Duration
}

func NewEvent(eventAsString string) Event {
	return NewDelayedEvent(eventAsString, nil)
}

func NewDelayedEvent(eventAsString string, delay *time.Duration) Event {
	var f interface{}
	err := json.Unmarshal([]byte(eventAsString), &f)
	if err != nil {
		return Event{EventAsString: eventAsString, Delay: delay}
	}
	m := f.(map[string]interface{})
	m[EventCreatedField] = time.Now().Unix()
	eventAsBytes, err := json.Marshal(m)
	if err != nil {
		return Event{EventAsString: eventAsString, Delay: delay}
	}
	return Event{EventAsString: string(eventAsBytes), Delay: delay}
}
