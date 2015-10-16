package Event

import "time"

//DelayedEvent is a Event which waits a certain time and then sends itself in the given queue
type DelayedEvent struct {
	*Event
	delay       time.Duration
	resultQueue chan Event
	quit        chan bool
	isWaiting   bool
}

//NewDelayedEvent creates a new DelayedEvent
func NewDelayedEvent(event *Event, delay time.Duration, resultQueue chan Event) *DelayedEvent {
	return &DelayedEvent{Event: event, delay: delay, resultQueue: resultQueue, quit: make(chan bool), isWaiting: true}
}

//Start starts waiting
func (delayed *DelayedEvent) Start() {
	go delayed.wait()
}

//Stop cancels the event
func (delayed DelayedEvent) Stop() {
	delayed.quit <- true
	<-delayed.quit
}

//IsWaiting returns true if the DelayedEvent has not been sent to the event queue
func (delayed DelayedEvent) IsWaiting() bool {
	return delayed.isWaiting
}

func (delayed DelayedEvent) wait() {
	select {
	case <-delayed.quit:
		delayed.quit <- true
		return
	case <-time.After(delayed.delay):
		delayed.resultQueue <- *delayed.Event
		delayed.isWaiting = false
	}
}
