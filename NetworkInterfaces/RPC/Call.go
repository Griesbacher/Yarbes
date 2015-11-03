package RPC

//Call represents a call on a external module
type Call struct {
	*Event
	Module string
}
