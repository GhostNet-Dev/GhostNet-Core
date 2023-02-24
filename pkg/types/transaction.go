package types

import (
	"bytes"
	"crypto/sha256"
	"unsafe"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
)

type TxOutputType uint32

const ( // tx output type
	None               TxOutputType = 0
	TxTypeCoinTransfer TxOutputType = 1
	TxTypeDataTransfer TxOutputType = 2
	TxTypeFSRoot       TxOutputType = 3
	TxTypeContract     TxOutputType = 4
)

const (
	TxIdSize      = gbytes.HashSize
	PublicKeySize = 25
	DummySize     = 4
	DataHash      = gbytes.HashSize
)

const ( //tx type
	AliceTx  = 0
	NormalTx = 1
)

type PrevOutputParam struct {
	TxType    TxOutputType
	VOutPoint TxOutPoint
	Vout      TxOutput
}

type NextOutputParam struct {
	TxType         TxOutputType
	RecvAddr       gbytes.HashBytes
	Broker         gbytes.HashBytes
	OutputScript   []byte
	OutputScriptEx []byte
	TransferCoin   uint64
}

type TxOutPoint struct {
	TxId       gbytes.HashBytes
	TxOutIndex uint32
}

type TxInput struct {
	PrevOut    TxOutPoint
	Sequence   uint32
	ScriptSize uint32
	ScriptSig  []byte
}

type TxOutput struct {
	Addr         gbytes.HashBytes
	BrokerAddr   gbytes.HashBytes
	Type         TxOutputType
	Value        uint64
	ScriptSize   uint32
	ScriptPubKey []byte
	ScriptExSize uint32
	ScriptEx     []byte
}

type TxBody struct {
	InputCounter  uint32
	Vin           []TxInput
	OutputCounter uint32
	Vout          []TxOutput
	Nonce         uint32
	LockTime      uint32
}

type GhostTransaction struct {
	TxId gbytes.HashBytes
	Body TxBody
}

func (txOutPoint *TxOutPoint) Size() uint32 {
	return uint32(unsafe.Sizeof(txOutPoint.TxOutIndex)) + gbytes.HashSize
}

func (input *TxInput) Size() uint32 {
	return input.PrevOut.Size() +
		uint32(unsafe.Sizeof(input.Sequence)) +
		uint32(unsafe.Sizeof(input.ScriptSize)) + input.ScriptSize
}

func (output *TxOutput) Size() uint32 {
	return uint32(gbytes.PubKeySize) + //address
		uint32(gbytes.PubKeySize) + // brokeraddress
		uint32(unsafe.Sizeof(output.Value)) +
		uint32(unsafe.Sizeof(output.Type)) +
		uint32(unsafe.Sizeof(output.ScriptSize)) +
		uint32(unsafe.Sizeof(output.ScriptExSize)) +
		output.ScriptSize + output.ScriptExSize
}

func (body *TxBody) Size() uint32 {
	var size uint32 = 0
	if body.InputCounter > 0 {
		for _, vin := range body.Vin {
			size += vin.Size()
		}
	}
	if body.OutputCounter > 0 {
		for _, vout := range body.Vout {
			size += vout.Size()
		}
	}
	return uint32(unsafe.Sizeof(body.InputCounter)) +
		uint32(unsafe.Sizeof(body.OutputCounter)) +
		uint32(unsafe.Sizeof(body.Nonce)) +
		uint32(unsafe.Sizeof(body.LockTime)) +
		size
}

func (tx *GhostTransaction) Size() uint32 {
	return tx.Body.Size() + uint32(gbytes.HashSize)
}

func (tx *GhostTransaction) GetHashKey() []byte {
	hash := sha256.New()
	hash.Write(tx.SerializeToByte())
	return hash.Sum(nil)
}

func (tx *GhostTransaction) TxCopy() (copy GhostTransaction) {
	data := tx.SerializeToByte()
	byteBuf := bytes.NewBuffer(data)
	copy.Deserialize(byteBuf)
	return copy
}
