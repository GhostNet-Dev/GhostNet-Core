package gvm

import (
	"bytes"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"

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
	if exec.Regs.Stack.Count() < 1 {
		return false
	}
	dupData := exec.Regs.Stack.Peek()
	exec.Regs.Stack.Push(dupData)
	return true
}

type OpHash160 struct {
	Regs *GVMRegister
}

func (exec *OpHash160) ExcuteOp(param interface{}) bool {
	if exec.Regs.Stack.Count() < 1 {
		return false
	}
	pubKey := exec.Regs.Stack.Pop().([]byte)

	publicSHA256 := sha256.Sum256(pubKey) // Public key를 SHA-256으로 해싱

	// RIPEMD-160으로 다시 해싱
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	exec.Regs.Stack.Push(publicRIPEMD160)

	return true
}

type OpEqualVerify struct {
	Regs *GVMRegister
}

func (exec *OpEqualVerify) ExcuteOp(param interface{}) bool {
	if exec.Regs.Stack.Count() < 3 {
		return false
	}
	src := exec.Regs.Stack.Pop().([]byte)
	dst := exec.Regs.Stack.Pop().([]byte)
	if bytes.Compare(src, dst) == 0 {
		return true
	}
	return false
}

type OpCheckSig struct {
	Regs *GVMRegister
}

func (exec *OpCheckSig) ExcuteOp(param interface{}) bool {
	if exec.Regs.Stack.Count() < 2 {
		return false
	}
	buf := param.([]byte)
	pubKey := exec.Regs.Stack.Pop().([]byte)
	sigRS := exec.Regs.Stack.Pop().([]byte)
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
	exec.Regs.Stack.Push(param)
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
	exec.Regs.Stack.Push(sigRS)
	exec.Regs.Stack.Push(sigHash.PubKey)
	return true
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
	if exec.Regs.Stack.Count() < 2 {
		return false
	}
	exec.Regs.R0 = exec.Regs.Stack.Pop().([]byte)
	exec.Regs.R1 = exec.Regs.Stack.Pop().([]byte)
	return true
}

type OpMapping struct {
	Regs *GVMRegister
}

func (exec *OpMapping) ExcuteOp(param interface{}) bool {
	if exec.Regs.Stack.Count() < 2 {
		return false
	}
	exec.Regs.R0 = exec.Regs.Stack.Pop().([]byte)
	exec.Regs.R1 = exec.Regs.Stack.Pop().([]byte)
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
	return false
}
