package parser

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/rollingopher/pkg/rollup"
)

var (
	EtherPortal         = common.HexToAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044")
	ERC20Portal         = common.HexToAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB")
	ERC721Portal        = common.HexToAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87")
	ERC1155SinglePortal = common.HexToAddress("0x7CFB0193Ca87eB6e48056885E026552c3A941FC4")
	ERC1155BatchPortal  = common.HexToAddress("0xedB53860A6B52bbb7561Ad596416ee9965B055Aa")
)

const (
	SelectorWithdrawEther         uint32 = 0x8cf70f0b
	SelectorWithdrawERC20         uint32 = 0x4f94d342
	SelectorWithdrawERC721        uint32 = 0x33acf293
	SelectorWithdrawERC1155Single uint32 = 0x8bb0a811
	SelectorWithdrawERC1155Batch  uint32 = 0x50c80019

	SelectorTransferEther         uint32 = 0xff67c903
	SelectorTransferERC20         uint32 = 0x03d61dcd
	SelectorTransferERC721        uint32 = 0xaf615a5a
	SelectorTransferERC1155Single uint32 = 0xe1c913ed
	SelectorTransferERC1155Batch  uint32 = 0x638ac6f9

	SelectorERC20Transfer            uint32 = 0xa9059cbb
	SelectorERC721SafeTransferFrom   uint32 = 0x42842e0e
	SelectorERC1155SafeTransferFrom  uint32 = 0xf242432a
	SelectorERC1155SafeBatchTransfer uint32 = 0x2eb2c2d6
)

var (
	erc20ABI   abi.ABI
	erc721ABI  abi.ABI
	erc1155ABI abi.ABI
)

func init() {
	var err error

	erc20ABI, err = abi.JSON(strings.NewReader(erc20ABIJson))
	if err != nil {
		panic("failed to parse ERC20 ABI: " + err.Error())
	}

	erc721ABI, err = abi.JSON(strings.NewReader(erc721ABIJson))
	if err != nil {
		panic("failed to parse ERC721 ABI: " + err.Error())
	}

	erc1155ABI, err = abi.JSON(strings.NewReader(erc1155ABIJson))
	if err != nil {
		panic("failed to parse ERC1155 ABI: " + err.Error())
	}
}

const erc20ABIJson = `[{
	"name": "transfer",
	"type": "function",
	"inputs": [
		{"name": "to", "type": "address"},
		{"name": "amount", "type": "uint256"}
	]
}]`

const erc721ABIJson = `[{
	"name": "safeTransferFrom",
	"type": "function",
	"inputs": [
		{"name": "from", "type": "address"},
		{"name": "to", "type": "address"},
		{"name": "tokenId", "type": "uint256"}
	]
}]`

const erc1155ABIJson = `[
	{
		"name": "safeTransferFrom",
		"type": "function",
		"inputs": [
			{"name": "from", "type": "address"},
			{"name": "to", "type": "address"},
			{"name": "id", "type": "uint256"},
			{"name": "amount", "type": "uint256"},
			{"name": "data", "type": "bytes"}
		]
	},
	{
		"name": "safeBatchTransferFrom",
		"type": "function",
		"inputs": [
			{"name": "from", "type": "address"},
			{"name": "to", "type": "address"},
			{"name": "ids", "type": "uint256[]"},
			{"name": "amounts", "type": "uint256[]"},
			{"name": "data", "type": "bytes"}
		]
	}
]`

func DecodeAdvance(advance *rollup.Advance) (interface{}, InputType, error) {
	sender := advance.MsgSender

	switch sender {
	case EtherPortal:
		deposit, err := DecodeEtherDeposit(advance.Payload)
		return deposit, InputTypeEtherDeposit, err

	case ERC20Portal:
		deposit, err := DecodeERC20Deposit(advance.Payload)
		return deposit, InputTypeERC20Deposit, err

	case ERC721Portal:
		deposit, err := DecodeERC721Deposit(advance.Payload)
		return deposit, InputTypeERC721Deposit, err

	case ERC1155SinglePortal:
		deposit, err := DecodeERC1155SingleDeposit(advance.Payload)
		return deposit, InputTypeERC1155SingleDeposit, err

	case ERC1155BatchPortal:
		deposit, err := DecodeERC1155BatchDeposit(advance.Payload)
		return deposit, InputTypeERC1155BatchDeposit, err

	default:
		return decodeBySelector(advance.Payload)
	}
}

func decodeBySelector(payload []byte) (interface{}, InputType, error) {
	if len(payload) < 4 {
		return nil, InputTypeNone, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])

	switch selector {
	case SelectorWithdrawEther:
		withdrawal, err := DecodeEtherWithdrawal(payload)
		return withdrawal, InputTypeEtherWithdrawal, err

	case SelectorWithdrawERC20:
		withdrawal, err := DecodeERC20Withdrawal(payload)
		return withdrawal, InputTypeERC20Withdrawal, err

	case SelectorWithdrawERC721:
		withdrawal, err := DecodeERC721Withdrawal(payload)
		return withdrawal, InputTypeERC721Withdrawal, err

	case SelectorWithdrawERC1155Single:
		withdrawal, err := DecodeERC1155SingleWithdrawal(payload)
		return withdrawal, InputTypeERC1155SingleWithdrawal, err

	case SelectorWithdrawERC1155Batch:
		withdrawal, err := DecodeERC1155BatchWithdrawal(payload)
		return withdrawal, InputTypeERC1155BatchWithdrawal, err

	case SelectorTransferEther:
		transfer, err := DecodeEtherTransfer(payload)
		return transfer, InputTypeEtherTransfer, err

	case SelectorTransferERC20:
		transfer, err := DecodeERC20Transfer(payload)
		return transfer, InputTypeERC20Transfer, err

	case SelectorTransferERC721:
		transfer, err := DecodeERC721Transfer(payload)
		return transfer, InputTypeERC721Transfer, err

	case SelectorTransferERC1155Single:
		transfer, err := DecodeERC1155SingleTransfer(payload)
		return transfer, InputTypeERC1155SingleTransfer, err

	case SelectorTransferERC1155Batch:
		transfer, err := DecodeERC1155BatchTransfer(payload)
		return transfer, InputTypeERC1155BatchTransfer, err

	default:
		return nil, InputTypeNone, ErrUnknownInputType
	}
}

func DecodeInspect(inspect *rollup.Inspect) (interface{}, InputType, error) {
	var req struct {
		Method string   `json:"method"`
		Params []string `json:"params"`
	}

	if err := json.Unmarshal(inspect.Payload, &req); err != nil {
		return nil, InputTypeNone, ErrMalformedInput
	}

	switch req.Method {
	case "ledger_getBalance":
		return decodeBalanceJSON(req.Params)
	case "ledger_getTotalSupply":
		return decodeSupplyJSON(req.Params)
	default:
		return nil, InputTypeNone, ErrUnknownInputType
	}
}

func decodeBalanceJSON(params []string) (*BalanceQuery, InputType, error) {
	query := &BalanceQuery{}

	if len(params) == 0 || len(params) > 4 {
		return nil, InputTypeNone, ErrMalformedInput
	}

	account := common.HexToHash(params[0])
	query.Account = account

	if len(params) == 1 {
		return query, InputTypeBalanceAccount, nil
	}

	query.Token = common.HexToAddress(params[1])

	if len(params) == 2 {
		return query, InputTypeBalanceAccountTokenAddress, nil
	}

	tokenID, ok := new(big.Int).SetString(params[2], 0)
	if !ok {
		return nil, InputTypeNone, ErrMalformedInput
	}
	query.TokenID = tokenID

	if len(params) == 4 {
		query.ExecLayerData = []byte(params[3])
	}

	return query, InputTypeBalanceAccountTokenAddressID, nil
}

func decodeSupplyJSON(params []string) (*SupplyQuery, InputType, error) {
	query := &SupplyQuery{}

	if len(params) == 0 {
		return query, InputTypeSupply, nil
	}

	if len(params) > 3 {
		return nil, InputTypeNone, ErrMalformedInput
	}

	query.Token = common.HexToAddress(params[0])

	if len(params) == 1 {
		return query, InputTypeSupplyTokenAddress, nil
	}

	tokenID, ok := new(big.Int).SetString(params[1], 0)
	if !ok {
		return nil, InputTypeNone, ErrMalformedInput
	}
	query.TokenID = tokenID

	if len(params) == 3 {
		query.ExecLayerData = []byte(params[2])
	}

	return query, InputTypeSupplyTokenAddressID, nil
}

func DecodeEtherDeposit(payload []byte) (*EtherDeposit, error) {
	if len(payload) < 52 {
		return nil, ErrMalformedInput
	}

	deposit := &EtherDeposit{
		Sender: common.BytesToAddress(payload[0:20]),
		Amount: new(big.Int).SetBytes(payload[20:52]),
	}

	if len(payload) > 52 {
		deposit.ExecLayerData = make([]byte, len(payload)-52)
		copy(deposit.ExecLayerData, payload[52:])
	}

	return deposit, nil
}

func DecodeERC20Deposit(payload []byte) (*ERC20Deposit, error) {
	if len(payload) < 72 {
		return nil, ErrMalformedInput
	}

	deposit := &ERC20Deposit{
		Token:  common.BytesToAddress(payload[0:20]),
		Sender: common.BytesToAddress(payload[20:40]),
		Amount: new(big.Int).SetBytes(payload[40:72]),
	}

	if len(payload) > 72 {
		deposit.ExecLayerData = make([]byte, len(payload)-72)
		copy(deposit.ExecLayerData, payload[72:])
	}

	return deposit, nil
}

func DecodeERC721Deposit(payload []byte) (*ERC721Deposit, error) {
	if len(payload) < 72 {
		return nil, ErrMalformedInput
	}

	deposit := &ERC721Deposit{
		Token:   common.BytesToAddress(payload[0:20]),
		Sender:  common.BytesToAddress(payload[20:40]),
		TokenID: new(big.Int).SetBytes(payload[40:72]),
	}

	if len(payload) > 72 {
		deposit.ExecLayerData = make([]byte, len(payload)-72)
		copy(deposit.ExecLayerData, payload[72:])
	}

	return deposit, nil
}

func DecodeERC1155SingleDeposit(payload []byte) (*ERC1155SingleDeposit, error) {
	if len(payload) < 104 {
		return nil, ErrMalformedInput
	}

	deposit := &ERC1155SingleDeposit{
		Token:   common.BytesToAddress(payload[0:20]),
		Sender:  common.BytesToAddress(payload[20:40]),
		TokenID: new(big.Int).SetBytes(payload[40:72]),
		Amount:  new(big.Int).SetBytes(payload[72:104]),
	}

	if len(payload) > 104 {
		deposit.ExecLayerData = make([]byte, len(payload)-104)
		copy(deposit.ExecLayerData, payload[104:])
	}

	return deposit, nil
}

func DecodeERC1155BatchDeposit(payload []byte) (*ERC1155BatchDeposit, error) {
	if len(payload) < 72 {
		return nil, ErrMalformedInput
	}

	deposit := &ERC1155BatchDeposit{
		Token:  common.BytesToAddress(payload[0:20]),
		Sender: common.BytesToAddress(payload[20:40]),
	}

	offset := 40
	idsOffset := new(big.Int).SetBytes(payload[offset:offset+32]).Uint64() + 40
	offset += 32

	if len(payload) < int(idsOffset)+32 {
		return nil, ErrMalformedInput
	}

	idsCount := new(big.Int).SetBytes(payload[idsOffset : idsOffset+32]).Uint64()
	idsStart := idsOffset + 32

	if len(payload) < int(idsStart)+int(idsCount)*32 {
		return nil, ErrMalformedInput
	}

	deposit.TokenIDs = make([]*big.Int, idsCount)
	for i := uint64(0); i < idsCount; i++ {
		start := idsStart + i*32
		deposit.TokenIDs[i] = new(big.Int).SetBytes(payload[start : start+32])
	}

	amountsOffset := new(big.Int).SetBytes(payload[offset:offset+32]).Uint64() + 40
	offset += 32

	if len(payload) < int(amountsOffset)+32 {
		return nil, ErrMalformedInput
	}

	amountsCount := new(big.Int).SetBytes(payload[amountsOffset : amountsOffset+32]).Uint64()
	amountsStart := amountsOffset + 32

	if len(payload) < int(amountsStart)+int(amountsCount)*32 {
		return nil, ErrMalformedInput
	}

	deposit.Amounts = make([]*big.Int, amountsCount)
	for i := uint64(0); i < amountsCount; i++ {
		start := amountsStart + i*32
		deposit.Amounts[i] = new(big.Int).SetBytes(payload[start : start+32])
	}

	baseLayerDataOffset := new(big.Int).SetBytes(payload[offset : offset+32]).Uint64()
	offset += 32

	if len(payload) >= int(baseLayerDataOffset)+32 {
		baseLayerDataLen := new(big.Int).SetBytes(payload[baseLayerDataOffset : baseLayerDataOffset+32]).Uint64()
		baseLayerDataStart := baseLayerDataOffset + 32
		if len(payload) >= int(baseLayerDataStart)+int(baseLayerDataLen) {
			deposit.BaseLayerData = make([]byte, baseLayerDataLen)
			copy(deposit.BaseLayerData, payload[baseLayerDataStart:baseLayerDataStart+baseLayerDataLen])
		}
	}

	execLayerDataOffset := new(big.Int).SetBytes(payload[offset : offset+32]).Uint64()

	if len(payload) >= int(execLayerDataOffset)+32 {
		execLayerDataLen := new(big.Int).SetBytes(payload[execLayerDataOffset : execLayerDataOffset+32]).Uint64()
		execLayerDataStart := execLayerDataOffset + 32
		if len(payload) >= int(execLayerDataStart)+int(execLayerDataLen) {
			deposit.ExecLayerData = make([]byte, execLayerDataLen)
			copy(deposit.ExecLayerData, payload[execLayerDataStart:execLayerDataStart+execLayerDataLen])
		}
	}

	return deposit, nil
}

func DecodeEtherWithdrawal(payload []byte) (*EtherWithdrawal, error) {
	if len(payload) < 36 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorWithdrawEther {
		return nil, ErrInvalidSelector
	}

	withdrawal := &EtherWithdrawal{
		Amount: new(big.Int).SetBytes(payload[4:36]),
	}

	if len(payload) > 36 {
		withdrawal.ExecLayerData = make([]byte, len(payload)-36)
		copy(withdrawal.ExecLayerData, payload[36:])
	}

	return withdrawal, nil
}

func DecodeERC20Withdrawal(payload []byte) (*ERC20Withdrawal, error) {
	if len(payload) < 68 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorWithdrawERC20 {
		return nil, ErrInvalidSelector
	}

	withdrawal := &ERC20Withdrawal{
		Token:  common.BytesToAddress(payload[16:36]),
		Amount: new(big.Int).SetBytes(payload[36:68]),
	}

	if len(payload) > 68 {
		withdrawal.ExecLayerData = make([]byte, len(payload)-68)
		copy(withdrawal.ExecLayerData, payload[68:])
	}

	return withdrawal, nil
}

func DecodeERC721Withdrawal(payload []byte) (*ERC721Withdrawal, error) {
	if len(payload) < 68 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorWithdrawERC721 {
		return nil, ErrInvalidSelector
	}

	withdrawal := &ERC721Withdrawal{
		Token:   common.BytesToAddress(payload[16:36]),
		TokenID: new(big.Int).SetBytes(payload[36:68]),
	}

	if len(payload) > 68 {
		withdrawal.ExecLayerData = make([]byte, len(payload)-68)
		copy(withdrawal.ExecLayerData, payload[68:])
	}

	return withdrawal, nil
}

func DecodeERC1155SingleWithdrawal(payload []byte) (*ERC1155SingleWithdrawal, error) {
	if len(payload) < 100 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorWithdrawERC1155Single {
		return nil, ErrInvalidSelector
	}

	withdrawal := &ERC1155SingleWithdrawal{
		Token:   common.BytesToAddress(payload[16:36]),
		TokenID: new(big.Int).SetBytes(payload[36:68]),
		Amount:  new(big.Int).SetBytes(payload[68:100]),
	}

	if len(payload) > 100 {
		withdrawal.ExecLayerData = make([]byte, len(payload)-100)
		copy(withdrawal.ExecLayerData, payload[100:])
	}

	return withdrawal, nil
}

func DecodeERC1155BatchWithdrawal(payload []byte) (*ERC1155BatchWithdrawal, error) {
	if len(payload) < 68 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorWithdrawERC1155Batch {
		return nil, ErrInvalidSelector
	}

	withdrawal := &ERC1155BatchWithdrawal{
		Token: common.BytesToAddress(payload[16:36]),
	}

	offset := 36
	idsOffset := new(big.Int).SetBytes(payload[offset:offset+32]).Uint64() + 4
	offset += 32

	if len(payload) < int(idsOffset)+32 {
		return nil, ErrMalformedInput
	}

	idsCount := new(big.Int).SetBytes(payload[idsOffset : idsOffset+32]).Uint64()
	idsStart := idsOffset + 32

	if len(payload) < int(idsStart)+int(idsCount)*32 {
		return nil, ErrMalformedInput
	}

	withdrawal.TokenIDs = make([]*big.Int, idsCount)
	for i := uint64(0); i < idsCount; i++ {
		start := idsStart + i*32
		withdrawal.TokenIDs[i] = new(big.Int).SetBytes(payload[start : start+32])
	}

	amountsOffset := new(big.Int).SetBytes(payload[offset:offset+32]).Uint64() + 4
	offset += 32

	if len(payload) < int(amountsOffset)+32 {
		return nil, ErrMalformedInput
	}

	amountsCount := new(big.Int).SetBytes(payload[amountsOffset : amountsOffset+32]).Uint64()
	amountsStart := amountsOffset + 32

	if len(payload) < int(amountsStart)+int(amountsCount)*32 {
		return nil, ErrMalformedInput
	}

	withdrawal.Amounts = make([]*big.Int, amountsCount)
	for i := uint64(0); i < amountsCount; i++ {
		start := amountsStart + i*32
		withdrawal.Amounts[i] = new(big.Int).SetBytes(payload[start : start+32])
	}

	return withdrawal, nil
}

func DecodeEtherTransfer(payload []byte) (*EtherTransfer, error) {
	if len(payload) < 68 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorTransferEther {
		return nil, ErrInvalidSelector
	}

	transfer := &EtherTransfer{
		Receiver: common.BytesToHash(payload[4:36]),
		Amount:   new(big.Int).SetBytes(payload[36:68]),
	}

	if len(payload) > 68 {
		transfer.ExecLayerData = make([]byte, len(payload)-68)
		copy(transfer.ExecLayerData, payload[68:])
	}

	return transfer, nil
}

func DecodeERC20Transfer(payload []byte) (*ERC20Transfer, error) {
	if len(payload) < 100 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorTransferERC20 {
		return nil, ErrInvalidSelector
	}

	transfer := &ERC20Transfer{
		Token:    common.BytesToAddress(payload[16:36]),
		Receiver: common.BytesToHash(payload[36:68]),
		Amount:   new(big.Int).SetBytes(payload[68:100]),
	}

	if len(payload) > 100 {
		transfer.ExecLayerData = make([]byte, len(payload)-100)
		copy(transfer.ExecLayerData, payload[100:])
	}

	return transfer, nil
}

func DecodeERC721Transfer(payload []byte) (*ERC721Transfer, error) {
	if len(payload) < 100 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorTransferERC721 {
		return nil, ErrInvalidSelector
	}

	transfer := &ERC721Transfer{
		Token:    common.BytesToAddress(payload[16:36]),
		Receiver: common.BytesToHash(payload[36:68]),
		TokenID:  new(big.Int).SetBytes(payload[68:100]),
	}

	if len(payload) > 100 {
		transfer.ExecLayerData = make([]byte, len(payload)-100)
		copy(transfer.ExecLayerData, payload[100:])
	}

	return transfer, nil
}

func DecodeERC1155SingleTransfer(payload []byte) (*ERC1155SingleTransfer, error) {
	if len(payload) < 132 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorTransferERC1155Single {
		return nil, ErrInvalidSelector
	}

	transfer := &ERC1155SingleTransfer{
		Token:    common.BytesToAddress(payload[16:36]),
		Receiver: common.BytesToHash(payload[36:68]),
		TokenID:  new(big.Int).SetBytes(payload[68:100]),
		Amount:   new(big.Int).SetBytes(payload[100:132]),
	}

	if len(payload) > 132 {
		transfer.ExecLayerData = make([]byte, len(payload)-132)
		copy(transfer.ExecLayerData, payload[132:])
	}

	return transfer, nil
}

func DecodeERC1155BatchTransfer(payload []byte) (*ERC1155BatchTransfer, error) {
	if len(payload) < 100 {
		return nil, ErrMalformedInput
	}

	selector := binary.BigEndian.Uint32(payload[0:4])
	if selector != SelectorTransferERC1155Batch {
		return nil, ErrInvalidSelector
	}

	transfer := &ERC1155BatchTransfer{
		Token:    common.BytesToAddress(payload[16:36]),
		Receiver: common.BytesToHash(payload[36:68]),
	}

	offset := 68
	idsOffset := new(big.Int).SetBytes(payload[offset:offset+32]).Uint64() + 4
	offset += 32

	if len(payload) < int(idsOffset)+32 {
		return nil, ErrMalformedInput
	}

	idsCount := new(big.Int).SetBytes(payload[idsOffset : idsOffset+32]).Uint64()
	idsStart := idsOffset + 32

	if len(payload) < int(idsStart)+int(idsCount)*32 {
		return nil, ErrMalformedInput
	}

	transfer.TokenIDs = make([]*big.Int, idsCount)
	for i := uint64(0); i < idsCount; i++ {
		start := idsStart + i*32
		transfer.TokenIDs[i] = new(big.Int).SetBytes(payload[start : start+32])
	}

	amountsOffset := new(big.Int).SetBytes(payload[offset:offset+32]).Uint64() + 4

	if len(payload) < int(amountsOffset)+32 {
		return nil, ErrMalformedInput
	}

	amountsCount := new(big.Int).SetBytes(payload[amountsOffset : amountsOffset+32]).Uint64()
	amountsStart := amountsOffset + 32

	if len(payload) < int(amountsStart)+int(amountsCount)*32 {
		return nil, ErrMalformedInput
	}

	transfer.Amounts = make([]*big.Int, amountsCount)
	for i := uint64(0); i < amountsCount; i++ {
		start := amountsStart + i*32
		transfer.Amounts[i] = new(big.Int).SetBytes(payload[start : start+32])
	}

	return transfer, nil
}

func DecodeBalanceQuery(payload []byte) (*BalanceQuery, InputType, error) {
	query := &BalanceQuery{}

	switch {
	case len(payload) == 0:
		return query, InputTypeBalance, nil

	case len(payload) == 32:
		query.Account = common.BytesToHash(payload[0:32])
		return query, InputTypeBalanceAccount, nil

	case len(payload) == 52:
		query.Account = common.BytesToHash(payload[0:32])
		query.Token = common.BytesToAddress(payload[32:52])
		return query, InputTypeBalanceAccountTokenAddress, nil

	case len(payload) >= 84:
		query.Account = common.BytesToHash(payload[0:32])
		query.Token = common.BytesToAddress(payload[32:52])
		query.TokenID = new(big.Int).SetBytes(payload[52:84])
		if len(payload) > 84 {
			query.ExecLayerData = make([]byte, len(payload)-84)
			copy(query.ExecLayerData, payload[84:])
		}
		return query, InputTypeBalanceAccountTokenAddressID, nil

	default:
		return nil, InputTypeNone, ErrMalformedInput
	}
}

func DecodeSupplyQuery(payload []byte) (*SupplyQuery, InputType, error) {
	query := &SupplyQuery{}

	switch {
	case len(payload) == 0:
		return query, InputTypeSupply, nil

	case len(payload) == 20:
		query.Token = common.BytesToAddress(payload[0:20])
		return query, InputTypeSupplyTokenAddress, nil

	case len(payload) >= 52:
		query.Token = common.BytesToAddress(payload[0:20])
		query.TokenID = new(big.Int).SetBytes(payload[20:52])
		if len(payload) > 52 {
			query.ExecLayerData = make([]byte, len(payload)-52)
			copy(query.ExecLayerData, payload[52:])
		}
		return query, InputTypeSupplyTokenAddressID, nil

	default:
		return nil, InputTypeNone, ErrMalformedInput
	}
}

func EncodeEtherVoucher(receiver common.Address, amount *big.Int) *Voucher {
	return &Voucher{
		Destination: receiver,
		Value:       amount,
		Payload:     nil,
	}
}

func EncodeERC20Voucher(token, receiver common.Address, amount *big.Int) (*Voucher, error) {
	payload, err := erc20ABI.Pack("transfer", receiver, amount)
	if err != nil {
		return nil, err
	}
	return &Voucher{
		Destination: token,
		Value:       big.NewInt(0),
		Payload:     payload,
	}, nil
}

func EncodeERC721Voucher(token, appAddress, receiver common.Address, tokenID *big.Int) (*Voucher, error) {
	payload, err := erc721ABI.Pack("safeTransferFrom", appAddress, receiver, tokenID)
	if err != nil {
		return nil, err
	}
	return &Voucher{
		Destination: token,
		Value:       big.NewInt(0),
		Payload:     payload,
	}, nil
}

func EncodeERC1155SingleVoucher(token, appAddress, receiver common.Address, tokenID, amount *big.Int) (*Voucher, error) {
	payload, err := erc1155ABI.Pack("safeTransferFrom", appAddress, receiver, tokenID, amount, []byte{})
	if err != nil {
		return nil, err
	}
	return &Voucher{
		Destination: token,
		Value:       big.NewInt(0),
		Payload:     payload,
	}, nil
}

func EncodeERC1155BatchVoucher(token, appAddress, receiver common.Address, tokenIDs, amounts []*big.Int) (*Voucher, error) {
	payload, err := erc1155ABI.Pack("safeBatchTransferFrom", appAddress, receiver, tokenIDs, amounts, []byte{})
	if err != nil {
		return nil, err
	}
	return &Voucher{
		Destination: token,
		Value:       big.NewInt(0),
		Payload:     payload,
	}, nil
}

func EncodeDelegateCallVoucher(target common.Address, payload []byte) *DelegateCallVoucher {
	return &DelegateCallVoucher{
		Destination: target,
		Payload:     payload,
	}
}
