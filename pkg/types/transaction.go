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
	TxTypeCoinTransfer TxOutputType = 1 // bitcoin
	TxTypeDataStore    TxOutputType = 2 // store user data
	TxTypeFSRoot       TxOutputType = 3 // create user
	TxTypeContract     TxOutputType = 4 // need change by code
	TxTypeShare        TxOutputType = 5 // ?
	TxTypeScript       TxOutputType = 6 // store gscript or glambda script
	TxTypeScriptStore  TxOutputType = 7 // store data by script
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

type GScript struct {
	Version uint32
	Type    uint32
	Param   string
	Script  []byte
}

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

func MakeTxOutputFromOutputParam(outputParam *NextOutputParam) *TxOutput {
	return &TxOutput{
		Addr:         outputParam.RecvAddr,
		BrokerAddr:   outputParam.Broker,
		Value:        outputParam.TransferCoin,
		Type:         outputParam.TxType,
		ScriptSize:   uint32(len(outputParam.OutputScript)),
		ScriptPubKey: outputParam.OutputScript,
		ScriptExSize: uint32(len(outputParam.OutputScriptEx)),
		ScriptEx:     outputParam.OutputScriptEx,
	}
}

func MakeTxInputFromOutputParam(outputParam *PrevOutputParam) *TxInput {
	return &TxInput{
		PrevOut:    outputParam.VOutPoint,
		Sequence:   0xFFFFFFFF,
		ScriptSize: outputParam.Vout.ScriptSize,
		ScriptSig:  outputParam.Vout.ScriptPubKey, // 서명후 새로 생성된 서명으로 교체된다.
	}
}

func MakeEmptyInput() *TxInput {
	dummyBuf4 := make([]byte, 4)
	dummyHash := make([]byte, gbytes.HashSize)
	return &TxInput{
		PrevOut: TxOutPoint{
			TxId: dummyHash,
		},
		Sequence:   0xFFFFFFFF,
		ScriptSize: uint32(len(dummyBuf4)),
		ScriptSig:  dummyBuf4,
	}
}
