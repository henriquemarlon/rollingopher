package ledger

type AssetID uint64

type InternalAccountID uint64

type AssetType int

const (
	AssetTypeID AssetType = iota
	AssetTypeTokenAddress
	AssetTypeTokenAddressID
)

type AccountType int

const (
	AccountTypeID AccountType = iota
	AccountTypeWalletAddress
	AccountTypeAccountID
)

type RetrieveOperation int

const (
	RetrieveOperationFind RetrieveOperation = iota
	RetrieveOperationCreate
	RetrieveOperationFindOrCreate
)
