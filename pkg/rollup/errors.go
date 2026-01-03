package rollup

import "errors"

var (
	ErrIOError         = errors.New("I/O error")
	ErrUnknown         = errors.New("unknown error")
	ErrAlreadyFinished = errors.New("already finished")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotInitialized  = errors.New("rollup not initialized")
)
