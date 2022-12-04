package gvm

import (
	"bytes"
	"encoding/binary"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/crypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GScript struct {
}

func NewGScript() *GScript {
	return &GScript{}
}

func (gScript *GScript) MakeScriptSigExecuteUnlock(tx *types.GhostTransaction, ghostAddr *crypto.GhostAddress) {
	inputParam := gScript.MakeInputParam(tx.SerializeToByte(), ghostAddr)
	for _, input := range tx.Body.Vin {
		input.ScriptSig = inputParam
		input.ScriptSize = uint32(len(input.ScriptSig))
	}
}

func (gScript *GScript) MakeInputParam(buf []byte, myAddress *crypto.GhostAddress) []byte {
	sig := gScript.makeSignature(buf, myAddress)
	sigBuf := sig.SerializeToByte()
	scriptBuf := new(bytes.Buffer)
	binary.Write(scriptBuf, binary.LittleEndian, OP_PUSHSIG)
	binary.Write(scriptBuf, binary.LittleEndian, byte(len(sigBuf)))
	resultBuf := append(scriptBuf.Bytes(), sigBuf...)
	return resultBuf
}

func (gScript *GScript) makeSignature(buf []byte, myAddress *crypto.GhostAddress) *types.SigHash {
	signPack := crypto.Signer(buf, myAddress)
	r, s := signPack.R.Bytes(), signPack.S.Bytes()
	sig := &types.SigHash{
		RBuf:          r,
		SBuf:          s,
		PubKey:        myAddress.PubKey,
		PubKeySize:    byte(len(myAddress.PubKey)),
		RSize:         byte(len(r)),
		SSize:         byte(len(s)),
		SignatureType: SIGHASH_ALL,
	}
	sig.SignatureSize = byte(sig.SigSize())
	return sig
}

func (gScript *GScript) MakeLockScriptOut(ToAddr []byte) []byte {
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
