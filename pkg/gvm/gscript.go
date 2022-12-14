package gvm

import (
	"bytes"
	"encoding/binary"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GScript struct {
}

func NewGScript() *GScript {
	return &GScript{}
}

func (gScript *GScript) MakeBlockSignature(block *types.GhostNetBlock, ghostAddr *gcrypto.GhostAddress) {
	block.Header.BlockSignature = types.SigHash{}
	block.Header.SignatureSize = uint32(block.Header.BlockSignature.Size())

	sig := makeSignature(block.Header.SerializeToByte(), ghostAddr)
	block.Header.SignatureSize = uint32(sig.Size())
	block.Header.BlockSignature = *sig
}

func (gScript *GScript) MakeScriptSigExecuteUnlock(tx *types.GhostTransaction, ghostAddr *gcrypto.GhostAddress) {
	inputParam := gScript.MakeInputParam(tx.SerializeToByte(), ghostAddr)
	for i := range tx.Body.Vin {
		tx.Body.Vin[i].ScriptSig = inputParam
		tx.Body.Vin[i].ScriptSize = uint32(len(tx.Body.Vin[i].ScriptSig))
	}
}

func (gScript *GScript) MakeInputParam(buf []byte, myAddress *gcrypto.GhostAddress) []byte {
	sig := makeSignature(buf, myAddress)
	sigBuf := sig.SerializeToByte()
	scriptBuf := new(bytes.Buffer)
	binary.Write(scriptBuf, binary.LittleEndian, OP_PUSHSIG)
	binary.Write(scriptBuf, binary.LittleEndian, byte(len(sigBuf)))
	resultBuf := append(scriptBuf.Bytes(), sigBuf...)
	return resultBuf
}

func makeSignature(buf []byte, myAddress *gcrypto.GhostAddress) *types.SigHash {
	signPack := gcrypto.Signer(buf, myAddress)
	r, s := signPack.R.Bytes(), signPack.S.Bytes()
	sig := &types.SigHash{
		RBuf:          r,
		SBuf:          s,
		PubKey:        myAddress.GetSignPubKey(),
		PubKeySize:    byte(len(myAddress.GetSignPubKey())),
		RSize:         byte(len(r)),
		SSize:         byte(len(s)),
		SignatureType: SIGHASH_ALL,
	}
	sig.SignatureSize = byte(sig.SigSize())
	return sig
}

func MakeLockScriptOut(ToAddr []byte) []byte {
	scriptBuf := new(bytes.Buffer)
	lockOutputScript(scriptBuf, ToAddr)
	binary.Write(scriptBuf, binary.LittleEndian, OP_PAY)
	binary.Write(scriptBuf, binary.LittleEndian, OP_RETURN)
	return scriptBuf.Bytes()
}

func lockOutputScript(scriptBuf *bytes.Buffer, ToAddr []byte) {
	toAddrUint8 := make([]uint8, len(ToAddr))
	copy(toAddrUint8[:], ToAddr[:])
	binary.Write(scriptBuf, binary.LittleEndian, OP_DUP)
	binary.Write(scriptBuf, binary.LittleEndian, OP_HASH160)
	binary.Write(scriptBuf, binary.LittleEndian, OP_PUSH)
	binary.Write(scriptBuf, binary.LittleEndian, uint8(len(ToAddr)))
	binary.Write(scriptBuf, binary.LittleEndian, toAddrUint8)
	binary.Write(scriptBuf, binary.LittleEndian, OP_EQUALVERIFY)
	binary.Write(scriptBuf, binary.LittleEndian, OP_CHECKSIG)
}

func MakeRootAccount(ToAddr []byte, Nickname string) []byte {
	nickname := []byte(Nickname)
	nickBuf := make([]uint8, len(nickname))
	copy(nickBuf[:], nickname[:])
	scriptBuf := new(bytes.Buffer)
	lockOutputScript(scriptBuf, ToAddr)
	binary.Write(scriptBuf, binary.LittleEndian, OP_PUSH)
	binary.Write(scriptBuf, binary.LittleEndian, uint8(len(nickBuf)))
	binary.Write(scriptBuf, binary.LittleEndian, nickBuf)
	binary.Write(scriptBuf, binary.LittleEndian, OP_RETURN)
	return scriptBuf.Bytes()
}
