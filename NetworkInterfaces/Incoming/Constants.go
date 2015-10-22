package Incoming

import "errors"

//ErrorInputWasNil will be returned if the RPC caller sends a nil object as input
var ErrorInputWasNil = errors.New("The given input was a nil pointer")

//ErrorResultWasNil will be returned if the RPC caller sends a nil object as result
var ErrorResultWasNil = errors.New("The given result was a nil pointer")
