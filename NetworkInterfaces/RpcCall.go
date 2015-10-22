package NetworkInterfaces

//RPCCall represents a call on a external module
type RPCCall struct {
	*RPCEvent
	Module string
}
