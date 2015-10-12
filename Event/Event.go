package Event

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

//Event represents a JSON event
type Event struct {
	dataRaw []byte
	data    map[string]interface{}
	created time.Time
}

func newEvent(jsonBytes []byte, jsonInterface interface{}) (*Event, error) {
	switch data := jsonInterface.(type) {
	case map[string]interface{}:
		return &Event{dataRaw: jsonBytes, data: data, created: time.Now()}, nil
	default:
		return nil, errors.New("Given Jsondata is not in the format: map[string]interface{}")
	}
}

//NewEventFromBytes creates an Event from a byte array
func NewEventFromBytes(jsonBytes []byte) (*Event, error) {
	var jsonInterface interface{}
	err := json.Unmarshal(jsonBytes, &jsonInterface)
	if err != nil {
		return nil, err
	}
	return newEvent(jsonBytes, jsonInterface)
}

//NewEventFromInterface creates an Event from a empty interface
func NewEventFromInterface(jsonInterface interface{}) (*Event, error) {
	jsonBytes, err := json.Marshal(jsonInterface)
	if err != nil {
		return nil, err
	}
	return newEvent(jsonBytes, jsonInterface)
}

//GetDataAsInterface returns the internal data as object
func (event Event) GetDataAsInterface() map[string]interface{} {
	return event.data
}

//String returns the internal data as string
func (event Event) String() string {
	return strings.TrimSpace(string(event.dataRaw))
}
