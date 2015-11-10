package Event

import (
	"reflect"
	"testing"
)

func TestNewEventFromBytes(t *testing.T) {
	t.Parallel()

	str := []byte(`{"a":"b"}`)
	event, err := NewEventFromBytes(str)

	if err != nil {
		t.Errorf("there should be no error: %s", err)
	}

	if event == nil {
		t.Errorf("there should be an event")
	}

	if !reflect.DeepEqual(event.GetDataAsBytes(), str) {
		t.Errorf("theses two should be equal: %s - %s", event.GetDataAsBytes(), str)
	}

	obj := map[string]interface{}{"a": "b"}
	if !reflect.DeepEqual(event.GetDataAsInterface(), obj) {
		t.Errorf("objects should be equal: %s - %s", event.GetDataAsInterface(), obj)
	}

	if event.String() != string(str) {
		t.Errorf("String method changed string, but shouldn't. %s - %s", event.String(), string(str))
	}

	wrong := []byte(`{"a":b}`)
	_, err = NewEventFromBytes(wrong)
	if err == nil {
		t.Errorf("there should be an error because it is not valid json")
	}
}

func TestNewEventFromInterface(t *testing.T) {
	t.Parallel()

	obj := map[string]interface{}{"a": "b"}
	event, err := NewEventFromInterface(obj)

	if err != nil {
		t.Errorf("there should be no error: %s", err)
	}

	if event == nil {
		t.Errorf("there should be an event")
	}

	str := []byte(`{"a":"b"}`)
	if !reflect.DeepEqual(event.GetDataAsBytes(), str) {
		t.Errorf("theses two should be equal: %s - %s", event.GetDataAsBytes(), str)
	}

	if !reflect.DeepEqual(event.GetDataAsInterface(), obj) {
		t.Errorf("objects should be equal: %s - %s", event.GetDataAsInterface(), obj)
	}

	if event.String() != string(str) {
		t.Errorf("String method changed string, but shouldn't. %s - %s", event.String(), string(str))
	}

	wrong := map[int]interface{}{1: "b"}
	_, err = NewEventFromInterface(wrong)
	if err == nil {
		t.Errorf("there should be an error because it is not valid json")
	}
}

func TestNewEvent(t *testing.T) {
	t.Parallel()
	wrong := map[int]interface{}{1: "b"}
	_, err := newEvent([]byte(`{1:"b"}`), wrong)
	if err == nil {
		t.Errorf("there should be an error because it is not valid json")
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	str := []byte(`{"a":"b","c":1}`)
	event, err := NewEventFromBytes(str)
	if err != nil {
		t.Error(err)
	}

	if event.String() != `{"a":"b", "c":1}` {
		t.Errorf("String should be formated. %s - %s", event, string(str))
	}
}
