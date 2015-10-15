package Client

//EventCreatable represents structs which can create events
type EventCreatable interface {
	CreateEvent(event []byte) error
}
