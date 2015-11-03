package RPC

//Result contains the result message as string and if given an error
type Result struct {
	Message string
	Err     error
}
