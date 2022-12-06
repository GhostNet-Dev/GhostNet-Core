package gcrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

type SignaturePackage struct {
	PubKey    []byte
	Signature []byte
	Text      []byte
	R         *big.Int
	S         *big.Int
}

func Signer(text []byte, ghostAddr *GhostAddress) *SignaturePackage {
	// Sign ecdsa style
	//var h hash.Hash
	//h = md5.New()
	//r := big.NewInt(0)
	//s := big.NewInt(0)

	//io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
	//signhash := h.Sum(nil)
	privateKey := &ghostAddr.PriKey
	r, s, serr := ecdsa.Sign(rand.Reader, privateKey, text)
	if serr != nil {
		fmt.Println(serr)
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	return &SignaturePackage{
		ghostAddr.PubKey, signature, text, r, s,
	}
}

func SignVerify(sig *SignaturePackage) bool {
	pubkeyCurve := elliptic.P256()
	sigLen := len(sig.Signature)
	keyLen := len(sig.PubKey)
	var r, s big.Int
	var x, y big.Int

	// Signature is a pair of numbers.
	r.SetBytes(sig.Signature[:(sigLen / 2)])
	s.SetBytes(sig.Signature[(sigLen / 2):])

	// PublicKey is a pair of coordinates.
	x.SetBytes(sig.PubKey[:(keyLen / 2)])
	y.SetBytes(sig.PubKey[(keyLen / 2):])

	rawPublicKey := &ecdsa.PublicKey{Curve: pubkeyCurve, X: &x, Y: &y}
	// Verify
	return ecdsa.Verify(rawPublicKey, sig.Text, sig.R, sig.S)
}
