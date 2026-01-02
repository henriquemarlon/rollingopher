//go:build !riscv64

package ledger

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type assetKey struct {
	tokenAddress common.Address
	tokenID      string
}

type accountKey struct {
	address   common.Address
	accountID common.Hash
	keyType   AccountType
}

type Ledger struct {
	mu sync.Mutex

	nextAssetID   AssetID
	nextAccountID InternalAccountID

	assets   map[assetKey]AssetID
	accounts map[accountKey]InternalAccountID

	balances map[AssetID]map[InternalAccountID]*big.Int
	supplies map[AssetID]*big.Int
}

func New() (*Ledger, error) {
	return &Ledger{
		nextAssetID:   1,
		nextAccountID: 1,
		assets:        make(map[assetKey]AssetID),
		accounts:      make(map[accountKey]InternalAccountID),
		balances:      make(map[AssetID]map[InternalAccountID]*big.Int),
		supplies:      make(map[AssetID]*big.Int),
	}, nil
}

func (l *Ledger) Close() error {
	return nil
}

func (l *Ledger) Reset() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.nextAssetID = 1
	l.nextAccountID = 1
	l.assets = make(map[assetKey]AssetID)
	l.accounts = make(map[accountKey]InternalAccountID)
	l.balances = make(map[AssetID]map[InternalAccountID]*big.Int)
	l.supplies = make(map[AssetID]*big.Int)
	return nil
}

func (l *Ledger) RetrieveAsset(tokenAddress common.Address, tokenID *big.Int, assetType AssetType, op RetrieveOperation) (AssetID, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	tokenIDStr := ""
	if tokenID != nil {
		tokenIDStr = tokenID.String()
	}

	key := assetKey{
		tokenAddress: tokenAddress,
		tokenID:      tokenIDStr,
	}

	if id, exists := l.assets[key]; exists {
		return id, nil
	}

	if op == RetrieveOperationFind {
		return 0, ErrAssetNotFound
	}

	id := l.nextAssetID
	l.nextAssetID++
	l.assets[key] = id
	l.balances[id] = make(map[InternalAccountID]*big.Int)
	l.supplies[id] = big.NewInt(0)
	return id, nil
}

func (l *Ledger) RetrieveAccountByAddress(address common.Address, op RetrieveOperation) (InternalAccountID, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := accountKey{
		address: address,
		keyType: AccountTypeWalletAddress,
	}

	if id, exists := l.accounts[key]; exists {
		return id, nil
	}

	if op == RetrieveOperationFind {
		return 0, ErrAccountNotFound
	}

	id := l.nextAccountID
	l.nextAccountID++
	l.accounts[key] = id
	return id, nil
}

func (l *Ledger) RetrieveAccountByID(accountID common.Hash, op RetrieveOperation) (InternalAccountID, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := accountKey{
		accountID: accountID,
		keyType:   AccountTypeAccountID,
	}

	if id, exists := l.accounts[key]; exists {
		return id, nil
	}

	if op == RetrieveOperationFind {
		return 0, ErrAccountNotFound
	}

	id := l.nextAccountID
	l.nextAccountID++
	l.accounts[key] = id
	return id, nil
}

func (l *Ledger) Deposit(assetID AssetID, accountID InternalAccountID, amount *big.Int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.balances[assetID]; !exists {
		return ErrAssetNotFound
	}

	balance, exists := l.balances[assetID][accountID]
	if !exists {
		balance = big.NewInt(0)
	}

	newBalance := new(big.Int).Add(balance, amount)
	l.balances[assetID][accountID] = newBalance

	l.supplies[assetID] = new(big.Int).Add(l.supplies[assetID], amount)
	return nil
}

func (l *Ledger) Withdraw(assetID AssetID, accountID InternalAccountID, amount *big.Int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.balances[assetID]; !exists {
		return ErrAssetNotFound
	}

	balance, exists := l.balances[assetID][accountID]
	if !exists {
		balance = big.NewInt(0)
	}

	if balance.Cmp(amount) < 0 {
		return ErrInsufficientFunds
	}

	newBalance := new(big.Int).Sub(balance, amount)
	l.balances[assetID][accountID] = newBalance

	l.supplies[assetID] = new(big.Int).Sub(l.supplies[assetID], amount)
	return nil
}

func (l *Ledger) Transfer(assetID AssetID, from, to InternalAccountID, amount *big.Int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.balances[assetID]; !exists {
		return ErrAssetNotFound
	}

	fromBalance, exists := l.balances[assetID][from]
	if !exists {
		fromBalance = big.NewInt(0)
	}

	if fromBalance.Cmp(amount) < 0 {
		return ErrInsufficientFunds
	}

	toBalance, exists := l.balances[assetID][to]
	if !exists {
		toBalance = big.NewInt(0)
	}

	l.balances[assetID][from] = new(big.Int).Sub(fromBalance, amount)
	l.balances[assetID][to] = new(big.Int).Add(toBalance, amount)
	return nil
}

func (l *Ledger) GetBalance(assetID AssetID, accountID InternalAccountID) (*big.Int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.balances[assetID]; !exists {
		return nil, ErrAssetNotFound
	}

	balance, exists := l.balances[assetID][accountID]
	if !exists {
		return big.NewInt(0), nil
	}

	return new(big.Int).Set(balance), nil
}

func (l *Ledger) GetTotalSupply(assetID AssetID) (*big.Int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	supply, exists := l.supplies[assetID]
	if !exists {
		return nil, ErrAssetNotFound
	}

	return new(big.Int).Set(supply), nil
}
