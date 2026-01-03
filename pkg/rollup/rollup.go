//go:build riscv64

package rollup

/*
#cgo LDFLAGS: -lcmt

#include <libcmt/rollup.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"math/big"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
)

type Rollup struct {
	rollup C.cmt_rollup_t
}

func New() (*Rollup, error) {
	r := &Rollup{}
	rc := C.cmt_rollup_init(&r.rollup)
	if rc != 0 {
		return nil, fmt.Errorf("cmt_rollup_init failed: %d", rc)
	}
	return r, nil
}

func (r *Rollup) Close() error {
	C.cmt_rollup_fini(&r.rollup)
	return nil
}

func (r *Rollup) EmitNotice(payload []byte) (uint64, error) {
	var cPayload C.cmt_abi_bytes_t
	if len(payload) > 0 {
		cPayload.length = C.size_t(len(payload))
		cPayload.data = C.CBytes(payload)
		defer C.free(cPayload.data)
	}

	var index C.uint64_t
	rc := C.cmt_rollup_emit_notice(&r.rollup, &cPayload, &index)
	if rc != 0 {
		return 0, fmt.Errorf("cmt_rollup_emit_notice failed: %d", rc)
	}
	return uint64(index), nil
}

func (r *Rollup) EmitVoucher(address common.Address, value *big.Int, data []byte) (uint64, error) {
	var cAddress C.cmt_abi_address_t
	addrPtr := (*[20]byte)(unsafe.Pointer(&cAddress.data[0]))
	copy(addrPtr[:], address[:])

	var cValue C.cmt_abi_u256_t
	if value != nil {
		valuePtr := (*[32]byte)(unsafe.Pointer(&cValue.data[0]))
		value.FillBytes(valuePtr[:])
	}

	var cData C.cmt_abi_bytes_t
	if len(data) > 0 {
		cData.length = C.size_t(len(data))
		cData.data = C.CBytes(data)
		defer C.free(cData.data)
	}

	var index C.uint64_t
	rc := C.cmt_rollup_emit_voucher(&r.rollup, &cAddress, &cValue, &cData, &index)
	if rc != 0 {
		return 0, fmt.Errorf("cmt_rollup_emit_voucher failed: %d", rc)
	}
	return uint64(index), nil
}

func (r *Rollup) EmitDelegateCallVoucher(address common.Address, data []byte) (uint64, error) {
	var cAddress C.cmt_abi_address_t
	addrPtr := (*[20]byte)(unsafe.Pointer(&cAddress.data[0]))
	copy(addrPtr[:], address[:])

	var cData C.cmt_abi_bytes_t
	if len(data) > 0 {
		cData.length = C.size_t(len(data))
		cData.data = C.CBytes(data)
		defer C.free(cData.data)
	}

	var index C.uint64_t
	rc := C.cmt_rollup_emit_delegate_call_voucher(&r.rollup, &cAddress, &cData, &index)
	if rc != 0 {
		return 0, fmt.Errorf("cmt_rollup_emit_delegate_call_voucher failed: %d", rc)
	}
	return uint64(index), nil
}

func (r *Rollup) EmitReport(payload []byte) error {
	var cPayload C.cmt_abi_bytes_t
	if len(payload) > 0 {
		cPayload.length = C.size_t(len(payload))
		cPayload.data = C.CBytes(payload)
		defer C.free(cPayload.data)
	}

	rc := C.cmt_rollup_emit_report(&r.rollup, &cPayload)
	if rc != 0 {
		return fmt.Errorf("cmt_rollup_emit_report failed: %d", rc)
	}
	return nil
}

func (r *Rollup) EmitException(payload []byte) error {
	var cPayload C.cmt_abi_bytes_t
	if len(payload) > 0 {
		cPayload.length = C.size_t(len(payload))
		cPayload.data = C.CBytes(payload)
		defer C.free(cPayload.data)
	}

	rc := C.cmt_rollup_emit_exception(&r.rollup, &cPayload)
	if rc != 0 {
		return fmt.Errorf("cmt_rollup_emit_exception failed: %d", rc)
	}
	return nil
}

func (r *Rollup) ReadAdvanceState() (*Advance, error) {
	var cAdvance C.cmt_rollup_advance_t
	rc := C.cmt_rollup_read_advance_state(&r.rollup, &cAdvance)
	if rc != 0 {
		return nil, fmt.Errorf("cmt_rollup_read_advance_state failed: %d", rc)
	}

	advance := &Advance{
		Metadata: Metadata{
			ChainID:        uint64(cAdvance.chain_id),
			BlockNumber:    uint64(cAdvance.block_number),
			BlockTimestamp: uint64(cAdvance.block_timestamp),
			Index:          uint64(cAdvance.index),
		},
	}

	appContractPtr := (*[20]byte)(unsafe.Pointer(&cAdvance.app_contract.data[0]))
	copy(advance.AppContract[:], appContractPtr[:])

	msgSenderPtr := (*[20]byte)(unsafe.Pointer(&cAdvance.msg_sender.data[0]))
	copy(advance.MsgSender[:], msgSenderPtr[:])

	prevRandaoPtr := (*[32]byte)(unsafe.Pointer(&cAdvance.prev_randao.data[0]))
	copy(advance.PrevRandao[:], prevRandaoPtr[:])

	if cAdvance.payload.length > 0 {
		advance.Payload = C.GoBytes(cAdvance.payload.data, C.int(cAdvance.payload.length))
	}

	return advance, nil
}

func (r *Rollup) ReadInspectState() (*Inspect, error) {
	var cInspect C.cmt_rollup_inspect_t
	rc := C.cmt_rollup_read_inspect_state(&r.rollup, &cInspect)
	if rc != 0 {
		return nil, fmt.Errorf("cmt_rollup_read_inspect_state failed: %d", rc)
	}

	inspect := &Inspect{}
	if cInspect.payload.length > 0 {
		inspect.Payload = C.GoBytes(cInspect.payload.data, C.int(cInspect.payload.length))
	}

	return inspect, nil
}

func (r *Rollup) Finish(accept bool) (RequestType, uint32, error) {
	var finish C.cmt_rollup_finish_t
	finish.accept_previous_request = C.bool(accept)

	rc := C.cmt_rollup_finish(&r.rollup, &finish)
	if rc != 0 {
		return 0, 0, fmt.Errorf("cmt_rollup_finish failed: %d", rc)
	}

	var reqType RequestType
	switch finish.next_request_type {
	case 0:
		reqType = RequestTypeAdvance
	case 1:
		reqType = RequestTypeInspect
	}

	return reqType, uint32(finish.next_request_payload_length), nil
}
