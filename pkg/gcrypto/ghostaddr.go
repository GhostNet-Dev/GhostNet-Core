package gcrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"

	"github.com/btcsuite/btcutil/base58"
)

type GhostAddress struct {
	PriKey    ecdsa.PrivateKey
	pubKey    []byte
	pubKey160 []byte
}

func GenerateKeyPair() *GhostAddress {
	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)

	if err != nil {
		fmt.Println(err)
	}

	//publicKey := append(privatekey.PublicKey.X.Bytes(), privatekey.PublicKey.Y.Bytes()...)
	//https://stackoverflow.com/questions/73721296/go-language-ecdsa-verify-the-valid-signature-to-invalid
	publicKey := elliptic.MarshalCompressed(privatekey.Curve, privatekey.PublicKey.X, privatekey.PublicKey.Y)
	return &GhostAddress{
		PriKey: *privatekey,
		pubKey: publicKey,
	}
}

func (ghostAddr *GhostAddress) GetPubAddress() string {
	pubKey := ghostAddr.pubKey
	publicSHA256 := sha256.Sum256(pubKey) // Public key를 SHA-256으로 해싱

	// RIPEMD-160으로 다시 해싱
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return base58.CheckEncode(publicRIPEMD160, 0)
}

func (ghostAddr *GhostAddress) GetSignPubKey() []byte {
	return ghostAddr.pubKey
}

func TranslateSigPubTo160PubKey(sigPubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(sigPubKey) // Public key를 SHA-256으로 해싱

	// RIPEMD-160으로 다시 해싱
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	return RIPEMD160Hasher.Sum(nil)
}

func (ghostAddr *GhostAddress) Get160PubKey() []byte {
	if ghostAddr.pubKey160 == nil {
		pubKey := ghostAddr.pubKey
		publicSHA256 := sha256.Sum256(pubKey) // Public key를 SHA-256으로 해싱

		// RIPEMD-160으로 다시 해싱
		RIPEMD160Hasher := ripemd160.New()
		_, err := RIPEMD160Hasher.Write(publicSHA256[:])
		if err != nil {
			log.Panic(err)
		}

		ghostAddr.pubKey160 = RIPEMD160Hasher.Sum(nil)
	}
	return ghostAddr.pubKey160
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
