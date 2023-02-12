package bootloader

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"google.golang.org/protobuf/proto"
)

type LoadWallet struct {
	db      *LiteStore
	ghostIp *ptypes.GhostIp
	table   string
}

func NewLoadWallet(table string, db *LiteStore, ghostIp *ptypes.GhostIp) *LoadWallet {
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
	der := loadWallet.Decryption(password, cipherPivateKey)
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
	cipherPivateKey := loadWallet.Encryption(password, data)
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

func (loadWallet *LoadWallet) Encryption(key []byte, privateKey []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	return encrypt(block, privateKey)
}

func (loadWallet *LoadWallet) Decryption(key []byte, cipherText []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decrypt(block, cipherText)
}

func encrypt(b cipher.Block, plaintext []byte) []byte {
	gcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func decrypt(b cipher.Block, ciphertext []byte) []byte {
	gcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Fatal(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Fatal("len(ciphertext) < nonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return plaintext
}
