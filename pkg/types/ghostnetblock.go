package types

import (
	"crypto/sha256"
	"unsafe"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	gvm "github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	mems "github.com/traherom/memstream"
)

type PairedBlock struct {
	Block     GhostNetBlock
	DataBlock GhostNetDataBlock
}

type GhostNetBlock struct {
	Header      GhostNetBlockHeader   `json:"Header"`
	HeaderEx    GhostNetBlockHeaderEx `json:"HeaderEx"`
	Alice       GhostTransaction      `json:"Alice"`
	Transaction []GhostTransaction    `json:"Transaction"`
}

type GhostNetBlockHeader struct {
	Id                      uint32               `json:"Id"`
	Version                 uint32               `json:"Version"`
	PreviousBlockHeaderHash ghostBytes.HashBytes `json:"PreviousBlockHeaderHash"`
	MerkleRoot              ghostBytes.HashBytes `json:"MerkleRoot"`
	DataBlockHeaderHash     ghostBytes.HashBytes `json:"DataBlockHeaderHash"`
	TimeStamp               uint64               `json:"TimeStamp"`
	Bits                    uint32               `json:"Bits"`
	Nonce                   uint32               `json:"Nonce"`
	TransactionCount        uint32               `json:"TransactionCount"`
}

type GhostNetBlockHeaderEx struct {
	HeaderSize     uint32      `json:"HeaderSize"`
	SignatureSize  uint32      `json:"SignatureSize"`
	BlockSignature gvm.SigHash `json:"BlockSignature"`
}

func (header *GhostNetBlockHeader) Size() uint32 {
	return uint32(unsafe.Sizeof(header.Id)) +
		uint32(unsafe.Sizeof(header.Version)) + ghostBytes.HashSize*3 +
		uint32(unsafe.Sizeof(header.TimeStamp)) +
		uint32(unsafe.Sizeof(header.Bits)) +
		uint32(unsafe.Sizeof(header.Nonce)) +
		uint32(unsafe.Sizeof(header.TransactionCount))
}

func (headerEx *GhostNetBlockHeaderEx) Size() uint32 {
	return uint32(unsafe.Sizeof(headerEx.HeaderSize)) +
		uint32(unsafe.Sizeof(headerEx.SignatureSize)) + uint32(headerEx.BlockSignature.Size())
}

func (block GhostNetBlock) GetHashKey() []byte {
	size := block.Header.Size()
	stream := mems.NewCapacity(int(size))
	hash := sha256.New()
	hash.Write(stream.Bytes())
	return hash.Sum(nil)
}
