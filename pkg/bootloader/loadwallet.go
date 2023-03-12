package bootloader

import (
	"errors"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"google.golang.org/protobuf/proto"
)

type LoadWallet struct {
	db      *store.LiteStore
	ghostIp *ptypes.GhostIp
	table   string
}

func NewLoadWallet(table string, db *store.LiteStore, ghostIp *ptypes.GhostIp) *LoadWallet {
	return &LoadWallet{db: db, ghostIp: ghostIp, table: table}
}

func (loadWallet *LoadWallet) CreateWallet(nickname string, password []byte) *gcrypto.Wallet {
	newGhostAddress := gcrypto.GenerateKeyPair()
	return gcrypto.NewWallet(nickname, newGhostAddress, loadWallet.ghostIp, nil)
}

func (loadWallet *LoadWallet) OpenWallet(nickname string, password []byte) (*gcrypto.Wallet, error) {
	cipherPivateKey, err := loadWallet.db.SelectEntry(loadWallet.table, []byte(nickname))
	if err != nil || cipherPivateKey == nil {
		return nil, err
	}
	der := gcrypto.Decryption(password, cipherPivateKey)
	if der == nil {
		return nil, errors.New("password is wrong")
	}
	keyPair := &ptypes.KeyPair{}
	if err := proto.Unmarshal(der, keyPair); err != nil {
		log.Fatal(err)
	}

	ghostAddr := &gcrypto.GhostAddress{}
	ghostAddr.PrivateKeyDeserialize(keyPair.PrivateKey)
	if ghostAddr.GetPubAddress() != keyPair.PubKey {
		return nil, errors.New("password is wrong")
	}

	return gcrypto.NewWallet(nickname, ghostAddr, loadWallet.ghostIp, nil), nil
}

func (loadWallet *LoadWallet) SaveWallet(w *gcrypto.Wallet, password []byte) {
	nickname := w.GetNickname()
	privateKey := w.GetGhostAddress().PrivateKeySerialize()
	keyPair := &ptypes.KeyPair{
		PubKey:     w.GetGhostAddress().GetPubAddress(),
		PrivateKey: privateKey,
	}
	data, err := proto.Marshal(keyPair)
	if err != nil {
		log.Fatal(err)
	}
	cipherPivateKey := gcrypto.Encryption(password, data)
	loadWallet.db.SaveEntry(loadWallet.table, []byte(nickname), cipherPivateKey)
}

func (loadWallet *LoadWallet) GetWalletList() (nicknames []string) {
	if k, _, err := loadWallet.db.LoadEntry(loadWallet.table); err == nil {
		for _, name := range k {
			nicknames = append(nicknames, string(name))
		}
		return nicknames
	} else {
		log.Print(err)
	}
	return nil
}
