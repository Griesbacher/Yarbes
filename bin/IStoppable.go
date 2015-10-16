package bin

//Stoppable a stoppable struct runs till this function is called
type Stoppable interface {
	Stop()
	IsRunning() bool
}
