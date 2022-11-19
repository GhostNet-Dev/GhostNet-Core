package types

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
)

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

type GhostTrasaction struct {
	TxId ghostBytes.HashBytes
	Body TxBody
}
