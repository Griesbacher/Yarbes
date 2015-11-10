package Event

import (
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	t.Parallel()
	eChan := make(chan Event)
	delay := time.Duration(25) * time.Millisecond
	event := &Event{}
	dEvent := NewDelayedEvent(event, delay, eChan)

	dEvent.Start()

	if !dEvent.IsWaiting() {
		t.Error("Event should be waiting after start")
	}

	select {
	case <-eChan:
	case <-time.After(delay * 150 / 100):
		t.Error("Could not recive delayed Event within time")
	}
}

func TestStop(t *testing.T) {
	t.Parallel()
	eChan := make(chan Event)
	delay := time.Duration(50) * time.Millisecond
	event := &Event{}
	dEvent := NewDelayedEvent(event, delay, eChan)

	dEvent.Start()

	if !dEvent.IsWaiting() {
		t.Error("Event should be waiting after start")
	}

	dEvent.Stop()
	if dEvent.IsWaiting() {
		t.Error("Event should be stopped")
	}

	select {
	case <-eChan:
		t.Error("Due to stop call, there should be no event left")
	case <-time.After(delay / 2):
	}
}
