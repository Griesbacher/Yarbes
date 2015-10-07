package Event

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Event struct {
	DataRaw []byte
	data    map[string]interface{}
	created time.Time
}

func NewEvent(jsonData []byte) (*Event, error) {
	var jsonInterface interface{}
	err := json.Unmarshal(jsonData, &jsonInterface)
	if err != nil {
		return nil, err
	}
	switch data := jsonInterface.(type) {
	case map[string]interface{}:
		return &Event{DataRaw: jsonData, data: data, created: time.Now()}, nil
	default:
		return nil, errors.New("Given Jsondata is not in the format: map[string]interface{}")
	}
}

func (event Event) GetDataAsInterface() map[string]interface{} {
	return event.data
}

func (event Event) String() string {
	return strings.TrimSpace(string(event.DataRaw))
}
