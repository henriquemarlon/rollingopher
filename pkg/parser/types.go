package parser

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Outputs

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

// Ledger

type InputType int

const (
	InputTypeNone InputType = iota
	InputTypeEtherDeposit
	InputTypeERC20Deposit
	InputTypeERC721Deposit
	InputTypeERC1155SingleDeposit
	InputTypeERC1155BatchDeposit
	InputTypeEtherWithdrawal
	InputTypeERC20Withdrawal
	InputTypeERC721Withdrawal
	InputTypeERC1155SingleWithdrawal
	InputTypeERC1155BatchWithdrawal
	InputTypeEtherTransfer
	InputTypeERC20Transfer
	InputTypeERC721Transfer
	InputTypeERC1155SingleTransfer
	InputTypeERC1155BatchTransfer
	InputTypeBalance
	InputTypeBalanceAccount
	InputTypeBalanceAccountTokenAddress
	InputTypeBalanceAccountTokenAddressID
	InputTypeSupply
	InputTypeSupplyTokenAddress
	InputTypeSupplyTokenAddressID
)

type EtherDeposit struct {
	Sender        common.Address
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC20Deposit struct {
	Token         common.Address
	Sender        common.Address
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC721Deposit struct {
	Token         common.Address
	Sender        common.Address
	TokenID       *big.Int
	ExecLayerData []byte
}

type ERC1155SingleDeposit struct {
	Token         common.Address
	Sender        common.Address
	TokenID       *big.Int
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC1155BatchDeposit struct {
	Token         common.Address
	Sender        common.Address
	TokenIDs      []*big.Int
	Amounts       []*big.Int
	BaseLayerData []byte
	ExecLayerData []byte
}

type EtherWithdrawal struct {
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC20Withdrawal struct {
	Token         common.Address
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC721Withdrawal struct {
	Token         common.Address
	TokenID       *big.Int
	ExecLayerData []byte
}

type ERC1155SingleWithdrawal struct {
	Token         common.Address
	TokenID       *big.Int
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC1155BatchWithdrawal struct {
	Token         common.Address
	TokenIDs      []*big.Int
	Amounts       []*big.Int
	ExecLayerData []byte
}

type EtherTransfer struct {
	Receiver      common.Hash
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC20Transfer struct {
	Token         common.Address
	Receiver      common.Hash
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC721Transfer struct {
	Token         common.Address
	Receiver      common.Hash
	TokenID       *big.Int
	ExecLayerData []byte
}

type ERC1155SingleTransfer struct {
	Token         common.Address
	Receiver      common.Hash
	TokenID       *big.Int
	Amount        *big.Int
	ExecLayerData []byte
}

type ERC1155BatchTransfer struct {
	Token         common.Address
	Receiver      common.Hash
	TokenIDs      []*big.Int
	Amounts       []*big.Int
	ExecLayerData []byte
}

type BalanceQuery struct {
	Account       common.Hash
	Token         common.Address
	TokenID       *big.Int
	ExecLayerData []byte
}

type SupplyQuery struct {
	Token         common.Address
	TokenID       *big.Int
	ExecLayerData []byte
}
