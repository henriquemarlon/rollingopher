package main

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/rollingopher/pkg/ledger"
	"github.com/henriquemarlon/rollingopher/pkg/parser"
	"github.com/henriquemarlon/rollingopher/pkg/rollup"
)

var etherAssetID ledger.AssetID

func handleAdvance(r *rollup.Rollup, l *ledger.Ledger) bool {
	advance, err := r.ReadAdvanceState()
	if err != nil {
		log.Printf("[handling-assets] failed to read advance: %v", err)
		return false
	}

	msgSender := advance.MsgSender
	log.Printf("[handling-assets] Received advance from %s", msgSender.Hex())

	decodedInput, inputType, err := parser.DecodeAdvance(advance)
	if err != nil {
		log.Printf("[handling-assets] failed to decode: %v", err)
		return false
	}

	switch inputType {
	// Ether
	case parser.InputTypeEtherDeposit:
		d := decodedInput.(*parser.EtherDeposit)
		accountID, _ := l.RetrieveAccountByAddress(d.Sender, ledger.RetrieveOperationFindOrCreate)
		l.Deposit(etherAssetID, accountID, d.Amount)
		log.Printf("[handling-assets] %s deposited %s ether", d.Sender.Hex(), d.Amount)
		return true

	case parser.InputTypeEtherWithdrawal:
		w := decodedInput.(*parser.EtherWithdrawal)
		accountID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		l.Withdraw(etherAssetID, accountID, w.Amount)
		v := parser.EncodeEtherVoucher(msgSender, w.Amount)
		r.EmitVoucher(v.Destination, v.Value, v.Payload)
		log.Printf("[handling-assets] %s withdrew %s ether", msgSender.Hex(), w.Amount)
		return true

	case parser.InputTypeEtherTransfer:
		t := decodedInput.(*parser.EtherTransfer)
		fromID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		toID, _ := l.RetrieveAccountByID(t.Receiver, ledger.RetrieveOperationFindOrCreate)
		l.Transfer(etherAssetID, fromID, toID, t.Amount)
		log.Printf("[handling-assets] %s transferred %s ether to %s", msgSender.Hex(), t.Amount, t.Receiver.Hex())
		return true

	// ERC20
	case parser.InputTypeERC20Deposit:
		d := decodedInput.(*parser.ERC20Deposit)
		assetID, _ := l.RetrieveAsset(d.Token, nil, ledger.AssetTypeTokenAddress, ledger.RetrieveOperationFindOrCreate)
		accountID, _ := l.RetrieveAccountByAddress(d.Sender, ledger.RetrieveOperationFindOrCreate)
		l.Deposit(assetID, accountID, d.Amount)
		log.Printf("[handling-assets] %s deposited %s of %s", d.Sender.Hex(), d.Amount, d.Token.Hex())
		return true

	case parser.InputTypeERC20Withdrawal:
		w := decodedInput.(*parser.ERC20Withdrawal)
		assetID, _ := l.RetrieveAsset(w.Token, nil, ledger.AssetTypeTokenAddress, ledger.RetrieveOperationFind)
		accountID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		l.Withdraw(assetID, accountID, w.Amount)
		v, _ := parser.EncodeERC20Voucher(w.Token, msgSender, w.Amount)
		r.EmitVoucher(v.Destination, v.Value, v.Payload)
		log.Printf("[handling-assets] %s withdrew %s of %s", msgSender.Hex(), w.Amount, w.Token.Hex())
		return true

	case parser.InputTypeERC20Transfer:
		t := decodedInput.(*parser.ERC20Transfer)
		assetID, _ := l.RetrieveAsset(t.Token, nil, ledger.AssetTypeTokenAddress, ledger.RetrieveOperationFind)
		fromID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		toID, _ := l.RetrieveAccountByID(t.Receiver, ledger.RetrieveOperationFindOrCreate)
		l.Transfer(assetID, fromID, toID, t.Amount)
		log.Printf("[handling-assets] %s transferred %s of %s to %s", msgSender.Hex(), t.Amount, t.Token.Hex(), t.Receiver.Hex())
		return true

	// ERC721
	case parser.InputTypeERC721Deposit:
		d := decodedInput.(*parser.ERC721Deposit)
		assetID, _ := l.RetrieveAsset(d.Token, d.TokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFindOrCreate)
		accountID, _ := l.RetrieveAccountByAddress(d.Sender, ledger.RetrieveOperationFindOrCreate)
		l.Deposit(assetID, accountID, big.NewInt(1))
		log.Printf("[handling-assets] %s deposited ERC721 %s #%s", d.Sender.Hex(), d.Token.Hex(), d.TokenID)
		return true

	case parser.InputTypeERC721Withdrawal:
		w := decodedInput.(*parser.ERC721Withdrawal)
		assetID, _ := l.RetrieveAsset(w.Token, w.TokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFind)
		accountID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		l.Withdraw(assetID, accountID, big.NewInt(1))
		v, _ := parser.EncodeERC721Voucher(w.Token, advance.AppContract, msgSender, w.TokenID)
		r.EmitVoucher(v.Destination, v.Value, v.Payload)
		log.Printf("[handling-assets] %s withdrew ERC721 %s #%s", msgSender.Hex(), w.Token.Hex(), w.TokenID)
		return true

	case parser.InputTypeERC721Transfer:
		t := decodedInput.(*parser.ERC721Transfer)
		assetID, _ := l.RetrieveAsset(t.Token, t.TokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFind)
		fromID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		toID, _ := l.RetrieveAccountByID(t.Receiver, ledger.RetrieveOperationFindOrCreate)
		l.Transfer(assetID, fromID, toID, big.NewInt(1))
		log.Printf("[handling-assets] %s transferred ERC721 %s #%s to %s", msgSender.Hex(), t.Token.Hex(), t.TokenID, t.Receiver.Hex())
		return true

	// ERC1155 Single
	case parser.InputTypeERC1155SingleDeposit:
		d := decodedInput.(*parser.ERC1155SingleDeposit)
		assetID, _ := l.RetrieveAsset(d.Token, d.TokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFindOrCreate)
		accountID, _ := l.RetrieveAccountByAddress(d.Sender, ledger.RetrieveOperationFindOrCreate)
		l.Deposit(assetID, accountID, d.Amount)
		log.Printf("[handling-assets] %s deposited %s of ERC1155 %s #%s", d.Sender.Hex(), d.Amount, d.Token.Hex(), d.TokenID)
		return true

	case parser.InputTypeERC1155SingleWithdrawal:
		w := decodedInput.(*parser.ERC1155SingleWithdrawal)
		assetID, _ := l.RetrieveAsset(w.Token, w.TokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFind)
		accountID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		l.Withdraw(assetID, accountID, w.Amount)
		v, _ := parser.EncodeERC1155SingleVoucher(w.Token, advance.AppContract, msgSender, w.TokenID, w.Amount)
		r.EmitVoucher(v.Destination, v.Value, v.Payload)
		log.Printf("[handling-assets] %s withdrew %s of ERC1155 %s #%s", msgSender.Hex(), w.Amount, w.Token.Hex(), w.TokenID)
		return true

	case parser.InputTypeERC1155SingleTransfer:
		t := decodedInput.(*parser.ERC1155SingleTransfer)
		assetID, _ := l.RetrieveAsset(t.Token, t.TokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFind)
		fromID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		toID, _ := l.RetrieveAccountByID(t.Receiver, ledger.RetrieveOperationFindOrCreate)
		l.Transfer(assetID, fromID, toID, t.Amount)
		log.Printf("[handling-assets] %s transferred %s of ERC1155 %s #%s to %s", msgSender.Hex(), t.Amount, t.Token.Hex(), t.TokenID, t.Receiver.Hex())
		return true

	// ERC1155 Batch
	case parser.InputTypeERC1155BatchDeposit:
		d := decodedInput.(*parser.ERC1155BatchDeposit)
		accountID, _ := l.RetrieveAccountByAddress(d.Sender, ledger.RetrieveOperationFindOrCreate)
		for i, tokenID := range d.TokenIDs {
			assetID, _ := l.RetrieveAsset(d.Token, tokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFindOrCreate)
			l.Deposit(assetID, accountID, d.Amounts[i])
		}
		log.Printf("[handling-assets] %s deposited ERC1155 batch from %s", d.Sender.Hex(), d.Token.Hex())
		return true

	case parser.InputTypeERC1155BatchWithdrawal:
		w := decodedInput.(*parser.ERC1155BatchWithdrawal)
		accountID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		for i, tokenID := range w.TokenIDs {
			assetID, _ := l.RetrieveAsset(w.Token, tokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFind)
			l.Withdraw(assetID, accountID, w.Amounts[i])
		}
		v, _ := parser.EncodeERC1155BatchVoucher(w.Token, advance.AppContract, msgSender, w.TokenIDs, w.Amounts)
		r.EmitVoucher(v.Destination, v.Value, v.Payload)
		log.Printf("[handling-assets] %s withdrew ERC1155 batch from %s", msgSender.Hex(), w.Token.Hex())
		return true

	case parser.InputTypeERC1155BatchTransfer:
		t := decodedInput.(*parser.ERC1155BatchTransfer)
		fromID, _ := l.RetrieveAccountByAddress(msgSender, ledger.RetrieveOperationFind)
		toID, _ := l.RetrieveAccountByID(t.Receiver, ledger.RetrieveOperationFindOrCreate)
		for i, tokenID := range t.TokenIDs {
			assetID, _ := l.RetrieveAsset(t.Token, tokenID, ledger.AssetTypeTokenAddressID, ledger.RetrieveOperationFind)
			l.Transfer(assetID, fromID, toID, t.Amounts[i])
		}
		log.Printf("[handling-assets] %s transferred ERC1155 batch of %s to %s", msgSender.Hex(), t.Token.Hex(), t.Receiver.Hex())
		return true

	default:
		log.Printf("[handling-assets] unknown input type: %v", inputType)
		return false
	}
}

func handleInspect(r *rollup.Rollup, l *ledger.Ledger) bool {
	inspect, err := r.ReadInspectState()
	if err != nil {
		log.Printf("[handling-assets] failed to read inspect: %v", err)
		return false
	}

	decoded, inputType, err := parser.DecodeInspect(inspect)
	if err != nil {
		log.Printf("[handling-assets] failed to decode inspect: %v", err)
		return false
	}

	switch inputType {
	case parser.InputTypeBalance, parser.InputTypeBalanceAccount, parser.InputTypeBalanceAccountTokenAddress, parser.InputTypeBalanceAccountTokenAddressID:
		query := decoded.(*parser.BalanceQuery)
		accountID, _ := l.RetrieveAccountByID(query.Account, ledger.RetrieveOperationFind)
		assetID := etherAssetID
		if query.Token != (common.Address{}) {
			assetType := ledger.AssetTypeTokenAddress
			if inputType == parser.InputTypeBalanceAccountTokenAddressID {
				assetType = ledger.AssetTypeTokenAddressID
			}
			assetID, _ = l.RetrieveAsset(query.Token, query.TokenID, assetType, ledger.RetrieveOperationFind)
		}
		balance, _ := l.GetBalance(assetID, accountID)
		report := make([]byte, 32)
		balance.FillBytes(report)
		r.EmitReport(report)
		log.Printf("[handling-assets] balance: %s (type %v)", balance, inputType)
		return true

	case parser.InputTypeSupply, parser.InputTypeSupplyTokenAddress, parser.InputTypeSupplyTokenAddressID:
		query := decoded.(*parser.SupplyQuery)
		assetID := etherAssetID
		if query.Token != (common.Address{}) {
			assetType := ledger.AssetTypeTokenAddress
			if inputType == parser.InputTypeSupplyTokenAddressID {
				assetType = ledger.AssetTypeTokenAddressID
			}
			assetID, _ = l.RetrieveAsset(query.Token, query.TokenID, assetType, ledger.RetrieveOperationFind)
		}
		supply, _ := l.GetTotalSupply(assetID)
		report := make([]byte, 32)
		supply.FillBytes(report)
		r.EmitReport(report)
		log.Printf("[handling-assets] supply: %s (type %v)", supply, inputType)
		return true

	default:
		log.Printf("[handling-assets] unknown inspect type: %v", inputType)
		return false
	}
}

func main() {
	r, _ := rollup.New()
	defer r.Close()

	l, _ := ledger.New()
	defer l.Close()

	etherAssetID, _ = l.RetrieveAsset(common.Address{}, nil, ledger.AssetTypeID, ledger.RetrieveOperationFindOrCreate)

	accept := true
	for {
		reqType, _, _ := r.Finish(accept)
		switch reqType {
		case rollup.RequestTypeAdvance:
			accept = handleAdvance(r, l)
		case rollup.RequestTypeInspect:
			accept = handleInspect(r, l)
		}
	}
}
