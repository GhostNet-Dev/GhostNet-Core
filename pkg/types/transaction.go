package types

import (
	"unsafe"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
)

const (
	TxIdSize      = ghostBytes.HashSize
	PublicKeySize = 25
	DummySize     = 4
	DataHash      = ghostBytes.HashSize
)

type PrevOutputPackage struct {
	TxType    uint32
	VOutPoint TxOutPoint
	Vout      TxOutput
}

type NewOutputPackage struct {
	TxType       uint32
	RecvAddr     ghostBytes.HashBytes
	Broker       ghostBytes.HashBytes
	OutputScript []byte
	TransferCoin int64
}

type TxOutPoint struct {
	TxId       ghostBytes.HashBytes
	TxOutIndex uint32
}

type TxInput struct {
	PrevOut    TxOutPoint
	Sequence   uint32
	ScriptSize uint32
	ScriptSig  []byte
}

type TxOutput struct {
	Addr         ghostBytes.HashBytes
	BrokerAddr   ghostBytes.HashBytes
	Value        uint64
	ScriptSize   uint32
	ScriptPubKey []byte
}

type TxBody struct {
	VinCounter  uint32
	Vin         []TxInput
	VoutCounter uint32
	Vout        []TxOutput
	Nonce       uint32
	LockTime    uint32
}

type GhostTransaction struct {
	TxId ghostBytes.HashBytes
	Body TxBody
}

func (txOutPoint *TxOutPoint) Size() uint32 {
	return uint32(unsafe.Sizeof(txOutPoint.TxOutIndex)) + ghostBytes.HashSize
}

func (input *TxInput) Size() uint32 {
	return input.PrevOut.Size() +
		uint32(unsafe.Sizeof(input.Sequence)) +
		uint32(unsafe.Sizeof(input.ScriptSize)) + input.ScriptSize
}

func (output *TxOutput) Size() uint32 {
	return uint32(ghostBytes.HashSize) + //address
		uint32(ghostBytes.HashSize) + // brokeraddress
		uint32(unsafe.Sizeof(output.Value)) +
		uint32(unsafe.Sizeof(output.ScriptSize)) +
		output.ScriptSize
}

func (body *TxBody) Size() uint32 {
	var size uint32 = 0
	if body.VinCounter > 0 {
		for _, vin := range body.Vin {
			size += vin.Size()
		}
	}
	if body.VoutCounter > 0 {
		for _, vout := range body.Vout {
			size += vout.Size()
		}
	}
	return uint32(unsafe.Sizeof(body.VinCounter)) +
		uint32(unsafe.Sizeof(body.VoutCounter)) +
		uint32(unsafe.Sizeof(body.Nonce)) +
		uint32(unsafe.Sizeof(body.LockTime)) +
		size
}

func (tx *GhostTransaction) Size() uint32 {
	return tx.Body.Size() + uint32(ghostBytes.HashSize)
}
