//go:build !riscv64

package rollup

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// Mock-specific types for capturing emitted outputs during testing.

type Notice struct {
	Payload []byte
}

type Voucher struct {
	Destination common.Address
	Value       *big.Int
	Payload     []byte
}

type Report struct {
	Payload []byte
}

type DelegateCallVoucher struct {
	Destination common.Address
	Payload     []byte
}

type Rollup struct {
	mu                   sync.Mutex
	vouchers             []Voucher
	delegateCallVouchers []DelegateCallVoucher
	notices              []Notice
	reports              []Report
	advances             []*Advance
	inspects             []*Inspect
	finished             bool
	rejected             bool
	advanceIdx           int
	inspectIdx           int
}

func New() (*Rollup, error) {
	return &Rollup{}, nil
}

func (r *Rollup) Close() error {
	return nil
}

func (r *Rollup) EmitVoucher(address common.Address, value *big.Int, data []byte) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.vouchers = append(r.vouchers, Voucher{
		Destination: address,
		Value:       new(big.Int).Set(value),
		Payload:     append([]byte(nil), data...),
	})
	return uint64(len(r.vouchers) - 1), nil
}

func (r *Rollup) EmitDelegateCallVoucher(address common.Address, data []byte) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.delegateCallVouchers = append(r.delegateCallVouchers, DelegateCallVoucher{
		Destination: address,
		Payload:     append([]byte(nil), data...),
	})
	return uint64(len(r.delegateCallVouchers) - 1), nil
}

func (r *Rollup) EmitNotice(payload []byte) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.notices = append(r.notices, Notice{Payload: append([]byte(nil), payload...)})
	return uint64(len(r.notices) - 1), nil
}

func (r *Rollup) EmitReport(payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.reports = append(r.reports, Report{Payload: append([]byte(nil), payload...)})
	return nil
}

func (r *Rollup) EmitException(payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rejected = true
	return nil
}

func (r *Rollup) Finish(accept bool) (RequestType, uint32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.finished = true

	if r.advanceIdx < len(r.advances) {
		return RequestTypeAdvance, uint32(len(r.advances[r.advanceIdx].Payload)), nil
	}
	if r.inspectIdx < len(r.inspects) {
		return RequestTypeInspect, uint32(len(r.inspects[r.inspectIdx].Payload)), nil
	}
	return RequestTypeAdvance, 0, nil
}

func (r *Rollup) ReadAdvanceState() (*Advance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.advanceIdx >= len(r.advances) {
		return nil, ErrNotInitialized
	}

	advance := r.advances[r.advanceIdx]
	r.advanceIdx++
	return advance, nil
}

func (r *Rollup) ReadInspectState() (*Inspect, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.inspectIdx >= len(r.inspects) {
		return nil, ErrNotInitialized
	}

	inspect := r.inspects[r.inspectIdx]
	r.inspectIdx++
	return inspect, nil
}

func (r *Rollup) Advance(advance *Advance) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.advances = append(r.advances, advance)
}

func (r *Rollup) Inspect(inspect *Inspect) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inspects = append(r.inspects, inspect)
}

func (r *Rollup) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.vouchers = nil
	r.delegateCallVouchers = nil
	r.notices = nil
	r.reports = nil
	r.advances = nil
	r.inspects = nil
	r.finished = false
	r.rejected = false
	r.advanceIdx = 0
	r.inspectIdx = 0
}
