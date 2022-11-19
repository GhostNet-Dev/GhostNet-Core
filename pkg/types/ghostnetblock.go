package types

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	gvm "github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
)

type GhostNetBlock struct {
	Header   GhostNetBlockHeader
	HeaderEx GhostNetBlockHeaderEx
}

type GhostNetBlockHeader struct {
	Id                      uint32
	Version                 uint32
	PreviousBlockHeaderHash ghostBytes.HashBytes
	MerkleRoot              ghostBytes.HashBytes
	DataBlockHeaderHash     ghostBytes.HashBytes
	TimeStamp               uint64
	Bits                    uint32
	Nonce                   uint32
	TransactionCount        uint32
}

type GhostNetBlockHeaderEx struct {
	HeaderSize     uint32
	SignatureSize  uint32
	BlockSignature gvm.SigHash
	Alice          GhostTrasaction
	Transaction    GhostTrasaction
}
