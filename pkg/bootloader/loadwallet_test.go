package bootloader

import (
	"bytes"
	"log"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/stretchr/testify/assert"
)

var (
	wallet   = NewLoadWallet(Tables[1], db, ghostIp)
	username = "User"
	password = "pass"
)

func TestOpenWallet(t *testing.T) {
	err := db.OpenStore()
	if err != nil {
		log.Fatal(err)
	}
	cipherText := grpc.PasswordToSha256(password)
	w := wallet.CreateWallet(username, cipherText)
	wallet.SaveWallet(w, cipherText)
	new, err := wallet.OpenWallet(username, cipherText)
	if err != nil {
		log.Fatal(err)
	}
	compare := w.GetPubAddress() == new.GetPubAddress()
	bytesCompare := bytes.Compare(w.GetGhostAddress().PrivateKeySerialize(),
		new.GetGhostAddress().PrivateKeySerialize())
	assert.Equal(t, true, compare, "not equal after decryption")
	assert.Equal(t, 0, bytesCompare, "privKey is not equal after decryption")
}

func TestEncryptDecrypt(t *testing.T) {
	cipherText := grpc.PasswordToSha256(password)
	keyPair := gcrypto.GenerateKeyPair()
	cipherPriv := wallet.Encryption(cipherText, keyPair.PrivateKeySerialize())
	privateKey := wallet.Decryption(cipherText, cipherPriv)
	comp := bytes.Compare(keyPair.PrivateKeySerialize(), privateKey)
	cipherComp := bytes.Compare(cipherText, grpc.PasswordToSha256(password))
	assert.Equal(t, 0, comp, "key not equal after decryption")
	assert.Equal(t, 0, cipherComp, "sha256 not working")
}
