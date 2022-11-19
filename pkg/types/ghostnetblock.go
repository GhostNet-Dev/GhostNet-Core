package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	gvm "github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
)

type GhostNetBlock struct {
	Header      GhostNetBlockHeader
	HeaderEx    GhostNetBlockHeaderEx
	Alice       GhostTrasaction
	Transaction GhostTrasaction
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
}

func (header GhostNetBlockHeader) GetHashKey() (key []byte) {
	var buf bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&buf)
	err := enc.Encode((header))
	if err != nil {
		log.Fatal("encode error:", err)
	}
	hash := sha256.New()

	hash.Write(buf.Bytes())

	return hash.Sum((nil))
}
