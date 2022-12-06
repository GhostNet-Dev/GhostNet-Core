package gcrypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

type GhostAddress struct {
	PriKey ecdsa.PrivateKey
	PubKey []byte
}

func GenerateKeyPair() *GhostAddress {
	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)

	if err != nil {
		fmt.Println(err)
	}

	publicKey := append(privatekey.PublicKey.X.Bytes(), privatekey.PublicKey.Y.Bytes()...)
	return &GhostAddress{
		*privatekey,
		publicKey,
	}
}

func (ghostAddr *GhostAddress) GetPubAddress() string {
	pubKey := ghostAddr.PubKey
	publicSHA256 := sha256.Sum256(pubKey) // Public key를 SHA-256으로 해싱

	// RIPEMD-160으로 다시 해싱
	RIPEMD160Hasher := crypto.RIPEMD160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return base58.CheckEncode(publicRIPEMD160, 0)
}

func (ghostAddr *GhostAddress) PrivateKeySerialize() []byte {
	privateByte, err := x509.MarshalECPrivateKey(&ghostAddr.PriKey)
	if err != nil {
		log.Panic(err)
	}
	return privateByte
}

func (ghostAddr *GhostAddress) PrivateKeyDeserialize(der []byte) {
	priKey, err := x509.ParseECPrivateKey(der)
	if err != nil {
		log.Panic(err)
	}
	ghostAddr.PriKey = *priKey
}
