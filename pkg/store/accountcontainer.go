package store

type AccountContainer struct {
	liteStore *LiteStore
}

func NewBcAccountContainer(liteStore *LiteStore) *AccountContainer {
	return &AccountContainer{
		liteStore: liteStore,
	}
}

func (account *AccountContainer) AddBcAccount(nickname []byte, txId []byte) bool {
	if v, err := account.liteStore.SelectEntry(DefaultNickInBlockChainTable, nickname); v == nil && err == nil {
		account.liteStore.SaveEntry(DefaultNickInBlockChainTable, nickname, txId)
		return true
	}
	return false
}

func (account *AccountContainer) GetBcAccount(nickname []byte) []byte {
	if v, err := account.liteStore.SelectEntry(DefaultNickInBlockChainTable, nickname); v != nil && err == nil {
		return v
	}
	return nil
}
