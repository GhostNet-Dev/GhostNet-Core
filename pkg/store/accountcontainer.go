package store

type AccountContainer struct {
	liteStore *LiteStore
}

func NewAccountContainer(liteStore *LiteStore) *AccountContainer {
	return &AccountContainer{
		liteStore: liteStore,
	}
}

func (account *AccountContainer) AddAccount(nickname []byte, txId []byte) bool {
	if v, err := account.liteStore.SelectEntry(DefaultNickTable, nickname); v == nil && err == nil {
		account.liteStore.SaveEntry(DefaultNickTable, nickname, txId)
		return true
	}
	return false
}

func (account *AccountContainer) GetAccount(nickname []byte) []byte {
	if v, err := account.liteStore.SelectEntry(DefaultNickTable, nickname); v != nil && err == nil {
		return v
	}
	return nil
}