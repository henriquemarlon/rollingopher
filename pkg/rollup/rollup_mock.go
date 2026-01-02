//go:build !riscv64

package rollup

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type Rollup struct {
	mu                   sync.Mutex
	vouchers             []Voucher
	delegateCallVouchers []DelegateCallVoucher
	notices              [][]byte
	reports              [][]byte
	advances             []*Advance
	inspects             []*Inspect
	finished             bool
	rejected             bool

	advanceIdx int
	inspectIdx int
}

type Voucher struct {
	Address common.Address
	Value   *big.Int
	Data    []byte
}

type DelegateCallVoucher struct {
	Address common.Address
	Data    []byte
}

func New() (*Rollup, error) {
	return &Rollup{
		vouchers:             make([]Voucher, 0),
		delegateCallVouchers: make([]DelegateCallVoucher, 0),
		notices:              make([][]byte, 0),
		reports:              make([][]byte, 0),
		advances:             make([]*Advance, 0),
		inspects:             make([]*Inspect, 0),
	}, nil
}

func (r *Rollup) Close() error {
	return nil
}

func (r *Rollup) EmitVoucher(address common.Address, value *big.Int, data []byte) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	r.vouchers = append(r.vouchers, Voucher{
		Address: address,
		Value:   new(big.Int).Set(value),
		Data:    dataCopy,
	})
	return uint64(len(r.vouchers) - 1), nil
}

func (r *Rollup) EmitDelegateCallVoucher(address common.Address, data []byte) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	r.delegateCallVouchers = append(r.delegateCallVouchers, DelegateCallVoucher{
		Address: address,
		Data:    dataCopy,
	})
	return uint64(len(r.delegateCallVouchers) - 1), nil
}

func (r *Rollup) EmitNotice(payload []byte) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	payloadCopy := make([]byte, len(payload))
	copy(payloadCopy, payload)

	r.notices = append(r.notices, payloadCopy)
	return uint64(len(r.notices) - 1), nil
}

func (r *Rollup) EmitReport(payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	payloadCopy := make([]byte, len(payload))
	copy(payloadCopy, payload)

	r.reports = append(r.reports, payloadCopy)
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
		advance := r.advances[r.advanceIdx]
		return RequestTypeAdvance, uint32(len(advance.Payload)), nil
	}

	if r.inspectIdx < len(r.inspects) {
		inspect := r.inspects[r.inspectIdx]
		return RequestTypeInspect, uint32(len(inspect.Payload)), nil
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

// Mock helpers

func (r *Rollup) AddAdvance(advance *Advance) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.advances = append(r.advances, advance)
}

func (r *Rollup) AddInspect(inspect *Inspect) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inspects = append(r.inspects, inspect)
}

func (r *Rollup) GetVouchers() []Voucher {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.vouchers
}

func (r *Rollup) GetDelegateCallVouchers() []DelegateCallVoucher {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.delegateCallVouchers
}

func (r *Rollup) GetNotices() [][]byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.notices
}

func (r *Rollup) GetReports() [][]byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reports
}

func (r *Rollup) IsFinished() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.finished
}

func (r *Rollup) IsRejected() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rejected
}

func (r *Rollup) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vouchers = make([]Voucher, 0)
	r.delegateCallVouchers = make([]DelegateCallVoucher, 0)
	r.notices = make([][]byte, 0)
	r.reports = make([][]byte, 0)
	r.advances = make([]*Advance, 0)
	r.inspects = make([]*Inspect, 0)
	r.finished = false
	r.rejected = false
	r.advanceIdx = 0
	r.inspectIdx = 0
}
