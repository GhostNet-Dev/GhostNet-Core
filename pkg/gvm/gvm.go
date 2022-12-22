package gvm

import (
	"bytes"
	"encoding/binary"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/container"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GFuncParam struct {
	InputSig      []byte
	ScriptPubbKey []byte
	TxType        types.TxOutputType
}

type GVMRegister struct {
	stack container.Stack
	r0    []byte
	r1    []byte
}

type GVM struct {
	Regs *GVMRegister
	exec map[uint16]GExecuter
}

func NewGVM() *GVM {
	regs := &GVMRegister{}
	gvm := GVM{
		Regs: regs,
		exec: map[uint16]GExecuter{
			OP_DUP:            &OpDuplication{Regs: regs},
			OP_HASH160:        &OpHash160{Regs: regs},
			OP_EQUALVERIFY:    &OpEqualVerify{Regs: regs},
			OP_CHECKSIG:       &OpCheckSig{Regs: regs},
			OP_PUSH:           &OpPush{Regs: regs},
			OP_PUSHTOKEN:      &OpPush{Regs: regs},
			OP_PUSHSIG:        &OpPushSig{Regs: regs},
			OP_PAY:            &OpPay{Regs: regs},
			OP_TRANSFER_TOKEN: &OpTransferToken{Regs: regs},
			OP_MAPPING:        &OpMapping{Regs: regs},
			OP_CREATE_TOKEN:   &OpCreateToken{Regs: regs},
			OP_RETURN:         &OpReturn{Regs: regs},
		},
	}

	return &gvm
}

func (gvm *GVM) ExecuteGFunction(buf []byte, params []GFuncParam) bool {
	verify := false
	for _, param := range params {
		gvm.Clear()

		if verify = gvm.PushParam(param.InputSig); verify == false {
			return false
		}

		if verify = gvm.ExecuteScript(param.ScriptPubbKey, buf); verify == false {
			return false
		}
	}

	return verify
}

func (gvm *GVM) PushParam(param []byte) bool {
	var op uint16
	result := false
	byteBuf := bytes.NewBuffer(param)
	for {
		if err := binary.Read(byteBuf, binary.LittleEndian, &op); err != nil {
			break
		}

		exec, ok := gvm.exec[op]
		if ok == false {
			break
		}

		if op == OP_PUSH || op == OP_PUSHSIG {
			var length byte
			binary.Read(byteBuf, binary.LittleEndian, &length)
			pushData := make([]byte, length)
			binary.Read(byteBuf, binary.LittleEndian, pushData)
			result = exec.ExcuteOp(pushData)
		}
	}
	return result
}

func (gvm *GVM) ExecuteScript(scriptPubKey []byte, param []byte) bool {
	var op uint16
	result := true
	byteBuf := bytes.NewBuffer(scriptPubKey)
	for result {
		if err := binary.Read(byteBuf, binary.LittleEndian, &op); err != nil {
			break
		}

		exec, ok := gvm.exec[op]
		if ok == false {
			break
		}

		if op == OP_PUSH || op == OP_PUSHTOKEN {
			var length byte
			binary.Read(byteBuf, binary.LittleEndian, &length)
			pushData := make([]byte, length)
			binary.Read(byteBuf, binary.LittleEndian, pushData)
			result = exec.ExcuteOp(pushData)
		} else if op == OP_CHECKSIG {
			result = exec.ExcuteOp(param)
		} else {
			result = exec.ExcuteOp(nil)
		}
	}
	return result
}

func (gvm *GVM) Clear() {
	gvm.Regs.stack.Clear()
}
