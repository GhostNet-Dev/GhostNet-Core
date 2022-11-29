package types

import (
	"unsafe"

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

func (header *GhostNetDataBlockHeader) Size() uint32 {
	return uint32(unsafe.Sizeof(header.Id)) +
		uint32(unsafe.Sizeof(header.Version)) + ghostBytes.HashSize*2 +
		uint32(unsafe.Sizeof(header.Nonce)) +
		uint32(unsafe.Sizeof(header.TransactionCount))
}

func (dataBlock *GhostNetDataBlock) Size() uint32 {
	var txSize uint32 = 0
	for _, tx := range dataBlock.Transaction {
		txSize += tx.Size()
	}
	return dataBlock.Header.Size() + txSize
}
