package ledger

import "errors"

var (
	ErrUnknown           = errors.New("unknown error")
	ErrException         = errors.New("exception")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
	ErrAssetNotFound     = errors.New("asset not found")
	ErrSupplyOverflow    = errors.New("supply overflow")
	ErrBalanceOverflow   = errors.New("balance overflow")
	ErrInvalidAccount    = errors.New("invalid account")
	ErrInsertionError    = errors.New("insertion error")
	ErrInvalidAmount     = errors.New("invalid amount")
)
