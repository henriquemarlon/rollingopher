package rollup

import (
	"github.com/ethereum/go-ethereum/common"
)

type RequestType int

const (
	RequestTypeAdvance RequestType = iota
	RequestTypeInspect
)

type Metadata struct {
	ChainID        uint64
	AppContract    common.Address
	MsgSender      common.Address
	BlockNumber    uint64
	BlockTimestamp uint64
	PrevRandao     common.Hash
	Index          uint64
}

type Advance struct {
	Metadata
	Payload []byte
}

type Inspect struct {
	Payload []byte
}
