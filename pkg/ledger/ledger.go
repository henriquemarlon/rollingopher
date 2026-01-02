//go:build riscv64

package ledger

/*
#cgo LDFLAGS: -lcma -lcmt

#include <libcma/ledger.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Ledger struct {
	ledger C.cma_ledger_t
}

func New() (*Ledger, error) {
	l := &Ledger{}
	rc := C.cma_ledger_init(&l.ledger)
	if rc != 0 {
		return nil, mapError(rc)
	}
	return l, nil
}

func (l *Ledger) Close() error {
	rc := C.cma_ledger_fini(&l.ledger)
	if rc != 0 {
		return mapError(rc)
	}
	return nil
}

func (l *Ledger) Reset() error {
	rc := C.cma_ledger_reset(&l.ledger)
	if rc != 0 {
		return mapError(rc)
	}
	return nil
}

func (l *Ledger) RetrieveAsset(tokenAddress common.Address, tokenID *big.Int, assetType AssetType, op RetrieveOperation) (AssetID, error) {
	var cAssetID C.cma_ledger_asset_id_t
	var cTokenAddress C.cma_token_address_t
	var cTokenID C.cma_token_id_t
	var cAssetType C.cma_ledger_asset_type_t

	copy(cTokenAddress[:], tokenAddress[:])

	if tokenID != nil {
		tokenIDBytes := tokenID.Bytes()
		if len(tokenIDBytes) > 32 {
			return 0, fmt.Errorf("tokenID too large")
		}
		copy(cTokenID[32-len(tokenIDBytes):], tokenIDBytes)
	}

	cAssetType = C.cma_ledger_asset_type_t(assetType)

	rc := C.cma_ledger_retrieve_asset(
		&l.ledger,
		&cAssetID,
		&cTokenAddress,
		&cTokenID,
		&cAssetType,
		C.cma_ledger_retrieve_operation_t(op),
	)
	if rc != 0 {
		return 0, mapError(rc)
	}
	return AssetID(cAssetID), nil
}

func (l *Ledger) RetrieveAccountByAddress(address common.Address, op RetrieveOperation) (InternalAccountID, error) {
	var cAccountID C.cma_ledger_account_id_t
	var cAccount C.cma_ledger_account_t
	var cAccountType C.cma_ledger_account_type_t = C.CMA_LEDGER_ACCOUNT_TYPE_WALLET_ADDRESS

	copy(cAccount.address[:], address[:])
	cAccount._type = C.CMA_LEDGER_ACCOUNT_TYPE_WALLET_ADDRESS

	rc := C.cma_ledger_retrieve_account(
		&l.ledger,
		&cAccountID,
		&cAccount,
		nil,
		&cAccountType,
		C.cma_ledger_retrieve_operation_t(op),
	)
	if rc != 0 {
		return 0, mapError(rc)
	}
	return InternalAccountID(cAccountID), nil
}

func (l *Ledger) RetrieveAccountByID(accountID common.Hash, op RetrieveOperation) (InternalAccountID, error) {
	var cAccountID C.cma_ledger_account_id_t
	var cAccount C.cma_ledger_account_t
	var cAccountType C.cma_ledger_account_type_t = C.CMA_LEDGER_ACCOUNT_TYPE_ACCOUNT_ID

	copy(cAccount.account_id[:], accountID[:])
	cAccount._type = C.CMA_LEDGER_ACCOUNT_TYPE_ACCOUNT_ID

	rc := C.cma_ledger_retrieve_account(
		&l.ledger,
		&cAccountID,
		&cAccount,
		nil,
		&cAccountType,
		C.cma_ledger_retrieve_operation_t(op),
	)
	if rc != 0 {
		return 0, mapError(rc)
	}
	return InternalAccountID(cAccountID), nil
}

func (l *Ledger) Deposit(assetID AssetID, accountID InternalAccountID, amount *big.Int) error {
	var cAmount C.cma_amount_t
	amountBytes := amount.Bytes()
	if len(amountBytes) > 32 {
		return fmt.Errorf("amount too large")
	}
	copy(cAmount[32-len(amountBytes):], amountBytes)

	rc := C.cma_ledger_deposit(
		&l.ledger,
		C.cma_ledger_asset_id_t(assetID),
		C.cma_ledger_account_id_t(accountID),
		&cAmount,
	)
	if rc != 0 {
		return mapError(rc)
	}
	return nil
}

func (l *Ledger) Withdraw(assetID AssetID, accountID InternalAccountID, amount *big.Int) error {
	var cAmount C.cma_amount_t
	amountBytes := amount.Bytes()
	if len(amountBytes) > 32 {
		return fmt.Errorf("amount too large")
	}
	copy(cAmount[32-len(amountBytes):], amountBytes)

	rc := C.cma_ledger_withdraw(
		&l.ledger,
		C.cma_ledger_asset_id_t(assetID),
		C.cma_ledger_account_id_t(accountID),
		&cAmount,
	)
	if rc != 0 {
		return mapError(rc)
	}
	return nil
}

func (l *Ledger) Transfer(assetID AssetID, from, to InternalAccountID, amount *big.Int) error {
	var cAmount C.cma_amount_t
	amountBytes := amount.Bytes()
	if len(amountBytes) > 32 {
		return fmt.Errorf("amount too large")
	}
	copy(cAmount[32-len(amountBytes):], amountBytes)

	rc := C.cma_ledger_transfer(
		&l.ledger,
		C.cma_ledger_asset_id_t(assetID),
		C.cma_ledger_account_id_t(from),
		C.cma_ledger_account_id_t(to),
		&cAmount,
	)
	if rc != 0 {
		return mapError(rc)
	}
	return nil
}

func (l *Ledger) GetBalance(assetID AssetID, accountID InternalAccountID) (*big.Int, error) {
	var cBalance C.cma_amount_t

	rc := C.cma_ledger_get_balance(
		&l.ledger,
		C.cma_ledger_asset_id_t(assetID),
		C.cma_ledger_account_id_t(accountID),
		&cBalance,
	)
	if rc != 0 {
		return nil, mapError(rc)
	}

	return new(big.Int).SetBytes(cBalance[:]), nil
}

func (l *Ledger) GetTotalSupply(assetID AssetID) (*big.Int, error) {
	var cSupply C.cma_amount_t

	rc := C.cma_ledger_get_total_supply(
		&l.ledger,
		C.cma_ledger_asset_id_t(assetID),
		&cSupply,
	)
	if rc != 0 {
		return nil, mapError(rc)
	}

	return new(big.Int).SetBytes(cSupply[:]), nil
}

func mapError(rc C.int) error {
	switch rc {
	case C.CMA_LEDGER_ERROR_UNKNOWN:
		return ErrUnknown
	case C.CMA_LEDGER_ERROR_EXCEPTION:
		return ErrException
	case C.CMA_LEDGER_ERROR_INSUFFICIENT_FUNDS:
		return ErrInsufficientFunds
	case C.CMA_LEDGER_ERROR_ACCOUNT_NOT_FOUND:
		return ErrAccountNotFound
	case C.CMA_LEDGER_ERROR_ASSET_NOT_FOUND:
		return ErrAssetNotFound
	case C.CMA_LEDGER_ERROR_SUPPLY_OVERFLOW:
		return ErrSupplyOverflow
	case C.CMA_LEDGER_ERROR_BALANCE_OVERFLOW:
		return ErrBalanceOverflow
	case C.CMA_LEDGER_ERROR_INVALID_ACCOUNT:
		return ErrInvalidAccount
	case C.CMA_LEDGER_ERROR_INSERTION_ERROR:
		return ErrInsertionError
	default:
		return fmt.Errorf("unknown error code: %d", rc)
	}
}
