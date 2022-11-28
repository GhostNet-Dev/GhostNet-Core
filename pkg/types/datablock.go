package types

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
)

type GhostNetDataBlock struct {
	Header      GhostNetDataBlockHeader
	Transaction []GhostDataTransaction
}

type GhostNetDataBlockHeader struct {
	Id                      uint32
	Version                 uint32
	PreviousBlockHeaderHash ghostBytes.HashBytes
	MerkleRoot              ghostBytes.HashBytes
	Nonce                   uint32
	TransactionCount        uint32
}
