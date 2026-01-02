package parser

import "errors"

var (
	ErrMalformedInput     = errors.New("malformed input")
	ErrDepositFailed      = errors.New("deposit failed")
	ErrUnknownInputType   = errors.New("unknown input type")
	ErrIncompatibleInput  = errors.New("incompatible input")
	ErrInvalidSelector    = errors.New("invalid function selector")
	ErrInvalidPayloadSize = errors.New("invalid payload size")
)
