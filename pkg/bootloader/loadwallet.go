package bootloader

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type LoadWallet struct {
	db      *LiteStore
	ghostIp *ptypes.GhostIp
	table   string
}

func NewLoadWallet(table string, db *LiteStore, ghostIp *ptypes.GhostIp) *LoadWallet {
	return &LoadWallet{db: db, ghostIp: ghostIp, table: table}
}

func (loadWallet *LoadWallet) CreateWallet(nickname string, password string) *gcrypto.Wallet {
	newGhostAddress := gcrypto.GenerateKeyPair()
	return gcrypto.NewWallet(nickname, newGhostAddress, loadWallet.ghostIp, nil)
}

func (loadWallet *LoadWallet) OpenWallet(nickname string, password string) (*gcrypto.Wallet, error) {
	der, err := loadWallet.db.SelectEntry(loadWallet.table, []byte(nickname))
	if err != nil {
		return nil, err
	}
	ghostAddr := &gcrypto.GhostAddress{}
	ghostAddr.PrivateKeyDeserialize(der)

	return gcrypto.NewWallet(nickname, ghostAddr, loadWallet.ghostIp, nil), nil
}

func (loadWallet *LoadWallet) SaveWallet(w *gcrypto.Wallet, password string) {
	nickname := w.GetNickname()
	privateKey := w.GetGhostAddress().PrivateKeySerialize()
	// TODO: need to encryption by passwd
	loadWallet.db.SaveEntry(loadWallet.table, []byte(nickname), privateKey)
}
