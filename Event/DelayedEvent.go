package Event

import "time"

//DelayedEvent is an Event which waits a certain time and then sends itself in the given queue
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
func (dEvent DelayedEvent) Start() {
	go dEvent.wait()
}

//Stop cancels the event
func (dEvent DelayedEvent) Stop() {
	dEvent.quit <- true
	<-dEvent.quit
}

//IsWaiting returns true if the DelayedEvent has not been sent to the event queue
func (dEvent DelayedEvent) IsWaiting() bool {
	return dEvent.isWaiting
}

func (dEvent *DelayedEvent) wait() {
	select {
	case <-dEvent.quit:
		dEvent.quit <- true
		return
	case <-time.After(dEvent.delay):
		dEvent.resultQueue <- *dEvent.Event
		dEvent.isWaiting = false
	}
}
