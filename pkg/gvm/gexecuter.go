package gvm

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GExecuter interface {
	ExcuteOp(interface{}) bool
}

type OpDuplication struct {
	Regs *GVMRegister
}

func (exec *OpDuplication) ExcuteOp(param interface{}) bool {
	if exec.Regs.stack.Count() < 1 {
		return false
	}
	dupData := exec.Regs.stack.Peek()
	exec.Regs.stack.Push(dupData)
	return true
}

type OpHash160 struct {
	Regs *GVMRegister
}

func (exec *OpHash160) ExcuteOp(param interface{}) bool {
	if exec.Regs.stack.Count() < 1 {
		return false
	}
	pubKey := exec.Regs.stack.Pop().([]byte)

	publicSHA256 := sha256.Sum256(pubKey) // Public key를 SHA-256으로 해싱

	// RIPEMD-160으로 다시 해싱
	RIPEMD160Hasher := crypto.RIPEMD160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	exec.Regs.stack.Push(publicRIPEMD160)

	return false
}

type OpEqualVerify struct {
	Regs *GVMRegister
}

func (exec *OpEqualVerify) ExcuteOp(param interface{}) bool {
	if exec.Regs.stack.Count() < 3 {
		return false
	}
	src := exec.Regs.stack.Pop().([]byte)
	dst := exec.Regs.stack.Pop().([]byte)
	if bytes.Compare(src, dst) == 0 {
		return true
	}
	return false
}

type OpCheckSig struct {
	Regs *GVMRegister
}

func (exec *OpCheckSig) ExcuteOp(param interface{}) bool {
	if exec.Regs.stack.Count() < 2 {
		return false
	}
	buf := param.([]byte)
	pubKey := exec.Regs.stack.Pop().([]byte)
	sigRS := exec.Regs.stack.Pop().([]byte)
	sig := gcrypto.SignaturePackage{
		PubKey:    pubKey,
		Text:      buf,
		Signature: sigRS,
	}

	return gcrypto.SignVerify(&sig)
}

type OpPush struct {
	Regs *GVMRegister
}

func (exec *OpPush) ExcuteOp(param interface{}) bool {
	exec.Regs.stack.Push(param)
	return true
}

type OpPushSig struct {
	Regs *GVMRegister
}

func (exec *OpPushSig) ExcuteOp(param interface{}) bool {
	buf := param.([]byte)
	sigHash := types.SigHash{}
	sigHash.DeserializeSigHashFromByte(buf)
	sigRS := append(sigHash.RBuf, sigHash.SBuf...)
	exec.Regs.stack.Push(sigRS)
	exec.Regs.stack.Push(sigHash.PubKey)
	return false
}

type OpPay struct {
	Regs *GVMRegister
}

func (exec *OpPay) ExcuteOp(param interface{}) bool {
	return true
}

type OpTransferToken struct {
	Regs *GVMRegister
}

func (exec *OpTransferToken) ExcuteOp(param interface{}) bool {
	if exec.Regs.stack.Count() < 2 {
		return false
	}
	exec.Regs.r0 = exec.Regs.stack.Pop().([]byte)
	exec.Regs.r1 = exec.Regs.stack.Pop().([]byte)
	return true
}

type OpMapping struct {
	Regs *GVMRegister
}

func (exec *OpMapping) ExcuteOp(param interface{}) bool {
	if exec.Regs.stack.Count() < 2 {
		return false
	}
	exec.Regs.r0 = exec.Regs.stack.Pop().([]byte)
	exec.Regs.r1 = exec.Regs.stack.Pop().([]byte)
	return true
}

type OpCreateToken struct {
	Regs *GVMRegister
}

func (exec *OpCreateToken) ExcuteOp(param interface{}) bool {
	return true
}

type OpReturn struct {
	Regs *GVMRegister
}

func (exec *OpReturn) ExcuteOp(param interface{}) bool {
	return true
}
