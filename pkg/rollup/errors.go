package rollup

import "errors"

var (
	ErrUnknown          = errors.New("unknown error")
	ErrIOError          = errors.New("I/O error")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrNotInitialized   = errors.New("rollup not initialized")
	ErrAlreadyFinished  = errors.New("already finished")
)
